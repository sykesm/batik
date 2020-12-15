// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package validator

import validationv1 "github.com/sykesm/batik/pkg/pb/validation/v1"

// A Validator is a pluggable instance responsible for validating resolved
// transactions of varying types. They take advantage of a go wasmtime host
// to handle validation via a web assembly module.
type Validator interface {
	Validate(*validationv1.ValidateRequest) (*validationv1.ValidateResponse, error)
}
