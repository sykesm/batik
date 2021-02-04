// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"net"
	"os"
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

	err := ioutil.WriteFile(ckp.CertFile, []byte(ckp.CertData), 0o644)
	gt.Expect(err).NotTo(HaveOccurred())
	err = ioutil.WriteFile(ckp.KeyFile, []byte(ckp.KeyData), 0o644)
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
	gt.Expect(ts).To(Equal(&ServerTLS{CertsDir: "tls-certs"}))
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

	gt.Expect(flags).To(HaveLen(3))
	gt.Expect(names).To(ConsistOf(
		"tls-cert-file",
		"tls-private-key-file",
		"tls-certs-dir",
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
		"certs dir": {
			[]string{"--tls-certs-dir", "some-dir"},
			ServerTLS{CertsDir: "some-dir"},
		},
		"cert and key files and certs dir": {
			[]string{
				"--tls-cert-file", "certificate.file",
				"--tls-private-key-file", "private.key",
				"--tls-certs-dir", "some-dir",
			},
			ServerTLS{
				ServerCert: CertKeyPair{
					CertFile: "certificate.file",
					KeyFile:  "private.key",
				},
				CertsDir: "some-dir",
			},
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
		"missing key pair": {errMatcher: MatchError(ErrServerTLSNotBootstrapped)},
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

func TestServerTLSBootstrap(t *testing.T) {
	tests := map[string]struct {
		errMatcher types.GomegaMatcher
		inspectTLS func(*GomegaWithT, *tls.Config)
		mangleDir  func(gt *GomegaWithT, certsDir string)
	}{
		"green path": {
			inspectTLS: func(gt *GomegaWithT, config *tls.Config) {
				gt.Expect(config.Certificates).To(HaveLen(1))
				gt.Expect(config.Certificates[0].Certificate).To(HaveLen(1))
				rawCert := config.Certificates[0].Certificate[0]

				cert, err := x509.ParseCertificate(rawCert)
				gt.Expect(err).NotTo(HaveOccurred())
				gt.Expect(cert).NotTo(BeNil())

				gt.Expect(cert.DNSNames).To(ContainElement("localhost"))
				gt.Expect(cert.IPAddresses).To(ConsistOf(
					net.ParseIP("::1"),
					// Weirdly, net.IPv4 produces an IPv6 length address and doesn't match
					net.IP([]byte{127, 0, 0, 1}),
				))
			},
		},
		"unwritable dir": {
			mangleDir: func(gt *GomegaWithT, certsDir string) {
				file, err := os.Create(certsDir)
				gt.Expect(err).NotTo(HaveOccurred())

				file.Close()
			},
			errMatcher: MatchError(MatchRegexp("failed to create tls-certs-dir.*")),
		},
		"unwritable cert": {
			mangleDir: func(gt *GomegaWithT, certsDir string) {
				err := os.MkdirAll(filepath.Join(certsDir, "server-cert.pem"), 0o755)
				gt.Expect(err).NotTo(HaveOccurred())
			},
			errMatcher: MatchError(MatchRegexp("failed to write cert to.*")),
		},
		"unwritable key": {
			mangleDir: func(gt *GomegaWithT, certsDir string) {
				err := os.MkdirAll(filepath.Join(certsDir, "server-key.pem"), 0o755)
				gt.Expect(err).NotTo(HaveOccurred())
			},
			errMatcher: MatchError(MatchRegexp("failed to write key to.*")),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			tempDir, cleanup := tested.TempDir(t, "", "options_tls_bootstrap")
			defer cleanup()

			serverTLSConf := &ServerTLS{
				CertsDir: filepath.Join(tempDir, "certs-dir"),
			}

			if tt.mangleDir != nil {
				tt.mangleDir(gt, serverTLSConf.CertsDir)
			}

			err := serverTLSConf.Bootstrap()
			if tt.errMatcher != nil {
				gt.Expect(err).To(tt.errMatcher)
				return
			}

			tlsConf, err := serverTLSConf.TLSConfig()
			gt.Expect(err).NotTo(HaveOccurred())

			if tt.inspectTLS != nil {
				tt.inspectTLS(gt, tlsConf)
			}
		})
	}
}
