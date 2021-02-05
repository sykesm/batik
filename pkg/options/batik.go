// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	cli "github.com/urfave/cli/v2"
)

// Batik exposes the configurable elements of the application.
type Batik struct {
	DataDir    string      `yaml:"data_dir,omitempty" batik:"relpath"`
	Server     Server      `yaml:"server,omitempty"`
	Namespaces []Namespace `yaml:"namespaces,omitempty"`
	Validators []Validator `yaml:"validators,omitempty"`
	Logging    Logging     `yaml:"logging,omitempty"`
}

// BatikDefaults returns the default configuration values for the app.
func BatikDefaults() *Batik {
	return &Batik{
		DataDir: "data",
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

	if c.DataDir == "" {
		c.DataDir = defaults.DataDir
	}

	c.Server.ApplyDefaults()

	for i := range c.Namespaces {
		(&c.Namespaces[i]).ApplyDefaults(c.DataDir)
	}

	if len(c.Validators) == 0 {
		c.Validators = defaults.Validators
	} else {
		for i := range c.Validators {
			(&c.Validators[i]).ApplyDefaults(c.DataDir)
		}
	}

	c.Logging.ApplyDefaults()
}

// Flags exposes configuration fields as flags. The current value of the
// receiver is used as the default value of the flag so a ApplyDefaults should
// be called before requesting flags.
func (c *Batik) Flags() []cli.Flag {
	def := BatikDefaults()
	return []cli.Flag{
		NewStringFlag(&cli.StringFlag{
			Name:        "data-dir",
			Value:       c.DataDir,
			Destination: &c.DataDir,
			Usage: flow(`Sets the base data directory where persisted data is stored.  Other component
					default paths such as for namespcaes and validators are created relative
					to this directory.`),
			DefaultText: def.DataDir,
		}),
	}
}
