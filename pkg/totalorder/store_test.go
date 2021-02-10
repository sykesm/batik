// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package totalorder

import (
	"context"
	"crypto"
	"crypto/sha256"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/sykesm/batik/pkg/store"
	"github.com/sykesm/batik/pkg/tested"
	"github.com/sykesm/batik/pkg/transaction"
)

func TestStore(t *testing.T) {
	gt := NewGomegaWithT(t)

	path, cleanup := tested.TempDir(t, "", "totalorder-store")
	defer cleanup()

	db, err := store.NewLevelDB(path)
	gt.Expect(err).NotTo(HaveOccurred())
	defer tested.Close(t, db)

	orderStore := NewStore(crypto.SHA256, db)

	err = orderStore.Append(TXIDAndHMAC{
		ID:   transaction.ID(sHash("tx1")),
		HMAC: sHash("tx1" + "secret"),
	})
	gt.Expect(err).NotTo(HaveOccurred())

	// Test that if already committed it returns immediately
	tah, err := orderStore.Get(context.Background(), 0)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(tah).To(Equal(TXIDAndHMAC{
		ID:   transaction.ID(sHash("tx1")),
		HMAC: sHash("tx1" + "secret"),
	}))

	type result struct {
		val TXIDAndHMAC
		err error
	}
	resultC := make(chan result)

	// Test a first waiter (which creates the waitC)
	go func() {
		tah, err := orderStore.Get(context.Background(), 1)
		resultC <- result{val: tah, err: err}
	}()

	// Test a second waiter (which re-uses the waitC)
	go func() {
		tah, err := orderStore.Get(context.Background(), 1)
		resultC <- result{val: tah, err: err}
	}()

	gt.Consistently(resultC).ShouldNot(Receive())

	// Test a third waiter (which cancels)
	errC := make(chan error)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		_, err := orderStore.Get(ctx, 1)
		errC <- err
	}()
	cancel()
	gt.Eventually(errC).Should(Receive(&err))
	gt.Expect(err).To(MatchError(context.Canceled))

	gt.Eventually(func() int {
		orderStore.mutex.Lock()
		defer orderStore.mutex.Unlock()
		return len(orderStore.waitCs)
	}).Should(Equal(1))

	err = orderStore.Append(TXIDAndHMAC{
		ID:   transaction.ID(sHash("tx2")),
		HMAC: sHash("tx2" + "secret"),
	})
	gt.Expect(err).NotTo(HaveOccurred())
	orderStore.mutex.Lock()
	gt.Expect(orderStore.waitCs).To(HaveLen(0))
	orderStore.mutex.Unlock()

	// Verify both results returned by the first two waiters are as expected
	var result1, result2 result
	gt.Eventually(resultC).Should(Receive(&result1))
	gt.Expect(result1.err).NotTo(HaveOccurred())
	gt.Expect(result1.val).To(Equal(TXIDAndHMAC{
		ID:   transaction.ID(sHash("tx2")),
		HMAC: sHash("tx2" + "secret"),
	}))

	gt.Eventually(resultC).Should(Receive(&result2))
	gt.Expect(result2.err).NotTo(HaveOccurred())
	gt.Expect(result2.val).To(Equal(TXIDAndHMAC{
		ID:   transaction.ID(sHash("tx2")),
		HMAC: sHash("tx2" + "secret"),
	}))

	val, err := db.Get(keyMetadataLastCommitted)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(bytesToUint64(val)).To(Equal(uint64(1)))

	accumulator, err := db.Get(keyMetadataAccumulator)
	gt.Expect(err).NotTo(HaveOccurred())
	initialAccumulator := []byte{}
	accumulatorAfterTx1 := bHash(
		append(
			initialAccumulator,
			(&TXIDAndHMAC{
				ID:   transaction.ID(sHash("tx1")),
				HMAC: sHash("tx1" + "secret"),
			}).serialize()...,
		),
	)
	acculatorAfterTx2 := bHash(
		append(
			accumulatorAfterTx1,
			(&TXIDAndHMAC{
				ID:   transaction.ID(sHash("tx2")),
				HMAC: sHash("tx2" + "secret"),
			}).serialize()...,
		),
	)
	gt.Expect(accumulator).To(Equal(acculatorAfterTx2))
}

func sHash(value string) []byte {
	return bHash([]byte(value))
}

func bHash(value []byte) []byte {
	h := sha256.New()
	h.Write(value)
	return h.Sum(nil)
}
