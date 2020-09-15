// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"flag"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestServerDefaults(t *testing.T) {
	gt := NewGomegaWithT(t)
	server := ServerDefaults()
	gt.Expect(server).To(Equal(&Server{
		ListenAddress: ":9443",
		GRPC:          *GRPCServerDefaults(),
		TLS:           *TLSServerDefaults(),
	}))
}

func TestServerApplyDefaults(t *testing.T) {
	tests := map[string]struct {
		setup func(*Server)
	}{
		"empty":          {setup: func(s *Server) { *s = Server{} }},
		"listen address": {setup: func(s *Server) { s.ListenAddress = "" }},
		"GRPC":           {setup: func(s *Server) { s.GRPC = GRPCServer{} }},
		"TLS":            {setup: func(s *Server) { s.TLS = TLSServer{} }},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			input := ServerDefaults()
			tt.setup(input)

			input.ApplyDefaults()
			gt.Expect(input).To(Equal(ServerDefaults()))
		})
	}
}

func TestServerFlagNames(t *testing.T) {
	gt := NewGomegaWithT(t)
	server := &Server{}
	flags := server.Flags()

	var names []string
	for _, f := range flags {
		names = append(names, f.Names()...)
	}

	gt.Expect(flags).To(HaveLen(6))
	gt.Expect(names).To(ConsistOf(
		"grpc-conn-timeout",
		"grpc-max-recv-message-size",
		"grpc-max-send-message-size",
		"listen-address",
		"tls-cert-file",
		"tls-private-key-file",
	))
}

func TestServerFlags(t *testing.T) {
	tests := map[string]struct {
		args     []string
		expected Server
	}{
		"no flags": {
			args:     []string{},
			expected: Server{},
		},
		"listen address": {
			args:     []string{"--listen-address", ":1234"},
			expected: Server{ListenAddress: ":1234"},
		},
		"grpc max send": {
			args:     []string{"--grpc-max-send-message-size", "233"},
			expected: Server{GRPC: GRPCServer{GRPC: GRPC{MaxSendMessageSize: 233}}},
		},
		"tls cert file": {
			args:     []string{"--tls-cert-file", "file.crt"},
			expected: Server{TLS: TLSServer{ServerCert: CertKeyPair{CertFile: "file.crt"}}},
		},
		"the works": {
			args: []string{
				"--listen-address", ":5678",
				"--grpc-conn-timeout", "90s",
				"--grpc-max-recv-message-size", "9999",
				"--grpc-max-send-message-size", "8888",
				"--tls-cert-file", "file.crt",
				"--tls-private-key-file", "private.key",
			},
			expected: Server{
				ListenAddress: ":5678",
				GRPC: GRPCServer{
					ConnTimeout: 90 * time.Second,
					GRPC:        GRPC{MaxRecvMessageSize: 9999, MaxSendMessageSize: 8888},
				},
				TLS: TLSServer{
					ServerCert: CertKeyPair{CertFile: "file.crt", KeyFile: "private.key"},
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			server := &Server{}
			flagSet := flag.NewFlagSet("server-test", flag.ContinueOnError)
			for _, f := range server.Flags() {
				err := f.Apply(flagSet)
				gt.Expect(err).NotTo(HaveOccurred())
			}

			err := flagSet.Parse(tt.args)
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(server).To(Equal(&tt.expected))
		})
	}
}
