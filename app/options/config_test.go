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

func TestConfigDefaults(t *testing.T) {
	gt := NewGomegaWithT(t)
	config := ConfigDefaults()
	gt.Expect(config).To(Equal(&Config{
		Server: *ServerDefaults(),
		Ledger: *LedgerDefaults(),
	}))
}

func TestConfigApplyDefaults(t *testing.T) {
	tests := map[string]struct {
		setup    func(*Config)
		matchErr types.GomegaMatcher
	}{
		"empty":  {setup: func(c *Config) { *c = Config{} }, matchErr: BeNil()},
		"server": {setup: func(c *Config) { c.Server = Server{} }, matchErr: BeNil()},
		"ledger": {setup: func(c *Config) { c.Ledger = Ledger{} }, matchErr: BeNil()},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			input := ConfigDefaults()
			tt.setup(input)

			err := input.ApplyDefaults()
			gt.Expect(err).To(tt.matchErr)
			if err != nil {
				return
			}
			gt.Expect(input).To(Equal(ConfigDefaults()))
		})
	}
}

func TestReadConfigApplyDefaults(t *testing.T) {
	gt := NewGomegaWithT(t)

	cf, err := os.Open("testdata/config.yaml")
	gt.Expect(err).NotTo(HaveOccurred())
	defer cf.Close()

	var config Config
	decoder := yaml.NewDecoder(cf)

	err = decoder.Decode(&config)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(config).To(Equal(Config{
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
	gt.Expect(config).To(Equal(Config{
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
