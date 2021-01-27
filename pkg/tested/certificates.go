// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tested

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"

	"github.com/sykesm/batik/pkg/tlscerts"
)

type CertKeyPair struct {
	Cert        []byte
	CertChain   []byte
	Key         []byte
	Certificate tls.Certificate
}

type CA CertKeyPair

func NewCA(t TestingT, subjectCN string) CA {
	caTemplate, err := tlscerts.NewCATemplate(subjectCN)
	assertNoError(t, err)

	cert, key, err := tlscerts.GenerateCA(caTemplate)
	assertNoError(t, err)

	certificate, err := tls.X509KeyPair(cert, key)
	assertNoError(t, err)
	return CA{Cert: cert, CertChain: cert, Key: key, Certificate: certificate}
}

func (ca CA) IssueServerCertificate(t TestingT, subjectCN string, sans ...string) CertKeyPair {
	cacert, err := x509.ParseCertificate(ca.Certificate.Certificate[0])
	assertNoError(t, err)

	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	assertNoError(t, err)
	publicKey := privateKey.Public()

	template, err := tlscerts.NewBaseTemplate(subjectCN, sans...)
	assertNoError(t, err)

	template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, cacert, publicKey, ca.Certificate.PrivateKey)
	assertNoError(t, err)

	cert, key, err := tlscerts.PemEncode(derBytes, privateKey)
	assertNoError(t, err)

	certChain := join(cert, ca.Cert)
	certificate, err := tls.X509KeyPair(certChain, key)
	assertNoError(t, err)

	return CertKeyPair{Cert: cert, CertChain: certChain, Key: key, Certificate: certificate}
}

func (ca CA) IssueClientCertificate(t TestingT, subjectCN string, sans ...string) CertKeyPair {
	cacert, err := x509.ParseCertificate(ca.Certificate.Certificate[0])
	assertNoError(t, err)

	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	assertNoError(t, err)
	publicKey := privateKey.Public()

	template, err := tlscerts.NewBaseTemplate(subjectCN, sans...)
	assertNoError(t, err)
	template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, cacert, publicKey, ca.Certificate.PrivateKey)
	assertNoError(t, err)

	cert, key, err := tlscerts.PemEncode(derBytes, privateKey)
	assertNoError(t, err)

	certChain := join(cert, ca.Cert)
	certificate, err := tls.X509KeyPair(certChain, key)
	assertNoError(t, err)

	return CertKeyPair{Cert: cert, CertChain: certChain, Key: key, Certificate: certificate}
}

func (ca CA) TLSConfig(t TestingT) *tls.Config {
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(ca.Cert)
	return &tls.Config{
		RootCAs: caCertPool,
	}
}

func (ckp CertKeyPair) ServerTLSConfig(t TestingT, clientCA *tls.Certificate) *tls.Config {
	caCertPool := x509.NewCertPool()
	for i := 1; i < len(ckp.Certificate.Certificate); i++ {
		caCertPool.AppendCertsFromPEM(ckp.Certificate.Certificate[i])
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{ckp.Certificate},
		RootCAs:      caCertPool,
	}
	if clientCA != nil {
		tlsConfig.ClientCAs = x509.NewCertPool()
		cacert, err := x509.ParseCertificate(clientCA.Certificate[0])
		assertNoError(t, err)
		tlsConfig.ClientCAs.AddCert(cacert)
		tlsConfig.ClientAuth = tls.VerifyClientCertIfGiven
	}
	return tlsConfig
}

func (ckp CertKeyPair) ClientTLSConfig(t TestingT, serverCA *tls.Certificate) *tls.Config {
	caCertPool := x509.NewCertPool()
	cacert, err := x509.ParseCertificate(serverCA.Certificate[0])
	assertNoError(t, err)
	caCertPool.AddCert(cacert)

	return &tls.Config{
		Certificates:       []tls.Certificate{ckp.Certificate},
		RootCAs:            caCertPool,
		ClientSessionCache: tls.NewLRUClientSessionCache(10),
	}
}

func join(b ...[]byte) []byte {
	var res []byte
	for _, b := range b {
		res = append(res, b...)
	}
	return res
}

func assertNoError(t TestingT, err error) {
	if err != nil {
		if h, ok := t.(tHelper); ok {
			h.Helper()
		}
		t.Fatalf("unexpected error: %s", err)
	}
}
