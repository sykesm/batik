// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package store

import (
	"bytes"
	"context"
	"crypto"
	"strconv"
	"sync"

	sb "github.com/sykesm/batik/pkg/pb/store"
	tb "github.com/sykesm/batik/pkg/pb/transaction"
	"github.com/sykesm/batik/pkg/protomsg"
	"github.com/sykesm/batik/pkg/transaction"
	"google.golang.org/protobuf/proto"
)

// StoreService implements the StoreAPIServer gRPC interface.
type StoreService struct {
	mu sync.Locker
	db KV
}

var _ sb.StoreAPIServer = (*StoreService)(nil)

func NewStoreService(db KV) *StoreService {
	return &StoreService{
		mu: &sync.Mutex{},
		db: db,
	}
}

// GetTransaction retrieves the associated transaction corresponding to the
// txid passed in the GetTransactionRequest.
func (s *StoreService) GetTransaction(ctx context.Context, req *sb.GetTransactionRequest) (*sb.GetTransactionResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	tx := &tb.Transaction{}

	key := transactionKey(req.Txid)
	data, err := s.db.Get(key)
	if err != nil {
		return nil, err
	}

	if err := proto.Unmarshal(data, tx); err != nil {
		return nil, err
	}

	return &sb.GetTransactionResponse{
		Transaction: tx,
	}, nil
}

// PutTransaction hashes the transaction to a txid and then stores
// the encoded transaction in the backing store.
func (s *StoreService) PutTransaction(ctx context.Context, req *sb.PutTransactionRequest) (*sb.PutTransactionResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id, encoded, err := transaction.Marshal(crypto.SHA256, req.Transaction)
	if err != nil {
		return nil, err
	}

	key := transactionKey(id)
	if err := s.db.Put(key, encoded); err != nil {
		return nil, err
	}

	return &sb.PutTransactionResponse{}, nil
}

// GetState retrieves the associated ResolvedState corresponding to the state reference
// passed in the GetStateRequest from the backing store indexed by a txid and
// output index that the State was originally created at in the transaction output
// list.
func (s *StoreService) GetState(ctx context.Context, req *sb.GetStateRequest) (*sb.GetStateResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	state := &tb.ResolvedState{}

	key := stateKey(req.StateRef)
	data, err := s.db.Get(key)
	if err != nil {
		return nil, err
	}

	if err := proto.Unmarshal(data, state); err != nil {
		return nil, err
	}

	state.Txid = req.StateRef.Txid
	state.OutputIndex = req.StateRef.OutputIndex

	return &sb.GetStateResponse{
		State: state,
	}, nil
}

// PutState stores the encoded resolved state in the backing store.
func (s *StoreService) PutState(ctx context.Context, req *sb.PutStateRequest) (*sb.PutStateResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	encodedState, err := protomsg.MarshalDeterministic(req.State)
	if err != nil {
		return nil, err
	}

	stateRef := &tb.StateReference{
		Txid:        req.State.Txid,
		OutputIndex: req.State.OutputIndex,
	}

	key := stateKey(stateRef)
	if err := s.db.Put(key, encodedState); err != nil {
		return nil, err
	}

	return &sb.PutStateResponse{}, nil
}

// transactionKey returns a byte slice key of the format "tx:<txid>"
func transactionKey(txid []byte) []byte {
	keySlice := [][]byte{
		[]byte("tx"),
		txid,
	}

	return bytes.Join(keySlice, []byte(":"))
}

// stateKey returns a byte slice key of the format "state:<txid>:<tx_output_index>"
func stateKey(stateRef *tb.StateReference) []byte {
	keySlice := [][]byte{
		[]byte("state"),
		stateRef.Txid,
		[]byte(strconv.FormatUint(stateRef.OutputIndex, 10)),
	}

	return bytes.Join(keySlice, []byte(":"))
}
