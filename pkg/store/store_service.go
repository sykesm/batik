// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package store

import (
	"bytes"
	"context"
	"crypto"
	"fmt"
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
	sync.RWMutex
	Db *LevelDBKV
}

var _ sb.StoreAPIServer = (*StoreService)(nil)

// GetTransaction retrieves the associated transaction corresponding to the
// txid passed in the GetTransactionRequest.
func (s *StoreService) GetTransaction(ctx context.Context, req *sb.GetTransactionRequest) (*sb.GetTransactionResponse, error) {
	s.Lock()
	defer s.Unlock()

	tx := &tb.Transaction{}

	key := transactionKey(req.Txid)
	data, err := s.Db.Get(key)
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

// PutTransaction first verifies that the transaction can be hashed to the
// provided txid and then stores the transaction in the backing store.
func (s *StoreService) PutTransaction(ctx context.Context, req *sb.PutTransactionRequest) (*sb.PutTransactionResponse, error) {
	s.Lock()
	defer s.Unlock()

	id, err := transaction.ID(crypto.SHA256, req.Transaction)
	if err != nil {
		return nil, err
	}

	if bytes.Compare(id, req.Txid) != 0 {
		return nil, fmt.Errorf("request txid [%x] does not match hashed tx: [%x]", req.Txid, id)
	}

	encodedTx, err := protomsg.MarshalDeterministic(req.Transaction)
	if err != nil {
		return nil, err
	}

	key := transactionKey(req.Txid)
	if err := s.Db.Put(key, encodedTx); err != nil {
		return nil, err
	}

	return &sb.PutTransactionResponse{}, nil
}

// GetState retrieves the associated ResolvedState corresponding to the state reference
// passed in the GetStateRequest from the backing store indexed by a txid and
// output index that the State was originally created at in the transaction output
// list.
func (s *StoreService) GetState(ctx context.Context, req *sb.GetStateRequest) (*sb.GetStateResponse, error) {
	s.Lock()
	defer s.Unlock()

	state := &tb.ResolvedState{}

	key := stateKey(req.StateRef)
	data, err := s.Db.Get(key)
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
