// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package tlscerts

import (
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"net"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestNewBaseTemplate(t *testing.T) {
	gt := NewGomegaWithT(t)

	template, err := NewBaseTemplate("some-cn", "addr1.com", "127.0.0.1")
	gt.Expect(err).NotTo(HaveOccurred())

	gt.Expect(template.Subject).To(Equal(pkix.Name{CommonName: "some-cn"}))
	gt.Expect(template.SerialNumber).NotTo(Equal(0))
	gt.Expect(template.NotAfter).To(BeTemporally(">", time.Now()))
	gt.Expect(template.NotBefore).To(BeTemporally("<", time.Now()))
	gt.Expect(template.KeyUsage).To(Equal(x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature))
	gt.Expect(template.ExtKeyUsage).To(HaveLen(0))
	gt.Expect(template.BasicConstraintsValid).To(BeTrue())

	gt.Expect(template.IPAddresses).To(Equal([]net.IP{net.ParseIP("127.0.0.1")}))
	gt.Expect(template.DNSNames).To(Equal([]string{"addr1.com"}))

}

func TestNewCATemplate(t *testing.T) {
	gt := NewGomegaWithT(t)

	template, err := NewCATemplate("some-cn", "addr1.com", "127.0.0.1")
	gt.Expect(err).NotTo(HaveOccurred())

	gt.Expect(template.Subject).To(Equal(pkix.Name{CommonName: "some-cn"}))
	gt.Expect(template.SerialNumber).NotTo(Equal(0))
	gt.Expect(time.Since(template.NotAfter)).To(BeNumerically("<", 0))
	gt.Expect(time.Since(template.NotBefore)).To(BeNumerically(">", 0))
	gt.Expect(template.KeyUsage).To(Equal(x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign))
	gt.Expect(template.ExtKeyUsage).To(HaveLen(0))
	gt.Expect(template.BasicConstraintsValid).To(BeTrue())

	gt.Expect(template.IPAddresses).To(Equal([]net.IP{net.ParseIP("127.0.0.1")}))
	gt.Expect(template.DNSNames).To(Equal([]string{"addr1.com"}))

}

func TestNewCAServerTemplate(t *testing.T) {
	gt := NewGomegaWithT(t)

	template, err := NewCAServerTemplate("some-cn", "addr1.com", "127.0.0.1")
	gt.Expect(err).NotTo(HaveOccurred())

	gt.Expect(template.Subject).To(Equal(pkix.Name{CommonName: "some-cn"}))
	gt.Expect(template.SerialNumber).NotTo(Equal(0))
	gt.Expect(time.Since(template.NotAfter)).To(BeNumerically("<", 0))
	gt.Expect(time.Since(template.NotBefore)).To(BeNumerically(">", 0))
	gt.Expect(template.KeyUsage).To(Equal(x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign))
	gt.Expect(template.ExtKeyUsage).To(Equal([]x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}))
	gt.Expect(template.BasicConstraintsValid).To(BeTrue())

	gt.Expect(template.IPAddresses).To(Equal([]net.IP{net.ParseIP("127.0.0.1")}))
	gt.Expect(template.DNSNames).To(Equal([]string{"addr1.com"}))

}

func TestGenerateCA(t *testing.T) {
	gt := NewGomegaWithT(t)

	template, err := NewCATemplate("some-cn", "addr1.com", "127.0.0.1")
	gt.Expect(err).NotTo(HaveOccurred())

	cert, key, err := GenerateCA(template)
	gt.Expect(err).NotTo(HaveOccurred())

	_, err = tls.X509KeyPair(cert, key)
	gt.Expect(err).NotTo(HaveOccurred())
}
