// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"path/filepath"
)

// Namespace exposes configuration for a namespace.
type TotalOrder struct {
	// Name is the human readable name for this total order and how other
	// parts of the configuration may reference it.
	Name string `yaml:"name"`

	// Type is the consensus type for this total order.  Depending on the
	// type, other configuration may be set.  Currently, the only type
	// is 'in-process', but other types including 'static-leader', and ultimately
	// more robust consensus types will be added.
	Type string `yaml:"type,omitempty"`

	// DataDir is the path where the db for this total order will be stored.  Note
	// the database will only be created if this peer is a consenter on order.
	DataDir string `yaml:"data_dir,omitempty" batik:"relpath"`
}

// ApplyDefaults applies default values for missing configuration fields.
func (n *TotalOrder) ApplyDefaults(baseDataDir string) {
	if n.Type == "" {
		n.Type = "in-process"
	}

	if n.DataDir == "" {
		n.DataDir = filepath.Join(baseDataDir, "totalorders", n.Name)
	}
}
