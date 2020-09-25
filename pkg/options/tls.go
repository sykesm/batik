// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"crypto/tls"
	"regexp"

	"github.com/pkg/errors"
	cli "github.com/urfave/cli/v2"
)

// A CertKeyPair references files containing the TLS server certificate
// and private key.
type CertKeyPair struct {
	// CertFile is the name of a file containing a PEM encoded server certificate
	// or server certificate chain. If the file contains a certificate chain, the
	// PEM blocks must be concatenated such that each certificate certifies the
	// preceding it; the root CA shall be the last certificate in the list.
	CertFile string `yaml:"cert_file,omitempty" batik:"relpath"`
	// KeyFile is the name of a file containing a PEM encoded private key for the
	// certificate provided in CertFile.
	KeyFile string `yaml:"key_file,omitempty" batik:"relpath"`
	// CertData is the PEM encoded server certificate or server certifiate chain.
	CertData string `yaml:"cert,omitempty"`
	// KeyData is the PEM encoded private key for the server certificate.
	KeyData string `yaml:"key,omitempty"`
}

func (ckp CertKeyPair) validate() error {
	// cannot set all four fields
	if ckp.CertData != "" && ckp.KeyData != "" && ckp.CertFile != "" && ckp.KeyFile != "" {
		return errors.New("certificate files and data were both provided but only one is allowed")
	}
	// need to specify cert file or data plus key data or file
	if (ckp.CertData != "" || ckp.CertFile != "") && (ckp.KeyData == "" && ckp.KeyFile == "") {
		return errors.New("private key data or file was not provided")
	}
	// need to specify key file or data plus cert data or file
	if (ckp.KeyData != "" || ckp.KeyFile != "") && (ckp.CertData == "" && ckp.CertFile == "") {
		return errors.New("certificate data or file was not provided")
	}
	return nil
}

func (ckp CertKeyPair) load() (cert tls.Certificate, err error) {
	if err = ckp.validate(); err != nil {
		return cert, err
	}

	if ckp.CertFile != "" {
		cert, err = tls.LoadX509KeyPair(ckp.CertFile, ckp.KeyFile)
	}

	if ckp.CertData != "" {
		cert, err = tls.X509KeyPair([]byte(ckp.CertData), []byte(ckp.KeyData))
	}

	return cert, err
}

// TLSServer exposes configuration options for network services secured
// by TLS.
type TLSServer struct {
	// ServerCert contains the TLS server certifcate and key.
	ServerCert CertKeyPair `yaml:",inline,omitempty"`
}

// TLSServerDefaults returns the default configuration values for TLS servers.
func TLSServerDefaults() *TLSServer {
	return &TLSServer{}
}

// ApplyDefaults applies default values for missing configuration fields.
func (t *TLSServer) ApplyDefaults() {}

// Flags exposes configuration fields as flags. The current value of the
// receiver is used as the default value of the flag so a ApplyDefaults should
// be called before requesting flags.
func (t *TLSServer) Flags() []cli.Flag {
	def := TLSServerDefaults()
	return []cli.Flag{
		NewStringFlag(&cli.StringFlag{
			Name:        "tls-cert-file",
			Value:       t.ServerCert.CertFile,
			Destination: &t.ServerCert.CertFile,
			TakesFile:   true,
			Usage: flow(`File containing the PEM encoded certificate (or chain) for the server.
				When providing a certificate chain, the chain must start with the server certificate
				and the remaining certificates must each certify the preceeding certificate.`),
			DefaultText: def.ServerCert.CertFile,
		}),
		NewStringFlag(&cli.StringFlag{
			Name:        "tls-private-key-file",
			Value:       t.ServerCert.KeyFile,
			Destination: &t.ServerCert.KeyFile,
			TakesFile:   true,
			Usage:       flow(`File containing the PEM encoded private key for the server.`),
			DefaultText: def.ServerCert.KeyFile,
		}),
	}
}

// TLSConfig returns a tls.Config based on the configuration options set for TLSServer
// and returns error for invalid configuration options.
func (t *TLSServer) TLSConfig() (*tls.Config, error) {
	if (t.ServerCert == CertKeyPair{}) {
		return nil, nil
	}

	cert, err := t.ServerCert.load()
	if err != nil {
		return nil, errors.Wrap(err, "options: failed to build TLS configuration")
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}, nil
}

var trimRegex = regexp.MustCompile(`\n\s*`)

func flow(s string) string {
	return trimRegex.ReplaceAllString(s, " ")
}
