// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"errors"

	"github.com/urfave/cli/v2"
	"github.com/sykesm/batik/pkg/log"
	"go.uber.org/zap"
)

type metadataKey string

const (
	configKey metadataKey = "config"
	loggerKey             = "logger"
	serverKey             = "server"
)

// GetConfig retrieves a Config object from the app Metadata.
func GetConfig(ctx *cli.Context) Config {
	config, ok := getMetadata(ctx, configKey).(Config)
	if !ok {
		return Config{}
	}

	return config
}

// SetConfig stores a Config object on the app Metadata.
func SetConfig(ctx *cli.Context, config Config) {
	setMetadata(ctx, configKey, config)
}

// GetLogger retrieves a logger from the app Metadata, if one
// does not exist it will return a new default logger.
func GetLogger(ctx *cli.Context) (*zap.Logger, error) {
	logger := getMetadata(ctx, loggerKey)
	if logger == nil {
		return log.NewLogger(log.Config{})
	}

	l, ok := logger.(*zap.Logger)
	if !ok {
		return nil, errors.New("logger not of type *zerolog.Logger")
	}

	return l, nil
}

// SetLogger stores a logger on the app Metadata.
func SetLogger(ctx *cli.Context, logger *zap.Logger) {
	setMetadata(ctx, loggerKey, logger)
}

// GetServer retrieves a server from the app Metadata if one exists.
func GetServer(ctx *cli.Context) (*BatikServer, error) {
	server := getMetadata(ctx, serverKey)
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
func SetServer(ctx *cli.Context, server *BatikServer) {
	setMetadata(ctx, serverKey, server)
}

func getMetadata(ctx *cli.Context, key metadataKey) interface{} {
	return ctx.App.Metadata[string(key)]
}

func setMetadata(ctx *cli.Context, key metadataKey, val interface{}) {
	ctx.App.Metadata[string(key)] = val
}
