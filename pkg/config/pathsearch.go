// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"os"
	"path/filepath"
	"runtime"
)

// SearchPath returns an ordered list of directories to search for a
// configuration file. This list is constructed from the following:
//   - current working directory
//   - os.UserConfigDir()
//   - $XDG_CONFIG_HOME || $HOME/.config (on Darwin)
func SearchPath(stem string) ([]string, error) {
	var paths []string

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	paths = append(paths, cwd)

	confDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	paths = append(paths, filepath.Join(confDir, stem))

	if runtime.GOOS == "darwin" {
		if home, ok := os.LookupEnv("HOME"); ok {
			paths = append(paths, filepath.Join(home, ".config", stem))
		}
	}

	return paths, nil
}
