// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package namespace

import (
	"github.com/pkg/errors"

	validationv1 "github.com/sykesm/batik/pkg/pb/validation/v1"
	"github.com/sykesm/batik/pkg/store"
	"github.com/sykesm/batik/pkg/transaction"
)

// A Repository abstracts the data persistence layer for transactions and
// states.
type Repository interface {
	PutCommitted(transaction.ID, *transaction.Committed) error
	GetCommitted(transaction.ID) (*transaction.Committed, error)
	PutReceipt(*transaction.Receipt) error
	GetReceipt([]byte) (*transaction.Receipt, error)
	PutTransaction(*transaction.Transaction) error
	GetTransaction(transaction.ID) (*transaction.Transaction, error)
	PutState(*transaction.State) error
	GetState(transaction.StateID, bool) (*transaction.State, error)
	ConsumeState(transaction.StateID) error
}

// TODO: proper error values with semantics

type committer struct {
	repo      Repository // repo is a reference to the transaction state repository.
	validator Validator  // validator the transaction Validator
}

func newCommitter(repo Repository, validator Validator) *committer {
	return &committer{
		repo:      repo,
		validator: validator,
	}
}

func (c *committer) commit(receiptID []byte) error {
	receipt, err := c.repo.GetReceipt(receiptID)
	if store.IsNotFound(err) {
		return newHaltError(err, "receipt should have been disseminated but was not found")
	}
	if err != nil {
		return newHaltError(err, "transaction store failure")
	}

	tx, err := c.repo.GetTransaction(receipt.TxID)
	if store.IsNotFound(err) {
		return newHaltError(err, "transaction should have been disseminated but was not found")
	}
	if err != nil {
		return newHaltError(err, "transaction store failure")
	}

	// resolve all inputs and references
	resolved, err := resolve(c.repo, tx, receipt.Signatures)
	if err != nil && store.IsNotFound(err) {
		return errors.WithMessagef(err, "missing state for transaction %s", tx.ID)
	}
	if err != nil {
		return newHaltError(err, "state resolution for transaction %s failed", tx.ID)
	}

	resp, err := c.validator.Validate(&validationv1.ValidateRequest{
		ResolvedTransaction: transaction.FromResolved(resolved),
	})
	if err != nil {
		return newHaltError(err, "validator failed")
	}
	if !resp.Valid && resp.ErrorMessage != "" {
		return errors.Errorf("validation failed: %s", resp.ErrorMessage)
	}
	if !resp.Valid {
		return errors.New("validation failed")
	}

	err = c.repo.PutCommitted(tx.ID, &transaction.Committed{
		// TODO, once ordered, include SeqNo
		ReceiptID: receipt.ID,
	})
	if err != nil {
		return newHaltError(err, "marking %s as committed failed", tx.ID)
	}

	for _, output := range resolved.Outputs {
		err = c.repo.PutState(output)
		if err != nil {
			return newHaltError(err, "storing transaction output %s failed", output.ID)
		}
	}

	for _, input := range resolved.Inputs {
		err = c.repo.ConsumeState(input.ID)
		if err != nil {
			return newHaltError(err, "consuming transaction state %s failed", input.ID)
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
