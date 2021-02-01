// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package namespace

import (
	"go.uber.org/zap"

	"github.com/sykesm/batik/pkg/store"
	"github.com/sykesm/batik/pkg/submit"
)

// Namespace carries all of the resources required
// for the operation of a namespace.
type Namespace struct {
	Logger        *zap.Logger
	TxRepo        *store.TransactionRepository
	SubmitService *submit.Service
}

func New(logger *zap.Logger, kvStore store.KV) *Namespace {
	repo := store.NewRepository(kvStore)
	return &Namespace{
		Logger:        logger,
		TxRepo:        repo,
		SubmitService: submit.NewService(repo),
	}
}
