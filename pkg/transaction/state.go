// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package transaction

import (
	"fmt"
)

type StateID struct {
	TxID        ID
	OutputIndex uint64
}

func (sid StateID) String() string {
	return fmt.Sprintf("%s:%016x", sid.TxID, sid.OutputIndex)
}

func (sid StateID) Equals(that StateID) bool {
	if sid.OutputIndex == that.OutputIndex && sid.TxID.Equals(that.TxID) {
		return true
	}
	return false
}

type Party struct {
	PublicKey []byte
}

type StateInfo struct {
	Kind   string
	Owners []*Party
}

type State struct {
	ID        StateID
	StateInfo *StateInfo
	Data      []byte
}

type Parameter struct {
	Name  string
	Value []byte
}

type Signature struct {
	PublicKey []byte
	Signature []byte
}

type Resolved struct {
	ID              ID
	Inputs          []*State
	References      []*State
	Outputs         []*State
	Parameters      []*Parameter
	RequiredSigners []*Party
	Signatures      []*Signature
}
