// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ecdsautil

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"io"

	"github.com/pkg/errors"
)

// GenerateKey simply delegates to ecdsa.GenerateKey.
func GenerateKey(curve elliptic.Curve, rand io.Reader) (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(curve, rand)
}

// MarshalPublicKey returns the PKIX, ASN.1 DER form of the ECDSA public key.
// This is the self-describing, distinguished encoding of the public key that
// includes the OID of the associated curve.
func MarshalPublicKey(pub *ecdsa.PublicKey) ([]byte, error) {
	return x509.MarshalPKIXPublicKey(pub)
}

// UnmarshalPublicKey parses a PKIX, ASN.1 DER form of an ECDSA public key. If
// any other key type is provided, an error is returned.
func UnmarshalPublicKey(pub []byte) (*ecdsa.PublicKey, error) {
	key, err := x509.ParsePKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}
	pk, ok := key.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.Errorf("unsupported public key type: %T", key)
	}
	return pk, nil
}

// Verify decodes the provided ASN.1 DER encoded signature and verifies the
// signature of the digest using the public key. An error is returned if the
// signature cannot be parsed or if the s component of the signature is not
// less than or equal to the half-order of the curve.
//
// Verify will return true, nil if the signature of the digest is valid.
func Verify(k *ecdsa.PublicKey, signature, digest []byte) (bool, error) {
	r, s, err := UnmarshalECDSASignature(signature)
	if err != nil {
		return false, err
	}
	lowS, err := IsLowS(k, s)
	if err != nil {
		return false, err
	}
	if !lowS {
		return false, errors.Errorf("s must be smaller than the half order of the curve")
	}
	return ecdsa.Verify(k, digest, r, s), nil
}

// Sign signs the digest using the private key, normalizes the s component of
// the signature to be less than or equal to the half-order of the curve, and
// ASN.1 DER encodes the signature.
func Sign(rand io.Reader, k *ecdsa.PrivateKey, digest []byte) ([]byte, error) {
	r, s, err := ecdsa.Sign(rand, k, digest)
	if err != nil {
		return nil, err
	}
	s, _, err = ToLowS(&k.PublicKey, s)
	if err != nil {
		return nil, err
	}
	return MarshalECDSASignature(r, s)
}

// A Signer is a wrapper around an ecdsa.PrivateKey that produces signatures
// normalized to a low-s value.
type Signer struct {
	*ecdsa.PrivateKey
}

// NewSigner creates a new Signer that generates normalized signatures.
func NewSigner(k *ecdsa.PrivateKey) *Signer {
	return &Signer{
		PrivateKey: k,
	}
}

// Sign implements the crypto.Signer interface.
func (s *Signer) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) ([]byte, error) {
	return Sign(rand, s.PrivateKey, digest)
}
