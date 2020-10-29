// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package transaction

import (
	"strconv"
	"testing"

	. "github.com/onsi/gomega"
)

func TestNewID(t *testing.T) {
	tests := []struct {
		in       []byte
		expected ID
	}{
		{in: nil, expected: ID(nil)},
		{in: []byte{}, expected: ID([]byte{})},
		{in: []byte{1, 2, 3, 4, 5}, expected: ID([]byte{1, 2, 3, 4, 5})},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			gt := NewGomegaWithT(t)
			id := NewID(tt.in)
			gt.Expect(id).To(Equal(tt.expected))

			if len(tt.in) > 0 {
				id[0] = 0
				gt.Expect([]byte(id)).NotTo(Equal(tt.in), "ID must copy the source bytes")
			}
		})
	}
}

func TestIDString(t *testing.T) {
	tests := []struct {
		in       ID
		expected string
	}{
		{in: NewID(nil), expected: ""},
		{in: NewID([]byte{}), expected: ""},
		{in: NewID([]byte{1}), expected: "01"},
		{in: NewID([]byte{255}), expected: "ff"},
		{in: NewID([]byte{1, 2, 3, 4, 5}), expected: "0102030405"},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			gt := NewGomegaWithT(t)
			gt.Expect(tt.in.String()).To(Equal(tt.expected))
		})
	}
}

func TestIDBytes(t *testing.T) {
	tests := []struct {
		in       ID
		expected []byte
	}{
		{in: NewID(nil), expected: nil},
		{in: NewID([]byte{}), expected: []byte{}},
		{in: NewID([]byte{1}), expected: []byte{1}},
		{in: NewID([]byte{255}), expected: []byte{255}},
		{in: NewID([]byte{1, 2, 3, 4, 5}), expected: []byte{1, 2, 3, 4, 5}},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			gt := NewGomegaWithT(t)
			gt.Expect(tt.in.Bytes()).To(Equal(tt.expected))
		})
	}
}

func TestIDEquals(t *testing.T) {
	tests := []struct {
		a, b  ID
		equal bool
	}{
		{a: nil, b: nil, equal: true},
		{a: nil, b: ID(nil), equal: true},
		{a: nil, b: NewID(nil), equal: true},
		{a: NewID(nil), b: nil, equal: true},
		{a: NewID(nil), b: ID(nil), equal: true},
		{a: NewID(nil), b: NewID(nil), equal: true},
		{a: NewID([]byte{1}), b: nil, equal: false},
		{a: NewID([]byte{1}), b: ID(nil), equal: false},
		{a: NewID([]byte{1}), b: NewID(nil), equal: false},
		{a: NewID([]byte{1}), b: ID([]byte{1}), equal: true},
		{a: NewID([]byte{1}), b: NewID([]byte{1}), equal: true},
		{a: NewID([]byte{1}), b: NewID([]byte{2}), equal: false},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			gt := NewGomegaWithT(t)
			gt.Expect(tt.a.Equals(tt.b)).To(Equal(tt.equal))
		})
	}
}
