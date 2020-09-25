// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"io/ioutil"
	"math/big"
	"path/filepath"
	"testing"

	"github.com/sykesm/batik/pkg/tested"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
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

func TestTLSConfig(t *testing.T) {
	gt := NewGomegaWithT(t)

	dir, cleanup := tested.TempDir(t, "", "options_tls")
	defer cleanup()

	keyPEM, certPEM, keyFile, certFile := genCKPPem(dir, gt)
	// create expected TLS Config from generated key/cert pair
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	gt.Expect(err).NotTo(HaveOccurred())
	expectedTLSConf := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	tests := map[string]struct {
		srv         TLSServer
		confMatcher *tls.Config
		errMatcher  types.GomegaMatcher
	}{
		"empty TLSServer": {
			srv:         TLSServer{},
			confMatcher: nil,
			errMatcher:  BeNil(),
		},
		"both data and files": {
			srv: TLSServer{
				ServerCert: CertKeyPair{
					CertData: string(certPEM),
					KeyData:  string(keyPEM),
					CertFile: certFile,
					KeyFile:  keyFile,
				},
			},
			confMatcher: nil,
			errMatcher:  MatchError("options: failed to build TLS configuration: certificate files and data were both provided but only one is allowed"),
		},
		"valid data": {
			srv: TLSServer{
				ServerCert: CertKeyPair{
					CertData: string(certPEM),
					KeyData:  string(keyPEM),
				},
			},
			confMatcher: expectedTLSConf,
			errMatcher:  BeNil(),
		},
		"valid files": {
			srv: TLSServer{
				ServerCert: CertKeyPair{
					CertFile: certFile,
					KeyFile:  keyFile,
				},
			},
			confMatcher: expectedTLSConf,
			errMatcher:  BeNil(),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tlsConf, err := tt.srv.TLSConfig()
			gt.Expect(tlsConf).To(Equal(tt.confMatcher))
			gt.Expect(err).To(tt.errMatcher)
		})
	}
}

func TestValidateCPK(t *testing.T) {
	tests := map[string]struct {
		ckp        CertKeyPair
		errMatcher types.GomegaMatcher
	}{
		"both data and files": {
			ckp: CertKeyPair{
				CertData: "test",
				KeyData:  "test",
				CertFile: "test",
				KeyFile:  "test",
			},
			errMatcher: MatchError("certificate files and data were both provided but only one is allowed"),
		},
		"key file provided - missing cert": {
			ckp: CertKeyPair{
				KeyFile: "test",
			},
			errMatcher: MatchError("certificate data or file was not provided"),
		},
		"cert file provided - missing key": {
			ckp: CertKeyPair{
				CertFile: "test",
			},
			errMatcher: MatchError("private key data or file was not provided"),
		},
		"key data provided - missing cert": {
			ckp: CertKeyPair{
				KeyData: "test",
			},
			errMatcher: MatchError("certificate data or file was not provided"),
		},
		"cert data provided - missing key": {
			ckp: CertKeyPair{
				CertData: "test",
			},
			errMatcher: MatchError("private key data or file was not provided"),
		},
		"cert data and key data": {
			ckp: CertKeyPair{
				CertData: "test",
				KeyData:  "test",
			},
			errMatcher: BeNil(),
		},
		"cert data and key file": {
			ckp: CertKeyPair{
				CertData: "test",
				KeyFile:  "test",
			},
			errMatcher: BeNil(),
		},
		"cert file and key file": {
			ckp: CertKeyPair{
				CertFile: "test",
				KeyFile:  "test",
			},
			errMatcher: BeNil(),
		},
		"cert file and key data": {
			ckp: CertKeyPair{
				CertFile: "test",
				KeyData:  "test",
			},
			errMatcher: BeNil(),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)
			err := tt.ckp.validate()
			gt.Expect(err).To(tt.errMatcher)
		})
	}
}

func TestLoadCKP(t *testing.T) {
	gt := NewGomegaWithT(t)

	dir, cleanup := tested.TempDir(t, "", "options_tls")
	defer cleanup()

	keyPEM, certPEM, keyFile, certFile := genCKPPem(dir, gt)
	// create certificate
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		t.Fatalf("failed to load certificate: %s", err)
	}

	tests := map[string]struct {
		ckp       CertKeyPair
		expectErr bool
		cert      tls.Certificate
	}{
		"valid files": {
			ckp: CertKeyPair{
				CertFile: certFile,
				KeyFile:  keyFile,
			},
			expectErr: false,
			cert:      cert,
		},
		"valid data": {
			ckp: CertKeyPair{
				CertData: string(certPEM),
				KeyData:  string(keyPEM),
			},
			expectErr: false,
			cert:      cert,
		},
		"invalid files": {
			ckp: CertKeyPair{
				CertFile: filepath.Join(dir, "doesnotexist.pem"),
				KeyFile:  filepath.Join(dir, "doesnotexist.pem"),
			},
			expectErr: true,
			cert:      tls.Certificate{},
		},
		"invalid data": {
			ckp: CertKeyPair{
				CertData: "not PEM encoded",
				KeyData:  "not PEM encoded",
			},
			expectErr: true,
			cert:      tls.Certificate{},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			cert, err := tt.ckp.load()
			gt.Expect(cert).To(Equal(tt.cert))
			if tt.expectErr {
				gt.Expect(err).To(HaveOccurred())
			}

		})
	}

}

// genCKPPem generates a PEM-encoded ECDSA private key and X.509 certificate pair.
// It returns both raw bytes as well as paths to files created in directory dir.
func genCKPPem(dir string, gt *GomegaWithT) (keyPEM, certPEM []byte, keyFile, certFile string) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	gt.Expect(err).NotTo(HaveOccurred())
	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	gt.Expect(err).NotTo(HaveOccurred())
	var privBuf bytes.Buffer
	err = pem.Encode(
		&privBuf,
		&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: privBytes,
		},
	)
	gt.Expect(err).NotTo(HaveOccurred())

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
	}
	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	gt.Expect(err).NotTo(HaveOccurred())

	var certBuf bytes.Buffer
	err = pem.Encode(
		&certBuf,
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: certBytes,
		},
	)
	gt.Expect(err).NotTo(HaveOccurred())

	keyPEM = privBuf.Bytes()
	certPEM = certBuf.Bytes()
	keyFile = filepath.Join(dir, "key.pem")
	certFile = filepath.Join(dir, "cert.pem")
	// create key/cert files
	err = ioutil.WriteFile(keyFile, keyPEM, 0644)
	gt.Expect(err).NotTo(HaveOccurred())
	err = ioutil.WriteFile(certFile, certPEM, 0644)
	gt.Expect(err).NotTo(HaveOccurred())

	return
}
