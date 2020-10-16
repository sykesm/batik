// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package transaction

import (
	"context"
	"crypto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sykesm/batik/pkg/merkle"
	tb "github.com/sykesm/batik/pkg/pb/transaction"
)

// SubmitService implements the EncodeTransactionAPIServer gRPC interface.
type SubmitService struct {
	// Unnsafe has been chosed to ensure there's a compilation failure when the
	// implementation diverges from the gRPC service.
	tb.UnsafeSubmitTransactionAPIServer
	// The hash algorithm used to build and validate the transaction ID.
	hasher merkle.Hasher
}

var _ tb.SubmitTransactionAPIServer = (*SubmitService)(nil)

// NewSubmitService creates a new instance of the SubmitService.
func NewSubmitService() *SubmitService {
	return &SubmitService{
		hasher: crypto.SHA256,
	}
}

// SubmitTransaction submits a transaction for validation and commit processing.
//
// NOTE: This is an implementation for prototyping.
func (s *SubmitService) SubmitTransaction(ctx context.Context, req *tb.SubmitTransactionRequest) (*tb.SubmitTransactionResponse, error) {
	inTx := req.GetTransaction()
	if inTx == nil {
		return nil, status.Errorf(codes.InvalidArgument, "transaction was not provided")
	}
	tx, err := Marshal(s.hasher, inTx)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	return &tb.SubmitTransactionResponse{Txid: tx.ID}, nil
}
