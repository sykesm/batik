// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
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

// NewBatikConfig returns a new Batik configuration based on reading configuration
// from the following sources in order:
//   1. batik.yaml config file
//   2. environment variable overrides
// If a path is specified but not found, it will return an error.
// If a path is not specified, it will attempt to load the configuration yaml from
// one of the following sources in order should one exist at that location:
//   1. $(pwd)/batik.yaml
//   2. os specific $XDG_CONFIG_HOME/batik/batik.yaml
//   3. $HOME/.config/batik/batik.yaml
// If a config file still does not already exist at any of the above paths, configuration
// parameters will need to be passed via command line flags or environment variables.
func NewBatikConfig(cfgPath string, m ...EnvMap) (BatikConfig, error) {
	batikConfig := BatikConfig{}

	if cfgPath == "" {
		// Config paths to check in order if they exist:
		//	1. cwd
		//	2. $XDG_CONFIG_HOME/batik/
		//	3. $HOME/.config/batik/
		cfgPaths := []string{
			".",
		}

		switch len(m) {
		case 0:
			if usrCfgDir, err := os.UserConfigDir(); err == nil {
				cfgPaths = append(cfgPaths, filepath.Join(usrCfgDir, "batik"))
			}

			if usrHomeDir, err := os.UserHomeDir(); err == nil {
				cfgPaths = append(cfgPaths, filepath.Join(usrHomeDir, ".config", "batik"))
			}
		case 1:
			if usrCfgDir, err := m[0].Getenv("XDG_CONFIG_HOME"); err == nil {
				cfgPaths = append(cfgPaths, filepath.Join(usrCfgDir, "batik"))
			}

			if usrHomeDir, err := m[0].Getenv("HOME"); err == nil {
				cfgPaths = append(cfgPaths, filepath.Join(usrHomeDir, ".config", "batik"))
			}
		default:
			return BatikConfig{}, errors.New("expected at most 1 optional EnvMap")
		}

		for _, p := range cfgPaths {
			path := filepath.Join(p, "batik.yaml")
			// fmt.Printf("Checking if config exists at %s\n", path)
			_, err := os.Stat(path)
			if err == nil {
				// fmt.Printf("Found config at %s\n", path)
				cfgPath = path
				break
			}
			// fmt.Printf("No config found at %s: %s\n", path, err)
		}
	}

	if cfgPath != "" {
		if err := readFile(&batikConfig, cfgPath); err != nil {
			return BatikConfig{}, fmt.Errorf("read file: %s", err)
		}
	}

	// TODO: Parse env and post-decode defaults
	// if err := readEnv(&batikConfig); err != nil {
	// 	return BatikConfig{}, fmt.Errorf("read env: %s", err)
	// }

	return batikConfig, nil
}

func readFile(cfg *BatikConfig, cfgPath string) error {
	f, err := os.Open(cfgPath)
	if err != nil {
		return err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(cfg); err != nil {
		return err
	}

	return nil
}

// func readEnv(cfg *BatikConfig) error {
// 	if err := env.Parse(cfg); err != nil {
// 		return err
// 	}
//
// 	return nil
// }
