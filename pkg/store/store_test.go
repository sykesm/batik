// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package store

import (
	"crypto"
	"encoding/hex"
	"testing"

	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/proto"

	txv1 "github.com/sykesm/batik/pkg/pb/tx/v1"
	"github.com/sykesm/batik/pkg/tested"
	"github.com/sykesm/batik/pkg/transaction"
)

func TestStoreReceipts(t *testing.T) {
	gt := NewGomegaWithT(t)

	store, cleanup := setupTestStore(t)
	defer cleanup()

	r := &transaction.Receipt{
		TxID: []byte("txid"),
		Signatures: []*transaction.Signature{
			{
				PublicKey: []byte("pk1"),
				Signature: []byte("sig1"),
			},
			{
				PublicKey: []byte("pk2"),
				Signature: []byte("sig2"),
			},
		},
		ID: []byte("receiptid"),
	}

	_, err := store.GetReceipt([]byte("receiptid"))
	gt.Expect(err).To(HaveOccurred())
	gt.Expect(IsNotFound(err)).To(BeTrue())

	err = store.PutReceipt(r)
	gt.Expect(err).NotTo(HaveOccurred())

	nr, err := store.GetReceipt([]byte("receiptid"))
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(nr).To(Equal(r))
}

func TestStoreCommits(t *testing.T) {
	gt := NewGomegaWithT(t)

	store, cleanup := setupTestStore(t)
	defer cleanup()

	txid := transaction.ID([]byte("tx-id"))

	c := &transaction.Committed{
		SeqNo:     7,
		ReceiptID: []byte("receiptid"),
	}

	_, err := store.GetCommitted(txid)
	gt.Expect(err).To(HaveOccurred())
	gt.Expect(IsNotFound(err)).To(BeTrue())

	err = store.PutCommitted(txid, c)
	gt.Expect(err).NotTo(HaveOccurred())

	nc, err := store.GetCommitted(txid)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(nc).To(Equal(c))
}

func TestStoreTransaction(t *testing.T) {
	gt := NewGomegaWithT(t)

	store, cleanup := setupTestStore(t)
	defer cleanup()

	tx, err := transaction.New(crypto.SHA256, newTestTransaction())
	gt.Expect(err).NotTo(HaveOccurred())

	_, err = store.GetTransaction(tx.ID)
	gt.Expect(err).To(HaveOccurred())
	gt.Expect(IsNotFound(err)).To(BeTrue())

	err = store.PutTransaction(tx)
	gt.Expect(err).NotTo(HaveOccurred())

	ntx, err := store.GetTransaction(tx.ID)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(proto.Equal(tx.Tx, ntx.Tx)).To(BeTrue())
	tx.Tx, ntx.Tx = nil, nil
	gt.Expect(tx).To(Equal(ntx))
}

func TestStoreState(t *testing.T) {
	gt := NewGomegaWithT(t)

	store, cleanup := setupTestStore(t)
	defer cleanup()

	tx, err := transaction.New(crypto.SHA256, newTestTransaction())
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(tx.Outputs).To(HaveLen(2))
	state := tx.Outputs[0]

	// Verify the state is not available as consumed or unconsumed, and fails to consume
	_, err = store.GetState(state.ID, false)
	gt.Expect(err).To(HaveOccurred())
	gt.Expect(IsNotFound(err)).To(BeTrue())

	_, err = store.GetState(state.ID, true)
	gt.Expect(err).To(HaveOccurred())
	gt.Expect(IsNotFound(err)).To(BeTrue())

	err = store.ConsumeState(state.ID)
	gt.Expect(err).To(HaveOccurred())
	gt.Expect(IsNotFound(err)).To(BeTrue())

	// Put the state
	err = store.PutState(state)
	gt.Expect(err).NotTo(HaveOccurred())

	// Verify it is reported not consumed
	nstate, err := store.GetState(state.ID, false)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(nstate).To(Equal(state))

	// Verify it is not reported consumed
	_, err = store.GetState(state.ID, true)
	gt.Expect(err).To(HaveOccurred())
	gt.Expect(IsNotFound(err)).To(BeTrue())

	// Consume it
	err = store.ConsumeState(state.ID)
	gt.Expect(err).NotTo(HaveOccurred())

	// Verify it is not reported not consumed
	_, err = store.GetState(state.ID, false)
	gt.Expect(err).To(HaveOccurred())
	gt.Expect(IsNotFound(err)).To(BeTrue())

	// Verify it is reported consumed
	nstate, err = store.GetState(state.ID, true)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(nstate).To(Equal(state))
}

func setupTestStore(t *testing.T) (*TransactionRepository, func()) {
	path, cleanup := tested.TempDir(t, "", "store")
	db, err := NewLevelDB(path)
	if err != nil {
		cleanup()
		t.Fatalf("could not create db: %s", err)
	}

	return NewRepository(db), func() {
		tested.Close(t, db)
		cleanup()
	}
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
					Owners: []*txv1.Party{{PublicKey: []byte("owner-1")}, {PublicKey: []byte("owner-2")}},
					Kind:   "state-kind-0",
				},
				State: []byte("state-0"),
			},
			{
				Info: &txv1.StateInfo{
					Owners: []*txv1.Party{{PublicKey: []byte("owner-1")}, {PublicKey: []byte("owner-2")}},
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
			{PublicKey: []byte("observer-1")},
			{PublicKey: []byte("observer-2")},
		},
	}
}

func fromHex(t *testing.T, s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatalf("failed to decode %q as hex string: %s", s, err)
	}
	return b
}

func TestKeyEncodings(t *testing.T) {
	gt := NewGomegaWithT(t)
	gt.Expect(true).To(BeTrue())

	txID := fromHex(t, "deadbeef")

	txKey := transactionKey(txID)
	gt.Expect(txKey).To(Equal(fromHex(t, "01deadbeef")))

	stateID := transaction.StateID{TxID: txID, OutputIndex: 1}

	sKey := stateKey(stateID)
	gt.Expect(sKey).To(Equal(fromHex(t, "02deadbeef0000000000000001")))

	siKey := stateInfoKey(stateID)
	gt.Expect(siKey).To(Equal(fromHex(t, "03deadbeef0000000000000001")))

	scKey := consumedStateKey(stateID)
	gt.Expect(scKey).To(Equal(fromHex(t, "04deadbeef0000000000000001")))
}
