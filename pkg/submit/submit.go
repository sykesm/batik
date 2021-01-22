// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package submit

import (
	"context"
	"path/filepath"

	"github.com/pkg/errors"

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

// TODO: proper error values with semantics

type Service struct {
	repo      Repository          // repo is a reference to the transaction state repository.
	validator validator.Validator // validator the transaction Validator
}

func NewService(repo Repository) *Service {
	return &Service{
		repo:      repo,
		validator: validator.NewSignature(),
	}
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

	// err = validateWASM(resolved)
	resp, err := s.validator.Validate(&validationv1.ValidateRequest{
		ResolvedTransaction: transaction.FromResolved(resolved),
	})
	if err != nil { // TODO: recovery needs to halt
		panic("validation failed with unrecoverable error")
	}
	if !resp.Valid && resp.ErrorMessage != "" {
		return errors.Errorf("validation failed: %s", resp.ErrorMessage)
	}
	if !resp.Valid {
		return errors.New("validation failed")
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
