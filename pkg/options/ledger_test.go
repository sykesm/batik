// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"flag"
	"testing"

	. "github.com/onsi/gomega"
)

func TestLedgerDefaults(t *testing.T) {
	gt := NewGomegaWithT(t)
	ledger := LedgerDefaults()
	gt.Expect(ledger).To(Equal(&Ledger{
		DataDir: "data",
	}))
}

func TestLedgerApplyDefaults(t *testing.T) {
	tests := map[string]struct {
		setup func(*Ledger)
	}{
		"empty":    {setup: func(l *Ledger) { *l = Ledger{} }},
		"data dir": {setup: func(l *Ledger) { l.DataDir = "" }},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			input := LedgerDefaults()
			tt.setup(input)

			input.ApplyDefaults()
			gt.Expect(input).To(Equal(LedgerDefaults()))
		})
	}
}

func TestLedgerFlagNames(t *testing.T) {
	gt := NewGomegaWithT(t)
	flags := (&Ledger{}).Flags()

	var names []string
	for _, f := range flags {
		names = append(names, f.Names()...)
	}

	gt.Expect(flags).To(HaveLen(1))
	gt.Expect(names).To(ConsistOf(
		"data-dir",
	))
}

func TestLedgerFlags(t *testing.T) {
	tests := map[string]struct {
		args     []string
		expected Ledger
	}{
		"no flags": {
			args:     []string{},
			expected: Ledger{},
		},
		"data dir": {
			args:     []string{"--data-dir=some/path/name"},
			expected: Ledger{DataDir: "some/path/name"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			ledger := &Ledger{}
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
