// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package grpcapi

import (
	"context"

	"github.com/pkg/errors"

	"github.com/sykesm/batik/pkg/namespace"
	"github.com/sykesm/batik/pkg/transaction"
)

type NamespaceMapAdapter map[string]*namespace.Namespace

func (nma NamespaceMapAdapter) Submitter(namespace string) Submitter {
	ns, ok := nma[namespace]
	if !ok {
		return notFoundSubmitter(namespace)
	}

	return ns
}

func (nma NamespaceMapAdapter) Repository(namespace string) Repository {
	ns, ok := nma[namespace]
	if !ok {
		return notFoundRepository(namespace)
	}

	return ns.Repo
}

var errNamespaceNotFound = errors.Errorf("namespace not found")

func isNamespaceNotFound(err error) bool {
	return errors.Is(err, errNamespaceNotFound)
}

type notFoundSubmitter string

func (nfs notFoundSubmitter) Submit(context.Context, *transaction.Signed) error {
	return errors.WithMessagef(errNamespaceNotFound, "bad namespace %q", nfs)
}

type notFoundRepository string

func (nfr notFoundRepository) PutTransaction(*transaction.Transaction) error {
	return errors.WithMessagef(errNamespaceNotFound, "bad namespace %q", nfr)
}

func (nfr notFoundRepository) GetTransaction(transaction.ID) (*transaction.Transaction, error) {
	return nil, errors.WithMessagef(errNamespaceNotFound, "bad namespace %q", nfr)
}

func (nfr notFoundRepository) PutState(*transaction.State) error {
	return errors.WithMessagef(errNamespaceNotFound, "bad namespace %q", nfr)
}

func (nfr notFoundRepository) GetState(transaction.StateID, bool) (*transaction.State, error) {
	return nil, errors.WithMessagef(errNamespaceNotFound, "bad namespace %q", nfr)
}
