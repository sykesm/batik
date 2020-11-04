// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package conf

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// File returns the first config file candidate that exists. If no candidate
// file exists, an empty path is returned with no error.
func File(stem string) (string, error) {
	candidates, err := candidateFiles(stem)
	if err != nil {
		return "", errors.Wrap(err, "unable to determine candidate configuration files")
	}
	for _, c := range candidates {
		if fileExists(c) {
			return c, nil
		}
	}
	return "", nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

// candidateFiles returns an ordered list of paths to candidate configuration
// files. The list contains the file references from the following directories:
//   - current working directory
//   - os.UserConfigDir()
//   - os.UserHomeDir() + /.config (if not already in the list)
//
// The provided stem is used as the XDG directory name and the filename
// stem. Candidates are returned for .yaml and .yml file suffixes.
func candidateFiles(stem string) ([]string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	confDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	var paths []string
	paths = append(paths, filepath.Join(cwd, stem+".yml"))
	paths = append(paths, filepath.Join(cwd, stem+".yaml"))
	paths = append(paths, filepath.Join(confDir, stem, stem+".yml"))
	paths = append(paths, filepath.Join(confDir, stem, stem+".yaml"))
	if filepath.Clean(confDir) != filepath.Join(homeDir, ".config") {
		paths = append(paths, filepath.Join(homeDir, ".config", stem, stem+".yml"))
		paths = append(paths, filepath.Join(homeDir, ".config", stem, stem+".yaml"))
	}
	return paths, nil
}
