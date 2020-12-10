// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ecdsautil

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"math/big"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

func TestECDSA(t *testing.T) {
	tests := map[string]struct {
		curve elliptic.Curve
	}{
		"P-224": {curve: elliptic.P224()},
		"P-256": {curve: elliptic.P256()},
		"P-384": {curve: elliptic.P384()},
		"P-521": {curve: elliptic.P521()},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)
			pk, err := GenerateKey(tt.curve, rand.Reader)
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(pk.Curve).To(Equal(tt.curve))

			t.Run("Marshaling", func(t *testing.T) {
				gt := NewGomegaWithT(t)
				mpk, err := MarshalPublicKey(&pk.PublicKey)
				gt.Expect(err).NotTo(HaveOccurred())
				upk, err := UnmarshalPublicKey(mpk)
				gt.Expect(err).NotTo(HaveOccurred())
				gt.Expect(upk).To(Equal(&pk.PublicKey))
			})

			t.Run("Signing", func(t *testing.T) {
				gt := NewGomegaWithT(t)
				otherKey, err := GenerateKey(tt.curve, rand.Reader)
				gt.Expect(err).NotTo(HaveOccurred())
				hash := sha256.Sum256([]byte("this-is-a-message"))

				sig1, err := Sign(rand.Reader, pk, hash[:])
				gt.Expect(err).NotTo(HaveOccurred())
				ok, err := Verify(&pk.PublicKey, sig1, hash[:])
				gt.Expect(err).NotTo(HaveOccurred())
				gt.Expect(ok).To(BeTrue())

				ok, err = Verify(&otherKey.PublicKey, sig1, hash[:])
				gt.Expect(err).NotTo(HaveOccurred())
				gt.Expect(ok).To(BeFalse(), "signature must be invalid for another key")

				sig2, err := NewSigner(pk).Sign(rand.Reader, hash[:], crypto.SHA256)
				gt.Expect(err).NotTo(HaveOccurred())
				ok, err = Verify(&pk.PublicKey, sig2, hash[:])
				gt.Expect(err).NotTo(HaveOccurred())
				gt.Expect(ok).To(BeTrue())

				ok, err = Verify(&otherKey.PublicKey, sig1, hash[:])
				gt.Expect(err).NotTo(HaveOccurred())
				gt.Expect(ok).To(BeFalse(), "signature must be invalid for another key")
			})

			t.Run("VerifyNotNormalized", func(t *testing.T) {
				gt := NewGomegaWithT(t)
				hash := sha256.Sum256([]byte("this-is-a-message"))

				sig, err := Sign(rand.Reader, pk, hash[:])
				gt.Expect(err).NotTo(HaveOccurred())
				ok, err := Verify(&pk.PublicKey, sig, hash[:])
				gt.Expect(err).NotTo(HaveOccurred())
				gt.Expect(ok).To(BeTrue())

				r, s, err := UnmarshalECDSASignature(sig)
				gt.Expect(err).NotTo(HaveOccurred())
				isLowS, err := IsLowS(&pk.PublicKey, s)
				gt.Expect(err).NotTo(HaveOccurred())
				gt.Expect(isLowS).To(BeTrue())

				highS := new(big.Int).Mod(new(big.Int).Neg(s), pk.Curve.Params().N)
				isLowS, err = IsLowS(&pk.PublicKey, highS)
				gt.Expect(err).NotTo(HaveOccurred())
				gt.Expect(isLowS).To(BeFalse())

				// Signature is valid with standard library
				gt.Expect(ecdsa.Verify(&pk.PublicKey, hash[:], r, highS)).To(BeTrue())

				sig, err = MarshalECDSASignature(r, highS)
				gt.Expect(err).NotTo(HaveOccurred())
				_, err = Verify(&pk.PublicKey, sig, hash[:])
				gt.Expect(err).To(MatchError("s must be smaller than the half order of the curve"))
			})
		})
	}
}

func TestUnmarshalPublicKeyErrors(t *testing.T) {
	gt := NewGomegaWithT(t)

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	gt.Expect(err).NotTo(HaveOccurred())
	rsaKeyBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	gt.Expect(err).NotTo(HaveOccurred())

	tests := map[string]struct {
		pk         []byte
		errMatcher types.GomegaMatcher
	}{
		"NoKey":    {pk: nil, errMatcher: MatchError(ContainSubstring("asn1: syntax error"))},
		"EmptyKey": {pk: []byte{}, errMatcher: MatchError(ContainSubstring("asn1: syntax error"))},
		"BadKey":   {pk: []byte{255, 255, 255, 255}, errMatcher: MatchError(ContainSubstring("asn1: syntax error"))},
		"RSAKey":   {pk: rsaKeyBytes, errMatcher: MatchError("unsupported public key type: *rsa.PublicKey")},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)
			_, err := UnmarshalPublicKey(tt.pk)
			gt.Expect(err).To(tt.errMatcher)
		})
	}
}

func TestUnmarshalECDSASignatureBadSign(t *testing.T) {
	t.Run("NegativeR", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		sig, err := MarshalECDSASignature(big.NewInt(-1), big.NewInt(1))
		gt.Expect(err).NotTo(HaveOccurred())

		_, _, err = UnmarshalECDSASignature(sig)
		gt.Expect(err).To(MatchError(ContainSubstring("invalid signature")))
	})

	t.Run("NegativeS", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		sig, err := MarshalECDSASignature(big.NewInt(1), big.NewInt(-1))
		gt.Expect(err).NotTo(HaveOccurred())

		_, _, err = UnmarshalECDSASignature(sig)
		gt.Expect(err).To(MatchError(ContainSubstring("invalid signature")))
	})
}

func TestVerifyBadSignatureFormat(t *testing.T) {
	t.Run("CorruptSignature", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		key, err := GenerateKey(elliptic.P256(), rand.Reader)
		gt.Expect(err).NotTo(HaveOccurred())

		_, err = Verify(&key.PublicKey, []byte("garbage"), []byte{})
		gt.Expect(err).To(MatchError(ContainSubstring("failed unmarshalling signature: asn1: structure error")))
	})

	t.Run("LowSFailure", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		key, err := GenerateKey(elliptic.P256(), rand.Reader)
		gt.Expect(err).NotTo(HaveOccurred())

		hash := sha256.Sum256([]byte("this-is-a-message"))
		sig, err := Sign(rand.Reader, key, hash[:])
		gt.Expect(err).NotTo(HaveOccurred())

		pk := key.PublicKey
		pk.Curve = &elliptic.CurveParams{Name: "BAD"}
		_, err = Verify(&pk, sig, hash[:])
		gt.Expect(err).To(MatchError("unrecognized ecdsa curve"))
	})
}

func TestToLowSBadCurve(t *testing.T) {
	gt := NewGomegaWithT(t)
	key, err := GenerateKey(elliptic.P256(), rand.Reader)
	gt.Expect(err).NotTo(HaveOccurred())
	pk := key.PublicKey
	pk.Curve = &elliptic.CurveParams{Name: "BAD"}

	_, _, err = ToLowS(&pk, big.NewInt(1))
	gt.Expect(err).To(MatchError("unrecognized ecdsa curve"))
}

func TestSignFailure(t *testing.T) {
	gt := NewGomegaWithT(t)
	pk, err := GenerateKey(elliptic.P256(), rand.Reader)
	gt.Expect(err).NotTo(HaveOccurred())

	hash := sha256.Sum256([]byte("this-is-a-message"))
	_, err = Sign(bytes.NewBuffer(nil), pk, hash[:])
	gt.Expect(err).To(HaveOccurred())
}
