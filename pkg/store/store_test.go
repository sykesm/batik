// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package store

import (
	"encoding/hex"
	"testing"

	. "github.com/onsi/gomega"

	txv1 "github.com/sykesm/batik/pkg/pb/tx/v1"
	"github.com/sykesm/batik/pkg/tested"
	"github.com/sykesm/batik/pkg/transaction"
)

func TestReceipts(t *testing.T) {
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

func TestPending(t *testing.T) {
	gt := NewGomegaWithT(t)
	gt.Expect(true).To(BeTrue())
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
