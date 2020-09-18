// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"flag"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v2"
)

func TestTLSServerDefaults(t *testing.T) {
	gt := NewGomegaWithT(t)
	ts := TLSServerDefaults()
	gt.Expect(ts).To(Equal(&TLSServer{}))
}

func TestTLSServerApplyDefaults(t *testing.T) {
	tests := map[string]struct {
		setup func(ts *TLSServer)
	}{
		"empty": {setup: func(ts *TLSServer) { *ts = TLSServer{} }},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			input := TLSServerDefaults()
			tt.setup(input)

			input.ApplyDefaults()
			gt.Expect(input).To(Equal(TLSServerDefaults()))
		})
	}
}

func TestTLSServerFlagNames(t *testing.T) {
	gt := NewGomegaWithT(t)
	flags := (&TLSServer{}).Flags()

	var names []string
	for _, f := range flags {
		names = append(names, f.Names()...)
	}

	gt.Expect(flags).To(HaveLen(2))
	gt.Expect(names).To(ConsistOf(
		"tls-cert-file",
		"tls-private-key-file",
	))
}

func TestTLSServerFlags(t *testing.T) {
	tests := map[string]struct {
		args     []string
		expected TLSServer
	}{
		"no flags": {
			[]string{},
			TLSServer{},
		},
		"cert file": {
			[]string{"--tls-cert-file", "certificate.file"},
			TLSServer{ServerCert: CertKeyPair{CertFile: "certificate.file"}},
		},
		"key file": {
			[]string{"--tls-private-key-file", "private.key"},
			TLSServer{ServerCert: CertKeyPair{KeyFile: "private.key"}},
		},
		"cert and key files": {
			[]string{
				"--tls-cert-file", "certificate.file",
				"--tls-private-key-file", "private.key",
			},
			TLSServer{ServerCert: CertKeyPair{
				CertFile: "certificate.file",
				KeyFile:  "private.key",
			}},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			opts := &TLSServer{}
			flagSet := flag.NewFlagSet("server-tls-test", flag.ContinueOnError)
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

func TestServerTLSFlagsDefaultText(t *testing.T) {
	flags := TLSServerDefaults().Flags()
	assertWrappedFlagWithDefaultText(t, flags...)
}

func TestServerTLSUsage(t *testing.T) {
	gt := NewGomegaWithT(t)

	opts := &TLSServer{}
	for _, f := range opts.Flags() {
		f := f.(cli.DocGenerationFlag)
		gt.Expect(f.GetUsage()).NotTo(ContainSubstring("\n"))
		gt.Expect(f.GetUsage()).NotTo(MatchRegexp(`\s{2}`))
	}
}
