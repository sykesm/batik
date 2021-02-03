// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	cli "github.com/urfave/cli/v2"
)

// Namespace exposes configuration for a namespace.
type Namespace struct {
	// Name is the human readable name for this namespace and how external
	// transactors reference
	Name string `yaml:"name"`

	// DataDir is the path to the directory where ledger artifacts will be
	// placed. If a relative path is used, the path is relative to the working
	// directory unless the path was retrieved from a configuration file, in
	// which case it is relative to the directory containing the configuration
	// file.  Note, the namespace name is always concatenated with the DataDir
	// to avoid collisions, so multiple namespaces may safely specify the same
	// DataDir.
	DataDir string `yaml:"data_dir,omitempty" batik:"relpath"`

	// Validator is the name of the validator used to validate transactions
	// in this namespace.  It must be defined in the top level Validators
	// section of the Batik configuration.
	Validator string `yaml:"validator,omitempty"`
}

// NamespaceDefaults returns the default configuration values for the ledger.
func NamespaceDefaults() *Namespace {
	return &Namespace{
		DataDir:   "data",
		Validator: "signature-builtin",
	}
}

// ApplyDefaults applies default values for missing configuration fields.
func (n *Namespace) ApplyDefaults() {
	defaults := NamespaceDefaults()
	if n.DataDir == "" {
		n.DataDir = defaults.DataDir
	}
	if n.Validator == "" {
		n.Validator = defaults.Validator
	}
}

// Flags exposes configuration fields as flags. The current value of the
// receiver is used as the default value of the flag so a ApplyDefaults should
// be called before requesting flags.
func (n *Namespace) Flags() []cli.Flag {
	def := NamespaceDefaults()
	return []cli.Flag{
		NewStringFlag(&cli.StringFlag{
			Name:        "data-dir",
			Value:       n.DataDir,
			Destination: &n.DataDir,
			Usage: flow(`Sets the data directory for all namespaces in the configuration.  Because
					the data directory is concatenated with the namespace name when creating
					artifacts inside the data directory, this value is safe to set with multiple
					namespaces without fear of collision.`),
			DefaultText: def.DataDir,
		}),
	}
}
