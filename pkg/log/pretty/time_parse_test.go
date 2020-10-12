// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pretty

import (
	"strconv"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

func TestParseUnixTime(t *testing.T) {
	tests := []struct {
		input      string
		expected   time.Time
		errMatcher types.GomegaMatcher
	}{
		{input: "", expected: time.Unix(0, 0)},
		{input: "1599593917", expected: time.Unix(1_599_593_917, 0)},
		{input: "0.159959391", expected: time.Unix(0, int64(159_959_391))},
		{input: "invalidtime", errMatcher: MatchError(ContainSubstring(`parsing "invalidtime": invalid syntax`))},
		{input: "1.2.3", errMatcher: MatchError(`ParseUnixTime: invalid syntax: "1.2.3"`)},
		{input: "0.fred", errMatcher: MatchError(ContainSubstring(`parsing "0.fred": invalid syntax`))},
		{input: "bob.0", errMatcher: MatchError(ContainSubstring(`parsing "bob": invalid syntax`))},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			gt := NewGomegaWithT(t)

			actual, err := ParseUnixTime(tt.input)
			if tt.errMatcher != nil {
				gt.Expect(err).To(tt.errMatcher)
				return
			}
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(actual).To(Equal(tt.expected))
		})
	}
}
