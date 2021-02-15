// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package namespace

import (
	"context"
	"crypto"
	"errors"
	"testing"

	. "github.com/onsi/gomega"
	"go.uber.org/zap"

	"github.com/sykesm/batik/pkg/namespace/fake"
	"github.com/sykesm/batik/pkg/store"
	"github.com/sykesm/batik/pkg/tested"
	"github.com/sykesm/batik/pkg/transaction"
	"github.com/sykesm/batik/pkg/validator"
)

func TestNamespace_New(t *testing.T) {
	gt := NewGomegaWithT(t)

	storeDB, cleanup := newKVDB(t)
	defer cleanup()

	logger := zap.NewExample()
	v := validator.NewSignature()

	ns := New(logger, crypto.SHA256, storeDB, v)
	gt.Expect(ns.Logger).To(Equal(logger))
	gt.Expect(ns.LevelDB).To(Equal(storeDB))
	gt.Expect(ns.Repo).NotTo(BeNil())
}

func newKVDB(t *testing.T) (*store.LevelDBKV, func()) {
	path, cleanup := tested.TempDir(t, "", "level")

	db, err := store.NewLevelDB(path)
	NewGomegaWithT(t).Expect(err).NotTo(HaveOccurred())

	return db, cleanup
}

func TestNamespace_Submit(t *testing.T) {
	var (
		ns       *Namespace
		fakeRepo *fake.Repository
	)

	setup := func() {
		fakeRepo = &fake.Repository{}

		ns = &Namespace{
			Logger:  zap.NewExample(),
			Hasher:  crypto.SHA256,
			LevelDB: nil, // We use a faked repo, so no db needed
			Repo:    fakeRepo,
			committer: &committer{
				validator: validator.NewSignature(),
				repo:      fakeRepo,
			},
		}
	}

	signed := &transaction.Signed{
		Transaction: &transaction.Transaction{
			ID: transaction.NewID([]byte("transaction-id")),
		},
	}

	t.Run("PutTransactionFails", func(t *testing.T) {
		setup()
		gt := NewGomegaWithT(t)

		fakeRepo.PutTransactionReturns(errors.New("put-tx-error"))

		err := ns.Submit(context.Background(), signed)
		gt.Expect(err).To(HaveOccurred())
		gt.Expect(err).To(MatchError("failed to store transaction: put-tx-error"))
	})

	t.Run("UnexpectedError", func(t *testing.T) {
		setup()
		gt := NewGomegaWithT(t)

		fakeRepo.PutReceiptReturns(errors.New("put-receipt-error"))

		err := ns.Submit(context.Background(), signed)
		gt.Expect(err).To(MatchError("failed to store transaction receipt: put-receipt-error"))
	})
}
