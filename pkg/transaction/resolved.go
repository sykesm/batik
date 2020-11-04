// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package transaction

// Resolved is a representation of a transaction where all input state
// references have been resolved. This is the information that a validator uses
// to validate a transaction.
type Resolved struct {
	ID              ID           `json:"id"`
	Inputs          []*State     `json:"inputs,omitempty"`
	References      []*State     `json:"references,omitempty"`
	Outputs         []*State     `json:"outputs,omitempty"`
	Parameters      []*Parameter `json:"parameters,omitempty"`
	RequiredSigners []*Party     `json:"required_signers,omitempty"`
	Signatures      []*Signature `json:"signatures,omitempty"`
}
