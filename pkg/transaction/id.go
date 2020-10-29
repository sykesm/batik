// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package transaction

import (
	"bytes"
	"encoding/hex"
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

// Equals returns true if this identity is equal to the argument.
func (id ID) Equals(that ID) bool {
	return bytes.Equal(id.Bytes(), that.Bytes())
}
