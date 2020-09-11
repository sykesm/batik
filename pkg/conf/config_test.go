// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package conf

import (
	"path/filepath"
	"testing"

	. "github.com/onsi/gomega"
)

// BatikConfig contains the configuration properties for a Batik instance.
type BatikConfig struct {
	// Server contains the batik grpc server configuration properties.
	Server Server `yaml:"server"`
}

// Server contains configuration properties for a Batik gRPC server.
type Server struct {
	// Address configures the listen address for the gRPC server.
	Address string `yaml:"address"`
}

func TestLoad(t *testing.T) {
	tests := []struct {
		testName       string
		cfgPath        string
		expectedConfig BatikConfig
	}{
		{
			testName: "load yaml from file",
			cfgPath:  filepath.Join("testdata", "batik-config.yaml"),
			expectedConfig: BatikConfig{
				Server: Server{
					Address: "127.0.0.1:9000",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			var batikConfig BatikConfig
			err := LoadFile(tt.cfgPath, &batikConfig)
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(batikConfig).To(Equal(tt.expectedConfig))
		})
	}
}

func TestLoadFailures(t *testing.T) {
	tests := []struct {
		testName    string
		cfgPath     string
		expectedErr string
	}{
		{
			testName:    "nonexistent dir",
			cfgPath:     filepath.Join("dne", "batik.yaml"),
			expectedErr: "read file: open dne/batik.yaml: no such file or directory",
		},
		{
			testName:    "invalid yaml",
			cfgPath:     filepath.Join("testdata", "invalid.yaml"),
			expectedErr: "read file: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!seq into conf.BatikConfig",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			var batikConfig BatikConfig
			err := LoadFile(tt.cfgPath, &batikConfig)
			gt.Expect(err).To(MatchError(tt.expectedErr))
		})
	}
}
