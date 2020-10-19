// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package store

import (
	"context"
	"crypto"
	"testing"

	. "github.com/onsi/gomega"

	storev1 "github.com/sykesm/batik/pkg/pb/store/v1"
	txv1 "github.com/sykesm/batik/pkg/pb/tx/v1"
	"github.com/sykesm/batik/pkg/protomsg"
	"github.com/sykesm/batik/pkg/tested"
	. "github.com/sykesm/batik/pkg/tested/matcher"
	"github.com/sykesm/batik/pkg/transaction"
)

func TestStoreService_GetTransaction(t *testing.T) {
	gt := NewGomegaWithT(t)

	path, cleanup := tested.TempDir(t, "", "level")
	defer cleanup()

	db, err := NewLevelDB(path)
	gt.Expect(err).NotTo(HaveOccurred())
	defer tested.Close(t, db)

	storeSvc := NewStoreService(db)

	testTx := newTestTransaction()
	intTx, err := transaction.Marshal(crypto.SHA256, testTx)
	gt.Expect(err).NotTo(HaveOccurred())

	key := transactionKey(intTx.ID)

	req := &storev1.GetTransactionRequest{
		Txid: intTx.ID,
	}
	resp, err := storeSvc.GetTransaction(context.Background(), req)
	gt.Expect(err).To(MatchError(MatchRegexp("leveldb: not found")))

	err = db.Put(key, intTx.Encoded)
	gt.Expect(err).NotTo(HaveOccurred())

	resp, err = storeSvc.GetTransaction(context.Background(), req)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(resp.Transaction).To(ProtoEqual(testTx))
}

func TestStoreService_PutTransaction(t *testing.T) {
	gt := NewGomegaWithT(t)

	path, cleanup := tested.TempDir(t, "", "level")
	defer cleanup()

	db, err := NewLevelDB(path)
	gt.Expect(err).NotTo(HaveOccurred())
	defer tested.Close(t, db)

	storeSvc := NewStoreService(db)

	testTx := newTestTransaction()
	intTx, err := transaction.Marshal(crypto.SHA256, testTx)
	gt.Expect(err).NotTo(HaveOccurred())

	key := transactionKey(intTx.ID)

	req := &storev1.PutTransactionRequest{
		Transaction: testTx,
	}
	_, err = storeSvc.PutTransaction(context.Background(), req)
	gt.Expect(err).NotTo(HaveOccurred())

	data, err := db.Get(key)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(data).To(Equal(intTx.Encoded))
}

func TestStoreService_GetState(t *testing.T) {
	gt := NewGomegaWithT(t)

	path, cleanup := tested.TempDir(t, "", "level")
	defer cleanup()

	db, err := NewLevelDB(path)
	gt.Expect(err).NotTo(HaveOccurred())
	defer tested.Close(t, db)

	storeSvc := NewStoreService(db)

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

	req := &storev1.GetStateRequest{
		StateRef: testStateRef,
	}
	resp, err := storeSvc.GetState(context.Background(), req)
	gt.Expect(err).To(MatchError(MatchRegexp("leveldb: not found")))

	err = db.Put(key, encodedState)
	gt.Expect(err).NotTo(HaveOccurred())

	resp, err = storeSvc.GetState(context.Background(), req)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(resp.State).To(ProtoEqual(testState))
}

func TestStoreService_PutState(t *testing.T) {
	gt := NewGomegaWithT(t)

	path, cleanup := tested.TempDir(t, "", "level")
	defer cleanup()

	db, err := NewLevelDB(path)
	gt.Expect(err).NotTo(HaveOccurred())
	defer tested.Close(t, db)

	storeSvc := NewStoreService(db)

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

	req := &storev1.PutStateRequest{
		State: testState,
	}
	_, err = storeSvc.PutState(context.Background(), req)
	gt.Expect(err).NotTo(HaveOccurred())

	data, err := db.Get(key)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(data).To(Equal(encodedState))
}
