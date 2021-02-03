// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	cli "github.com/urfave/cli/v2"
)

// Validator exposes configuration for a transaction validator.
type Validator struct {
	// Name is the human readable name for this validator and how external
	// other elements in the configuration may reference it.  Note, a derived
	// identifier will be used internally within namespaces to ensure consistent
	// validation across peers.
	Name string `yaml:"name"`

	// Type is the type of validator and may be one of "builtin" or "wasm".
	// If omitted, defaults to "wasm".
	Type string `yaml:"type"`

	// CodeDir is the path to the directory where WASM binaries are stored.
	// For validators of type "builtin", this field is ignored.
	// For validators of type "wasm", compiled code  must exist at <code_dir>/<name>.wasm
	// If a relative path is used, the path is relative to the working
	// directory unless the path was retrieved from a configuration file, in
	// which case it is relative to the directory containing the configuration
	// file.
	CodeDir string `yaml:"code_dir,omitempty" batik:"relpath"`
}

// ValidatorDefaults returns the default configuration values for the ledger.
func ValidatorDefaults() *Validator {
	return &Validator{
		Type:    "wasm",
		CodeDir: "validators",
	}
}

// ApplyDefaults applies default values for missing configuration fields.
func (n *Validator) ApplyDefaults() {
	defaults := ValidatorDefaults()
	if n.Type == "" {
		n.Type = defaults.Type
	}
	if n.Type == "wasm" && n.CodeDir == "" {
		n.CodeDir = defaults.CodeDir
	}
}

// Flags exposes configuration fields as flags. The current value of the
// receiver is used as the default value of the flag so a ApplyDefaults should
// be called before requesting flags.
func (n *Validator) Flags() []cli.Flag {
	def := ValidatorDefaults()
	return []cli.Flag{
		NewStringFlag(&cli.StringFlag{
			Name:        "validators-dir",
			Value:       n.CodeDir,
			Destination: &n.CodeDir,
			Usage: flow(`Sets the code directory for all validators in the configuration.  Because
					the code directory is concatenated with the validator name when locating
					wasm artifacts inside the code directory, this value is safe to set with multiple
					validators without fear of collision.`),
			DefaultText: def.CodeDir,
		}),
	}
}
