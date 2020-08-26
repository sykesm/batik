// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

// Config contains the configuration properties for a Batik instance.
type Config struct {
	// Server contains the batik grpc server configuration properties.
	Server Server `yaml:"server"`

	// DBPath configures the level db instance filepath. If empty, will default to
	// in memory storage.
	DBPath string `yaml:"dbpath" default:"" env:"DB_PATH"`
}

// Server contains configuration properties for a Batik gRPC server.
type Server struct {
	// Address configures the listen address for the gRPC server.
	Address string `yaml:"address" default:"127.0.0.1:9053" env:"BATIK_ADDRESS"`
}
