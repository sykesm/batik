// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package validator

import (
	"bytes"
	"crypto/sha256"

	"github.com/pkg/errors"

	"github.com/sykesm/batik/pkg/ecdsautil"
	validationv1 "github.com/sykesm/batik/pkg/pb/validation/v1"
	"github.com/sykesm/batik/pkg/transaction"
)

type Signature struct{}

func NewSignature() *Signature {
	return &Signature{}
}

func (s *Signature) Validate(req *validationv1.ValidateRequest) (*validationv1.ValidateResponse, error) {
	err := validate(transaction.ToResolved(req.ResolvedTransaction))
	if err != nil {
		return &validationv1.ValidateResponse{Valid: false, ErrorMessage: err.Error()}, nil
	}
	return &validationv1.ValidateResponse{Valid: true}, nil
}

func validate(resolved *transaction.Resolved) error {
	requiredSigners := requiredSigners(resolved)
	for _, signer := range requiredSigners {
		if signer.PublicKey == nil {
			return errors.New("required signer missing public key")
		}
		sig := signature(signer.PublicKey, resolved.Signatures)
		if sig == nil {
			return errors.Errorf("missing signature from %x", signer.PublicKey)
		}
		pk, err := ecdsautil.UnmarshalPublicKey(signer.PublicKey)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal public key")
		}
		txidHash := digest(resolved.ID.Bytes())
		ok, err := ecdsautil.Verify(pk, sig.Signature, txidHash[:])
		if err != nil {
			return err
		}
		if !ok {
			return errors.New("signature verification failed")
		}
	}
	return nil
}

func digest(preImage []byte) [32]byte {
	return sha256.Sum256(preImage)
}

func signature(publicKey []byte, signatures []*transaction.Signature) *transaction.Signature {
	for _, sig := range signatures {
		if bytes.Equal(sig.PublicKey, publicKey) {
			return sig
		}
	}
	return nil
}

// TODO(mjs): Consider duplicate removal
func requiredSigners(resolved *transaction.Resolved) []*transaction.Party {
	var required []*transaction.Party
	for _, input := range resolved.Inputs {
		if input.StateInfo != nil {
			required = append(required, input.StateInfo.Owners...)
		}
	}
	required = append(required, resolved.RequiredSigners...)
	return required
}
