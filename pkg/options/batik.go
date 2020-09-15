// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

// Batik exposes the conigurable elements of the application.
type Batik struct {
	Server  Server  `yaml:"server,omitempty"`
	Ledger  Ledger  `yaml:"ledger,omitempty"`
	Logging Logging `yaml:"logging,omitempty"`
}

// BatikDefaults returns the default configuration values for the app.
func BatikDefaults() *Batik {
	return &Batik{
		Server:  *ServerDefaults(),
		Ledger:  *LedgerDefaults(),
		Logging: *LoggingDefaults(),
	}
}

// ApplyDefaults applies default values for missing configuration fields.
func (c *Batik) ApplyDefaults() {
	c.Server.ApplyDefaults()
	c.Ledger.ApplyDefaults()
	c.Logging.ApplyDefaults()
}
