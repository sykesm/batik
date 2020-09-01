// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"time"

	"github.com/urfave/cli/v2"
)

// GRPCOptions exposes configurationf or gRPC servers and clients.
type GRPCOptions struct {
	MaxRecvMessageSize uint
	MaxSendMessageSize uint
}

func DefaultGRPCOptions() *GRPCOptions {
	return &GRPCOptions{
		MaxRecvMessageSize: 100 * 1024 * 1024,
		MaxSendMessageSize: 100 * 1024 * 1024,
	}
}

func (g *GRPCOptions) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.UintFlag{
			Name:        "max-recv-message-size",
			Value:       g.MaxRecvMessageSize,
			Destination: &g.MaxRecvMessageSize,
			Usage:       "FIXME: max-recv-message-size",
		},
		&cli.UintFlag{
			Name:        "max-send-message-size",
			Value:       g.MaxSendMessageSize,
			Destination: &g.MaxSendMessageSize,
			Usage:       "FIXME: max-send-message-size",
		},
	}
}

// GRPCServerOptions exposes configuration for gRPC servers.
type GRPCServerOptions struct {
	GRPCOptions
	ConnectionTimeout time.Duration
}

func DefaultGRPCServerOptions() *GRPCServerOptions {
	return &GRPCServerOptions{
		GRPCOptions:       *DefaultGRPCOptions(),
		ConnectionTimeout: 5 * time.Second,
	}
}

func (g *GRPCServerOptions) Flags() []cli.Flag {
	flags := []cli.Flag{
		&cli.DurationFlag{
			Name:        "grpc-conn-timeout",
			Value:       g.ConnectionTimeout,
			Destination: &g.ConnectionTimeout,
			Usage:       "FIXME: connection-timeout",
		},
	}
	return append(flags, g.GRPCOptions.Flags()...)
}
