// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

// Batik exposes the configurable elements of the application.
type Batik struct {
	Server     Server      `yaml:"server,omitempty"`
	Namespaces []Namespace `yaml:"namespaces,omitempty"`
	Logging    Logging     `yaml:"logging,omitempty"`
}

// BatikDefaults returns the default configuration values for the app.
func BatikDefaults() *Batik {
	return &Batik{
		Server:  *ServerDefaults(),
		Logging: *LoggingDefaults(),
	}
}

// ApplyDefaults applies default values for missing configuration fields.
func (c *Batik) ApplyDefaults() {
	c.Server.ApplyDefaults()
	for i := range c.Namespaces {
		(&c.Namespaces[i]).ApplyDefaults()
	}
	c.Logging.ApplyDefaults()
}
