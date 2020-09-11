// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

// Config exposes the conigurable elements of the application.
type Config struct {
	Server Server `yaml:"server,omitempty"`
	Ledger Ledger `yaml:"ledger,omitempty"`
}

// ConfigDefaults returns the default configuration values for the app.
func ConfigDefaults() *Config {
	return &Config{
		Server: *ServerDefaults(),
		Ledger: *LedgerDefaults(),
	}
}

// ApplyDefaults applies default values for missing configuration fields.
func (c *Config) ApplyDefaults() error {
	if err := c.Server.ApplyDefaults(); err != nil {
		return err
	}
	if err := c.Ledger.ApplyDefaults(); err != nil {
		return err
	}
	return nil
}
