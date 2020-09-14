// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"os"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	yaml "gopkg.in/yaml.v3"
)

func TestBatikDefaults(t *testing.T) {
	gt := NewGomegaWithT(t)
	config := BatikDefaults()
	gt.Expect(config).To(Equal(&Batik{
		Server: *ServerDefaults(),
		Ledger: *LedgerDefaults(),
	}))
}

func TestBatikApplyDefaults(t *testing.T) {
	tests := map[string]struct {
		setup    func(*Batik)
		matchErr types.GomegaMatcher
	}{
		"empty":  {setup: func(c *Batik) { *c = Batik{} }, matchErr: BeNil()},
		"server": {setup: func(c *Batik) { c.Server = Server{} }, matchErr: BeNil()},
		"ledger": {setup: func(c *Batik) { c.Ledger = Ledger{} }, matchErr: BeNil()},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			input := BatikDefaults()
			tt.setup(input)

			err := input.ApplyDefaults()
			gt.Expect(err).To(tt.matchErr)
			if err != nil {
				return
			}
			gt.Expect(input).To(Equal(BatikDefaults()))
		})
	}
}

func TestReadConfigFileApplyDefaults(t *testing.T) {
	gt := NewGomegaWithT(t)

	cf, err := os.Open("testdata/batik.yaml")
	gt.Expect(err).NotTo(HaveOccurred())
	defer cf.Close()

	var config Batik
	decoder := yaml.NewDecoder(cf)

	err = decoder.Decode(&config)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(config).To(Equal(Batik{
		Server: Server{
			ListenAddress: "127.0.0.1:7879",
			GRPC: GRPCServer{
				GRPC: GRPC{MaxSendMessageSize: 104857600},
			},
			TLS: TLSServer{
				ServerCert: CertKeyPair{CertData: "PEM ME\n", KeyData: "PEM ME\n"},
			},
		},
		Ledger: Ledger{
			DataDir: "relative/path",
		},
	}))

	err = config.ApplyDefaults()
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(config).To(Equal(Batik{
		Server: Server{
			ListenAddress: "127.0.0.1:7879",
			GRPC: GRPCServer{
				ConnTimeout: 30 * time.Second,
				GRPC: GRPC{
					MaxRecvMessageSize: 104857600,
					MaxSendMessageSize: 104857600,
				},
			},
			TLS: TLSServer{
				ServerCert: CertKeyPair{
					CertData: "PEM ME\n",
					KeyData:  "PEM ME\n",
				},
			},
		},
		Ledger: Ledger{
			DataDir: "relative/path",
		},
	}))

}
