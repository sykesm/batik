// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tlscerts

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"time"

	"github.com/pkg/errors"
)

// NewBaseTemplate returns an x509.Certificate template with basic fields set.
func NewBaseTemplate(subjectCN string, sans ...string) (x509.Certificate, error) {
	notBefore := time.Now().Add(-1 * time.Minute)
	notAfter := time.Now().Add(time.Duration(365 * 24 * time.Hour))

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return x509.Certificate{}, err
	}

	template := x509.Certificate{
		Subject:               pkix.Name{CommonName: subjectCN},
		SerialNumber:          serialNumber,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}
	for _, s := range sans {
		if ip := net.ParseIP(s); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, s)
		}
	}

	return template, nil
}

// NewCATemplate returns an x509.Certificate template capable of acting as a CA.
func NewCATemplate(subjectCN string, sans ...string) (x509.Certificate, error) {
	template, err := NewBaseTemplate(subjectCN, sans...)
	if err != nil {
		return x509.Certificate{}, err
	}
	template.KeyUsage |= x509.KeyUsageCertSign
	template.IsCA = true
	return template, nil
}

// NewCAServerTemplate returns an x509.Certificate template capable of acting
// both as a CA as well as server auth.  This is typical for a self-signed server
// certificate.
func NewCAServerTemplate(subjectCN string, sans ...string) (x509.Certificate, error) {
	template, err := NewCATemplate(subjectCN, sans...)
	if err != nil {
		return x509.Certificate{}, err
	}
	template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	return template, nil
}

// GenerateCA takes a template, generates an ecdsa key, and creates a certificate
// self-signed with the generated key.
func GenerateCA(template x509.Certificate) (pemCert, pemKey []byte, err error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	publicKey := privateKey.Public()

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey, privateKey)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not create certificate from template")
	}

	return PemEncode(derBytes, privateKey)
}

// PemEncode takes a der encoded certificate, and an ecdsa private key, and returns each
// as Pem encoded blocks.
func PemEncode(derCert []byte, key *ecdsa.PrivateKey) (pemCert, pemKey []byte, err error) {
	keyBytes, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return nil, nil, err
	}

	pemCert = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derCert})
	pemKey = pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes})
	return pemCert, pemKey, nil
}
