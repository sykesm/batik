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
func Load(cfgPath string, l Lookuper, out interface{}) error {
	if l == nil {
		return errors.New("empty lookuper")
	}

	if cfgPath == "" {
		// Config paths to check in order if they exist:
		//	1. cwd
		//	2. $XDG_CONFIG_HOME/batik/
		//	3. $HOME/.config/batik/
		cfgPaths := []string{
			".",
		}

		// switch l.(type) {
		// case OsEnv:
		// 	if usrCfgDir, err := os.UserConfigDir(); err == nil {
		// 		cfgPaths = append(cfgPaths, filepath.Join(usrCfgDir, "batik"))
		// 	}

		// 	if usrHomeDir, err := os.UserHomeDir(); err == nil {
		// 		cfgPaths = append(cfgPaths, filepath.Join(usrHomeDir, ".config", "batik"))
		// 	}
		// case EnvMap:
		// 	if usrCfgDir, err := l.Lookup("XDG_CONFIG_HOME"); err == nil {
		// 		cfgPaths = append(cfgPaths, filepath.Join(usrCfgDir, "batik"))
		// 	}

		// 	if usrHomeDir, err := l.Lookup("HOME"); err == nil {
		// 		cfgPaths = append(cfgPaths, filepath.Join(usrHomeDir, ".config", "batik"))
		// 	}
		// default:
		// 	return fmt.Errorf("unsupported lookuper of type: %T", l)
		// }

		for _, p := range cfgPaths {
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
			return fmt.Errorf("read file: %s", err)
		}
	}

	d := Decoder{
		lookuper:   l,
		defaultTag: "example",
		parseTag:   "env",
	}

	if err := d.Parse(out); err != nil {
		return fmt.Errorf("decode: %+w", err)
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
