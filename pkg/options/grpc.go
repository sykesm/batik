// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"strconv"
	"time"

	cli "github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

// GRPC exposes configuration for gRPC servers and clients.
type GRPC struct {
	// MaxRecvMessageSize limits the size of messages that can be received.
	MaxRecvMessageSize uint `yaml:"max_recv_message_size,omitempty"`
	// MaxSendMessageSize limits the size of messages that can be sent.
	MaxSendMessageSize uint `yaml:"max_send_message_size,omitempty"`
}

// GRPCDefaults returns the default configuration values for gRPC servers and
// clients.
func GRPCDefaults() *GRPC {
	return &GRPC{
		MaxRecvMessageSize: 100 * 1024 * 1024,
		MaxSendMessageSize: 100 * 1024 * 1024,
	}
}

// ApplyDefaults applies default values for missing configuration fields.
func (g *GRPC) ApplyDefaults() {
	defaults := GRPCDefaults()
	if g.MaxRecvMessageSize == 0 {
		g.MaxRecvMessageSize = defaults.MaxRecvMessageSize
	}
	if g.MaxSendMessageSize == 0 {
		g.MaxSendMessageSize = defaults.MaxSendMessageSize
	}
}

// Flags exposes configuration fields as flags. The current value of the
// receiver is used as the default value of the flag so a ApplyDefaults should
// be called before requesting flags.
func (g *GRPC) Flags() []cli.Flag {
	def := GRPCDefaults()
	return []cli.Flag{
		NewUintFlag(&cli.UintFlag{
			Name:        "grpc-max-recv-message-size",
			Value:       g.MaxRecvMessageSize,
			Destination: &g.MaxRecvMessageSize,
			Usage:       "FIXME: max-recv-message-size",
			DefaultText: strconv.Itoa(int(def.MaxRecvMessageSize)),
		}),
		NewUintFlag(&cli.UintFlag{
			Name:        "grpc-max-send-message-size",
			Value:       g.MaxSendMessageSize,
			Destination: &g.MaxSendMessageSize,
			Usage:       "FIXME: max-send-message-size",
			DefaultText: strconv.Itoa(int(def.MaxSendMessageSize)),
		}),
	}
}

// GRPCServer exposes configuration for gRPC servers.
type GRPCServer struct {
	GRPC `yaml:",inline,omitempty"`
	// ConnTimeout limits the time a server will wait for client connections to
	// be established.
	ConnTimeout time.Duration `yaml:"conn_timeout,omitempty"`
	// ListenAddress determines the addresses that the server will listen on. The
	// address should be in a form that is compatible with net.Listen from the Go
	// standard library.
	ListenAddress string `yaml:"listen_address,omitempty"`
}

// GRPCServerDefaults returns the default configuration values for gRPC
// servers.
func GRPCServerDefaults() *GRPCServer {
	return &GRPCServer{
		GRPC:          *GRPCDefaults(),
		ConnTimeout:   30 * time.Second,
		ListenAddress: ":9443",
	}
}

// ApplyDefaults applies default values for missing configuration fields.
func (g *GRPCServer) ApplyDefaults() {
	defaults := GRPCServerDefaults()
	g.GRPC.ApplyDefaults()
	if g.ConnTimeout == 0 {
		g.ConnTimeout = defaults.ConnTimeout
	}
	if g.ListenAddress == "" {
		g.ListenAddress = defaults.ListenAddress
	}
}

// Flags exposes configuration fields as flags. The current value of the
// receiver is used as the default value of the flag so ApplyDefaults should be
// called before requesting flags.
func (g *GRPCServer) Flags() []cli.Flag {
	def := GRPCServerDefaults()
	flags := []cli.Flag{
		NewDurationFlag(&cli.DurationFlag{
			Name:        "grpc-conn-timeout",
			Value:       g.ConnTimeout,
			Destination: &g.ConnTimeout,
			Usage:       "FIXME: connection-timeout",
			DefaultText: def.ConnTimeout.String(),
		}),
		NewStringFlag(&cli.StringFlag{
			Name:        "grpc-listen-address",
			Value:       g.ListenAddress,
			Destination: &g.ListenAddress,
			Usage:       "FIXME: grpc-listen-address",
			DefaultText: def.ListenAddress,
		}),
	}
	return append(flags, g.GRPC.Flags()...)
}

// BuildServerOptions returns an array of grpc.ServerOptions based on the configuration
// fields.
func (g *GRPCServer) BuildServerOptions() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.MaxRecvMsgSize(int(g.GRPC.MaxRecvMessageSize)),
		grpc.MaxSendMsgSize(int(g.GRPC.MaxSendMessageSize)),
		grpc.ConnectionTimeout(g.ConnTimeout),
	}
}
