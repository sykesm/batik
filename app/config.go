// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

// Config contains the configuration properties for a Batik instance.
type Config struct {
	// Server contains the batik grpc server configuration properties.
	Server Server `yaml:"server"`
}

// Server contains configuration properties for a Batik gRPC server.
type Server struct {
	// Address configures the listen address for the gRPC server.
	Address string `yaml:"address" example:"127.0.0.1:9053" env:"BATIK_ADDRESS"`
}
