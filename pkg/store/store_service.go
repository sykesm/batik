// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package store

import (
	"context"

	storev1 "github.com/sykesm/batik/pkg/pb/store/v1"
	txv1 "github.com/sykesm/batik/pkg/pb/transaction/v1"
)

// StoreService implements the StoreAPIServer gRPC interface.
type StoreService struct {
	// Unnsafe has been chosed to ensure there's a compilation failure when the
	// implementation diverges from the gRPC service.
	storev1.UnsafeStoreAPIServer

	db KV
}

var _ storev1.StoreAPIServer = (*StoreService)(nil)

func NewStoreService(db KV) *StoreService {
	return &StoreService{
		db: db,
	}
}

// GetTransaction retrieves the associated transaction corresponding to the
// txid passed in the GetTransactionRequest.
func (s *StoreService) GetTransaction(ctx context.Context, req *storev1.GetTransactionRequest) (*storev1.GetTransactionResponse, error) {
	txs, err := LoadTransactions(s.db, [][]byte{req.Txid})
	if err != nil {
		return nil, err
	}

	return &storev1.GetTransactionResponse{
		Transaction: txs[0],
	}, nil
}

// PutTransaction hashes the transaction to a txid and then stores
// the encoded transaction in the backing store.
func (s *StoreService) PutTransaction(ctx context.Context, req *storev1.PutTransactionRequest) (*storev1.PutTransactionResponse, error) {
	if err := StoreTransactions(s.db, []*txv1.Transaction{req.Transaction}); err != nil {
		return nil, err
	}

	return &storev1.PutTransactionResponse{}, nil
}

// GetState retrieves the associated ResolvedState corresponding to the state reference
// passed in the GetStateRequest from the backing store indexed by a txid and
// output index that the State was originally created at in the transaction output
// list.
func (s *StoreService) GetState(ctx context.Context, req *storev1.GetStateRequest) (*storev1.GetStateResponse, error) {
	states, err := LoadStates(s.db, []*txv1.StateReference{req.StateRef})
	if err != nil {
		return nil, err
	}

	return &storev1.GetStateResponse{
		State: states[0],
	}, nil
}

// PutState stores the encoded resolved state in the backing store.
func (s *StoreService) PutState(ctx context.Context, req *storev1.PutStateRequest) (*storev1.PutStateResponse, error) {
	if err := StoreStates(s.db, []*txv1.ResolvedState{req.State}); err != nil {
		return nil, err
	}

	return &storev1.PutStateResponse{}, nil
}
