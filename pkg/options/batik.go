// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

// Batik exposes the configurable elements of the application.
type Batik struct {
	Server     Server      `yaml:"server,omitempty"`
	Namespaces []Namespace `yaml:"namespaces,omitempty"`
	Validators []Validator `yaml:"validators,omitempty"`
	Logging    Logging     `yaml:"logging,omitempty"`
}

// BatikDefaults returns the default configuration values for the app.
func BatikDefaults() *Batik {
	return &Batik{
		Server:  *ServerDefaults(),
		Logging: *LoggingDefaults(),
		Validators: []Validator{
			{
				Name: "signature-builtin",
				Type: "builtin",
			},
		},
	}
}

// ApplyDefaults applies default values for missing configuration fields.
func (c *Batik) ApplyDefaults() {
	defaults := BatikDefaults()
	c.Server.ApplyDefaults()
	for i := range c.Namespaces {
		(&c.Namespaces[i]).ApplyDefaults()
	}
	if len(c.Validators) == 0 {
		c.Validators = defaults.Validators
	} else {
		for i := range c.Validators {
			(&c.Validators[i]).ApplyDefaults()
		}
	}
	c.Logging.ApplyDefaults()
}
