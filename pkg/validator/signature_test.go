// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package validator

import (
	"crypto"
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/sykesm/batik/pkg/ecdsautil"
	"github.com/sykesm/batik/pkg/transaction"
)

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
