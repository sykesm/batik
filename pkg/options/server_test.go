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
		GRPC: *GRPCServerDefaults(),
		HTTP: *HTTPServerDefaults(),
		TLS:  *ServerTLSDefaults(),
	}))
}

func TestServerApplyDefaults(t *testing.T) {
	tests := map[string]struct {
		setup func(*Server)
	}{
		"empty": {setup: func(s *Server) { *s = Server{} }},
		"GRPC":  {setup: func(s *Server) { s.GRPC = GRPCServer{} }},
		"TLS":   {setup: func(s *Server) { s.TLS = ServerTLS{} }},
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

	gt.Expect(flags).To(HaveLen(11))
	gt.Expect(names).To(ConsistOf(
		"grpc-conn-timeout",
		"grpc-listen-address",
		"grpc-max-recv-message-size",
		"grpc-max-send-message-size",
		"http-listen-address",
		"http-read-timeout",
		"http-read-header-timeout",
		"http-write-timeout",
		"http-idle-timeout",
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
		"grpc max send": {
			args:     []string{"--grpc-max-send-message-size", "233"},
			expected: Server{GRPC: GRPCServer{GRPC: GRPC{MaxSendMessageSize: 233}}},
		},
		"tls cert file": {
			args:     []string{"--tls-cert-file", "file.crt"},
			expected: Server{TLS: ServerTLS{ServerCert: CertKeyPair{CertFile: "file.crt"}}},
		},
		"the works": {
			args: []string{
				"--grpc-conn-timeout", "90s",
				"--grpc-max-recv-message-size", "9999",
				"--grpc-max-send-message-size", "8888",
				"--tls-cert-file", "file.crt",
				"--tls-private-key-file", "private.key",
			},
			expected: Server{
				GRPC: GRPCServer{
					ConnTimeout: 90 * time.Second,
					GRPC:        GRPC{MaxRecvMessageSize: 9999, MaxSendMessageSize: 8888},
				},
				TLS: ServerTLS{
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

func TestServerFlagsDefaultText(t *testing.T) {
	flags := ServerDefaults().Flags()
	assertWrappedFlagWithDefaultText(t, flags...)
}
