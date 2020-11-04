// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package transaction

import (
	"strconv"
	"testing"

	. "github.com/onsi/gomega"
)

func TestToHex(t *testing.T) {
	tests := []struct {
		size     int
		expected string
	}{
		{size: 0, expected: ""},
		{size: 1, expected: "00"},
		{size: 2, expected: "0001"},
		{size: 31, expected: "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e"},
		{size: 32, expected: "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"},
		{size: 33, expected: "000102030405060708090a0b0c0d0e....12131415161718191a1b1c1d1e1f20"},
		{size: 127, expected: "000102030405060708090a0b0c0d0e....707172737475767778797a7b7c7d7e"},
	}
	for _, tt := range tests {
		t.Run(strconv.Itoa(tt.size), func(t *testing.T) {
			gt := NewGomegaWithT(t)

			b := make([]byte, tt.size, tt.size)
			for i := range b {
				b[i] = byte(i)
			}
			gt.Expect(toHexString(b)).To(Equal(tt.expected))
		})
	}
}
