// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package store

import (
	"crypto"
	"testing"

	. "github.com/onsi/gomega"

	tb "github.com/sykesm/batik/pkg/pb/transaction"
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

	err = StoreTransactions(db, []*tb.Transaction{testTx})
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
	gt.Expect(err).To(MatchError(MatchRegexp("leveldb: not found")))

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

	testState := &tb.ResolvedState{
		Txid:        intTx.ID,
		OutputIndex: 0,
		Info:        testTx.Outputs[0].Info,
		State:       testTx.Outputs[0].State,
	}

	encodedState, err := protomsg.MarshalDeterministic(testState)
	gt.Expect(err).NotTo(HaveOccurred())

	testStateRef := &tb.StateReference{
		Txid:        intTx.ID,
		OutputIndex: 0,
	}

	key := stateKey(testStateRef)

	err = StoreStates(db, []*tb.ResolvedState{testState})
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

	testState := &tb.ResolvedState{
		Txid:        intTx.ID,
		OutputIndex: 0,
		Info:        testTx.Outputs[0].Info,
		State:       testTx.Outputs[0].State,
	}

	testStateRef := &tb.StateReference{
		Txid:        intTx.ID,
		OutputIndex: 0,
	}

	encodedState, err := protomsg.MarshalDeterministic(testState)
	gt.Expect(err).NotTo(HaveOccurred())

	key := stateKey(testStateRef)

	_, err = LoadStates(db, []*tb.StateReference{testStateRef})
	gt.Expect(err).To(MatchError(MatchRegexp("leveldb: not found")))

	err = db.Put(key, encodedState)
	gt.Expect(err).NotTo(HaveOccurred())

	states, err := LoadStates(db, []*tb.StateReference{testStateRef})
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(states[0]).To(ProtoEqual(testState))
}

func newTestTransaction() *tb.Transaction {
	return &tb.Transaction{
		Inputs: []*tb.StateReference{
			{Txid: []byte("input-transaction-id-0"), OutputIndex: 1},
			{Txid: []byte("input-transaction-id-1"), OutputIndex: 0},
		},
		References: []*tb.StateReference{
			{Txid: []byte("ref-transaction-id-0"), OutputIndex: 1},
			{Txid: []byte("ref-transaction-id-1"), OutputIndex: 0},
		},
		Outputs: []*tb.State{
			{
				Info: &tb.StateInfo{
					Owners: []*tb.Party{{Credential: []byte("owner-1")}, {Credential: []byte("owner-2")}},
					Kind:   "state-kind-0",
				},
				State: []byte("state-0"),
			},
			{
				Info: &tb.StateInfo{
					Owners: []*tb.Party{{Credential: []byte("owner-1")}, {Credential: []byte("owner-2")}},
					Kind:   "state-kind-1",
				},
				State: []byte("state-1"),
			},
		},
		Parameters: []*tb.Parameter{
			{Name: "name-0", Value: []byte("value-0")},
			{Name: "name-1", Value: []byte("value-1")},
		},
		RequiredSigners: []*tb.Party{
			{Credential: []byte("observer-1")},
			{Credential: []byte("observer-2")},
		},
		Salt: []byte("NaCl"),
	}
}
