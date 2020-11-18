// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package submit

import (
	"context"
	"crypto"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	"github.com/sykesm/batik/pkg/ecdsautil"
	txv1 "github.com/sykesm/batik/pkg/pb/tx/v1"
	"github.com/sykesm/batik/pkg/store"
	"github.com/sykesm/batik/pkg/submit/fake"
	"github.com/sykesm/batik/pkg/transaction"
)

//go:generate counterfeiter -o fake/repository.go --fake-name Repository . fakeRepository
type fakeRepository Repository // private to prevent an import cycle in generated fake

var _ fakeRepository = (*fake.Repository)(nil)

func TestSubmit(t *testing.T) {
	var (
		fakeRepo      *fake.Repository
		submitService *Service
		signed        *transaction.Signed
		activeStates  []*transaction.State
	)

	setup := func(t *testing.T) {
		gt := NewGomegaWithT(t)
		sk, err := ecdsautil.GenerateKey(elliptic.P256(), rand.Reader)
		gt.Expect(err).NotTo(HaveOccurred())
		pk, err := ecdsautil.MarshalPublicKey(&sk.PublicKey)
		gt.Expect(err).NotTo(HaveOccurred())
		sig, err := ecdsautil.NewSigner(sk).Sign(rand.Reader, []byte("transaction-id"), crypto.SHA256)
		gt.Expect(err).NotTo(HaveOccurred())

		fakeRepo = &fake.Repository{}
		submitService = NewService(fakeRepo)
		activeStates = []*transaction.State{
			{
				ID: transaction.StateID{TxID: transaction.ID("transaction-id-1"), OutputIndex: 1},
				StateInfo: &transaction.StateInfo{
					Kind: "dummy-state",
					Owners: []*transaction.Party{
						{PublicKey: pk},
					},
				},
				Data: []byte("state-data-1"),
			},
			{
				ID:   transaction.StateID{TxID: transaction.ID("transaction-id-2"), OutputIndex: 0},
				Data: []byte("state-data-2"),
			},
		}

		signed = &transaction.Signed{
			Transaction: &transaction.Transaction{
				ID: transaction.NewID([]byte("transaction-id")),
				Inputs: []*transaction.StateID{
					{TxID: transaction.NewID([]byte("transaction-id-1")), OutputIndex: 1},
				},
				References: []*transaction.StateID{
					{TxID: transaction.NewID([]byte("transaction-id-2")), OutputIndex: 0},
				},
				Outputs: []*transaction.State{
					{
						ID: transaction.StateID{
							TxID: transaction.NewID([]byte("transaction-id")), OutputIndex: 0,
						},
						StateInfo: &transaction.StateInfo{},
						Data:      []byte("output-0-data"),
					},
				},
				Tx:      &txv1.Transaction{},
				Encoded: []byte("encoded-transaction"),
			},
			Signatures: []*transaction.Signature{
				{PublicKey: pk, Signature: sig},
			},
		}
		fakeRepo.GetTransactionReturns(nil, &store.NotFoundError{Err: errors.New("not-found-error")})
		fakeRepo.GetStateStub = func(sid transaction.StateID, consumed bool) (*transaction.State, error) {
			for _, s := range activeStates {
				if s.ID.Equals(sid) {
					return s, nil
				}
			}
			return nil, &store.NotFoundError{Err: errors.Errorf("missing-state %s", sid)}
		}
	}

	t.Run("BasicTransaction", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		setup(t)

		err := submitService.Submit(context.Background(), signed)
		gt.Expect(err).NotTo(HaveOccurred())

		gt.Expect(fakeRepo.GetTransactionCallCount()).To(Equal(1))
		gt.Expect(fakeRepo.GetTransactionArgsForCall(0)).To(Equal(transaction.ID([]byte("transaction-id"))))

		gt.Expect(fakeRepo.GetStateCallCount()).To(Equal(2))
		sID, consumed := fakeRepo.GetStateArgsForCall(0)
		gt.Expect(sID).To(Equal(*signed.Transaction.Inputs[0]))
		gt.Expect(consumed).To(BeFalse())
		sID, consumed = fakeRepo.GetStateArgsForCall(1)
		gt.Expect(sID).To(Equal(*signed.Transaction.References[0]))
		gt.Expect(consumed).To(BeFalse())

		gt.Expect(fakeRepo.PutTransactionCallCount()).To(Equal(1))
		gt.Expect(fakeRepo.PutTransactionArgsForCall(0)).To(Equal(signed.Transaction)) // TODO: This should be a signed transaction

		gt.Expect(fakeRepo.PutStateCallCount()).To(Equal(1))
		gt.Expect(fakeRepo.PutStateArgsForCall(0)).To(Equal(signed.Transaction.Outputs[0]))

		gt.Expect(fakeRepo.ConsumeStateCallCount()).To(Equal(1))
		gt.Expect(fakeRepo.ConsumeStateArgsForCall(0)).To(Equal(*signed.Transaction.Inputs[0]))
	})

	t.Run("MissingInput", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		setup(t)
		signed.Inputs = append(signed.Inputs, &transaction.StateID{
			TxID: transaction.NewID([]byte("missing-input-txid")),
		})

		err := submitService.Submit(context.Background(), signed)
		gt.Expect(err).To(HaveOccurred())
		gt.Expect(store.IsNotFound(err)).To(BeTrue())
	})

	t.Run("MissingReference", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		setup(t)
		signed.References = append(signed.References, &transaction.StateID{
			TxID: transaction.NewID([]byte("missing-ref-txid")),
		})

		err := submitService.Submit(context.Background(), signed)
		gt.Expect(err).To(HaveOccurred())
		gt.Expect(store.IsNotFound(err)).To(BeTrue())
	})

	t.Run("PutTransactionFailure", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		setup(t)
		fakeRepo.PutTransactionReturns(errors.New("put-transaction-store-failure"))

		err := submitService.Submit(context.Background(), signed)
		gt.Expect(err).To(HaveOccurred())
		gt.Expect(err).To(MatchError("storing transaction " + signed.Transaction.ID.String() + " failed: put-transaction-store-failure"))
	})

	t.Run("PutStateFailure", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		setup(t)
		fakeRepo.PutStateReturns(errors.New("put-state-store-failure"))

		err := submitService.Submit(context.Background(), signed)
		gt.Expect(err).To(HaveOccurred())
		gt.Expect(err).To(MatchError("storing transaction output " + signed.Transaction.Outputs[0].ID.String() + " failed: put-state-store-failure"))
	})

	t.Run("ConsumeStateFailure", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		setup(t)
		fakeRepo.ConsumeStateReturns(errors.New("consume-state-store-failure"))

		err := submitService.Submit(context.Background(), signed)
		gt.Expect(err).To(HaveOccurred())
		gt.Expect(err).To(MatchError("consuming transaction state " + signed.Inputs[0].String() + " failed: consume-state-store-failure"))
	})

	t.Run("ValidateFailure", func(t *testing.T) {
		t.Skip("Implement after validate")
	})
}

func TestSubmitGetTransaction(t *testing.T) {
	signed := &transaction.Signed{
		Transaction: &transaction.Transaction{
			ID: transaction.NewID([]byte("transaction-id")),
		},
	}

	t.Run("ExistingTransaction", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		fakeRepo := &fake.Repository{}
		submitService := NewService(fakeRepo)

		err := submitService.Submit(context.Background(), signed)
		gt.Expect(err).To(HaveOccurred())
		gt.Expect(store.IsAlreadyExists(err)).To(BeTrue())
	})

	t.Run("UnexpectedError", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		fakeRepo := &fake.Repository{}
		fakeRepo.GetTransactionReturns(nil, errors.New("unexpected-error"))
		submitService := NewService(fakeRepo)

		err := submitService.Submit(context.Background(), signed)
		gt.Expect(err).To(MatchError("unexpected-error"))
	})
}

func TestValidate(t *testing.T) {
}
