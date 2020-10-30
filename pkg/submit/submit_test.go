// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package submit

import (
	"context"
	"errors"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/sykesm/batik/pkg/store"
	"github.com/sykesm/batik/pkg/submit/fakes"
	"github.com/sykesm/batik/pkg/transaction"
)

//go:generate counterfeiter -o fakes/repository.go --fake-name Repository . fakeRepository
type fakeRepository Repository // private import to prevent import cycle in generated fake

func TestSubmit(t *testing.T) {
	gt := NewGomegaWithT(t)
	fakeRepo := &fakes.Repository{}
	submitService := NewService(fakeRepo)

	fakeRepo.GetTransactionReturns(nil, &store.NotFoundError{Err: errors.New("not-found-error")})
	signed := &transaction.Signed{
		Transaction: &transaction.Transaction{
			ID: transaction.NewID([]byte("transaction-id")),
		},
	}

	err := submitService.Submit(context.Background(), signed)
	gt.Expect(err).NotTo(HaveOccurred())

	gt.Expect(fakeRepo.GetTransactionCallCount()).To(Equal(1))
	getID := fakeRepo.GetTransactionArgsForCall(0)
	gt.Expect(getID).To(Equal(transaction.ID([]byte("transaction-id"))))

	gt.Expect(fakeRepo.GetStateCallCount()).To(Equal(0))

	gt.Expect(fakeRepo.PutTransactionCallCount()).To(Equal(1))
	putTx := fakeRepo.PutTransactionArgsForCall(0)
	gt.Expect(putTx).To(Equal(signed.Transaction)) // TODO: This should be a signed transaction

	gt.Expect(fakeRepo.PutStateCallCount()).To(Equal(0))
	gt.Expect(fakeRepo.ConsumeStateCallCount()).To(Equal(0))
}

func TestSubmitGetTransaction(t *testing.T) {
	signed := &transaction.Signed{
		Transaction: &transaction.Transaction{
			ID: transaction.NewID([]byte("transaction-id")),
		},
	}

	t.Run("ExistingTransaction", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		fakeRepo := &fakes.Repository{}
		submitService := NewService(fakeRepo)

		err := submitService.Submit(context.Background(), signed)
		gt.Expect(err).To(HaveOccurred())
		gt.Expect(store.IsAlreadyExists(err)).To(BeTrue())
	})

	t.Run("UnexpectedError", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		fakeRepo := &fakes.Repository{}
		fakeRepo.GetTransactionReturns(nil, errors.New("unexpected-error"))
		submitService := NewService(fakeRepo)

		err := submitService.Submit(context.Background(), signed)
		gt.Expect(err).To(Equal(errors.New("unexpected-error")))
	})
}
