// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

type Config struct {
	Server Server `yaml:"server,omitempty"`
	Ledger Ledger `yaml:"ledger,omitempty"`
}

func ConfigDefaults() *Config {
	return &Config{
		Server: *ServerDefaults(),
		Ledger: *LedgerDefaults(),
	}
}

func (c *Config) ApplyDefaults() error {
	if err := c.Server.ApplyDefaults(); err != nil {
		return err
	}
	if err := c.Ledger.ApplyDefaults(); err != nil {
		return err
	}

	return nil
}
