// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"errors"

	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	"github.com/sykesm/batik/pkg/log"
)

type metadataKey string

const (
	configKey    metadataKey = "config"
	loggerKey                = "logger"
	errLoggerKey             = "errLogger"
	serverKey                = "server"
)

// GetConfig retrieves a Config object from the app Metadata.
func GetConfig(c *cli.Context) Config {
	config, ok := getMetadata(c, configKey).(Config)
	if !ok {
		return Config{}
	}

	return config
}

// SetConfig stores a Config object on the app Metadata.
func SetConfig(c *cli.Context, config Config) {
	setMetadata(c, configKey, config)
}

// GetLogger retrieves a logger from the app Metadata, if one
// does not exist it will return a new default logger.
func GetLogger(c *cli.Context) (*zerolog.Logger, error) {
	logger := getMetadata(c, loggerKey)
	if logger == nil {
		return log.NewDefaultLogger(), nil
	}

	l, ok := logger.(*zerolog.Logger)
	if !ok {
		return nil, errors.New("logger not of type *zerolog.Logger")
	}

	return l, nil
}

// SetLogger stores a logger on the app Metadata.
func SetLogger(c *cli.Context, logger *zerolog.Logger) {
	setMetadata(c, loggerKey, logger)
}

// GetErrLogger retrieves an errLogger from the app Metadata, if one
// does not exist it will return a new default errLogger.
func GetErrLogger(c *cli.Context) (*zerolog.Logger, error) {
	logger := getMetadata(c, errLoggerKey)
	if logger == nil {
		return log.NewDefaultErrLogger(), nil
	}

	l, ok := logger.(*zerolog.Logger)
	if !ok {
		return nil, errors.New("errLogger not of type *zerolog.Logger")
	}

	return l, nil
}

// SetErrLogger stores an errLogger on the app Metadata.
func SetErrLogger(c *cli.Context, logger *zerolog.Logger) {
	setMetadata(c, errLoggerKey, logger)
}

// GetServer retrieves a server from the app Metadata if one exists.
func GetServer(c *cli.Context) (*BatikServer, error) {
	server := getMetadata(c, serverKey)
	if server == nil {
		return nil, nil
	}

	s, ok := server.(*BatikServer)
	if !ok {
		return nil, errors.New("server not of type *BatikServer")
	}

	return s, nil
}

// SetServer stores a server on the app Metadata.
func SetServer(c *cli.Context, server *BatikServer) {
	setMetadata(c, serverKey, server)
}

func getMetadata(c *cli.Context, key metadataKey) interface{} {
	return c.App.Metadata[string(key)]
}

func setMetadata(c *cli.Context, key metadataKey, val interface{}) {
	c.App.Metadata[string(key)] = val
}
