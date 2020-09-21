// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	cli "github.com/urfave/cli/v2"
)

// Logging exposes configuration for logging.
type Logging struct {
	// LogSpec defines the level to log at.
	LogSpec string `yaml:"log_spec,omitempty"`
	// Color can be either "yes", "no", or "auto" and defines different modes for
	// configuring colored log output.
	Color string `yaml:"color,omitempty"`
	// Format can be either "logfmt" or "json" and defines the encoding format for
	// non-colorized log output.
	Format string `yaml:"format,omitempty"`
}

// LoggingDefaults returns the default configuration values for logging.
func LoggingDefaults() *Logging {
	return &Logging{
		LogSpec: "info",
		Color:   "auto",
		Format:  "logfmt",
	}
}

// ApplyDefaults applies default values for missing configuration fields.
func (l *Logging) ApplyDefaults() {
	defaults := LoggingDefaults()
	if l.LogSpec == "" {
		l.LogSpec = defaults.LogSpec
	}
	if l.Color == "" {
		l.Color = defaults.Color
	}
	if l.Format == "" {
		l.Format = defaults.Format
	}
}

// Flags exposes configuration fields as flags. The current value of the
// receiver is used as the default value of the flag so a ApplyDefaults should
// be called before requesting flags.
func (l *Logging) Flags() []cli.Flag {
	def := LoggingDefaults()
	return []cli.Flag{
		NewStringFlag(&cli.StringFlag{
			Name:        "log-spec",
			Value:       l.LogSpec,
			Destination: &l.LogSpec,
			Usage:       "FIXME: log spec",
			DefaultText: def.LogSpec,
		}),
		NewStringFlag(&cli.StringFlag{
			Name:        "color",
			Value:       l.Color,
			Destination: &l.Color,
			Usage:       "FIXME: color",
			DefaultText: def.Color,
		}),
		NewStringFlag(&cli.StringFlag{
			Name:        "format",
			Value:       l.Format,
			Destination: &l.Format,
			Usage:       "FIXME: format",
			DefaultText: def.Format,
		}),
	}
}
