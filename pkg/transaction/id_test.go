// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package transaction

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
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

func TestIDMarshalingJSON(t *testing.T) {
	tests := []struct {
		id ID
	}{
		{id: NewID([]byte{})},
		{id: NewID([]byte{1})},
		{id: NewID([]byte{255})},
		{id: NewID([]byte{1, 2, 3, 4, 5})},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			gt := NewGomegaWithT(t)
			out, err := json.Marshal(tt.id)
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(out).To(Equal([]byte(fmt.Sprintf(`"%s"`, tt.id))))

			var unmarshaled ID
			err = json.Unmarshal([]byte(fmt.Sprintf(`"%s"`, tt.id)), &unmarshaled)
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(unmarshaled).To(Equal(tt.id))
		})
	}
}

func TestIDUnmarshalErrors(t *testing.T) {
	tests := []struct {
		in         string
		errMatcher types.GomegaMatcher
	}{
		{in: `""`, errMatcher: BeNil()},
		{in: `"0"`, errMatcher: MatchError(ContainSubstring("odd length hex string"))},
		{in: `"00`, errMatcher: MatchError(ContainSubstring("unexpected end of JSON input"))},
		{in: `00"`, errMatcher: HaveOccurred()},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			gt := NewGomegaWithT(t)

			var id ID
			err := json.Unmarshal([]byte(tt.in), &id)
			gt.Expect(err).To(tt.errMatcher)
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

func TestStateIDMarshalingJSON(t *testing.T) {
	tests := []struct {
		id       StateID
		expected string
	}{
		{id: StateID{}, expected: `{"txid":"", "output_index": 0}`},
		{id: StateID{OutputIndex: 33}, expected: `{"txid":"", "output_index": 33}`},
		{id: StateID{TxID: ID([]byte{1, 2, 3, 4})}, expected: `{"txid":"01020304", "output_index": 0}`},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			gt := NewGomegaWithT(t)
			enc, err := json.Marshal(tt.id)
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(enc).To(MatchJSON(tt.expected))
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
