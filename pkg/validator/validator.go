// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package validator

import validationv1 "github.com/sykesm/batik/pkg/pb/validation/v1"

type Validator interface {
	Init() error
	Validate(*validationv1.ValidateRequest) (*validationv1.ValidateResponse, error)
}
