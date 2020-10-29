// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package transaction

import (
	"strconv"
	"testing"

	. "github.com/onsi/gomega"
)

func TestStateIDString(t *testing.T) {
	tests := []struct {
		sid      StateID
		expected string
	}{
		{sid: StateID{}, expected: ":0000000000000000"},
		{sid: StateID{TxID: ID([]byte{1})}, expected: "01:0000000000000000"},
		{sid: StateID{TxID: ID([]byte{255}), OutputIndex: 255}, expected: "ff:00000000000000ff"},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			gt := NewGomegaWithT(t)
			gt.Expect(tt.sid.String()).To(Equal(tt.expected))
		})
	}
}

func TestStateIDEquals(t *testing.T) {
	tests := []struct {
		a, b  StateID
		equal bool
	}{
		{a: StateID{}, b: StateID{}, equal: true},
		{a: StateID{OutputIndex: 1}, b: StateID{}, equal: false},
		{a: StateID{}, b: StateID{OutputIndex: 1}, equal: false},
		{a: StateID{OutputIndex: 1}, b: StateID{OutputIndex: 1}, equal: true},
		{a: StateID{TxID: ID([]byte{}), OutputIndex: 1}, b: StateID{OutputIndex: 1}, equal: true},
		{a: StateID{OutputIndex: 1}, b: StateID{TxID: ID([]byte{}), OutputIndex: 1}, equal: true},
		{a: StateID{TxID: ID([]byte{1}), OutputIndex: 1}, b: StateID{TxID: ID([]byte{1}), OutputIndex: 1}, equal: true},
		{a: StateID{TxID: ID([]byte{0}), OutputIndex: 1}, b: StateID{TxID: ID([]byte{1}), OutputIndex: 1}, equal: false},
		{a: StateID{TxID: ID([]byte{1}), OutputIndex: 2}, b: StateID{TxID: ID([]byte{0}), OutputIndex: 1}, equal: false},
		{a: StateID{TxID: ID([]byte{1}), OutputIndex: 0}, b: StateID{TxID: ID([]byte{1}), OutputIndex: 1}, equal: false},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			gt := NewGomegaWithT(t)
			gt.Expect(tt.a.Equals(tt.b)).To(Equal(tt.equal))
		})
	}
}
