// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package transaction

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

// ID is a transaction identifier. A transaction ID is a merkle hash that
// uniquely identifies a transaction based on its contents.
type ID []byte

// NewID creates a new ID from a byte slice.
func NewID(b []byte) ID {
	var id []byte
	if b != nil {
		id = make([]byte, len(b), len(b))
		copy(id, b)
	}
	return ID(id)
}

// ID implements fmt.Stringer and returns the hex encoded representation of ID.
func (id ID) String() string {
	return hex.EncodeToString(id)
}

// Bytes implements an explicit conversion of the ID to a byte slice.
func (id ID) Bytes() []byte {
	return []byte(id)
}

// Equals returns true if this identifer is equal to the argument.
func (id ID) Equals(that ID) bool {
	return bytes.Equal(id.Bytes(), that.Bytes())
}

// MarshalJSON satisfies the json.Marshaler interface.
func (id ID) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%x"`, id.Bytes())), nil
}

// UnmarshalJSON satisfies the json.Unmarshaler interface.
func (id *ID) UnmarshalJSON(b []byte) error {
	// hex encoded string in quotes
	decoded, err := hex.DecodeString(string(b[1 : len(b)-1]))
	if err != nil {
		return err
	}
	*id = decoded
	return nil
}

// StateID is a state identifier. A state ID uniquely identfies an output of a
// transaction.
type StateID struct {
	TxID        ID     `json:"txid"`
	OutputIndex uint64 `json:"output_index"`
}

// StateID implements fmt.Stringer and returns hex encoded transaction
// identifier and the hex encoded output index separated by a colon (':').
func (sid StateID) String() string {
	return fmt.Sprintf("%s:%016x", sid.TxID, sid.OutputIndex)
}

// Equals returns true if this state identifier is equal to the argument.
func (sid StateID) Equals(that StateID) bool {
	if sid.OutputIndex == that.OutputIndex && sid.TxID.Equals(that.TxID) {
		return true
	}
	return false
}
