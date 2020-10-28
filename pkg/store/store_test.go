// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package store

import (
	"crypto"
	"encoding/hex"
	"testing"

	. "github.com/onsi/gomega"

	txv1 "github.com/sykesm/batik/pkg/pb/tx/v1"
	"github.com/sykesm/batik/pkg/protomsg"
	"github.com/sykesm/batik/pkg/tested"
	. "github.com/sykesm/batik/pkg/tested/matcher"
	"github.com/sykesm/batik/pkg/transaction"
)

func TestStoreTransactions(t *testing.T) {
	gt := NewGomegaWithT(t)

	path, cleanup := tested.TempDir(t, "", "level")
	defer cleanup()

	db, err := NewLevelDB(path)
	gt.Expect(err).NotTo(HaveOccurred())
	defer tested.Close(t, db)

	testTx := newTestTransaction()
	intTx, err := transaction.Marshal(crypto.SHA256, testTx)
	gt.Expect(err).NotTo(HaveOccurred())

	key := transactionKey(intTx.ID)

	err = StoreTransactions(db, []*txv1.Transaction{testTx})
	gt.Expect(err).NotTo(HaveOccurred())

	data, err := db.Get(key)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(data).To(Equal(intTx.Encoded))
}

func TestLoadTransactions(t *testing.T) {
	gt := NewGomegaWithT(t)

	path, cleanup := tested.TempDir(t, "", "level")
	defer cleanup()

	db, err := NewLevelDB(path)
	gt.Expect(err).NotTo(HaveOccurred())
	defer tested.Close(t, db)

	testTx := newTestTransaction()
	intTx, err := transaction.Marshal(crypto.SHA256, testTx)
	gt.Expect(err).NotTo(HaveOccurred())

	key := transactionKey(intTx.ID)

	_, err = LoadTransactions(db, [][]byte{intTx.ID})
	gt.Expect(err).To(MatchError(ContainSubstring("leveldb: not found")))

	err = db.Put(key, intTx.Encoded)
	gt.Expect(err).NotTo(HaveOccurred())

	txs, err := LoadTransactions(db, [][]byte{intTx.ID})
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(txs[0]).To(ProtoEqual(testTx))
}

func TestStoreStates(t *testing.T) {
	gt := NewGomegaWithT(t)

	path, cleanup := tested.TempDir(t, "", "level")
	defer cleanup()

	db, err := NewLevelDB(path)
	gt.Expect(err).NotTo(HaveOccurred())
	defer tested.Close(t, db)

	testTx := newTestTransaction()
	intTx, err := transaction.Marshal(crypto.SHA256, testTx)
	gt.Expect(err).NotTo(HaveOccurred())

	testState := &txv1.ResolvedState{
		Txid:        intTx.ID,
		OutputIndex: 0,
		Info:        testTx.Outputs[0].Info,
		State:       testTx.Outputs[0].State,
	}

	encodedState, err := protomsg.MarshalDeterministic(testState)
	gt.Expect(err).NotTo(HaveOccurred())

	testStateRef := &txv1.StateReference{
		Txid:        intTx.ID,
		OutputIndex: 0,
	}

	key := stateKey(testStateRef)

	err = StoreStates(db, []*txv1.ResolvedState{testState})
	gt.Expect(err).NotTo(HaveOccurred())

	data, err := db.Get(key)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(data).To(Equal(encodedState))
}

func TestLoadStates(t *testing.T) {
	gt := NewGomegaWithT(t)

	path, cleanup := tested.TempDir(t, "", "level")
	defer cleanup()

	db, err := NewLevelDB(path)
	gt.Expect(err).NotTo(HaveOccurred())
	defer tested.Close(t, db)

	testTx := newTestTransaction()
	intTx, err := transaction.Marshal(crypto.SHA256, testTx)
	gt.Expect(err).NotTo(HaveOccurred())

	testState := &txv1.ResolvedState{
		Txid:        intTx.ID,
		OutputIndex: 0,
		Info:        testTx.Outputs[0].Info,
		State:       testTx.Outputs[0].State,
	}

	testStateRef := &txv1.StateReference{
		Txid:        intTx.ID,
		OutputIndex: 0,
	}

	encodedState, err := protomsg.MarshalDeterministic(testState)
	gt.Expect(err).NotTo(HaveOccurred())

	key := stateKey(testStateRef)

	_, err = LoadStates(db, []*txv1.StateReference{testStateRef})
	gt.Expect(err).To(MatchError(ContainSubstring("leveldb: not found")))

	err = db.Put(key, encodedState)
	gt.Expect(err).NotTo(HaveOccurred())

	states, err := LoadStates(db, []*txv1.StateReference{testStateRef})
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(states[0]).To(ProtoEqual(testState))
}

func newTestTransaction() *txv1.Transaction {
	return &txv1.Transaction{
		Salt: []byte("NaCl - abcdefghijklmnopqrstuvwxyz"),
		Inputs: []*txv1.StateReference{
			{Txid: []byte("input-transaction-id-0"), OutputIndex: 1},
			{Txid: []byte("input-transaction-id-1"), OutputIndex: 0},
		},
		References: []*txv1.StateReference{
			{Txid: []byte("ref-transaction-id-0"), OutputIndex: 1},
			{Txid: []byte("ref-transaction-id-1"), OutputIndex: 0},
		},
		Outputs: []*txv1.State{
			{
				Info: &txv1.StateInfo{
					Owners: []*txv1.Party{{Credential: []byte("owner-1")}, {Credential: []byte("owner-2")}},
					Kind:   "state-kind-0",
				},
				State: []byte("state-0"),
			},
			{
				Info: &txv1.StateInfo{
					Owners: []*txv1.Party{{Credential: []byte("owner-1")}, {Credential: []byte("owner-2")}},
					Kind:   "state-kind-1",
				},
				State: []byte("state-1"),
			},
		},
		Parameters: []*txv1.Parameter{
			{Name: "name-0", Value: []byte("value-0")},
			{Name: "name-1", Value: []byte("value-1")},
		},
		RequiredSigners: []*txv1.Party{
			{Credential: []byte("observer-1")},
			{Credential: []byte("observer-2")},
		},
	}
}

func fromHex(t *testing.T, s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatalf("failed to decode %q as hex string", s)
	}
	return b
}
