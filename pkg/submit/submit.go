// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package submit

import (
	"context"
	"crypto"

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
	ConsumeStates(...transaction.StateID) error
}

type Service struct {
	repo Repository // repo is a reference to the transaction state repository.
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// TODO: build submitter instance

func (s *Service) Submit(ctx context.Context, tx *transaction.Transaction) error {
	// Check if the transaction already exists
	_, err := s.repo.GetTransaction(tx.ID)
	if err == nil || !store.IsNotFound(err) {
		return &store.AlreadyExistsError{Err: errors.Errorf("transaction %s already exists", tx.ID)}
	}

	// resolve all inputs and references
	_, err = resolve(s.repo, tx)
	if err != nil {
		return errors.WithMessagef(err, "state resolution for transaction %s failed", tx.ID)
	}

	err = s.repo.PutTransaction(tx)
	if err != nil {
		return errors.WithMessagef(err, "storing transaction %s failed", tx.ID)
	}

	for i := range tx.Tx.Outputs {
		state := transaction.ToState(tx.Tx.Outputs[i], tx.ID, uint64(i))
		err = s.repo.PutState(state)
		if err != nil {
			return errors.WithMessagef(err, "storing transaction %s failed", tx.ID)
		}
	}

	var inputs []transaction.StateID
	for _, input := range tx.Tx.Inputs {
		inputs = append(inputs, transaction.StateID{TxID: input.Txid, OutputIndex: input.OutputIndex})
	}
	err = s.repo.ConsumeStates(inputs...)
	if err != nil {
		return err
	}

	return nil
}

func resolve(repo Repository, tx *transaction.Transaction) (*transaction.Resolved, error) {
	var inputs []*transaction.State
	for _, input := range tx.Tx.Inputs {
		stateID := transaction.StateID{TxID: input.Txid, OutputIndex: input.OutputIndex}
		state, err := repo.GetState(stateID)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, state)
	}
	var refs []*transaction.State
	for _, ref := range tx.Tx.References {
		stateID := transaction.StateID{TxID: ref.Txid, OutputIndex: ref.OutputIndex}
		state, err := repo.GetState(stateID)
		if err != nil {
			return nil, err
		}
		refs = append(refs, state)
	}
	var outputs []*transaction.State
	for i, out := range tx.Tx.Outputs {
		outputs = append(outputs, transaction.ToState(out, tx.ID, uint64(i)))
	}

	resolved := &transaction.Resolved{
		Tx:              tx,
		Inputs:          inputs,
		References:      refs,
		Outputs:         outputs,
		Parameters:      transaction.ToParameters(tx.Tx.Parameters...),
		RequiredSigners: transaction.ToParties(tx.Tx.RequiredSigners...),
	}

	return resolved, nil
}

func verifySignature(pk crypto.PublicKey, signature, hash []byte) bool {
	return false
}
