// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"crypto/tls"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/pkg/errors"
	cli "github.com/urfave/cli/v2"
	"github.com/sykesm/batik/pkg/tlscerts"
)

const (
	// TLSServerCertName is the filename within the TLS certs dir to load
	// (or store) the server TLS cert.
	TLSServerCertName = "server-cert.pem"

	// TLSServerKeyName is the filename within the TLS certs dir to load
	// (or store) the server TLS key.
	TLSServerKeyName = "server-key.pem"
)

var (
	// ErrServerTLSNotBootstrapped is returned when no server TLS configuration has been specified
	// and the filesystem does not contain certs or keys in the default paths.
	ErrServerTLSNotBootstrapped = errors.Errorf("no server certificate or key found in yaml or tls certs dir")
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

	// CertsDir is the location to generate or find certificates if
	// the ServerCert does not explicitly specify them.
	CertsDir string `yaml:"certs_dir,omitempty" batik:"relpath"`
}

// ServerTLSDefaults returns the default configuration values for TLS servers.
func ServerTLSDefaults() *ServerTLS {
	return &ServerTLS{
		CertsDir: "tls-certs",
	}
}

// ApplyDefaults applies default values for missing configuration fields.
func (s *ServerTLS) ApplyDefaults() {
	defaults := ServerTLSDefaults()
	if s.CertsDir == "" {
		s.CertsDir = defaults.CertsDir
	}
}

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
		NewStringFlag(&cli.StringFlag{
			Name:        "tls-certs-dir",
			Value:       s.CertsDir,
			Destination: &s.CertsDir,
			Usage:       flow(`Path to the directory containing TLS certificates if not otherwise specified.`),
			DefaultText: def.CertsDir,
		}),
	}
}

func (s *ServerTLS) Bootstrap() error {
	err := os.MkdirAll(s.CertsDir, 0700)
	if err != nil {
		return errors.Wrap(err, "failed to create tls-certs-dir")
	}

	sanList := []string{"127.0.0.1", "::1", "localhost"}

	hostname, err := os.Hostname()
	if err == nil {
		sanList = append(sanList, hostname)
	}

	template, err := tlscerts.NewCAServerTemplate(hostname, sanList...)
	if err != nil {
		return errors.WithMessage(err, "failed to generate cert template")
	}

	cert, key, err := tlscerts.GenerateCA(template)
	if err != nil {
		return errors.WithMessage(err, "failed to generate cert")
	}

	certPath := filepath.Join(s.CertsDir, TLSServerCertName)
	keyPath := filepath.Join(s.CertsDir, TLSServerKeyName)

	err = ioutil.WriteFile(certPath, cert, 0644)
	if err != nil {
		return errors.Wrapf(err, "failed to write cert to %s", certPath)
	}

	err = ioutil.WriteFile(keyPath, key, 0600)
	if err != nil {
		return errors.Wrapf(err, "failed to write key to %s", keyPath)
	}

	return nil
}

// TLSConfig returns a *tls.Config from the ServerTLS options or an error if
// the configuration is not valid.
func (s *ServerTLS) TLSConfig() (*tls.Config, error) {
	if (s.ServerCert == CertKeyPair{}) {
		certPath := filepath.Join(s.CertsDir, TLSServerCertName)
		_, certPathErr := os.Stat(certPath)
		keyPath := filepath.Join(s.CertsDir, TLSServerKeyName)
		_, keyPathErr := os.Stat(keyPath)
		if os.IsNotExist(certPathErr) && os.IsNotExist(keyPathErr) {
			return nil, ErrServerTLSNotBootstrapped
		}

		s.ServerCert.CertFile = certPath
		s.ServerCert.KeyFile = keyPath
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
