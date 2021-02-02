// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"context"

	"github.com/pkg/errors"
	zaplogfmt "github.com/sykesm/zap-logfmt"
	cli "github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/sykesm/batik/pkg/log"
	"github.com/sykesm/batik/pkg/namespace"
)

type contextKey int

const (
	configKey contextKey = iota
	loggerKey
	levelerKey
	serverKey
	namespacesKey
)

// GetLogger retrieves a zap.Logger from the *cli.Context if one exists.
// If no logger exists on the context a default one is created and returned.
func GetLogger(ctx *cli.Context) (*zap.Logger, error) {
	logger, ok := retrieveFromCtx(ctx, loggerKey).(*zap.Logger)
	if ok {
		return logger, nil
	}

	logger = log.NewLogger(
		zaplogfmt.NewEncoder(zap.NewProductionEncoderConfig()),
		log.NewWriteSyncer(ctx.App.ErrWriter),
		log.NewLeveler("info"),
	)
	return logger.Named(ctx.App.Name), nil
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

// GetNamespaces retrieves the namespaces map from the *cli.Context if one exists.
func GetNamespaces(ctx *cli.Context) map[string]*namespace.Namespace {
	val := retrieveFromCtx(ctx, namespacesKey)
	if val == nil {
		return nil
	}

	namespaces, ok := val.(map[string]*namespace.Namespace)
	if !ok {
		return nil
	}

	return namespaces
}

// SetNamespaces stores a map of namespaces on the *cli.Context.
func SetNamespaces(ctx *cli.Context, namespaces map[string]*namespace.Namespace) {
	setOnCtx(ctx, namespacesKey, namespaces)
}

func GetCurrentNamespace(ctx *cli.Context) (*namespace.Namespace, error) {
	namespaces := GetNamespaces(ctx)
	if namespaces == nil {
		return nil, errors.Errorf("could not find namespaces from context")
	}

	namespaceName := ctx.String("namespace")
	if namespaceName == "" {
		return nil, errors.Errorf("target namespace is not set")
	}

	ns, ok := namespaces[namespaceName]
	if !ok {
		return nil, errors.Errorf("namespace %q is not defined", namespaceName)
	}

	return ns, nil
}

func retrieveFromCtx(ctx *cli.Context, key contextKey) interface{} {
	return ctx.Context.Value(key)
}

func setOnCtx(ctx *cli.Context, key contextKey, val interface{}) {
	ctx.Context = context.WithValue(ctx.Context, key, val)
}
