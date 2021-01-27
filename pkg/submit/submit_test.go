// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package submit

import (
	"context"
	"crypto/sha256"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	txv1 "github.com/sykesm/batik/pkg/pb/tx/v1"
	validationv1 "github.com/sykesm/batik/pkg/pb/validation/v1"
	"github.com/sykesm/batik/pkg/store"
	"github.com/sykesm/batik/pkg/submit/fake"
	. "github.com/sykesm/batik/pkg/tested/matcher"
	"github.com/sykesm/batik/pkg/transaction"
)

//go:generate counterfeiter -o fake/repository.go --fake-name Repository . fakeRepository
type fakeRepository Repository // private to prevent an import cycle in generated fake

var _ fakeRepository = (*fake.Repository)(nil)

type validatorFunc func(req *validationv1.ValidateRequest) (*validationv1.ValidateResponse, error)

func (v validatorFunc) Validate(req *validationv1.ValidateRequest) (*validationv1.ValidateResponse, error) {
	return v(req)
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
		gt.Expect(err).To(MatchError(ErrHalt))
		gt.Expect(err).To(MatchError("transaction store failure: halt processing: unexpected-error"))
	})
}

func TestSubmitTxResolve(t *testing.T) {
	var (
		fakeRepo *fake.Repository
		tx       transaction.Transaction
		signed   *transaction.Signed
	)

	noopValidator := func(r *validationv1.ValidateRequest) (*validationv1.ValidateResponse, error) {
		return &validationv1.ValidateResponse{Valid: true}, nil
	}

	setup := func(t *testing.T) {
		tx = transaction.Transaction{
			ID: transaction.NewID([]byte("txid-3")),
			Inputs: []*transaction.StateID{
				newStateID("txid-1", 1),
			},
			References: []*transaction.StateID{
				newStateID("txid-2", 0),
			},
			Outputs: []*transaction.State{{
				ID:        *newStateID("txid-3", 0),
				StateInfo: &transaction.StateInfo{},
				Data:      []byte("output-0-data"),
			}},
			Parameters: []*transaction.Parameter{{
				Name:  "name-1",
				Value: []byte("value-1"),
			}},
			RequiredSigners: []*transaction.Party{{
				PublicKey: []byte("public-key-signer"),
			}},
			Tx:      &txv1.Transaction{},
			Encoded: []byte("encoded-transaction"),
		}
		signed = &transaction.Signed{
			Transaction: &tx,
			Signatures: []*transaction.Signature{{
				PublicKey: []byte("public-key"),
				Signature: []byte("signature"),
			}},
		}

		fakeRepo = &fake.Repository{}
		fakeRepo.GetTransactionReturns(nil, &store.NotFoundError{Err: errors.New("missing-transaction-error")})
		fakeRepo.GetStateStub = func(sid transaction.StateID, consumed bool) (*transaction.State, error) {
			switch {
			case sid.Equals(*newStateID("txid-1", 1)):
				return &transaction.State{
					ID:        sid,
					StateInfo: &transaction.StateInfo{Kind: "kind-1"},
					Data:      []byte("txid-1:1-state-data"),
				}, nil
			case sid.Equals(*newStateID("txid-2", 0)):
				return &transaction.State{
					ID:        sid,
					StateInfo: &transaction.StateInfo{Kind: "kind-2"},
					Data:      []byte("txid-2:0-state-data"),
				}, nil
			default:
				return nil, &store.NotFoundError{Err: errors.Errorf("missing-state %s", sid)}
			}
		}
	}

	t.Run("Success", func(t *testing.T) {
		setup(t)
		gt := NewGomegaWithT(t)

		var req *validationv1.ValidateRequest
		validator := func(r *validationv1.ValidateRequest) (*validationv1.ValidateResponse, error) {
			req = r
			return &validationv1.ValidateResponse{Valid: true}, nil
		}
		submitService := &Service{
			repo:      fakeRepo,
			validator: validatorFunc(validator),
		}

		err := submitService.Submit(context.Background(), signed)
		gt.Expect(err).NotTo(HaveOccurred())

		gt.Expect(fakeRepo.GetStateCallCount()).To(Equal(2))
		input, consumed := fakeRepo.GetStateArgsForCall(0)
		gt.Expect(input).To(Equal(*tx.Inputs[0]))
		gt.Expect(consumed).To(BeFalse())
		ref, consumed := fakeRepo.GetStateArgsForCall(1)
		gt.Expect(ref).To(Equal(*tx.References[0]))
		gt.Expect(consumed).To(BeFalse())

		gt.Expect(req).NotTo(BeNil())
		gt.Expect(req).To(ProtoEqual(&validationv1.ValidateRequest{
			ResolvedTransaction: &validationv1.ResolvedTransaction{
				Txid: []byte("txid-3"),
				Inputs: []*validationv1.ResolvedState{{
					Reference: &txv1.StateReference{
						Txid:        []byte("txid-1"),
						OutputIndex: 1,
					},
					State: &txv1.State{
						Info:  &txv1.StateInfo{Kind: "kind-1"},
						State: []byte("txid-1:1-state-data"),
					},
				}},
				References: []*validationv1.ResolvedState{{
					Reference: &txv1.StateReference{
						Txid:        []byte("txid-2"),
						OutputIndex: 0,
					},
					State: &txv1.State{
						Info:  &txv1.StateInfo{Kind: "kind-2"},
						State: []byte("txid-2:0-state-data"),
					},
				}},
				Outputs: []*txv1.State{{
					Info:  &txv1.StateInfo{},
					State: []byte("output-0-data"),
				}},
				Parameters: []*txv1.Parameter{{
					Name:  "name-1",
					Value: []byte("value-1"),
				}},
				RequiredSigners: []*txv1.Party{{
					PublicKey: []byte("public-key-signer"),
				}},
				Signatures: []*txv1.Signature{{
					PublicKey: []byte("public-key"),
					Signature: []byte("signature"),
				}},
			},
		}))
	})

	t.Run("WhenInputMissing", func(t *testing.T) {
		setup(t)
		gt := NewGomegaWithT(t)

		submitService := &Service{
			repo:      fakeRepo,
			validator: validatorFunc(noopValidator),
		}

		fakeRepo.GetStateReturnsOnCall(0, nil, &store.NotFoundError{Err: errors.New("missing-input-state")})

		err := submitService.Submit(context.Background(), signed)
		gt.Expect(err).To(MatchError(ContainSubstring("missing-input-state")))
		gt.Expect(store.IsNotFound(err)).To(BeTrue())
		input, _ := fakeRepo.GetStateArgsForCall(0)
		gt.Expect(input).To(Equal(*tx.Inputs[0]))
	})

	t.Run("WhenGetInputFailure", func(t *testing.T) {
		setup(t)
		gt := NewGomegaWithT(t)

		submitService := &Service{
			repo:      fakeRepo,
			validator: validatorFunc(noopValidator),
		}

		fakeRepo.GetStateReturnsOnCall(0, nil, errors.New("get-input-state-failed"))

		err := submitService.Submit(context.Background(), signed)
		gt.Expect(err).To(MatchError(ErrHalt))
		gt.Expect(err).To(MatchError(MatchRegexp("state resolution for transaction [[:xdigit:]]+ failed: halt processing: get-input-state-failed")), err.Error())
		input, _ := fakeRepo.GetStateArgsForCall(0)
		gt.Expect(input).To(Equal(*tx.Inputs[0]))
	})

	t.Run("WhenReferenceMissing", func(t *testing.T) {
		setup(t)
		gt := NewGomegaWithT(t)

		submitService := &Service{
			repo:      fakeRepo,
			validator: validatorFunc(noopValidator),
		}

		fakeRepo.GetStateReturnsOnCall(1, nil, &store.NotFoundError{Err: errors.New("missing-reference-state")})

		err := submitService.Submit(context.Background(), signed)
		gt.Expect(err).To(MatchError(ContainSubstring("missing-reference-state")))
		gt.Expect(store.IsNotFound(err)).To(BeTrue())
		ref, _ := fakeRepo.GetStateArgsForCall(1)
		gt.Expect(ref).To(Equal(*tx.References[0]))
	})

	t.Run("WhenGetReferenceFailure", func(t *testing.T) {
		setup(t)
		gt := NewGomegaWithT(t)

		submitService := &Service{
			repo:      fakeRepo,
			validator: validatorFunc(noopValidator),
		}

		fakeRepo.GetStateReturnsOnCall(1, nil, errors.New("get-reference-state-failed"))

		err := submitService.Submit(context.Background(), signed)
		gt.Expect(err).To(MatchError(ErrHalt))
		gt.Expect(err).To(MatchError(MatchRegexp("state resolution for transaction [[:xdigit:]]+ failed: halt processing: get-reference-state-failed")), err.Error())
		ref, _ := fakeRepo.GetStateArgsForCall(1)
		gt.Expect(ref).To(Equal(*tx.References[0]))
	})
}

