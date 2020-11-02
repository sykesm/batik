// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package submit

import (
	"context"

	"github.com/pkg/errors"

	"github.com/sykesm/batik/pkg/store"
	"github.com/sykesm/batik/pkg/transaction"
)

// A Repository abstracts the data persistence layer for transactions and
// states.
type Repository interface {
	PutTransaction(*transaction.Transaction) error
	GetTransaction(transaction.ID) (*transaction.Transaction, error)
	PutState(*transaction.State) error
	GetState(transaction.StateID) (*transaction.State, error)
	ConsumeState(transaction.StateID) error
}

type Service struct {
	repo Repository // repo is a reference to the transaction state repository.
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// TODO: build submitter instance as a unit of work

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

	err = verify(resolved)
	if err != nil {
		panic(err)
	}

	err = s.repo.PutTransaction(resolved.Tx)
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
		state, err := repo.GetState(*input)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, state)
	}
	var refs []*transaction.State
	for _, ref := range tx.References {
		state, err := repo.GetState(*ref)
		if err != nil {
			return nil, err
		}
		refs = append(refs, state)
	}

	resolved := &transaction.Resolved{
		Tx:              tx,
		Inputs:          inputs,
		References:      refs,
		Outputs:         tx.Outputs,
		Parameters:      tx.Parameters,
		RequiredSigners: tx.RequiredSigners,
		Signatures:      sigs,
	}

	return resolved, nil
}

func verify(resolved *transaction.Resolved) error {
	return nil
}
