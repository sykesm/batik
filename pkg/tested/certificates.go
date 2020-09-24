// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tested

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"time"
)

type CertKeyPair struct {
	Cert        []byte
	Key         []byte
	Certificate tls.Certificate
}

type CA CertKeyPair

func NewCA(t TestingT, subjectCN string) CA {
	cert, key := generateCA(t, subjectCN)
	certificate, err := tls.X509KeyPair(cert, key)
	assertNoError(t, err)
	return CA{Cert: cert, Key: key, Certificate: certificate}
}

func (ca CA) IssueServerCertificate(t TestingT, subjectCN string, sans ...string) CertKeyPair {
	cacert, err := x509.ParseCertificate(ca.Certificate.Certificate[0])
	assertNoError(t, err)

	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	assertNoError(t, err)
	publicKey := privateKey.Public()

	template := newTemplate(t, subjectCN, sans...)
	template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, cacert, publicKey, ca.Certificate.PrivateKey)
	assertNoError(t, err)

	cert, key := pemEncode(t, derBytes, privateKey)
	certificate, err := tls.X509KeyPair(join(cert, ca.Cert), key)
	assertNoError(t, err)

	return CertKeyPair{Cert: cert, Key: key, Certificate: certificate}
}

func (ca CA) IssueClientCertificate(t TestingT, subjectCN string, sans ...string) CertKeyPair {
	cacert, err := x509.ParseCertificate(ca.Certificate.Certificate[0])
	assertNoError(t, err)

	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	assertNoError(t, err)
	publicKey := privateKey.Public()

	template := newTemplate(t, subjectCN, sans...)
	template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, cacert, publicKey, ca.Certificate.PrivateKey)
	assertNoError(t, err)

	cert, key := pemEncode(t, derBytes, privateKey)
	certificate, err := tls.X509KeyPair(join(cert, ca.Cert), key)
	assertNoError(t, err)

	return CertKeyPair{Cert: cert, Key: key, Certificate: certificate}
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

func generateCA(t TestingT, subjectCN string) (pemCert, pemKey []byte) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	assertNoError(t, err)
	publicKey := privateKey.Public()

	template := newTemplate(t, subjectCN)
	template.KeyUsage |= x509.KeyUsageCertSign
	template.IsCA = true

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey, privateKey)
	assertNoError(t, err)

	return pemEncode(t, derBytes, privateKey)
}

func newTemplate(t TestingT, subjectCN string, sans ...string) x509.Certificate {
	notBefore := time.Now().Add(-1 * time.Minute)
	notAfter := time.Now().Add(time.Duration(365 * 24 * time.Hour))

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	assertNoError(t, err)

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

	return template
}

func pemEncode(t TestingT, derCert []byte, key *ecdsa.PrivateKey) (pemCert, pemKey []byte) {
	certBuf := &bytes.Buffer{}
	err := pem.Encode(certBuf, &pem.Block{Type: "CERTIFICATE", Bytes: derCert})
	assertNoError(t, err)

	keyBytes, err := x509.MarshalECPrivateKey(key)
	assertNoError(t, err)

	keyBuf := &bytes.Buffer{}
	err = pem.Encode(keyBuf, &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes})
	assertNoError(t, err)

	return certBuf.Bytes(), keyBuf.Bytes()
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
