// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package store

import (
	"crypto"
	"encoding/binary"
	"encoding/json"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	txv1 "github.com/sykesm/batik/pkg/pb/tx/v1"
	"github.com/sykesm/batik/pkg/protomsg"
	"github.com/sykesm/batik/pkg/transaction"
)

// TODO: Determine how to model the hasher required to restore a transaction.
// TODO: Standarize on binary mashaling and unmarshaling to remove proto
// TODO: Unit of work / atomicity / snapshot isolation

type TransactionRepository struct {
	kv KV
}

func NewRepository(kv KV) *TransactionRepository {
	return &TransactionRepository{
		kv: kv,
	}
}

func (t *TransactionRepository) PutReceipt(receipt *transaction.Receipt) error {
	serialized, err := json.Marshal(receipt)
	if err != nil {
		return errors.WithMessage(err, "could not serialize receipt to JSON")
	}

	err = t.kv.Put(receiptKey(receipt.ID), serialized)
	if err != nil {
		return errors.WithMessage(err, "failed to store receipt")
	}

	return nil
}

func (t *TransactionRepository) GetReceipt(id []byte) (*transaction.Receipt, error) {
	data, err := t.kv.Get(receiptKey(id))
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get receipt from db")
	}

	var r transaction.Receipt
	err = json.Unmarshal(data, &r)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to unmarshal retreived receipt")
	}

	r.ID = id

	return &r, nil
}

func (t *TransactionRepository) PutCommitted(id transaction.ID, commit *transaction.Committed) error {
	serialized, err := json.Marshal(commit)
	if err != nil {
		return errors.WithMessage(err, "could not serialize commit to JSON")
	}

	err = t.kv.Put(commitKey(id), serialized)
	if err != nil {
		return errors.WithMessage(err, "failed to store commit")
	}

	return nil
}

func (t *TransactionRepository) GetCommitted(id transaction.ID) (*transaction.Committed, error) {
	data, err := t.kv.Get(commitKey(id))
	if err != nil {
		return nil, errors.WithMessage(err, "failed to get tx commitment from db")
	}

	var r transaction.Committed
	err = json.Unmarshal(data, &r)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to unmarshal retreived commit")
	}

	return &r, nil
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
	var owners []*txv1.Party
	si := state.StateInfo
	for i := range si.Owners {
		owners = append(owners, &txv1.Party{PublicKey: si.Owners[i].PublicKey})
	}
	stateInfo := &txv1.StateInfo{
		Owners: owners,
		Kind:   si.Kind,
	}

	info, err := protomsg.MarshalDeterministic(stateInfo)
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

func (t *TransactionRepository) GetState(stateID transaction.StateID, consumed bool) (*transaction.State, error) {
	infoPayload, err := t.kv.Get(stateInfoKey(stateID))
	if err != nil {
		return nil, errors.WithMessagef(err, "error getting state info for %s from db", stateID)
	}
	var stateInfo txv1.StateInfo
	if err := proto.Unmarshal(infoPayload, &stateInfo); err != nil {
		return nil, errors.WithMessagef(err, "error unmarshaling state info for ref %s", stateID)
	}

	var payload []byte
	switch {
	case consumed:
		payload, err = t.kv.Get(consumedStateKey(stateID))
	default:
		payload, err = t.kv.Get(stateKey(stateID))
	}
	if err != nil {
		return nil, errors.WithMessagef(err, "error getting state %s from db", stateID)
	}

	var owners []*transaction.Party
	for i := range stateInfo.Owners {
		owners = append(owners, &transaction.Party{PublicKey: stateInfo.Owners[i].PublicKey})
	}

	state := &transaction.State{
		ID: stateID,
		StateInfo: &transaction.StateInfo{
			Kind:   stateInfo.Kind,
			Owners: owners,
		},
		Data: payload,
	}

	return state, nil
}

func (t *TransactionRepository) ConsumeState(stateID transaction.StateID) error {
	return t.consumeStates(stateID)
}

func (t *TransactionRepository) consumeStates(stateIDs ...transaction.StateID) error {
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
	keyReceipts       = [...]byte{0x5}
	keyCommits        = [...]byte{0x6}
)

// transactionKey returns a db key for a transaction
func transactionKey(txid []byte) []byte {
	key := make([]byte, len(keyTransactions)+len(txid))
	copy(key, keyTransactions[:])
	copy(key[len(keyTransactions):], txid)
	return key
}

// buildKey builds a fixed length key of the form:
//  <prefix><txid><big-endian-uint64>
func buildKey(prefix, txid []byte, outputIndex uint64) []byte {
	key := make([]byte, len(prefix)+len(txid)+8)
	copy(key, prefix)
	copy(key[len(prefix):], txid)
	binary.BigEndian.PutUint64(key[len(prefix)+len(txid):], outputIndex)
	return key
}

// stateKey returns a db key for a state
func stateKey(id transaction.StateID) []byte {
	return buildKey(keyStates[:], id.TxID, id.OutputIndex)
}

// stateInfoKey returns a db key for a stateInfo
func stateInfoKey(id transaction.StateID) []byte {
	return buildKey(keyStateInfos[:], id.TxID, id.OutputIndex)
}

// stateIfoKey returns a db key for a stateInfo
func consumedStateKey(id transaction.StateID) []byte {
	return buildKey(keyConsumedStates[:], id.TxID, id.OutputIndex)
}

func receiptKey(receiptID []byte) []byte {
	key := make([]byte, len(keyReceipts)+len(receiptID))
	copy(key, keyReceipts[:])
	copy(key[len(keyReceipts):], receiptID)
	return key
}

func commitKey(txid []byte) []byte {
	key := make([]byte, len(keyCommits)+len(txid))
	copy(key, keyCommits[:])
	copy(key[len(keyCommits):], txid)
	return key
}
