// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"crypto/tls"
	"io/ioutil"
	"regexp"

	"github.com/pkg/errors"
	cli "github.com/urfave/cli/v2"
)

// A CertKeyPair references files containing the TLS certificate and private
// key.
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
	// If CertFile is set, CertData is ignored.
	CertData string `yaml:"cert,omitempty"`
	// KeyData is the PEM encoded private key for the server certificate. If
	// KeyFile is set, Keydata is ignored.
	KeyData string `yaml:"key,omitempty"`
}

func (c *CertKeyPair) certData() (data []byte, err error) {
	switch {
	case c.CertFile != "":
		data, err = ioutil.ReadFile(c.CertFile)
		if err != nil {
			return nil, errors.Wrap(err, "unable to read certificate file")
		}
	case c.CertData != "":
		data = []byte(c.CertData)
	}
	return data, nil
}

func (c *CertKeyPair) keyData() (data []byte, err error) {
	switch {
	case c.KeyFile != "":
		data, err = ioutil.ReadFile(c.KeyFile)
		if err != nil {
			return nil, errors.Wrap(err, "unable to read private key file")
		}
	case c.KeyData != "":
		data = []byte(c.KeyData)
	}
	return data, nil
}

// TLSCertificate returns a tls.Certificate from a CertKeyPair.
func (c *CertKeyPair) TLSCertificate() (tls.Certificate, error) {
	cert, err := c.certData()
	if err != nil {
		return tls.Certificate{}, err
	}
	key, err := c.keyData()
	if err != nil {
		return tls.Certificate{}, err
	}
	return tls.X509KeyPair(cert, key)
}

// ServerTLS exposes configuration options for network services secured
// by TLS.
type ServerTLS struct {
	// ServerCert contains the TLS certifcate and key for a server.
	ServerCert CertKeyPair `yaml:",inline,omitempty"`
}

// ServerTLSDefaults returns the default configuration values for TLS servers.
func ServerTLSDefaults() *ServerTLS {
	return &ServerTLS{}
}

// ApplyDefaults applies default values for missing configuration fields.
func (s *ServerTLS) ApplyDefaults() {}

// Flags exposes configuration fields as flags. The current value of the
// receiver is used as the default value of the flag so a ApplyDefaults should
// be called before requesting flags.
func (s *ServerTLS) Flags() []cli.Flag {
	def := ServerTLSDefaults()
	return []cli.Flag{
		NewStringFlag(&cli.StringFlag{
			Name:        "tls-cert-file",
			Value:       s.ServerCert.CertFile,
			Destination: &s.ServerCert.CertFile,
			TakesFile:   true,
			Usage: flow(`File containing the PEM encoded certificate (or chain) for the server.
				When providing a certificate chain, the chain must start with the server certificate
				and the remaining certificates must each certify the preceeding certificate.`),
			DefaultText: def.ServerCert.CertFile,
		}),
		NewStringFlag(&cli.StringFlag{
			Name:        "tls-private-key-file",
			Value:       s.ServerCert.KeyFile,
			Destination: &s.ServerCert.KeyFile,
			TakesFile:   true,
			Usage:       flow(`File containing the PEM encoded private key for the server.`),
			DefaultText: def.ServerCert.KeyFile,
		}),
	}
}

// TLSConfig returns a *tls.Config from the ServerTLS options or an error if
// the configuration is not valid.
func (s *ServerTLS) TLSConfig() (*tls.Config, error) {
	if (s.ServerCert == CertKeyPair{}) {
		return nil, nil
	}

	cert, err := s.ServerCert.TLSCertificate()
	if err != nil {
		return nil, err
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
