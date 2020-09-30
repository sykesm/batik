// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"crypto/tls"
	"flag"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/sykesm/batik/pkg/tested"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"github.com/urfave/cli/v2"
)

func TestCertKeyPairTLSCertificate(t *testing.T) {
	gt := NewGomegaWithT(t)
	tempDir, cleanup := tested.TempDir(t, "", "options_tls")
	defer cleanup()

	ca := tested.NewCA(t, "ca")
	skp := ca.IssueServerCertificate(t, "server", "127.0.0.1")

	ckp := CertKeyPair{
		KeyData:  string(skp.Key),
		KeyFile:  filepath.Join(tempDir, "server.crt"),
		CertData: string(skp.CertChain),
		CertFile: filepath.Join(tempDir, "server.key"),
	}

	err := ioutil.WriteFile(ckp.CertFile, []byte(ckp.CertData), 0644)
	gt.Expect(err).NotTo(HaveOccurred())
	err = ioutil.WriteFile(ckp.KeyFile, []byte(ckp.KeyData), 0644)
	gt.Expect(err).NotTo(HaveOccurred())

	tests := map[string]struct {
		setup      func(*CertKeyPair)
		expected   tls.Certificate
		errMatcher types.GomegaMatcher
	}{
		"KeyData": {
			setup:    func(c *CertKeyPair) { c.KeyFile = "" },
			expected: skp.Certificate,
		},
		"CertData": {
			setup:    func(c *CertKeyPair) { c.CertFile = "" },
			expected: skp.Certificate,
		},
		"KeyFile": {
			setup:    func(c *CertKeyPair) { c.KeyData = "" },
			expected: skp.Certificate,
		},
		"CertFile": {
			setup:    func(c *CertKeyPair) { c.CertData = "" },
			expected: skp.Certificate,
		},
		"CertDataIgnored": {
			setup:    func(c *CertKeyPair) { c.CertData = "bogus-data" },
			expected: skp.Certificate,
		},
		"KeyDataIgnored": {
			setup:    func(c *CertKeyPair) { c.KeyData = "bogus-data" },
			expected: skp.Certificate,
		},
		"Empty": {
			setup:      func(c *CertKeyPair) { *c = CertKeyPair{} },
			errMatcher: MatchError(MatchRegexp("tls:.*certificate input")),
		},
		"BadKeyFile": {
			setup:      func(c *CertKeyPair) { c.KeyFile = "missing.txt" },
			errMatcher: MatchError(MatchRegexp("unable to read private key file.*missing.txt")),
		},
		"BadCertFile": {
			setup:      func(c *CertKeyPair) { c.CertFile = "missing.txt" },
			errMatcher: MatchError(MatchRegexp("unable to read certificate file.*missing.txt")),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)
			ckp := ckp
			tt.setup(&ckp)

			cert, err := ckp.TLSCertificate()
			if tt.errMatcher != nil {
				gt.Expect(err).To(tt.errMatcher)
				return
			}
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(cert).To(Equal(tt.expected))
		})
	}
}

func TestServerTLSDefaults(t *testing.T) {
	gt := NewGomegaWithT(t)
	ts := ServerTLSDefaults()
	gt.Expect(ts).To(Equal(&ServerTLS{}))
}

func TestServerTLSApplyDefaults(t *testing.T) {
	tests := map[string]struct {
		setup func(ts *ServerTLS)
	}{
		"empty": {setup: func(ts *ServerTLS) { *ts = ServerTLS{} }},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			input := ServerTLSDefaults()
			tt.setup(input)

			input.ApplyDefaults()
			gt.Expect(input).To(Equal(ServerTLSDefaults()))
		})
	}
}

func TestServerTLSFlagNames(t *testing.T) {
	gt := NewGomegaWithT(t)
	flags := (&ServerTLS{}).Flags()

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

func TestServerTLSFlags(t *testing.T) {
	tests := map[string]struct {
		args     []string
		expected ServerTLS
	}{
		"no flags": {
			[]string{},
			ServerTLS{},
		},
		"cert file": {
			[]string{"--tls-cert-file", "certificate.file"},
			ServerTLS{ServerCert: CertKeyPair{CertFile: "certificate.file"}},
		},
		"key file": {
			[]string{"--tls-private-key-file", "private.key"},
			ServerTLS{ServerCert: CertKeyPair{KeyFile: "private.key"}},
		},
		"cert and key files": {
			[]string{
				"--tls-cert-file", "certificate.file",
				"--tls-private-key-file", "private.key",
			},
			ServerTLS{ServerCert: CertKeyPair{
				CertFile: "certificate.file",
				KeyFile:  "private.key",
			}},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			opts := &ServerTLS{}
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
	flags := ServerTLSDefaults().Flags()
	assertWrappedFlagWithDefaultText(t, flags...)
}

func TestServerTLSUsage(t *testing.T) {
	gt := NewGomegaWithT(t)

	opts := &ServerTLS{}
	for _, f := range opts.Flags() {
		f := f.(cli.DocGenerationFlag)
		gt.Expect(f.GetUsage()).NotTo(ContainSubstring("\n"))
		gt.Expect(f.GetUsage()).NotTo(MatchRegexp(`\s{2}`))
	}
}

func TestServerTLSConfig(t *testing.T) {
	ca := tested.NewCA(t, "ca")
	skp := ca.IssueClientCertificate(t, "server", "127.0.0.1")
	serverCert := skp.Certificate
	serverCert.Certificate = serverCert.Certificate[0:1]

	expectedTLSConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		MinVersion:   tls.VersionTLS12,
	}

	tests := map[string]struct {
		srv        ServerTLS
		expected   *tls.Config
		errMatcher types.GomegaMatcher
	}{
		"missing key pair": {expected: nil},
		"valid key pair": {
			srv: ServerTLS{
				ServerCert: CertKeyPair{CertData: string(skp.Cert), KeyData: string(skp.Key)},
			},
			expected: expectedTLSConfig,
		},
		"invalid key pair": {
			srv: ServerTLS{
				ServerCert: CertKeyPair{CertFile: "missing.txt"},
			},
			errMatcher: MatchError(MatchRegexp("unable to read certificate file.*missing.txt")),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			tlsConf, err := tt.srv.TLSConfig()
			if tt.errMatcher != nil {
				gt.Expect(err).To(tt.errMatcher)
				return
			}
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(tlsConf).To(Equal(tt.expected))
		})
	}
}
