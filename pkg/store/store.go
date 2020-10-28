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
		ID:      transaction.NewID(id),
		Tx:      &tx,
		Encoded: payload,
	}, nil
}

func (t *TransactionRepository) PutState(state *transaction.State) error {
	info, err := protomsg.MarshalDeterministic(state.StateInfo)
	if err != nil {
		return errors.WithMessage(err, "error marshalling state info")
	}

	batch := t.KV.NewWriteBatch()
	if err := batch.Put(stateKey(state.ID), state.Data); err != nil {
		return err
	}
	if err := batch.Put(stateInfoKey(state.ID), info); err != nil {
		return err
	}

	return errors.WithMessage(batch.Commit(), "error committing resolved states batch")
}

func GetState(kv KV, stateID transaction.StateID) (*txv1.ResolvedState, error) {
	infoPayload, err := kv.Get(stateInfoKey(stateID))
	if err != nil {
		return nil, errors.WithMessagef(err, "error getting state info for %s from db", stateID)
	}
	var stateInfo txv1.StateInfo
	if err := proto.Unmarshal(infoPayload, &stateInfo); err != nil {
		return nil, errors.WithMessagef(err, "error unmarshaling state info for ref %s", stateID)
	}

	payload, err := kv.Get(stateKey(stateID))
	if err != nil {
		return nil, errors.WithMessagef(err, "error getting state %s from db", stateID)
	}

	state := &txv1.ResolvedState{
		Txid:        stateID.TxID,
		OutputIndex: stateID.OutputIndex,
		Info:        &stateInfo,
		State:       payload,
	}

	return state, nil
}

func ConsumeStates(kv KV, refs []*txv1.StateReference) error {
	batch := kv.NewWriteBatch()
	for _, ref := range refs {
		state, err := kv.Get(stateKey(transaction.StateID{TxID: ref.Txid, OutputIndex: ref.OutputIndex}))
		if err != nil {
			return err
		}
		err = kv.Delete(stateKey(transaction.StateID{TxID: ref.Txid, OutputIndex: ref.OutputIndex}))
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
	keyStateInfos     = [...]byte{0x3}
	keyConsumedStates = [...]byte{0x4}
)

// transactionKey returns a db key for a transaction
func transactionKey(txid []byte) []byte {
	return append(keyTransactions[:], txid[:]...)
}

// stateKey returns a db key for a state
func stateKey(id transaction.StateID) []byte {
	return append(
		append(keyStates[:], id.TxID[:]...),
		strconv.Itoa(int(id.OutputIndex))...,
	)
}

// stateKey returns a db key for a state
func stateInfoKey(id transaction.StateID) []byte {
	return append(
		append(keyStateInfos[:], id.TxID[:]...),
		strconv.Itoa(int(id.OutputIndex))...,
	)
}

func consumedStateKey(stateRef *txv1.StateReference) []byte {
	return append(
		append(keyConsumedStates[:], stateRef.Txid[:]...),
		strconv.Itoa(int(stateRef.OutputIndex))...,
	)
}
