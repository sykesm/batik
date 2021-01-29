// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package ecdsautil

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/asn1"
	"math/big"

	"github.com/pkg/errors"
)

var (
	// curveHalfOrders contains the precomputed curve group orders halved.
	// It is used to ensure that signature' S value is lower or equal to the
	// curve group order halved. We accept only low-S signatures.
	// They are precomputed for efficiency reasons.
	curveHalfOrders = map[elliptic.Curve]*big.Int{
		elliptic.P224(): new(big.Int).Rsh(elliptic.P224().Params().N, 1),
		elliptic.P256(): new(big.Int).Rsh(elliptic.P256().Params().N, 1),
		elliptic.P384(): new(big.Int).Rsh(elliptic.P384().Params().N, 1),
		elliptic.P521(): new(big.Int).Rsh(elliptic.P521().Params().N, 1),
	}
)

type ecdsaSignature struct {
	R, S *big.Int
}

func MarshalECDSASignature(r, s *big.Int) ([]byte, error) {
	return asn1.Marshal(ecdsaSignature{r, s})
}

func UnmarshalECDSASignature(raw []byte) (*big.Int, *big.Int, error) {
	sig := &ecdsaSignature{}
	_, err := asn1.Unmarshal(raw, sig)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to unmarshal signature")
	}
	if sig.R == nil || sig.S == nil || sig.R.Sign() != 1 || sig.S.Sign() != 1 {
		return nil, nil, errors.Errorf("invalid signature: %x", raw)
	}
	return sig.R, sig.S, nil
}

// IsLow checks if s is a low-S, less than or equal to the half-order of the
// curve.
func IsLowS(k *ecdsa.PublicKey, s *big.Int) (bool, error) {
	halfOrder, ok := curveHalfOrders[k.Curve]
	if !ok {
		return false, errors.New("unrecognized ecdsa curve")
	}
	return s.Cmp(halfOrder) != 1, nil
}

// ToLowS normalizes the s value to be less than or equsl to the half-order of
// the curve.
func ToLowS(k *ecdsa.PublicKey, s *big.Int) (*big.Int, bool, error) {
	lowS, err := IsLowS(k, s)
	if err != nil {
		return nil, false, err
	}
	if !lowS {
		// Set s to N - s that will be then in the lower part of signature space
		// less or equal to half order
		return s.Sub(k.Params().N, s), true, nil
	}
	return s, false, nil
}
