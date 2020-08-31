// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"regexp"

	"github.com/urfave/cli/v2"
)

// A CertKeyPair references files containing the TLS server certificate
// and private key.
type CertKeyPair struct {
	// CertFile is the name of a file containing a PEM encoded server certificate
	// or server certificate chain. If the file contains a certificate chain, the
	// PEM blocks must be concatenated such that each certificate certifies the
	// preceding it; the root CA shall be the last certificate in the list.
	CertFile string
	// KeyFile is the name of a file containing a PEM encoded private key for the
	// certificate provided in CertFile.
	KeyFile string
}

// TLSOptions exposes configuration options for network services secured by
// TLS.
type TLSOptions struct {
	// ServerCert contains the TLS server certifcate and key.
	ServerCert CertKeyPair
}

// Flags returns flags that can be applied to commands to configure network
// services secured by TLS.
func (o *TLSOptions) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "tls-cert-file",
			Value:       o.ServerCert.CertFile,
			Destination: &o.ServerCert.CertFile,
			TakesFile:   true,
			Usage: flow(`File containing the PEM encoded certificate (or chain) for the server.
				When providing a certificate chain, the chain must start with the server certificate
				and the remaining certificates must each certify the preceeding certificate.`),
		},
		&cli.StringFlag{
			Name:        "tls-private-key-file",
			Value:       o.ServerCert.KeyFile,
			Destination: &o.ServerCert.KeyFile,
			TakesFile:   true,
			Usage:       flow(`File containing the PEM encoded private key for the server.`),
		},
	}
}

var trimRegex = regexp.MustCompile(`\n\s*`)

func flow(s string) string {
	return trimRegex.ReplaceAllString(s, " ")
}
