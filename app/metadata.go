// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"github.com/urfave/cli/v2"
)

type MetadataKey string

const (
	CONFIG MetadataKey = "config"
)

// GetConfig retrieves a Config object from the app Metadata
func GetConfig(c *cli.Context) Config {
	config, ok := getMetadata(c, CONFIG).(Config)
	if !ok {
		return Config{}
	}

	return config
}

// SetConfig stores a Config object on the app Metadata
func SetConfig(c *cli.Context, config Config) {
	setMetadata(c, CONFIG, config)
}

func getMetadata(c *cli.Context, key MetadataKey) interface{} {
	return c.App.Metadata[string(key)]
}

func setMetadata(c *cli.Context, key MetadataKey, val interface{}) {
	c.App.Metadata[string(key)] = val
}
