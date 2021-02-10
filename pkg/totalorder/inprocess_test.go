// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package totalorder

import (
	"context"
	"crypto"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/sykesm/batik/pkg/store"
	"github.com/sykesm/batik/pkg/tested"
	"github.com/sykesm/batik/pkg/transaction"
)

func TestInProcess(t *testing.T) {
	gt := NewGomegaWithT(t)

	path, cleanup := tested.TempDir(t, "", "totalorder-inprocess")
	defer cleanup()

	db, err := store.NewLevelDB(path)
	gt.Expect(err).NotTo(HaveOccurred())
	defer tested.Close(t, db)

	ip := &InProcess{
		store: NewStore(crypto.SHA256, db),
		doneC: make(chan struct{}),
		queue: make(chan TXIDAndHMAC),
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err = ip.Broadcast(ctx, TXIDAndHMAC{
		ID:   transaction.ID(sHash("tx1")),
		HMAC: sHash("tx1" + "secret"),
	})
	gt.Expect(err).To(MatchError(context.Canceled))

	_, err = ip.Deliver(ctx, 1)
	gt.Expect(err).To(MatchError(context.Canceled))

	exitC := make(chan struct{})

	go func() {
		ip.run()
		close(exitC)
	}()

	err = ip.Broadcast(context.Background(), TXIDAndHMAC{
		ID:   transaction.ID(sHash("tx1")),
		HMAC: sHash("tx1" + "secret"),
	})
	gt.Expect(err).NotTo(HaveOccurred())

	tah, err := ip.Deliver(context.Background(), 0)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(tah).To(Equal(
		TXIDAndHMAC{
			ID:   transaction.ID(sHash("tx1")),
			HMAC: sHash("tx1" + "secret"),
		},
	))

	ip.Stop()

	gt.Eventually(exitC).Should(BeClosed())

	err = ip.Broadcast(context.Background(), TXIDAndHMAC{
		ID:   transaction.ID(sHash("tx1")),
		HMAC: sHash("tx1" + "secret"),
	})
	gt.Expect(err).To(MatchError("told to exit"))

}
