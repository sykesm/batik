// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package validator

import (
	"crypto"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/bytecodealliance/wasmtime-go"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"

	"github.com/sykesm/batik/pkg/ecdsautil"
	validationv1 "github.com/sykesm/batik/pkg/pb/validation/v1"
	"github.com/sykesm/batik/pkg/transaction"
)

type validator interface {
	Validate(*validationv1.ValidateRequest) (*validationv1.ValidateResponse, error)
}

// These tests attempt to keep the Rust/WASM and native Go validation behavior
// in sync.
func TestValidate(t *testing.T) {
	gt := NewGomegaWithT(t)

	sk, err := ecdsautil.GenerateKey(elliptic.P256(), rand.Reader)
	gt.Expect(err).NotTo(HaveOccurred())
	pk, err := ecdsautil.MarshalPublicKey(&sk.PublicKey)
	gt.Expect(err).NotTo(HaveOccurred())
	signer := ecdsautil.NewSigner(sk)
	txidHash := digest([]byte("transaction-id"))
	sig, err := signer.Sign(rand.Reader, txidHash[:], crypto.SHA256)
	gt.Expect(err).NotTo(HaveOccurred())

	sk2, err := ecdsautil.GenerateKey(elliptic.P256(), rand.Reader)
	gt.Expect(err).NotTo(HaveOccurred())
	pk2, err := ecdsautil.MarshalPublicKey(&sk2.PublicKey)
	gt.Expect(err).NotTo(HaveOccurred())

	tests := []struct {
		desc       string
		validator  validator
		setupTx    func(*transaction.Resolved)
		valid      bool
		errMessage types.GomegaMatcher
		errMatcher types.GomegaMatcher
	}{
		{
			desc:       "SuccessfulValidation",
			setupTx:    nil,
			valid:      true,
			errMessage: BeEmpty(),
			errMatcher: nil,
		},
		{
			desc: "MissingRequiredSignature",
			setupTx: func(tx *transaction.Resolved) {
				tx.RequiredSigners = append(tx.RequiredSigners, &transaction.Party{PublicKey: pk2})
			},
			valid:      false,
			errMessage: Equal("missing signature from " + hex.EncodeToString(pk2)),
			errMatcher: nil,
		},
		{
			desc: "InvalidPublicKey",
			setupTx: func(tx *transaction.Resolved) {
				tx.RequiredSigners[0].PublicKey = []byte("invalid-public-key")
				tx.Signatures[0].PublicKey = []byte("invalid-public-key")
			},
			valid:      false,
			errMessage: ContainSubstring("unmarshal public key failed"),
			errMatcher: nil,
		},
		{
			desc: "NilPublicKey",
			setupTx: func(tx *transaction.Resolved) {
				tx.RequiredSigners[0].PublicKey = nil
			},
			valid:      false,
			errMessage: ContainSubstring("required signer missing public key"),
			errMatcher: nil,
		},
		{
			desc: "InvalidSignatureFormat",
			setupTx: func(tx *transaction.Resolved) {
				tx.Signatures[0].Signature = []byte("bad-signature")
			},
			valid:      false,
			errMessage: ContainSubstring("failed unmarshalling signature"),
			errMatcher: nil,
		},
		{
			desc: "BadSignature",
			setupTx: func(tx *transaction.Resolved) {
				newTxidHash := digest([]byte("this-is-a-different-message"))
				sig, err := signer.Sign(rand.Reader, newTxidHash[:], crypto.SHA256)
				gt.Expect(err).NotTo(HaveOccurred())
				tx.Signatures[0].Signature = sig
			},
			valid:      false,
			errMessage: ContainSubstring("signature verification failed"),
			errMatcher: nil,
		},
	}

	engine := wasmtime.NewEngine()
	modfile := filepath.Join(moddir, "target", "wasm32-unknown-unknown", "debug", "utxotx.wasm")
	gt.Expect(modfile).To(BeAnExistingFile())
	module, err := ioutil.ReadFile(modfile)
	gt.Expect(err).NotTo(HaveOccurred())

	validators := []struct {
		name string
		ctor func() (validator, error)
	}{
		{
			name: "Native",
			ctor: func() (validator, error) { return NewSignature(), nil },
		},
		{
			name: "WASM",
			ctor: func() (validator, error) { return NewWASM(engine, module) },
		},
	}

	for _, v := range validators {
		t.Run(v.name, func(t *testing.T) {
			gt := NewGomegaWithT(t)
			validator, err := v.ctor()
			gt.Expect(err).NotTo(HaveOccurred())

			for _, tt := range tests {
				t.Run(tt.desc, func(t *testing.T) {
					gt := NewGomegaWithT(t)

					resolvedTx := &transaction.Resolved{
						ID:              []byte("transaction-id"),
						RequiredSigners: []*transaction.Party{{PublicKey: pk}},
						Signatures:      []*transaction.Signature{{PublicKey: pk, Signature: sig}},
					}
					if tt.setupTx != nil {
						tt.setupTx(resolvedTx)
					}
					resp, err := validator.Validate(&validationv1.ValidateRequest{
						ResolvedTransaction: transaction.FromResolved(resolvedTx),
					})
					if tt.errMatcher != nil {
						gt.Expect(err).To(tt.errMatcher)
						return
					}
					gt.Expect(resp.Valid).To(Equal(tt.valid))
					gt.Expect(resp.ErrorMessage).To(tt.errMessage)
				})
			}
		})
	}
}
