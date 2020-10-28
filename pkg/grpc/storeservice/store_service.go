// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package storeservice

import (
	"context"
	"crypto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sykesm/batik/pkg/merkle"
	storev1 "github.com/sykesm/batik/pkg/pb/store/v1"
	txv1 "github.com/sykesm/batik/pkg/pb/tx/v1"
	"github.com/sykesm/batik/pkg/store"
	"github.com/sykesm/batik/pkg/transaction"
)

type Repository interface {
	PutTransaction(*transaction.Transaction) error
	GetTransaction(transaction.ID) (*transaction.Transaction, error)
	PutState(*transaction.State) error
	GetState(transaction.StateID) (*transaction.State, error)
}

// StoreService implements the StoreAPIServer gRPC interface.
type StoreService struct {
	// Unnsafe has been chosed to ensure there's a compilation failure when the
	// implementation diverges from the gRPC service.
	storev1.UnsafeStoreAPIServer

	db     store.KV
	hasher merkle.Hasher
	repo   Repository
}

var _ storev1.StoreAPIServer = (*StoreService)(nil)

func NewStoreService(db store.KV) *StoreService {
	return &StoreService{
		db:     db,
		hasher: crypto.SHA256,
		repo:   &store.TransactionRepository{KV: db},
	}
}

// GetTransaction retrieves the associated transaction corresponding to the
// txid passed in the GetTransactionRequest.
func (s *StoreService) GetTransaction(ctx context.Context, req *storev1.GetTransactionRequest) (*storev1.GetTransactionResponse, error) {
	tx, err := s.repo.GetTransaction(req.Txid)
	if err != nil {
		return nil, err
	}

	return &storev1.GetTransactionResponse{
		Transaction: tx.Tx,
	}, nil
}

// PutTransaction hashes the transaction to a txid and then stores
// the encoded transaction in the backing store.
func (s *StoreService) PutTransaction(ctx context.Context, req *storev1.PutTransactionRequest) (*storev1.PutTransactionResponse, error) {
	tx, err := transaction.New(s.hasher, req.Transaction)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := s.repo.PutTransaction(tx); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &storev1.PutTransactionResponse{}, nil
}

// GetState retrieves the associated ResolvedState corresponding to the state reference
// passed in the GetStateRequest from the backing store indexed by a txid and
// output index that the State was originally created at in the transaction output
// list.
func (s *StoreService) GetState(ctx context.Context, req *storev1.GetStateRequest) (*storev1.GetStateResponse, error) {
	stateID := transaction.StateID{TxID: req.StateRef.Txid, OutputIndex: req.StateRef.OutputIndex}
	state, err := s.repo.GetState(stateID)
	if err != nil {
		return nil, err
	}

	return &storev1.GetStateResponse{
		State: &txv1.ResolvedState{
			Txid:        state.ID.TxID,
			OutputIndex: state.ID.OutputIndex,
			Info:        state.StateInfo,
			State:       state.Data,
		},
	}, nil
}

// PutState stores the encoded resolved state in the backing store.
func (s *StoreService) PutState(ctx context.Context, req *storev1.PutStateRequest) (*storev1.PutStateResponse, error) {
	state := &transaction.State{
		ID:        transaction.StateID{TxID: req.State.Txid, OutputIndex: req.State.OutputIndex},
		StateInfo: req.State.Info,
		Data:      req.State.State,
	}
	if err := s.repo.PutState(state); err != nil {
		return nil, err
	}

	return &storev1.PutStateResponse{}, nil
}
