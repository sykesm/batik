// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package store

import (
	"strconv"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	txv1 "github.com/sykesm/batik/pkg/pb/tx/v1"
	"github.com/sykesm/batik/pkg/protomsg"
	"github.com/sykesm/batik/pkg/transaction"
)

type TransactionRepository struct {
	KV KV
}

func (t *TransactionRepository) PutTransaction(tx *transaction.Transaction) error {
	err := t.KV.Put(transactionKey(tx.ID), tx.Encoded)
	if err != nil {
		return errors.WithMessagef(err, "failed to put transaction %s", tx.ID)
	}
	return nil
}

func (t *TransactionRepository) GetTransaction(id transaction.ID) (*transaction.Transaction, error) {
	payload, err := t.KV.Get(transactionKey(id))
	if err != nil {
		return nil, errors.WithMessagef(err, "error getting tx %x from db", id)
	}
	var tx txv1.Transaction
	if err := proto.Unmarshal(payload, &tx); err != nil {
		return nil, errors.WithMessagef(err, "error unmarshaling tx %x", id)
	}
	return &transaction.Transaction{
		Tx:      &tx,
		ID:      transaction.NewID(id),
		Encoded: payload,
	}, nil
}

func StoreStates(kv KV, states []*txv1.ResolvedState) error {
	batch := kv.NewWriteBatch()

	for _, state := range states {
		encodedState, err := protomsg.MarshalDeterministic(state)
		if err != nil {
			return errors.WithMessage(err, "error marshalling resolved state")
		}

		if err := batch.Put(stateKey(&txv1.StateReference{Txid: state.Txid, OutputIndex: state.OutputIndex}), encodedState); err != nil {
			return err
		}
	}

	return errors.WithMessage(batch.Commit(), "error committing resolved states batch")
}

func LoadStates(kv KV, refs []*txv1.StateReference) ([]*txv1.ResolvedState, error) {
	result := make([]*txv1.ResolvedState, 0, len(refs))

	for _, ref := range refs {
		payload, err := kv.Get(stateKey(ref))
		if err != nil {
			return nil, errors.WithMessagef(err, "error getting state %x from db", ref)
		}

		state := &txv1.ResolvedState{
			Txid:        ref.Txid,
			OutputIndex: ref.OutputIndex,
		}
		if err := proto.Unmarshal(payload, state); err != nil {
			return nil, errors.WithMessagef(err, "error unmarshaling state for ref %x", ref)
		}

		result = append(result, state)
	}

	return result, nil
}

func ConsumeStates(kv KV, refs []*txv1.StateReference) error {
	batch := kv.NewWriteBatch()
	for _, ref := range refs {
		state, err := kv.Get(stateKey(ref))
		if err != nil {
			return err
		}
		err = batch.Delete(stateKey(ref))
		if err != nil {
			return err
		}
		err = batch.Put(consumedStateKey(ref), state)
		if err != nil {
			return err
		}
	}
	return errors.WithMessage(batch.Commit(), "error consuming states batch")
}

var (
	// Global prefixes.
	keyTransactions   = [...]byte{0x1}
	keyStates         = [...]byte{0x2}
	keyConsumedStates = [...]byte{0x3}
)

// transactionKey returns a db key for a transaction
func transactionKey(txid []byte) []byte {
	return append(keyTransactions[:], txid[:]...)
}

// TODO: Use variable length encoding for index to save ~7 bytes per key

// stateKey returns a db key for a state
func stateKey(stateRef *txv1.StateReference) []byte {
	return append(
		append(keyStates[:], stateRef.Txid[:]...),
		strconv.Itoa(int(stateRef.OutputIndex))...,
	)
}

func consumedStateKey(stateRef *txv1.StateReference) []byte {
	return append(
		append(keyConsumedStates[:], stateRef.Txid[:]...),
		strconv.Itoa(int(stateRef.OutputIndex))...,
	)
}
