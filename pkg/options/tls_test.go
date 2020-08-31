// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"flag"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v2"
)

func TestTLSOptionsFlags(t *testing.T) {
	tests := map[string]struct {
		args     []string
		expected TLSOptions
	}{
		"no flags": {
			[]string{},
			TLSOptions{},
		},
		"cert flag": {
			[]string{"--tls-cert-file", "certificate.file"},
			TLSOptions{ServerCert: CertKeyPair{CertFile: "certificate.file"}},
		},
		"key flag": {
			[]string{"--tls-private-key-file", "private.key"},
			TLSOptions{ServerCert: CertKeyPair{KeyFile: "private.key"}},
		},
		"cert and key flags": {
			[]string{
				"--tls-cert-file", "certificate.file",
				"--tls-private-key-file", "private.key",
			},
			TLSOptions{ServerCert: CertKeyPair{
				CertFile: "certificate.file",
				KeyFile:  "private.key",
			}},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			opts := &TLSOptions{}
			flagSet := flag.NewFlagSet("tls-options-test", flag.ContinueOnError)
			for _, f := range opts.Flags() {
				err := f.Apply(flagSet)
				gt.Expect(err).NotTo(HaveOccurred())
			}

			err := flagSet.Parse(tt.args)
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(opts).To(Equal(&tt.expected))
		})
	}
}

func TestTLSOptionsUsage(t *testing.T) {
	gt := NewGomegaWithT(t)

	opts := &TLSOptions{}
	for _, f := range opts.Flags() {
		f := f.(cli.DocGenerationFlag)
		gt.Expect(f.GetUsage()).NotTo(ContainSubstring("\n"))
		gt.Expect(f.GetUsage()).NotTo(MatchRegexp(`\s{2}`))
	}
}
