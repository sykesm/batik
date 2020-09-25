// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tested

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"testing"

	. "github.com/onsi/gomega"
)

func TestCertificates(t *testing.T) {
	ca := NewCA(t, "my-ca")

	t.Run("CA", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		cert, err := x509.ParseCertificate(ca.Certificate.Certificate[0])
		gt.Expect(err).NotTo(HaveOccurred())

		gt.Expect(cert.IsCA).To(BeTrue())
		gt.Expect(cert.Subject.CommonName).To(Equal("my-ca"))
		gt.Expect(cert.Subject.CommonName).To(Equal(cert.Issuer.CommonName))
		gt.Expect(cert.KeyUsage).To(Equal(x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign))
		gt.Expect(cert.ExtKeyUsage).To(HaveLen(0))
	})

	t.Run("Server", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		srv := ca.IssueServerCertificate(t, "server", "127.0.0.1", "::1", "localhost")
		cert, err := x509.ParseCertificate(srv.Certificate.Certificate[0])
		gt.Expect(err).NotTo(HaveOccurred())

		gt.Expect(cert.IsCA).To(BeFalse())
		gt.Expect(cert.Subject.CommonName).To(Equal("server"))
		gt.Expect(cert.Issuer.CommonName).To(Equal("my-ca"))
		gt.Expect(cert.KeyUsage).To(Equal(x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature))
		gt.Expect(cert.ExtKeyUsage).To(ConsistOf(x509.ExtKeyUsageServerAuth))
		gt.Expect(cert.IPAddresses).To(ConsistOf(
			net.IP([]byte{127, 0, 0, 1}),
			net.IPv6loopback,
		))
	})

	t.Run("Client", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		clt := ca.IssueClientCertificate(t, "client", "127.0.0.1")
		cert, err := x509.ParseCertificate(clt.Certificate.Certificate[0])
		gt.Expect(err).NotTo(HaveOccurred())

		gt.Expect(cert.IsCA).To(BeFalse())
		gt.Expect(cert.Subject.CommonName).To(Equal("client"))
		gt.Expect(cert.Issuer.CommonName).To(Equal("my-ca"))
		gt.Expect(cert.KeyUsage).To(Equal(x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature))
		gt.Expect(cert.ExtKeyUsage).To(ConsistOf(x509.ExtKeyUsageClientAuth))
		gt.Expect(cert.IPAddresses).To(ConsistOf([]net.IP{
			net.IP([]byte{127, 0, 0, 1}),
		}))
	})
}

func TestCATLSConfig(t *testing.T) {
	gt := NewGomegaWithT(t)

	ca := NewCA(t, "ca")
	cert, err := x509.ParseCertificate(ca.Certificate.Certificate[0])
	gt.Expect(err).NotTo(HaveOccurred())

	tlsConfig := ca.TLSConfig(t)
	gt.Expect(tlsConfig.RootCAs.Subjects()).To(HaveLen(1))
	gt.Expect(tlsConfig.RootCAs.Subjects()[0]).To(Equal(cert.RawSubject))
}

func TestCertificatesTLSConfig(t *testing.T) {
	gt := NewGomegaWithT(t)
	serverCA := NewCA(t, "server-ca")
	clientCA := NewCA(t, "client-ca")

	srv := serverCA.IssueServerCertificate(t, "server", "127.0.0.1", "::1", "localhost")
	clt := clientCA.IssueClientCertificate(t, "client", "127.0.0.1")

	serverTLSConfig := srv.ServerTLSConfig(t, &clientCA.Certificate)
	clientTLSConfig := clt.ClientTLSConfig(t, &serverCA.Certificate)

	lis, err := tls.Listen("tcp", "127.0.0.1:0", serverTLSConfig)
	gt.Expect(err).NotTo(HaveOccurred())

	serverResultCh := make(chan error, 1)
	go func() {
		defer lis.Close()
		conn, err := lis.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		_, err = conn.Write([]byte("bye"))
		serverResultCh <- err
	}()

	conn, err := tls.Dial("tcp", lis.Addr().String(), clientTLSConfig)
	gt.Expect(err).NotTo(HaveOccurred())
	defer conn.Close()

	res, err := ioutil.ReadAll(conn)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(string(res)).To(Equal("bye"))

	gt.Eventually(serverResultCh).Should(Receive(BeNil()))
}
