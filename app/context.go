// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"context"
	"errors"

	zaplogfmt "github.com/sykesm/zap-logfmt"
	cli "github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/sykesm/batik/pkg/log"
)

type contextKey int

const (
	configKey contextKey = iota
	loggerKey
	serverKey
)

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
		return nil, errors.New("logger not of type *zerolog.Logger")
	}

	return l, nil
}

// SetLogger stores a *zap.Logger on the *cli.Context.
func SetLogger(ctx *cli.Context, logger *zap.Logger) {
	setOnCtx(ctx, loggerKey, logger)
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
