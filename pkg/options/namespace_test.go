// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"flag"
	"testing"

	. "github.com/onsi/gomega"
)

func TestNamespaceDefaults(t *testing.T) {
	gt := NewGomegaWithT(t)
	ledger := NamespaceDefaults()
	gt.Expect(ledger).To(Equal(&Namespace{
		DataDir:   "data",
		Validator: "signature-builtin",
	}))
}

func TestNamespaceApplyDefaults(t *testing.T) {
	tests := map[string]struct {
		setup func(*Namespace)
	}{
		"empty":     {setup: func(l *Namespace) { *l = Namespace{} }},
		"data dir":  {setup: func(l *Namespace) { l.DataDir = "" }},
		"validator": {setup: func(l *Namespace) { l.Validator = "" }},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			input := NamespaceDefaults()
			tt.setup(input)

			input.ApplyDefaults()
			gt.Expect(input).To(Equal(NamespaceDefaults()))
		})
	}
}

func TestNamespaceFlagNames(t *testing.T) {
	gt := NewGomegaWithT(t)
	flags := (&Namespace{}).Flags()

	var names []string
	for _, f := range flags {
		names = append(names, f.Names()...)
	}

	gt.Expect(flags).To(HaveLen(1))
	gt.Expect(names).To(ConsistOf(
		"data-dir",
	))
}

func TestNamespaceFlags(t *testing.T) {
	tests := map[string]struct {
		args     []string
		expected Namespace
	}{
		"no flags": {
			args:     []string{},
			expected: Namespace{},
		},
		"data dir": {
			args:     []string{"--data-dir=some/path/name"},
			expected: Namespace{DataDir: "some/path/name"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			ledger := &Namespace{}
			flagSet := flag.NewFlagSet("ledger-test", flag.ContinueOnError)
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

func TestNamespaceFlagsDefaultText(t *testing.T) {
	flags := NamespaceDefaults().Flags()
	assertWrappedFlagWithDefaultText(t, flags...)
}
