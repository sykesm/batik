// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestValidatorApplyDefaults(t *testing.T) {
	defaults := Validator{
		Name: "name",
		Type: "wasm",
		Path: "data/validators/name.wasm",
	}

	tests := map[string]struct {
		setup    func(*Validator)
		expected Validator
	}{
		"empty": {
			setup:    func(l *Validator) { *l = Validator{Name: "name"} },
			expected: defaults,
		},
		"path empty": {
			setup:    func(l *Validator) { l.Path = "" },
			expected: defaults,
		},
		"path specified": {
			setup: func(l *Validator) { l.Path = "some/path" },
			expected: Validator{
				Name: "name",
				Type: "wasm",
				Path: "some/path",
			},
		},
		"type": {
			setup:    func(l *Validator) { l.Type = "" },
			expected: defaults,
		},
		"not wasm": {
			setup: func(l *Validator) {
				l.Type = "builtin"
				l.Path = ""
			},
			expected: Validator{
				Name: "name",
				Type: "builtin",
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
