// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package store

import (
	"crypto"
	"strconv"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	txv1 "github.com/sykesm/batik/pkg/pb/tx/v1"
	"github.com/sykesm/batik/pkg/protomsg"
	"github.com/sykesm/batik/pkg/transaction"
)

// TODO: Determine how to model the hasher required to restore a transaction.
// TODO: Standarize on binary mashaling and unmarshaling to remove proto

type TransactionRepository struct {
	kv KV
}

func NewRepository(kv KV) *TransactionRepository {
	return &TransactionRepository{
		kv: kv,
	}
}

func (t *TransactionRepository) PutTransaction(tx *transaction.Transaction) error {
	err := t.kv.Put(transactionKey(tx.ID), tx.Encoded)
	if err != nil {
		return errors.WithMessagef(err, "failed to put transaction %s", tx.ID)
	}
	return nil
}

func (t *TransactionRepository) GetTransaction(id transaction.ID) (*transaction.Transaction, error) {
	payload, err := t.kv.Get(transactionKey(id))
	if err != nil {
		return nil, errors.WithMessagef(err, "error getting tx %x from db", id)
	}
	tx, err := transaction.NewFromBytes(crypto.SHA256, payload)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to reconstruct the transaction")
	}
	if !tx.ID.Equals(id) {
		return nil, errors.Errorf("requested transaction %s but retrieved %s", id, tx.ID)
	}
	return tx, nil
}

func (t *TransactionRepository) PutState(state *transaction.State) error {
	info, err := protomsg.MarshalDeterministic(state.StateInfo)
	if err != nil {
		return errors.WithMessage(err, "error marshalling state info")
	}

	batch := t.kv.NewWriteBatch()
	if err := batch.Put(stateKey(state.ID), state.Data); err != nil {
		return err
	}
	if err := batch.Put(stateInfoKey(state.ID), info); err != nil {
		return err
	}

	return errors.WithMessage(batch.Commit(), "error committing resolved states batch")
}

func (t *TransactionRepository) GetState(stateID transaction.StateID) (*transaction.State, error) {
	infoPayload, err := t.kv.Get(stateInfoKey(stateID))
	if err != nil {
		return nil, errors.WithMessagef(err, "error getting state info for %s from db", stateID)
	}
	var stateInfo txv1.StateInfo
	if err := proto.Unmarshal(infoPayload, &stateInfo); err != nil {
		return nil, errors.WithMessagef(err, "error unmarshaling state info for ref %s", stateID)
	}

	payload, err := t.kv.Get(stateKey(stateID))
	if err != nil {
		return nil, errors.WithMessagef(err, "error getting state %s from db", stateID)
	}

	state := &transaction.State{
		ID:        stateID,
		StateInfo: &stateInfo,
		Data:      payload,
	}

	return state, nil
}

func (t *TransactionRepository) ConsumeStates(stateIDs ...transaction.StateID) error {
	batch := t.kv.NewWriteBatch()
	for _, id := range stateIDs {
		state, err := t.kv.Get(stateKey(id))
		if err != nil {
			return err
		}
		err = batch.Put(consumedStateKey(id), state)
		if err != nil {
			return err
		}
		err = t.kv.Delete(stateKey(id))
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

func consumedStateKey(id transaction.StateID) []byte {
	return append(
		append(keyConsumedStates[:], id.TxID[:]...),
		strconv.Itoa(int(id.OutputIndex))...,
	)
}
