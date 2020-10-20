// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package transaction

import (
	"crypto"
	"strconv"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	txv1 "github.com/sykesm/batik/pkg/pb/tx/v1"
	"github.com/sykesm/batik/pkg/protomsg"
	"github.com/sykesm/batik/pkg/store"
)

var (
	// Global prefixes.
	keyTransactions = [...]byte{0x1}
	keyStates       = [...]byte{0x2}
)

func StoreTransactions(kv store.KV, txs []*txv1.Transaction) error {
	batch := kv.NewWriteBatch()

	for _, tx := range txs {
		intTx, err := Marshal(crypto.SHA256, tx)
		if err != nil {
			return errors.WithMessage(err, "error marshaling transaction")
		}

		if err := batch.Put(transactionKey(intTx.ID), intTx.Encoded); err != nil {
			return err
		}
	}

	return errors.WithMessage(batch.Commit(), "error committing transactions batch")
}

func LoadTransactions(kv store.KV, ids [][]byte) ([]*txv1.Transaction, error) {
	result := make([]*txv1.Transaction, 0, len(ids))

	for _, id := range ids {
		payload, err := kv.Get(transactionKey(id))
		if err != nil {
			return nil, errors.WithMessagef(err, "error getting tx %x from db", id)
		}

		tx := &txv1.Transaction{}
		if err := proto.Unmarshal(payload, tx); err != nil {
			return nil, errors.WithMessagef(err, "error unmarshaling tx %x", id)
		}

		result = append(result, tx)
	}

	return result, nil
}

func StoreStates(kv store.KV, states []*txv1.ResolvedState) error {
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

func LoadStates(kv store.KV, refs []*txv1.StateReference) ([]*txv1.ResolvedState, error) {
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

// transactionKey returns a db key for a transaction
func transactionKey(txid []byte) []byte {
	return append(keyTransactions[:], txid[:]...)
}

// stateKey returns a db key for a state
func stateKey(stateRef *txv1.StateReference) []byte {
	return append(
		append(keyStates[:], stateRef.Txid[:]...),
		strconv.Itoa(int(stateRef.OutputIndex))...,
	)
}
