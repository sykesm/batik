// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package utxo

import (
	"crypto"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/bytecodealliance/wasmtime-go"
	. "github.com/onsi/gomega"

	"github.com/sykesm/batik/pkg/ecdsautil"
	txv1 "github.com/sykesm/batik/pkg/pb/tx/v1"
	validationv1 "github.com/sykesm/batik/pkg/pb/validation/v1"
)

var modulePath = filepath.Join("..", "..", "..", "wasm", "modules", "utxotx")

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

func TestUTXOValidator(t *testing.T) {
	var (
		validateRequest *validationv1.ValidateRequest
		validator       *UTXOValidator
		signer          *ecdsautil.Signer
	)

	gt := NewGomegaWithT(t)
	engine := wasmtime.NewEngine()
	store := wasmtime.NewStore(engine)
	modulefile := filepath.Join(modulePath, "target", "wasm32-unknown-unknown", "debug", "utxotx.wasm")
	gt.Expect(modulefile).To(BeAnExistingFile())
	module, err := wasmtime.NewModuleFromFile(engine, modulefile)
	gt.Expect(err).NotTo(HaveOccurred())

	setup := func(t *testing.T) {
		gt := NewGomegaWithT(t)

		var err error
		validator, err = NewValidator(store, module)
		gt.Expect(err).NotTo(HaveOccurred())

		sk, err := ecdsautil.GenerateKey(elliptic.P256(), rand.Reader)
		gt.Expect(err).NotTo(HaveOccurred())
		pk, err := ecdsautil.MarshalPublicKey(&sk.PublicKey)
		gt.Expect(err).NotTo(HaveOccurred())
		signer = ecdsautil.NewSigner(sk)
		txidHash := digest([]byte("transaction-id"))
		sig, err := signer.Sign(rand.Reader, txidHash[:], crypto.SHA256)
		gt.Expect(err).NotTo(HaveOccurred())

		validateRequest = &validationv1.ValidateRequest{
			ResolvedTransaction: &validationv1.ResolvedTransaction{
				Txid: []byte("transaction-id"),
				RequiredSigners: []*txv1.Party{
					{
						PublicKey: pk,
					},
				},
				Signatures: []*txv1.Signature{
					{
						PublicKey: pk,
						Signature: sig,
					},
				},
			},
		}
	}

	t.Run("SuccessfulValidation", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		setup(t)
		_, err := validator.Validate(validateRequest)
		gt.Expect(err).NotTo(HaveOccurred())
	})

	t.Run("BadSignature", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		setup(t)
		newTxidHash := digest([]byte("this-is-a-different-message"))
		sig, err := signer.Sign(rand.Reader, newTxidHash[:], crypto.SHA256)
		gt.Expect(err).NotTo(HaveOccurred())
		validateRequest.ResolvedTransaction.Signatures[0].Signature = sig

		res, err := validator.Validate(validateRequest)
		gt.Expect(err).NotTo(HaveOccurred())
		gt.Expect(res.Valid).To(BeFalse())
		gt.Expect(res.ErrorMessage).To(Equal("signature error"))
	})
}

func digest(preImage []byte) [32]byte {
	return sha256.Sum256(preImage)
}
