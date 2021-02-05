// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"path/filepath"
)

// Namespace exposes configuration for a namespace.
type Namespace struct {
	// Name is the human readable name for this namespace and how external
	// transactors reference
	Name string `yaml:"name"`

	// DataDir is the location of the ledger artifacts for this namespace.
	// If this field is not specified, its default is generated from
	// the BaseDir in the Namespaces configuration.
	DataDir string `yaml:"data_dir,omitempty" batik:"relpath"`

	// Validator is the name of the validator used to validate transactions
	// in this namespace.  It must be defined in the top level Validators
	// section of the Batik configuration.
	Validator string `yaml:"validator,omitempty"`
}

// ApplyDefaults applies default values for missing configuration fields.
func (n *Namespace) ApplyDefaults(baseDataDir string) {
	if n.DataDir == "" {
		n.DataDir = filepath.Join(baseDataDir, "namespaces", n.Name)
	}
	if n.Validator == "" {
		n.Validator = "signature-builtin"
	}
}