func TestSubmitPostValidation(t *testing.T) {
	var (
		fakeRepo *fake.Repository
		tx       transaction.Transaction
		signed   *transaction.Signed
	)

	noopValidator := func(req *validationv1.ValidateRequest) (*validationv1.ValidateResponse, error) {
		return &validationv1.ValidateResponse{Valid: true}, nil
	}

	setup := func(t *testing.T) {
		tx = transaction.Transaction{
			ID: transaction.NewID([]byte("txid-3")),
			Inputs: []*transaction.StateID{
				newStateID("txid-1", 1),
			},
			References: []*transaction.StateID{
				newStateID("txid-2", 0),
			},
			Outputs: []*transaction.State{{
				ID:        *newStateID("txid-3", 0),
				StateInfo: &transaction.StateInfo{},
				Data:      []byte("output-0-data"),
			}},
			Tx:      &txv1.Transaction{},
			Encoded: []byte("encoded-transaction"),
		}
		signed = &transaction.Signed{
			Transaction: &tx,
			Signatures: []*transaction.Signature{{
				PublicKey: []byte("public-key"),
				Signature: []byte("signature"),
			}},
		}

		fakeRepo = &fake.Repository{}
		fakeRepo.GetTransactionReturns(nil, &store.NotFoundError{Err: errors.New("missing-transaction-error")})
		fakeRepo.GetStateStub = func(sid transaction.StateID, consumed bool) (*transaction.State, error) {
			switch {
			case sid.Equals(*newStateID("txid-1", 1)):
				return &transaction.State{ID: sid, Data: []byte("txid-1:1-state-data")}, nil
			case sid.Equals(*newStateID("txid-2", 0)):
				return &transaction.State{ID: sid, Data: []byte("txid-2:0-state-data")}, nil
			default:
				return nil, &store.NotFoundError{Err: errors.Errorf("missing-state %s", sid)}
			}
		}
	}

	t.Run("Success", func(t *testing.T) {
		setup(t)
		gt := NewGomegaWithT(t)

		submitService := &Service{
			repo:      fakeRepo,
			validator: validatorFunc(noopValidator),
		}

		err := submitService.Submit(context.Background(), signed)
		gt.Expect(err).NotTo(HaveOccurred())

		gt.Expect(fakeRepo.PutTransactionCallCount()).To(Equal(1))
		gt.Expect(fakeRepo.PutTransactionArgsForCall(0)).To(Equal(&tx))
		gt.Expect(fakeRepo.PutStateCallCount()).To(Equal(1))
		gt.Expect(fakeRepo.PutStateArgsForCall(0)).To(Equal(tx.Outputs[0]))
		gt.Expect(fakeRepo.ConsumeStateCallCount()).To(Equal(1))
		gt.Expect(fakeRepo.ConsumeStateArgsForCall(0)).To(Equal(*tx.Inputs[0]))
	})

	t.Run("WhenInvalid", func(t *testing.T) {
		setup(t)
		gt := NewGomegaWithT(t)

		submitService := &Service{
			repo: fakeRepo,
			validator: validatorFunc(func(req *validationv1.ValidateRequest) (*validationv1.ValidateResponse, error) {
				return &validationv1.ValidateResponse{Valid: false}, nil
			}),
		}

		err := submitService.Submit(context.Background(), signed)
		gt.Expect(err).To(MatchError(ContainSubstring("validation failed")))
		gt.Expect(err).NotTo(MatchError(ErrHalt))

		gt.Expect(fakeRepo.PutTransactionCallCount()).To(Equal(0))
		gt.Expect(fakeRepo.PutStateCallCount()).To(Equal(0))
		gt.Expect(fakeRepo.ConsumeStateCallCount()).To(Equal(0))
	})

	t.Run("WhenInvalidWithMessage", func(t *testing.T) {
		setup(t)
		gt := NewGomegaWithT(t)

		submitService := &Service{
			repo: fakeRepo,
			validator: validatorFunc(func(req *validationv1.ValidateRequest) (*validationv1.ValidateResponse, error) {
				return &validationv1.ValidateResponse{Valid: false, ErrorMessage: "texas-toast"}, nil
			}),
		}

		err := submitService.Submit(context.Background(), signed)
		gt.Expect(err).To(MatchError(ContainSubstring("validation failed: texas-toast")))
		gt.Expect(err).NotTo(MatchError(ErrHalt))

		gt.Expect(fakeRepo.PutTransactionCallCount()).To(Equal(0))
		gt.Expect(fakeRepo.PutStateCallCount()).To(Equal(0))
		gt.Expect(fakeRepo.ConsumeStateCallCount()).To(Equal(0))
	})

	t.Run("WhenValidationFails", func(t *testing.T) {
		setup(t)
		gt := NewGomegaWithT(t)

		submitService := &Service{
			repo: fakeRepo,
			validator: validatorFunc(func(req *validationv1.ValidateRequest) (*validationv1.ValidateResponse, error) {
				return nil, errors.New("boom!")
			}),
		}

		err := submitService.Submit(context.Background(), signed)
		gt.Expect(err).To(MatchError(ErrHalt))
		gt.Expect(err).To(MatchError("validator failed: halt processing: boom!"))

		gt.Expect(fakeRepo.PutTransactionCallCount()).To(Equal(0))
		gt.Expect(fakeRepo.PutStateCallCount()).To(Equal(0))
		gt.Expect(fakeRepo.ConsumeStateCallCount()).To(Equal(0))
	})

	t.Run("WhenPutTransactionFails", func(t *testing.T) {
		setup(t)
		gt := NewGomegaWithT(t)

		submitService := &Service{
			repo:      fakeRepo,
			validator: validatorFunc(noopValidator),
		}

		fakeRepo.PutTransactionReturns(errors.New("put-transaction-failed"))

		err := submitService.Submit(context.Background(), signed)
		gt.Expect(err).To(MatchError(ErrHalt))
		gt.Expect(err).To(MatchError(MatchRegexp("storing transaction [[:xdigit:]]+ failed: halt processing: put-transaction-failed")))

		gt.Expect(fakeRepo.PutTransactionCallCount()).To(Equal(1))
		gt.Expect(fakeRepo.PutStateCallCount()).To(Equal(0))
		gt.Expect(fakeRepo.ConsumeStateCallCount()).To(Equal(0))
	})

	t.Run("WhenPutStateFails", func(t *testing.T) {
		setup(t)
		gt := NewGomegaWithT(t)

		submitService := &Service{
			repo:      fakeRepo,
			validator: validatorFunc(noopValidator),
		}

		fakeRepo.PutStateReturns(errors.New("put-state-failed"))

		err := submitService.Submit(context.Background(), signed)
		gt.Expect(err).To(MatchError(ErrHalt))
		gt.Expect(err).To(MatchError(MatchRegexp("storing transaction output [[:xdigit:]]+:[[:xdigit:]]+ failed: halt processing: put-state-failed")))

		gt.Expect(fakeRepo.PutTransactionCallCount()).To(Equal(1))
		gt.Expect(fakeRepo.PutStateCallCount()).To(Equal(1))
		gt.Expect(fakeRepo.ConsumeStateCallCount()).To(Equal(0))
	})

	t.Run("WhenConsumeStateFails", func(t *testing.T) {
		setup(t)
		gt := NewGomegaWithT(t)

		submitService := &Service{
			repo:      fakeRepo,
			validator: validatorFunc(noopValidator),
		}

		fakeRepo.ConsumeStateReturns(errors.New("consume-state-failed"))

		err := submitService.Submit(context.Background(), signed)
		gt.Expect(err).To(MatchError(ErrHalt))
		gt.Expect(err).To(MatchError(MatchRegexp("consuming transaction state [[:xdigit:]]+:[[:xdigit:]]+ failed: halt processing: consume-state-failed")))

		gt.Expect(fakeRepo.PutTransactionCallCount()).To(Equal(1))
		gt.Expect(fakeRepo.PutStateCallCount()).To(Equal(1))
		gt.Expect(fakeRepo.ConsumeStateCallCount()).To(Equal(1))
	})
}

func digest(preImage []byte) []byte {
	sum := sha256.Sum256(preImage)
	return sum[:]
}

func newStateID(txid string, index uint64) *transaction.StateID {
	return &transaction.StateID{
		TxID:        transaction.NewID([]byte(txid)),
		OutputIndex: index,
	}
}
