// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	cli "github.com/urfave/cli/v2"
)

// Ledger exposes configuration for the ledger.
type Ledger struct {
	// DataDir is the path to the directory where ledger artifacts will be
	// placed. If a relative path is used, the path is relative to the working
	// directory unless the path was retrieved from a configuration file, in
	// which case it is relative to the directory containing the configuration
	// file.
	DataDir string `yaml:"data_dir,omitempty" batik:"relpath"`
}

// LedgerDefaults returns the default configuration values for the ledger.
func LedgerDefaults() *Ledger {
	return &Ledger{
		DataDir: "data",
	}
}

// ApplyDefaults applies default values for missing configuration fields.
func (l *Ledger) ApplyDefaults() error {
	defaults := LedgerDefaults()
	if l.DataDir == "" {
		l.DataDir = defaults.DataDir
	}
	return nil
}

// Flags exposes configuration fields as flags. The current value of the
// receiver is used as the default value of the flag so a ApplyDefaults should
// be called before requesting flags.
func (l *Ledger) Flags(commandName string) []cli.Flag {
	return []cli.Flag{
		NewStringFlag(&cli.StringFlag{
			Name:        "data-dir",
			Value:       l.DataDir,
			Destination: &l.DataDir,
			Usage:       "FIXME: data directory",
		}),
	}
}
