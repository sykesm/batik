// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package transaction

// A State represents the output of a transaction.
type State struct {
	ID        StateID    `json:"id"`
	StateInfo *StateInfo `json:"info"`
	Data      []byte     `json:"data"`
}

// A StateInfo holds metadata about a State.
type StateInfo struct {
	Kind   string   `json:"kind"`
	Owners []*Party `json:"owners"`
}

// A Party represents a state owner or transaction signatory.
type Party struct {
	PublicKey []byte `json:"public_key,omitempty"`
}

// String satisfies the fmt.Stringer interface.
func (p Party) String() string {
	return toHexString(p.PublicKey)
}

// A Parameter holds information intended to parameterize the execution of a
// transaction.
type Parameter struct {
	Name  string `json:"name"`
	Value []byte `json:"value"`
}

// A Signature is used to hold a transaction endorsement by the party
// represented by the public key.
type Signature struct {
	PublicKey []byte `json:"public_key"`
	Signature []byte `json:"signature"`
}
