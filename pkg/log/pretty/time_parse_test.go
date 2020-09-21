// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pretty

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestTryParseTime(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected int64
	}{
		{
			name:     "string float",
			value:    "1599593917.548589",
			expected: 1599593917 * 1e9,
		},
		{
			name:     "time.RFC3339",
			value:    "2002-10-02T10:00:00-05:00",
			expected: 10335708 * 1e11,
		},
		{
			name:     "time.RFC3339Nano",
			value:    "2002-10-02T15:00:00.05Z",
			expected: 103357080005 * 1e7,
		},
		{
			name:     "time.RFC822",
			value:    "Wed, 02 Oct 2002 15:00:00 MST",
			expected: 10335708 * 1e11,
		},
		{
			name:     "time.RFC822Z",
			value:    "Wed, 02 Oct 2002 15:00:00 +0200",
			expected: 10335636 * 1e11,
		},
		{
			name:     "time.RFC850",
			value:    "Monday, 02-Jan-06 15:04:05 MST",
			expected: 1136214245 * 1e9,
		},
		{
			name:     "time.RFC1123",
			value:    "Mon, 02 Jan 2006 15:04:05 MST",
			expected: 1136214245 * 1e9,
		},
		{
			name:     "time.RFC1123Z",
			value:    "Mon, 02 Jan 2006 15:04:05 -0700",
			expected: 1136239445 * 1e9,
		},
		{
			name:     "time.UnixDate",
			value:    "Mon Jan 02 15:04:05 MST 2006",
			expected: 1136214245 * 1e9,
		},
		{
			name:     "time.RubyDate",
			value:    "Mon Jan 02 15:04:05 -0700 2006",
			expected: 1136239445 * 1e9,
		},
		{
			name:     "time.ANSIC",
			value:    "Mon Jan 02 15:04:05 2006",
			expected: 1136214245 * 1e9,
		},
		{
			name:     "time.Kitchen",
			value:    "3:04PM",
			expected: -6826932738871345152,
		},
		{
			name:     "time.Stamp",
			value:    "Jan 02 15:04:05",
			expected: -6826846333871345152,
		},
		{
			name:     "time.StampMilli",
			value:    "Jan 02 15:04:05.000",
			expected: -6826846333871345152,
		},
		{
			name:     "time.StampMicro",
			value:    "Jan 02 15:04:05.000000",
			expected: -6826846333871345152,
		},
		{
			name:     "time.StampNano",
			value:    "Jan 02 15:04:05.000000000",
			expected: -6826846333871345152,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			tm, isTime := tryParseTime(tt.value)
			gt.Expect(tm.UnixNano()).To(Equal(tt.expected))
			gt.Expect(isTime).To(BeTrue())
		})
	}

	t.Run("invalid format", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		_, isTime := tryParseTime("invalidtime")
		gt.Expect(isTime).To(BeFalse())
	})
}
