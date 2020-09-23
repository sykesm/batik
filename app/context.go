// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"context"
	"errors"

	zaplogfmt "github.com/sykesm/zap-logfmt"
	cli "github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/sykesm/batik/pkg/log"
)

type contextKey int

const (
	configKey contextKey = iota
	loggerKey
	levelerKey
	serverKey
)

// GetLogger retrieves a zap.Logger from the *cli.Context if one exists.
// If no logger exists on the context a default one is created and returned.
func GetLogger(ctx *cli.Context) (*zap.Logger, error) {
	logger := retrieveFromCtx(ctx, loggerKey)
	if logger == nil {
		return log.NewLogger(
			zaplogfmt.NewEncoder(zap.NewProductionEncoderConfig()),
			log.NewWriteSyncer(ctx.App.ErrWriter),
			log.NewLeveler("info"),
		).Named(ctx.App.Name), nil
	}

	l, ok := logger.(*zap.Logger)
	if !ok {
		return nil, errors.New("logger not of type *zap.Logger")
	}

	return l, nil
}

// SetLogger stores a *zap.Logger on the *cli.Context.
func SetLogger(ctx *cli.Context, logger *zap.Logger) {
	setOnCtx(ctx, loggerKey, logger)
}

// GetLeveler retrieves a zapcore.LevelEnabler from the *cli.Context if one exists.
// This leveler should be the one used by the enabled batik logger. If one does not
// exist, it will error.
func GetLeveler(ctx *cli.Context) (zapcore.LevelEnabler, error) {
	leveler := retrieveFromCtx(ctx, levelerKey)
	if leveler == nil {
		return nil, errors.New("leveler does not exist")
	}

	l, ok := leveler.(zapcore.LevelEnabler)
	if !ok {
		return nil, errors.New("leveler not of type zapcore.LevelEnabler")
	}

	return l, nil
}

// SetLeveler stores a zapcore.LevelEnabler on the *cli.Context.
func SetLeveler(ctx *cli.Context, leveler zapcore.LevelEnabler) {
	setOnCtx(ctx, levelerKey, leveler)
}

// // GetServer retrieves a server from the *cli.Context if one exists.
// func GetServer(ctx *cli.Context) (*BatikServer, error) {
// 	server := retrieveFromCtx(ctx, serverKey)
// 	if server == nil {
// 		return nil, nil
// 	}

// 	s, ok := server.(*BatikServer)
// 	if !ok {
// 		return nil, errors.New("server not of type *BatikServer")
// 	}

// 	return s, nil
// }

// // SetServer stores a *BatikServer on the *cli.Context.
// func SetServer(ctx *cli.Context, server *BatikServer) {
// 	setOnCtx(ctx, serverKey, server)
// }

func retrieveFromCtx(ctx *cli.Context, key contextKey) interface{} {
	return ctx.Context.Value(key)
}

func setOnCtx(ctx *cli.Context, key contextKey, val interface{}) {
	ctx.Context = context.WithValue(ctx.Context, key, val)
}
