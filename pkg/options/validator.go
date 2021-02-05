// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"fmt"
	"path/filepath"
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

	// Path is the location that the WASM binary are stored.  If not specified,
	// and the type is "wasm", it defaults to <data_dir>/validators/<validator_name>.wasm
	Path string `yaml:"path,omitempty" batik:"relpath"`
}

// ApplyDefaults applies default values for missing configuration fields.
func (n *Validator) ApplyDefaults(dataDir string) {
	if n.Type == "" {
		n.Type = "wasm"
	}
	if n.Type == "wasm" && n.Path == "" {
		n.Path = filepath.Join(dataDir, "validators", fmt.Sprintf("%s.wasm", n.Name))
	}
}
