// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package namespace

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/sykesm/batik/pkg/merkle"
	"github.com/sykesm/batik/pkg/store"
	"github.com/sykesm/batik/pkg/transaction"
)

// Namespace carries all of the resources required
// for the operation of a namespace.
type Namespace struct {
	Logger *zap.Logger
	Hasher merkle.Hasher

	LevelDB   *store.LevelDBKV
	Repo      Repository
	committer *committer
}

func New(
	logger *zap.Logger,
	hasher merkle.Hasher,
	level *store.LevelDBKV,
	validator Validator,
) *Namespace {
	repo := store.NewRepository(level)

	return &Namespace{
		Logger:  logger,
		Hasher:  hasher,
		LevelDB: level,
		Repo:    repo,
		committer: &committer{
			repo:      repo,
			validator: validator,
		},
	}
}

func (ns *Namespace) Submit(ctx context.Context, signed *transaction.Signed) error {
	// TODO, optimization, check if this transaction exists and if it's already been
	// committed.

	// TODO, mark in the store that this tx has been disseminated (by us).
	err := ns.Repo.PutTransaction(signed.Transaction)
	if err != nil {
		return errors.WithMessage(err, "failed to store transaction")
	}

	receipt := transaction.NewReceipt(ns.Hasher, signed.Transaction.ID, signed.Signatures)
	err = ns.Repo.PutReceipt(receipt)
	if err != nil {
		return errors.WithMessage(err, "failed to store transaction receipt")
	}

	// TODO, actually disseminate once we have some notion of
	// other peers in this namespace.

	// TODO, order the receipt ID

	return ns.committer.commit(receipt.ID)
}
