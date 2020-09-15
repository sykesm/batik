// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"github.com/pkg/errors"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/sigmon"
	cli "github.com/urfave/cli/v2"

	"github.com/sykesm/batik/pkg/grpccomm"
	"github.com/sykesm/batik/pkg/options"
	sb "github.com/sykesm/batik/pkg/pb/store"
	tb "github.com/sykesm/batik/pkg/pb/transaction"
	"github.com/sykesm/batik/pkg/store"
	"github.com/sykesm/batik/pkg/transaction"
)

func startCommand(config *options.Batik, interactive bool) *cli.Command {
	return &cli.Command{
		Name:        "start",
		Description: "start the grpc server",
		Flags: append(
			config.Server.Flags(),
			config.Ledger.Flags()...,
		),
		Action: func(ctx *cli.Context) error {
			logger, err := GetLogger(ctx)
			if err != nil {
				return cli.Exit(err, exitServerStartFailed)
			}

			grpcServer := grpccomm.NewServer(
				grpccomm.ServerConfig{
					ListenAddress: config.Server.ListenAddress,
					Logger:        logger,
				},
			)

			db, err := store.NewLevelDB("")
			if err != nil {
				return cli.Exit(errors.Wrap(err, "failed to create server"), exitServerCreateFailed)
			}

			encodeService := &transaction.EncodeService{}
			tb.RegisterEncodeTransactionAPIServer(grpcServer.Server, encodeService)

			storeService := store.NewStoreService(db)
			sb.RegisterStoreAPIServer(grpcServer.Server, storeService)

			// SetServer(ctx, server)

			logger.Info("Starting server")
			process := ifrit.Invoke(sigmon.New(grpcServer))
			logger.Info("Server started")
			if !interactive {
				return <-process.Wait()
			}

			return nil
		},
	}
}
