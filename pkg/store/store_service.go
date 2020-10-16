// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package store

import (
	"context"

	sb "github.com/sykesm/batik/pkg/pb/store"
	tb "github.com/sykesm/batik/pkg/pb/transaction"
)

// StoreService implements the StoreAPIServer gRPC interface.
type StoreService struct {
	// Unnsafe has been chosed to ensure there's a compilation failure when the
	// implementation diverges from the gRPC service.
	sb.UnsafeStoreAPIServer

	db KV
}

var _ sb.StoreAPIServer = (*StoreService)(nil)

func NewStoreService(db KV) *StoreService {
	return &StoreService{
		db: db,
	}
}

// GetTransaction retrieves the associated transaction corresponding to the
// txid passed in the GetTransactionRequest.
func (s *StoreService) GetTransaction(ctx context.Context, req *sb.GetTransactionRequest) (*sb.GetTransactionResponse, error) {
	txs, err := LoadTransactions(s.db, [][]byte{req.Txid})
	if err != nil {
		return nil, err
	}

	return &sb.GetTransactionResponse{
		Transaction: txs[0],
	}, nil
}

// PutTransaction hashes the transaction to a txid and then stores
// the encoded transaction in the backing store.
func (s *StoreService) PutTransaction(ctx context.Context, req *sb.PutTransactionRequest) (*sb.PutTransactionResponse, error) {
	if err := StoreTransactions(s.db, []*tb.Transaction{req.Transaction}); err != nil {
		return nil, err
	}

	return &sb.PutTransactionResponse{}, nil
}

// GetState retrieves the associated ResolvedState corresponding to the state reference
// passed in the GetStateRequest from the backing store indexed by a txid and
// output index that the State was originally created at in the transaction output
// list.
func (s *StoreService) GetState(ctx context.Context, req *sb.GetStateRequest) (*sb.GetStateResponse, error) {
	states, err := LoadStates(s.db, []*tb.StateReference{req.StateRef})
	if err != nil {
		return nil, err
	}

	return &sb.GetStateResponse{
		State: states[0],
	}, nil
}

// PutState stores the encoded resolved state in the backing store.
func (s *StoreService) PutState(ctx context.Context, req *sb.PutStateRequest) (*sb.PutStateResponse, error) {
	if err := StoreStates(s.db, []*tb.ResolvedState{req.State}); err != nil {
		return nil, err
	}

	return &sb.PutStateResponse{}, nil
}
