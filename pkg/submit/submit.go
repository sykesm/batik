// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package submit

import (
	"bytes"
	"context"
	"crypto/sha256"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/sykesm/batik/pkg/ecdsautil"
	validationv1 "github.com/sykesm/batik/pkg/pb/validation/v1"
	"github.com/sykesm/batik/pkg/store"
	"github.com/sykesm/batik/pkg/transaction"
	"github.com/sykesm/batik/pkg/validator"
	"github.com/sykesm/batik/pkg/validator/utxo"
)

// A Repository abstracts the data persistence layer for transactions and
// states.
type Repository interface {
	PutTransaction(*transaction.Transaction) error
	GetTransaction(transaction.ID) (*transaction.Transaction, error)
	PutState(*transaction.State) error
	GetState(transaction.StateID, bool) (*transaction.State, error)
	ConsumeState(transaction.StateID) error
}

// TODO: build submitter instance as a unit of work
// TODO: proper error values with semantics

type Service struct {
	repo Repository // repo is a reference to the transaction state repository.
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Submit(ctx context.Context, signed *transaction.Signed) error {
	txid := signed.Transaction.ID

	// Transaction must have been processed before
	_, err := s.repo.GetTransaction(txid)
	if err == nil {
		return &store.AlreadyExistsError{Err: errors.Errorf("transaction %s already exists", txid)}
	}
	if !store.IsNotFound(err) {
		return err
	}

	// resolve all inputs and references
	resolved, err := resolve(s.repo, signed.Transaction, signed.Signatures)
	if err != nil {
		return errors.WithMessagef(err, "state resolution for transaction %s failed", txid)
	}

	// TODO: Validator associated with namespace
	// err = validateWASM(resolved)
	err = validate(resolved)
	if err != nil {
		return errors.Wrap(err, "validation failed")
	}

	err = s.repo.PutTransaction(signed.Transaction)
	if err != nil {
		return errors.WithMessagef(err, "storing transaction %s failed", txid)
	}

	for _, output := range resolved.Outputs {
		err = s.repo.PutState(output)
		if err != nil {
			return errors.WithMessagef(err, "storing transaction output %s failed", output.ID)
		}
	}

	for _, input := range resolved.Inputs {
		err = s.repo.ConsumeState(input.ID)
		if err != nil {
			return errors.WithMessagef(err, "consuming transaction state %s failed", input.ID)
		}
	}

	return nil
}

func resolve(repo Repository, tx *transaction.Transaction, sigs []*transaction.Signature) (*transaction.Resolved, error) {
	var inputs []*transaction.State
	for _, input := range tx.Inputs {
		state, err := repo.GetState(*input, false)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, state)
	}
	var refs []*transaction.State
	for _, ref := range tx.References {
		state, err := repo.GetState(*ref, false)
		if err != nil {
			return nil, err
		}
		refs = append(refs, state)
	}

	resolved := &transaction.Resolved{
		ID:              tx.ID,
		Inputs:          inputs,
		References:      refs,
		Outputs:         tx.Outputs,
		Parameters:      tx.Parameters,
		RequiredSigners: tx.RequiredSigners,
		Signatures:      sigs,
	}

	return resolved, nil
}

func digest(preImage []byte) [32]byte {
	return sha256.Sum256(preImage)
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
			return errors.Wrap(err, "unmarshal public key failed")
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

func validateWASM(resolved *transaction.Resolved) error {
	var (
		validator validator.Validator
		err       error
	)

	// Check and choose type of validator
	modulePath := filepath.Join("..", "..", "wasm", "modules", "utxotx", "target", "wasm32-unknown-unknown", "debug", "utxotx.wasm")
	validator, err = utxo.NewValidator(modulePath)
	if err != nil {
		return errors.WithMessage(err, "failed creating validator")
	}

	validateRequest := &validationv1.ValidateRequest{
		ResolvedTransaction: transaction.FromResolved(resolved),
	}

	validateResponse, err := validator.Validate(validateRequest)
	if err != nil {
		return err
	}
	if !validateResponse.Valid {
		if msg := validateResponse.ErrorMessage; msg != "" {
			return errors.Errorf("validation failed for transaction %s: %s", resolved.ID, msg)
		}
		return errors.Errorf("validation failed for transaction %s", resolved.ID)
	}

	return nil
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
