// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	cli "github.com/urfave/cli/v2"
)

// Server exposes configuration for the server component.
type Server struct {
	// GRPC maintains the gRPC server configuration for a server.
	GRPC GRPCServer `yaml:"grpc,omitempty"`
	// HTTP maintains the HTTP server configuration for a server.
	HTTP HTTPServer `yaml:"http,omitempty"`
	// TLS references the TLS configuration for a server.
	TLS ServerTLS `yaml:"tls,omitempty"`
}

// ServerDefault returns the default configuration values for the server component.
func ServerDefaults() *Server {
	return &Server{
		GRPC: *GRPCServerDefaults(),
		HTTP: *HTTPServerDefaults(),
		TLS:  *ServerTLSDefaults(),
	}
}

// ApplyDefaults applies default values for missing configuration fields.
func (s *Server) ApplyDefaults() {
	s.GRPC.ApplyDefaults()
	s.HTTP.ApplyDefaults()
	s.TLS.ApplyDefaults()
}

// Flags exposes configuration fields as flags. The current value of the
// receiver is used as the default value of the flag so a ApplyDefaults should
// be called before requesting flags.
func (s *Server) Flags() []cli.Flag {
	var flags []cli.Flag
	flags = append(flags, s.GRPC.Flags()...)
	flags = append(flags, s.HTTP.Flags()...)
	flags = append(flags, s.TLS.Flags()...)
	return flags
}
