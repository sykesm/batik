// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package submitservice

import (
	"context"
	"crypto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sykesm/batik/pkg/merkle"
	txv1 "github.com/sykesm/batik/pkg/pb/tx/v1"
	"github.com/sykesm/batik/pkg/store"
	"github.com/sykesm/batik/pkg/transaction"
)

// SubmitService implements the EncodeAPIServer gRPC interface.
type SubmitService struct {
	// Unnsafe has been chosed to ensure there's a compilation failure when the
	// implementation diverges from the gRPC service.
	txv1.UnsafeSubmitAPIServer
	// hasher implements the hash algorithm used to build and validate the
	// transaction ID.
	hasher merkle.Hasher
	// kv is a reference to the key value store backing this service
	kv store.KV
}

var _ txv1.SubmitAPIServer = (*SubmitService)(nil)

// NewSubmitService creates a new instance of the SubmitService.
func NewSubmitService(kv store.KV) *SubmitService {
	return &SubmitService{
		hasher: crypto.SHA256,
		kv:     kv,
	}
}

// Submit submits a transaction for validation and commit processing.
//
// NOTE: This is an implementation for prototyping.
func (s *SubmitService) Submit(ctx context.Context, req *txv1.SubmitRequest) (*txv1.SubmitResponse, error) {
	signedTx := req.GetSignedTransaction()
	if signedTx == nil {
		return nil, status.Errorf(codes.InvalidArgument, "signed transaction was not provided")
	}
	tx := signedTx.GetTransaction()
	if tx == nil {
		return nil, status.Errorf(codes.InvalidArgument, "transaction was not provided")
	}
	itx, err := transaction.Marshal(s.hasher, tx)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	// All inputs must exist
	for _, in := range tx.Inputs {
		if in == nil {
			continue
		}
	}
	// All references must exist
	for _, ref := range tx.References {
		if ref == nil {
			continue
		}
	}

	return &txv1.SubmitResponse{Txid: itx.ID}, nil
}
