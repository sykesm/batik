// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package submit

import (
	"context"
	"crypto"
	"crypto/elliptic"
	"crypto/rand"
	"os"
	"os/exec"
	"path/filepath"
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

var modulePath = filepath.Join("..", "..", "wasm", "modules", "utxotx")

func TestMain(m *testing.M) {
	cmd := exec.Command("cargo", "build", "--target", "wasm32-unknown-unknown")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = modulePath

	if err := cmd.Run(); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

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
		txidHash := digest([]byte("transaction-id"))
		sig, err := ecdsautil.NewSigner(sk).Sign(rand.Reader, txidHash[:], crypto.SHA256)
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

		err := submitService.Submit(context.Background(), signed, "")
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

		err := submitService.Submit(context.Background(), signed, "")
		gt.Expect(err).To(HaveOccurred())
		gt.Expect(store.IsNotFound(err)).To(BeTrue())
	})

	t.Run("MissingReference", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		setup(t)
		signed.References = append(signed.References, &transaction.StateID{
			TxID: transaction.NewID([]byte("missing-ref-txid")),
		})

		err := submitService.Submit(context.Background(), signed, "")
		gt.Expect(err).To(HaveOccurred())
		gt.Expect(store.IsNotFound(err)).To(BeTrue())
	})

	t.Run("PutTransactionFailure", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		setup(t)
		fakeRepo.PutTransactionReturns(errors.New("put-transaction-store-failure"))

		err := submitService.Submit(context.Background(), signed, "")
		gt.Expect(err).To(HaveOccurred())
		gt.Expect(err).To(MatchError("storing transaction " + signed.Transaction.ID.String() + " failed: put-transaction-store-failure"))
	})

	t.Run("PutStateFailure", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		setup(t)
		fakeRepo.PutStateReturns(errors.New("put-state-store-failure"))

		err := submitService.Submit(context.Background(), signed, "")
		gt.Expect(err).To(HaveOccurred())
		gt.Expect(err).To(MatchError("storing transaction output " + signed.Transaction.Outputs[0].ID.String() + " failed: put-state-store-failure"))
	})

	t.Run("ConsumeStateFailure", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		setup(t)
		fakeRepo.ConsumeStateReturns(errors.New("consume-state-store-failure"))

		err := submitService.Submit(context.Background(), signed, "")
		gt.Expect(err).To(HaveOccurred())
		gt.Expect(err).To(MatchError("consuming transaction state " + signed.Inputs[0].String() + " failed: consume-state-store-failure"))
	})

	t.Run("ValidateFailure", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		setup(t)
		signed.Signatures[0].Signature = []byte("bad-signature")

		err := submitService.Submit(context.Background(), signed, "")
		gt.Expect(err).To(MatchError(ContainSubstring("validation failed: ")))
	})

	t.Run("BasicTransactionWASMValidation", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		setup(t)

		err := submitService.Submit(context.Background(), signed, "utxo-wasm")
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

	t.Run("ValidateWASMFailure", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		setup(t)
		signed.Signatures[0].Signature = []byte("bad-signature")

		err := submitService.Submit(context.Background(), signed, "utxo-wasm")
		gt.Expect(err).To(MatchError(ContainSubstring("validation failed: ")))
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

		err := submitService.Submit(context.Background(), signed, "")
		gt.Expect(err).To(HaveOccurred())
		gt.Expect(store.IsAlreadyExists(err)).To(BeTrue())
	})

	t.Run("UnexpectedError", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		fakeRepo := &fake.Repository{}
		fakeRepo.GetTransactionReturns(nil, errors.New("unexpected-error"))
		submitService := NewService(fakeRepo)

		err := submitService.Submit(context.Background(), signed, "")
		gt.Expect(err).To(MatchError("unexpected-error"))
	})
}

func TestValidate(t *testing.T) {
	var (
		resolvedTx *transaction.Resolved
		signer     *ecdsautil.Signer
	)

	setup := func(t *testing.T) {
		gt := NewGomegaWithT(t)

		sk, err := ecdsautil.GenerateKey(elliptic.P256(), rand.Reader)
		gt.Expect(err).NotTo(HaveOccurred())
		pk, err := ecdsautil.MarshalPublicKey(&sk.PublicKey)
		gt.Expect(err).NotTo(HaveOccurred())
		signer = ecdsautil.NewSigner(sk)
		txidHash := digest([]byte("transaction-id"))
		sig, err := signer.Sign(rand.Reader, txidHash[:], crypto.SHA256)
		gt.Expect(err).NotTo(HaveOccurred())

		resolvedTx = &transaction.Resolved{
			ID: []byte("transaction-id"),
			RequiredSigners: []*transaction.Party{
				{
					PublicKey: pk,
				},
			},
			Signatures: []*transaction.Signature{
				{
					PublicKey: pk,
					Signature: sig,
				},
			},
		}
	}

	t.Run("SucessfulValidation", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		setup(t)
		err := validate(resolvedTx)
		gt.Expect(err).NotTo(HaveOccurred())
	})

	t.Run("MissingRequiredSignature", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		setup(t)
		resolvedTx.RequiredSigners = append(
			resolvedTx.RequiredSigners,
			&transaction.Party{PublicKey: []byte("absent-required-signer")},
		)

		err := validate(resolvedTx)
		gt.Expect(err).To(MatchError("missing signature from 616273656e742d72657175697265642d7369676e6572"))
	})

	t.Run("InvalidPublicKey", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		setup(t)
		resolvedTx.RequiredSigners[0].PublicKey = []byte("invalid-public-key")
		resolvedTx.Signatures[0].PublicKey = []byte("invalid-public-key")

		err := validate(resolvedTx)
		gt.Expect(err).To(MatchError(ContainSubstring("unmarshal public key failed")))
	})

	t.Run("NilPublicKey", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		setup(t)
		resolvedTx.RequiredSigners[0].PublicKey = nil

		err := validate(resolvedTx)
		gt.Expect(err).To(MatchError(ContainSubstring("required signer missing public key")))
	})

	t.Run("InvalidSignatureFormat", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		setup(t)
		resolvedTx.Signatures[0].Signature = []byte("bad-signature")

		err := validate(resolvedTx)
		gt.Expect(err).To(MatchError(ContainSubstring("failed unmarshalling signature")))
	})

	t.Run("BadSignature", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		setup(t)
		newTxidHash := digest([]byte("this-is-a-different-message"))
		sig, err := signer.Sign(rand.Reader, newTxidHash[:], crypto.SHA256)
		gt.Expect(err).NotTo(HaveOccurred())
		resolvedTx.Signatures[0].Signature = sig

		err = validate(resolvedTx)
		gt.Expect(err).To(MatchError("signature verification failed"))
	})
}

