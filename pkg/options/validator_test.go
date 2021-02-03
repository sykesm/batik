// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"flag"
	"testing"

	. "github.com/onsi/gomega"
)

func TestValidatorDefaults(t *testing.T) {
	gt := NewGomegaWithT(t)
	ledger := ValidatorDefaults()
	gt.Expect(ledger).To(Equal(&Validator{
		CodeDir: "validators",
		Type:    "wasm",
	}))
}

func TestValidatorApplyDefaults(t *testing.T) {
	tests := map[string]struct {
		setup func(*Validator)
	}{
		"empty":    {setup: func(l *Validator) { *l = Validator{} }},
		"code dir": {setup: func(l *Validator) { l.CodeDir = "" }},
		"type":     {setup: func(l *Validator) { l.Type = "" }},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			input := ValidatorDefaults()
			tt.setup(input)

			input.ApplyDefaults()
			gt.Expect(input).To(Equal(ValidatorDefaults()))
		})
	}
}

func TestValidatorFlagNames(t *testing.T) {
	gt := NewGomegaWithT(t)
	flags := (&Validator{}).Flags()

	var names []string
	for _, f := range flags {
		names = append(names, f.Names()...)
	}

	gt.Expect(flags).To(HaveLen(1))
	gt.Expect(names).To(ConsistOf(
		"validators-dir",
	))
}

func TestValidatorFlags(t *testing.T) {
	tests := map[string]struct {
		args     []string
		expected Validator
	}{
		"no flags": {
			args:     []string{},
			expected: Validator{},
		},
		"code dir": {
			args:     []string{"--validators-dir=some/path/name"},
			expected: Validator{CodeDir: "some/path/name"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			ledger := &Validator{}
			flagSet := flag.NewFlagSet("validator-test", flag.ContinueOnError)
			for _, f := range ledger.Flags() {
				err := f.Apply(flagSet)
				gt.Expect(err).NotTo(HaveOccurred())
			}

			err := flagSet.Parse(tt.args)
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(ledger).To(Equal(&tt.expected))
		})
	}
}

func TestValidatorFlagsDefaultText(t *testing.T) {
	flags := ValidatorDefaults().Flags()
	assertWrappedFlagWithDefaultText(t, flags...)
}
