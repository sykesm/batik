// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"flag"
	"testing"

	. "github.com/onsi/gomega"
)

func TestLoggingDefaults(t *testing.T) {
	gt := NewGomegaWithT(t)
	ledger := LoggingDefaults()
	gt.Expect(ledger).To(Equal(&Logging{
		LogSpec: "info",
		Color:   "auto",
	}))
}

func TestLoggingApplyDefaults(t *testing.T) {
	tests := map[string]struct {
		setup func(*Logging)
	}{
		"empty":    {setup: func(l *Logging) { *l = Logging{} }},
		"log spec": {setup: func(l *Logging) { l.LogSpec = "" }},
		"color":    {setup: func(l *Logging) { l.Color = "" }},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			input := LoggingDefaults()
			tt.setup(input)

			input.ApplyDefaults()
			gt.Expect(input).To(Equal(LoggingDefaults()))
		})
	}
}

func TestLoggingFlagNames(t *testing.T) {
	gt := NewGomegaWithT(t)
	flags := (&Logging{}).Flags()

	var names []string
	for _, f := range flags {
		names = append(names, f.Names()...)
	}

	gt.Expect(flags).To(HaveLen(2))
	gt.Expect(names).To(ConsistOf(
		"log-spec",
		"color",
	))
}

func TestLoggingFlags(t *testing.T) {
	tests := map[string]struct {
		args     []string
		expected Logging
	}{
		"no flags": {
			args:     []string{},
			expected: Logging{},
		},
		"log spec": {
			args:     []string{"--log-spec=debug"},
			expected: Logging{LogSpec: "debug"},
		},
		"color": {
			args:     []string{"--color=yes"},
			expected: Logging{Color: "yes"},
		},
		"both": {
			args:     []string{"--log-spec=debug", "--color=yes"},
			expected: Logging{LogSpec: "debug", Color: "yes"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			logging := &Logging{}
			flagSet := flag.NewFlagSet("logging-test", flag.ContinueOnError)
			for _, f := range logging.Flags() {
				err := f.Apply(flagSet)
				gt.Expect(err).NotTo(HaveOccurred())
			}

			err := flagSet.Parse(tt.args)
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(logging).To(Equal(&tt.expected))
		})
	}
}

func TestLoggingFlagsDefaultText(t *testing.T) {
	flags := LoggingDefaults().Flags()
	assertWrappedFlagWithDefaultText(t, flags...)
}
