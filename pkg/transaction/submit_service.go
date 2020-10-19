// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package transaction

import (
	"context"
	"crypto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sykesm/batik/pkg/merkle"
	txv1 "github.com/sykesm/batik/pkg/pb/transaction/v1"
)

// SubmitService implements the EncodeTransactionAPIServer gRPC interface.
type SubmitService struct {
	// Unnsafe has been chosed to ensure there's a compilation failure when the
	// implementation diverges from the gRPC service.
	txv1.UnsafeSubmitTransactionAPIServer
	// The hash algorithm used to build and validate the transaction ID.
	hasher merkle.Hasher
}

var _ txv1.SubmitTransactionAPIServer = (*SubmitService)(nil)

// NewSubmitService creates a new instance of the SubmitService.
func NewSubmitService() *SubmitService {
	return &SubmitService{
		hasher: crypto.SHA256,
	}
}

// SubmitTransaction submits a transaction for validation and commit processing.
//
// NOTE: This is an implementation for prototyping.
func (s *SubmitService) SubmitTransaction(ctx context.Context, req *txv1.SubmitTransactionRequest) (*txv1.SubmitTransactionResponse, error) {
	inTx := req.GetTransaction()
	if inTx == nil {
		return nil, status.Errorf(codes.InvalidArgument, "transaction was not provided")
	}
	tx, err := Marshal(s.hasher, inTx)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &txv1.SubmitTransactionResponse{Txid: tx.ID}, nil
}
