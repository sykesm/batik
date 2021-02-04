// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package conf

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
)

// BatikConfig contains the configuration properties for a Batik instance.
type BatikConfig struct {
	// defApplied is set if defaults have been applied
	defApplied bool
	// Server contains the batik grpc server configuration properties.
	Server Server `yaml:"server"`
}

func (b *BatikConfig) ApplyDefaults() { b.defApplied = true }

// Server contains configuration properties for a Batik gRPC server.
type Server struct {
	// Address configures the listen address for the gRPC server.
	Address string `yaml:"address"`
	Relpath string `batik:"relpath"`
}

func TestLoadFile(t *testing.T) {
	tests := []struct {
		name           string
		cfgPath        string
		expectedConfig BatikConfig
	}{
		{
			name:    "load yaml from file",
			cfgPath: filepath.Join("testdata", "batik-config.yaml"),
			expectedConfig: BatikConfig{
				defApplied: true,
				Server: Server{
					Address: "127.0.0.1:9000",
					Relpath: filepath.Join("testdata", "relative.txt"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			var batikConfig BatikConfig
			err := LoadFile(tt.cfgPath, &batikConfig)
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(batikConfig).To(Equal(tt.expectedConfig))
		})
	}
}

func TestLoad(t *testing.T) {
	gt := NewGomegaWithT(t)

	contents, err := ioutil.ReadFile(filepath.Join("testdata", "batik-config.yaml"))
	gt.Expect(err).NotTo(HaveOccurred())

	var batikConfig BatikConfig
	err = Load(bytes.NewBuffer(contents), &batikConfig)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(batikConfig).To(Equal(BatikConfig{
		defApplied: true,
		Server: Server{
			Address: "127.0.0.1:9000",
			Relpath: filepath.Join("relative.txt"),
		},
	}))
}

func TestLoadFileFailures(t *testing.T) {
	tests := []struct {
		name        string
		cfgPath     string
		expectedErr string
	}{
		{
			name:        "nonexistent dir",
			cfgPath:     filepath.Join("dne", "batik.yaml"),
			expectedErr: "conf: open dne/batik.yaml: no such file or directory",
		},
		{
			name:        "invalid yaml",
			cfgPath:     filepath.Join("testdata", "invalid.yaml"),
			expectedErr: "yaml: unmarshal errors:\n  line 1: cannot unmarshal !!seq into conf.BatikConfig",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			var batikConfig BatikConfig
			err := LoadFile(tt.cfgPath, &batikConfig)
			gt.Expect(err).To(MatchError(tt.expectedErr))
		})
	}
}

func TestLoadNonApplier(t *testing.T) {
	gt := NewGomegaWithT(t)

	data := map[string]interface{}{}
	err := Load(strings.NewReader("---\n{key: value}\n"), &data)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(data).To(Equal(map[string]interface{}{
		"key": "value",
	}))
}
