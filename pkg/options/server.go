// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	cli "github.com/urfave/cli/v2"
)

// Server exposes configuration for the server component.
type Server struct {
	// ListenAddress determines the addresses that the server will listen on. The
	// address should be in a form that is compatible with net.Listen from the Go
	// standard library.
	ListenAddress string `yaml:"listen_address,omitempty"`
	// GRPC maintains the gRPC server configuration for a server.
	GRPC GRPCServer `yaml:"grpc,omitempty"`
	// TLS references the TLS configuration for a server.
	TLS TLSServer `yaml:"tls,omitempty"`
}

// ServerDefault returns the default configuration values for the server component.
func ServerDefaults() *Server {
	return &Server{
		ListenAddress: ":9443",
		GRPC:          *GRPCServerDefaults(),
		TLS:           *TLSServerDefaults(),
	}
}

// ApplyDefaults applies default values for missing configuration fields.
func (s *Server) ApplyDefaults() {
	defaults := ServerDefaults()
	if s.ListenAddress == "" {
		s.ListenAddress = defaults.ListenAddress
	}
	s.GRPC.ApplyDefaults()
	s.TLS.ApplyDefaults()
}

// Flags exposes configuration fields as flags. The current value of the
// receiver is used as the default value of the flag so a ApplyDefaults should
// be called before requesting flags.
func (s *Server) Flags() []cli.Flag {
	def := ServerDefaults()
	flags := []cli.Flag{
		NewStringFlag(&cli.StringFlag{
			Name:        "listen-address",
			Value:       s.ListenAddress,
			Destination: &s.ListenAddress,
			EnvVars:     []string{"BATIK_LISTEN_ADDR"},
			Usage:       "FIXME: listen-address",
			DefaultText: def.ListenAddress,
		}),
	}

	flags = append(flags, s.GRPC.Flags()...)
	flags = append(flags, s.TLS.Flags()...)
	return flags
}
