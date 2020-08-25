// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"os"
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
	Address string `yaml:"address" example:"127.0.0.1:9053" env:"BATIK_ADDRESS"`
}

type fakeLookuper struct{}

func (f fakeLookuper) Lookup(key string) (string, error) {
	return "", nil
}

func TestLoad(t *testing.T) {
	tests := []struct {
		testName       string
		cfgPath        string
		envMap         EnvMap
		expectedConfig BatikConfig
	}{
		{
			testName: "load yaml from file",
			cfgPath:  filepath.Join("testdata", "batik-config.yaml"),
			envMap:   EnvMap{},
			expectedConfig: BatikConfig{
				Server: Server{
					Address: "127.0.0.1:9000",
				},
			},
		},
		{
			testName: "load yaml from $HOME/.config/batik.yaml",
			cfgPath:  "",
			envMap: EnvMap{
				"HOME": filepath.Join("testdata", "home"),
			},
			expectedConfig: BatikConfig{
				Server: Server{
					Address: "127.0.0.1:9002",
				},
			},
		},
		{
			testName: "load yaml from $XDG_CONFIG_HOME/batik/batik.yaml",
			cfgPath:  "",
			envMap: EnvMap{
				"XDG_CONFIG_HOME": filepath.Join("testdata", "xdg_home"),
			},
			expectedConfig: BatikConfig{
				Server: Server{
					Address: "127.0.0.1:9003",
				},
			},
		},
		{
			testName: "load env vars override",
			cfgPath:  "",
			envMap: EnvMap{
				"HOME":          filepath.Join("testdata", "home"),
				"BATIK_ADDRESS": "127.0.0.1:9004",
			},
			expectedConfig: BatikConfig{
				Server: Server{
					Address: "127.0.0.1:9004",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			var batikConfig BatikConfig
			err := Load(tt.cfgPath, tt.envMap, &batikConfig)
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(batikConfig).To(Equal(tt.expectedConfig))
		})
	}

	t.Run("load yaml from cwd", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		tempFile, err := os.Create("batik.yaml")
		gt.Expect(err).NotTo(HaveOccurred())

		_, err = tempFile.WriteString(`server:
  address: 127.0.0.1:9001
`)
		gt.Expect(err).NotTo(HaveOccurred())
		defer func() {
			tempFile.Close()
			os.Remove(tempFile.Name())
		}()

		expectedConfig := BatikConfig{
			Server: Server{
				Address: "127.0.0.1:9001",
			},
		}

		var batikConfig BatikConfig
		err = Load("", EnvMap{}, &batikConfig)
		gt.Expect(err).NotTo(HaveOccurred())
		gt.Expect(batikConfig).To(Equal(expectedConfig))
	})
}

func TestLoadFailures(t *testing.T) {
	tests := []struct {
		testName    string
		cfgPath     string
		envMap      EnvMap
		expectedErr string
	}{
		{
			testName: "nonexistent dir",
			cfgPath:  filepath.Join("dne", "batik.yaml"),
			envMap: EnvMap{
				"HOME": filepath.Join("testdata", "home"),
			},
			expectedErr: "read file: open dne/batik.yaml: no such file or directory",
		},
		{
			testName: "invalid yaml",
			cfgPath:  filepath.Join("testdata", "invalid.yaml"),
			envMap: EnvMap{
				"HOME": filepath.Join("testdata", "home"),
			},
			expectedErr: "read file: yaml: unmarshal errors:\n  line 1: cannot unmarshal !!seq into config.BatikConfig",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			var batikConfig BatikConfig
			err := Load(tt.cfgPath, tt.envMap, &batikConfig)
			gt.Expect(err).To(MatchError(tt.expectedErr))
		})
	}

	t.Run("unsupported lookuper", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		var batikConfig BatikConfig
		err := Load("", fakeLookuper{}, &batikConfig)

		gt.Expect(err).To(MatchError("unsupported lookuper of type: config.fakeLookuper"))
	})
}
