// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestTotalOrderApplyDefaults(t *testing.T) {
	defaults := TotalOrder{
		Name: "name", // TODO / note, if name is unspecified, what behavior do we want?
		Type: "in-process",
	}

	tests := map[string]struct {
		setup func(*TotalOrder)
	}{
		"empty": {
			setup: func(l *TotalOrder) { *l = TotalOrder{Name: "name"} },
		},
		"type": {
			setup: func(l *TotalOrder) { l.Type = "" },
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			input := defaults
			tt.setup(&input)

			input.ApplyDefaults()
			gt.Expect(input).To(Equal(defaults))
		})
	}
}