func TestValidateWASM(t *testing.T) {
	setup := func(t *testing.T) (*transaction.Resolved, *ecdsautil.Signer) {
		gt := NewGomegaWithT(t)

		sk, err := ecdsautil.GenerateKey(elliptic.P256(), rand.Reader)
		gt.Expect(err).NotTo(HaveOccurred())
		pk, err := ecdsautil.MarshalPublicKey(&sk.PublicKey)
		gt.Expect(err).NotTo(HaveOccurred())
		signer := ecdsautil.NewSigner(sk)
		txidHash := digest([]byte("transaction-id"))
		sig, err := signer.Sign(rand.Reader, txidHash[:], crypto.SHA256)
		gt.Expect(err).NotTo(HaveOccurred())

		return &transaction.Resolved{
			ID: []byte("transaction-id"),
			RequiredSigners: []*transaction.Party{
				{
					PublicKey: pk,
				},
			},
			Signatures: []*transaction.Signature{
				{
					PublicKey: pk,
					Signature: sig,
				},
			},
		}, signer
	}

	t.Run("SucessfulValidation", func(t *testing.T) {
		t.Parallel()

		gt := NewGomegaWithT(t)

		resolvedTx, _ := setup(t)
		err := validateWASM(resolvedTx)
		gt.Expect(err).NotTo(HaveOccurred())
	})

	t.Run("MissingRequiredSignature", func(t *testing.T) {
		t.Parallel()

		gt := NewGomegaWithT(t)

		resolvedTx, _ := setup(t)
		resolvedTx.RequiredSigners = append(
			resolvedTx.RequiredSigners,
			&transaction.Party{PublicKey: []byte("absent-required-signer")},
		)

		err := validateWASM(resolvedTx)
		gt.Expect(err).To(MatchError("validation failed for transaction 7472616e73616374696f6e2d6964: Incomplete data or invalid ASN1"))
	})

	t.Run("InvalidPublicKey", func(t *testing.T) {
		t.Parallel()

		gt := NewGomegaWithT(t)

		resolvedTx, _ := setup(t)
		resolvedTx.RequiredSigners[0].PublicKey = []byte("invalid-public-key")
		resolvedTx.Signatures[0].PublicKey = []byte("invalid-public-key")

		err := validateWASM(resolvedTx)
		gt.Expect(err).To(MatchError("validation failed for transaction 7472616e73616374696f6e2d6964: Incomplete data or invalid ASN1"))
	})

	t.Run("NilPublicKey", func(t *testing.T) {
		t.Parallel()

		gt := NewGomegaWithT(t)

		resolvedTx, _ := setup(t)
		resolvedTx.RequiredSigners[0].PublicKey = nil

		err := validateWASM(resolvedTx)
		gt.Expect(err).To(MatchError("validation failed for transaction 7472616e73616374696f6e2d6964: Encountered an empty buffer decoding ASN1 block."))
	})

	t.Run("InvalidSignatureFormat", func(t *testing.T) {
		t.Parallel()

		gt := NewGomegaWithT(t)

		resolvedTx, _ := setup(t)
		resolvedTx.Signatures[0].Signature = []byte("bad-signature")

		err := validateWASM(resolvedTx)
		gt.Expect(err).To(MatchError("validation failed for transaction 7472616e73616374696f6e2d6964: signature error"))
	})

	t.Run("BadSignature", func(t *testing.T) {
		t.Parallel()

		gt := NewGomegaWithT(t)

		resolvedTx, signer := setup(t)
		newTxidHash := digest([]byte("this-is-a-different-message"))
		sig, err := signer.Sign(rand.Reader, newTxidHash[:], crypto.SHA256)
		gt.Expect(err).NotTo(HaveOccurred())
		resolvedTx.Signatures[0].Signature = sig

		err = validateWASM(resolvedTx)
		gt.Expect(err).To(MatchError("validation failed for transaction 7472616e73616374696f6e2d6964: signature error"))
	})
}
