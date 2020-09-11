// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// Load populates configuration by searching configuration sources and assigns the
// discovered values into the object referenced by out.
//
// The configuration sources are examined in the following order:
//   1. batik.yaml config file
//   2. environment variable overrides
//
// If a path is specified but not found, it will return an error.
// If a path is not specified, it will attempt to load the configuration yaml from
// one of the following sources in order should one exist at that location:
//   1. $(pwd)/batik.yaml
//   2. os specific $XDG_CONFIG_HOME/batik/batik.yaml
//   3. $HOME/.config/batik/batik.yaml
//
// If a config file still does not already exist at any of the above paths, configuration
// parameters will need to be passed via command line flags or environment variables.
func Load(cfgPath string, out interface{}) error {
	if cfgPath == "" {
		paths, err := SearchPath("batik")
		if err != nil {
			return err
		}

		for _, p := range paths {
			path := filepath.Join(p, "batik.yaml")
			_, err := os.Stat(path)
			if err == nil {
				cfgPath = path
				break
			}
		}
	}

	if cfgPath != "" {
		if err := readFile(cfgPath, out); err != nil {
			return errors.Wrap(err, "read file")
		}
	}

	return nil
}

func readFile(cfgPath string, cfg interface{}) error {
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
