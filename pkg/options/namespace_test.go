// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestNamespaceApplyDefaults(t *testing.T) {
	defaults := Namespace{
		Name:      "name",
		DataDir:   "data/namespaces/name",
		Validator: "signature-builtin",
	}

	tests := map[string]struct {
		setup    func(*Namespace)
		expected Namespace
	}{
		"almost empty": {
			setup:    func(l *Namespace) { *l = Namespace{Name: "name"} },
			expected: defaults,
		},
		"data dir": {
			setup:    func(l *Namespace) { l.DataDir = "" },
			expected: defaults,
		},
		"validator": {
			setup:    func(l *Namespace) { l.Validator = "" },
			expected: defaults,
		},
		"overridden data dir": {
			setup: func(l *Namespace) { l.DataDir = "some/path" },
			expected: Namespace{
				Name:      "name",
				DataDir:   "some/path",
				Validator: "signature-builtin",
			},
		},
		"overridden validator": {
			setup: func(l *Namespace) { l.Validator = "custom" },
			expected: Namespace{
				Name:      "name",
				DataDir:   "data/namespaces/name",
				Validator: "custom",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			input := defaults
			tt.setup(&input)

			input.ApplyDefaults("data")
			gt.Expect(input).To(Equal(tt.expected))
		})
	}
}
