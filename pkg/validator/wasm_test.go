// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package validator

import (
	"crypto"
	"crypto/elliptic"
	"crypto/rand"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/bytecodealliance/wasmtime-go"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	"github.com/sykesm/batik/pkg/ecdsautil"
	validationv1 "github.com/sykesm/batik/pkg/pb/validation/v1"
	"github.com/sykesm/batik/pkg/transaction"
)

var moddir = filepath.Join("..", "..", "wasm", "modules", "utxotx")

func TestMain(m *testing.M) {
	cmd := exec.Command("cargo", "build", "--target", "wasm32-unknown-unknown")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = moddir

	if err := cmd.Run(); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
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
			RequiredSigners: []*transaction.Party{{
				PublicKey: pk,
			}},
			Signatures: []*transaction.Signature{{
				PublicKey: pk,
				Signature: sig,
			}},
		}, signer
	}

	gt := NewGomegaWithT(t)
	engine := wasmtime.NewEngine()
	modfile := filepath.Join(moddir, "target", "wasm32-unknown-unknown", "debug", "utxotx.wasm")
	gt.Expect(modfile).To(BeAnExistingFile())
	module, err := ioutil.ReadFile(modfile)
	gt.Expect(err).NotTo(HaveOccurred())

	validateWASM := func(in *transaction.Resolved) error {
		wasm, err := NewWASM(engine, module)
		if err != nil {
			return err
		}
		resp, err := wasm.Validate(
			&validationv1.ValidateRequest{ResolvedTransaction: transaction.FromResolved(in)},
		)
		if err != nil {
			return err
		}
		if !resp.Valid {
			if msg := resp.ErrorMessage; msg != "" {
				return errors.Errorf("validation failed for transaction %s: %s", in.ID, msg)
			}
			return errors.Errorf("validation failed for transaction %s", in.ID)
		}
		return nil
	}

	t.Run("SucessfulValidation", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		resolvedTx, _ := setup(t)
		err := validateWASM(resolvedTx)
		gt.Expect(err).NotTo(HaveOccurred())
	})

	t.Run("MissingRequiredSignature", func(t *testing.T) {
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
		gt := NewGomegaWithT(t)

		resolvedTx, _ := setup(t)
		resolvedTx.RequiredSigners[0].PublicKey = []byte("invalid-public-key")
		resolvedTx.Signatures[0].PublicKey = []byte("invalid-public-key")

		err := validateWASM(resolvedTx)
		gt.Expect(err).To(MatchError("validation failed for transaction 7472616e73616374696f6e2d6964: Incomplete data or invalid ASN1"))
	})

	t.Run("NilPublicKey", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		resolvedTx, _ := setup(t)
		resolvedTx.RequiredSigners[0].PublicKey = nil

		err := validateWASM(resolvedTx)
		gt.Expect(err).To(MatchError("validation failed for transaction 7472616e73616374696f6e2d6964: Encountered an empty buffer decoding ASN1 block."))
	})

	t.Run("InvalidSignatureFormat", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		resolvedTx, _ := setup(t)
		resolvedTx.Signatures[0].Signature = []byte("bad-signature")

		err := validateWASM(resolvedTx)
		gt.Expect(err).To(MatchError("validation failed for transaction 7472616e73616374696f6e2d6964: signature error"))
	})

	t.Run("BadSignature", func(t *testing.T) {
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
