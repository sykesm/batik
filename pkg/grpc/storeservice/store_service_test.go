// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package storeservice

import (
	"context"
	"crypto"
	"testing"

	. "github.com/onsi/gomega"

	storev1 "github.com/sykesm/batik/pkg/pb/store/v1"
	txv1 "github.com/sykesm/batik/pkg/pb/tx/v1"
	"github.com/sykesm/batik/pkg/store"
	"github.com/sykesm/batik/pkg/tested"
	. "github.com/sykesm/batik/pkg/tested/matcher"
	"github.com/sykesm/batik/pkg/transaction"
)

func TestStoreService_GetTransaction(t *testing.T) {
	gt := NewGomegaWithT(t)
	storeSvc, cleanup := newStoreService(t)
	defer cleanup()

	testTx := newTestTransaction()
	intTx, err := transaction.New(crypto.SHA256, testTx)
	gt.Expect(err).NotTo(HaveOccurred())

	req := &storev1.GetTransactionRequest{
		Txid: intTx.ID,
	}
	resp, err := storeSvc.GetTransaction(context.Background(), req)
	gt.Expect(err).To(MatchError(ContainSubstring("leveldb: not found")))

	err = storeSvc.repo.PutTransaction(intTx)
	gt.Expect(err).NotTo(HaveOccurred())

	resp, err = storeSvc.GetTransaction(context.Background(), req)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(resp.Transaction).To(ProtoEqual(testTx))
}

func TestStoreService_PutTransaction(t *testing.T) {
	gt := NewGomegaWithT(t)
	storeSvc, cleanup := newStoreService(t)
	defer cleanup()

	testTx := newTestTransaction()
	intTx, err := transaction.New(crypto.SHA256, testTx)
	gt.Expect(err).NotTo(HaveOccurred())

	req := &storev1.PutTransactionRequest{
		Transaction: testTx,
	}
	_, err = storeSvc.PutTransaction(context.Background(), req)
	gt.Expect(err).NotTo(HaveOccurred())

	result, err := storeSvc.GetTransaction(context.Background(), &storev1.GetTransactionRequest{
		Txid: intTx.ID,
	})
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(result.Transaction).To(ProtoEqual(testTx))
}

func TestStoreService_GetState(t *testing.T) {
	gt := NewGomegaWithT(t)
	storeSvc, cleanup := newStoreService(t)
	defer cleanup()

	testTx := newTestTransaction()
	intTx, err := transaction.New(crypto.SHA256, testTx)
	gt.Expect(err).NotTo(HaveOccurred())

	stateRef := &txv1.StateReference{
		Txid:        intTx.ID,
		OutputIndex: 0,
	}

	req := &storev1.GetStateRequest{
		StateRef: stateRef,
	}
	resp, err := storeSvc.GetState(context.Background(), req)
	gt.Expect(err).To(MatchError(ContainSubstring("leveldb: not found")))

	err = storeSvc.repo.PutState(&transaction.State{
		ID:        transaction.StateID{TxID: stateRef.Txid, OutputIndex: stateRef.OutputIndex},
		StateInfo: testTx.Outputs[0].Info,
		Data:      testTx.Outputs[0].State,
	})
	gt.Expect(err).NotTo(HaveOccurred())

	resp, err = storeSvc.GetState(context.Background(), req)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(resp.State).To(ProtoEqual(testTx.Outputs[0]))
}

func TestStoreService_PutState(t *testing.T) {
	gt := NewGomegaWithT(t)
	storeSvc, cleanup := newStoreService(t)
	defer cleanup()

	testTx := newTestTransaction()
	intTx, err := transaction.New(crypto.SHA256, testTx)
	gt.Expect(err).NotTo(HaveOccurred())

	testState := testTx.Outputs[0]

	req := &storev1.PutStateRequest{
		StateReference: &txv1.StateReference{
			Txid:        intTx.ID,
			OutputIndex: 0,
		},
		State: testState,
	}
	_, err = storeSvc.PutState(context.Background(), req)
	gt.Expect(err).NotTo(HaveOccurred())

	stateID := transaction.StateID{
		TxID:        req.StateReference.Txid,
		OutputIndex: req.StateReference.OutputIndex,
	}
	state, err := storeSvc.repo.GetState(stateID)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(state.StateInfo).To(ProtoEqual(testState.Info))
	gt.Expect(state.Data).To(Equal(testState.State))
}

func newStoreService(t *testing.T) (*StoreService, func()) {
	path, cleanup := tested.TempDir(t, "", "level")

	db, err := store.NewLevelDB(path)
	NewGomegaWithT(t).Expect(err).NotTo(HaveOccurred())

	repo := store.NewRepository(db)
	storeSvc := NewStoreService(repo)

	return storeSvc, func() {
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
