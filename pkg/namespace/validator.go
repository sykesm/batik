// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package namespace

import validationv1 "github.com/sykesm/batik/pkg/pb/validation/v1"

// A Validator is responsible for validating resolved transactions.
//
// Implementations must be iplemented as pure functions and any
// parameterization should be statically defined within the validator or be
// derived from the contents of the resolved transactions.
type Validator interface {
	// Validate is used to ensure that a transaction is semantically valid prior
	// to committing it to the ledger. This operation is only invoked if a
	// transaction has not already been committed and transaction inputs and
	// references have not been consumed.
	//
	// If the Valid field of the ValidateResponse is true, the transaction is
	// considered valid and will enter commit processing. Valid transactions may
	// fail to commit if the invariants established prior to validation have
	// changed or any invariants raised by the validator are not satisfied.
	//
	// If an error is returned, transaction processing is halted.
	Validate(*validationv1.ValidateRequest) (*validationv1.ValidateResponse, error)
}
