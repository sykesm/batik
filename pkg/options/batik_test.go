// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"flag"
	"os"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	yaml "gopkg.in/yaml.v3"
)

func TestBatikDefaults(t *testing.T) {
	gt := NewGomegaWithT(t)
	config := BatikDefaults()
	gt.Expect(config).To(Equal(&Batik{
		DataDir: "data",
		Server:  *ServerDefaults(),
		Validators: []Validator{
			{
				Name: "signature-builtin",
				Type: "builtin",
			},
		},
		Logging: *LoggingDefaults(),
	}))
}

func TestBatikFlagName(t *testing.T) {
	gt := NewGomegaWithT(t)
	flags := (&Batik{}).Flags()

	var names []string
	for _, f := range flags {
		names = append(names, f.Names()...)
	}

	gt.Expect(flags).To(HaveLen(1))
	gt.Expect(names).To(ConsistOf(
		"data-dir",
	))
}

func TestBatikFlags(t *testing.T) {
	tests := map[string]struct {
		args     []string
		expected Batik
	}{
		"no flags": {
			args:     []string{},
			expected: Batik{},
		},
		"data dir": {
			args:     []string{"--data-dir=some/path/name"},
			expected: Batik{DataDir: "some/path/name"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			batik := &Batik{}
			flagSet := flag.NewFlagSet("batik-flag-test", flag.ContinueOnError)
			for _, f := range batik.Flags() {
				err := f.Apply(flagSet)
				gt.Expect(err).NotTo(HaveOccurred())
			}

			err := flagSet.Parse(tt.args)
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(batik).To(Equal(&tt.expected))
		})
	}
}

func TestBatikFlagsDefaultText(t *testing.T) {
	flags := BatikDefaults().Flags()
	assertWrappedFlagWithDefaultText(t, flags...)
}

func TestBatikApplyDefaults(t *testing.T) {
	tests := map[string]struct {
		setup func(*Batik)
	}{
		"empty":   {setup: func(c *Batik) { *c = Batik{} }},
		"server":  {setup: func(c *Batik) { c.Server = Server{} }},
		"logging": {setup: func(c *Batik) { c.Logging = Logging{} }},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			input := BatikDefaults()
			tt.setup(input)

			input.ApplyDefaults()
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
		DataDir: "relative/path",
		Server: Server{
			GRPC: GRPCServer{
				GRPC:          GRPC{MaxSendMessageSize: 104857600},
				ListenAddress: "127.0.0.1:7878",
			},
			HTTP: HTTPServer{
				ListenAddress: "127.0.0.1:7879",
			},
			TLS: ServerTLS{
				ServerCert: CertKeyPair{CertData: "PEM ME\n", KeyData: "PEM ME\n"},
				CertsDir:   "relative/certs-dir-path",
			},
		},
		Namespaces: []Namespace{
			{
				Name:    "ns1",
				DataDir: "override/path",
			},
			{
				Name:      "ns2",
				Validator: "wasm-validator1",
			},
		},
		Validators: []Validator{
			{
				Name: "builtin-validator",
				Type: "builtin",
			},
			{
				Name: "wasm-validator1",
				Type: "wasm",
			},
			{
				Name: "wasm-validator2",
				Path: "custom/relative/path",
			},
		},
		Logging: Logging{
			LogSpec: "debug",
		},
	}))

	config.ApplyDefaults()
	gt.Expect(config).To(Equal(Batik{
		DataDir: "relative/path",
		Server: Server{
			GRPC: GRPCServer{
				ConnTimeout: 30 * time.Second,
				GRPC: GRPC{
					MaxRecvMessageSize: 104857600,
					MaxSendMessageSize: 104857600,
				},
				ListenAddress: "127.0.0.1:7878",
			},
			HTTP: HTTPServer{
				ListenAddress:     "127.0.0.1:7879",
				ReadHeaderTimeout: 30 * time.Second,
			},
			TLS: ServerTLS{
				ServerCert: CertKeyPair{
					CertData: "PEM ME\n",
					KeyData:  "PEM ME\n",
				},
				CertsDir: "relative/certs-dir-path",
			},
		},
		Namespaces: []Namespace{
			{
				Name:      "ns1",
				DataDir:   "override/path",
				Validator: "signature-builtin",
			},
			{
				Name:      "ns2",
				DataDir:   "relative/path/namespaces/ns2",
				Validator: "wasm-validator1",
			},
		},
		Validators: []Validator{
			{
				Name: "builtin-validator",
				Type: "builtin",
			},
			{
				Name: "wasm-validator1",
				Type: "wasm",
				Path: "relative/path/validators/wasm-validator1.wasm",
			},
			{
				Name: "wasm-validator2",
				Type: "wasm",
				Path: "custom/relative/path",
			},
		},
		Logging: Logging{
			LogSpec: "debug",
			Color:   "auto",
			Format:  "logfmt",
		},
	}))
}
