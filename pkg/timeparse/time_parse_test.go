// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package timeparse

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestParseUnixTime(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected int64
	}{
		{
			name:     "string",
			value:    "1599593917.548589",
			expected: 1599593917 * 1e9,
		},
		{
			name:     "float64",
			value:    1599593917.548589,
			expected: 1599593917 * 1e9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			tm, err := ParseUnixTime(tt.value)
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(tm.UnixNano()).To(Equal(tt.expected))
		})
	}

	t.Run("invalid format", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		_, err := ParseUnixTime("invalidtime")
		gt.Expect(err.Error()).To(Equal("strconv.ParseFloat: parsing \"invalidtime\": invalid syntax"))

		_, err = ParseUnixTime(int64(159953917))
		gt.Expect(err).To(MatchError("unix time is not string or float64"))
	})
}
