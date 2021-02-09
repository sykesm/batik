// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestTotalOrderApplyDefaults(t *testing.T) {
	defaults := TotalOrder{
		Name:    "name", // TODO / note, if name is unspecified, what behavior do we want?
		Type:    "in-process",
		DataDir: "base/dir/totalorders/name",
	}

	tests := map[string]struct {
		setup    func(*TotalOrder)
		expected TotalOrder
	}{
		"empty": {
			setup:    func(l *TotalOrder) { *l = TotalOrder{Name: "name"} },
			expected: defaults,
		},
		"type": {
			setup:    func(l *TotalOrder) { l.Type = "" },
			expected: defaults,
		},
		"data dir": {
			setup:    func(l *TotalOrder) { l.DataDir = "" },
			expected: defaults,
		},
		"overridden data dir": {
			setup: func(l *TotalOrder) { l.DataDir = "some/path" },
			expected: TotalOrder{
				Name:    "name",
				Type:    "in-process",
				DataDir: "some/path",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			input := defaults
			tt.setup(&input)

			input.ApplyDefaults("base/dir")
			gt.Expect(input).To(Equal(tt.expected))
		})
	}
}
