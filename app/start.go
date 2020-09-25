// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"github.com/pkg/errors"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/sigmon"
	cli "github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/sykesm/batik/pkg/grpccomm"
	"github.com/sykesm/batik/pkg/grpclogging"
	"github.com/sykesm/batik/pkg/options"
	sb "github.com/sykesm/batik/pkg/pb/store"
	tb "github.com/sykesm/batik/pkg/pb/transaction"
	"github.com/sykesm/batik/pkg/store"
	"github.com/sykesm/batik/pkg/transaction"
)

func startCommand(config *options.Batik, interactive bool) *cli.Command {
	return &cli.Command{
		Name:        "start",
		Description: "Establish network connections and begin processing.",
		Usage:       "start the server",
		Flags:       config.Server.Flags(),
		Action: func(ctx *cli.Context) error {
			logger, err := GetLogger(ctx)
			if err != nil {
				return cli.Exit(err, exitServerStartFailed)
			}

			grpcLogger := logger.Named("grpc")
			gRPCOpts := config.Server.GRPC.BuildServerOptions()

			tlsConf, err := options.TLSConfig(config.Server.TLS)
			if err != nil {
				return cli.Exit(errors.WithMessage(err, "failed to create server"), exitServerCreateFailed)
			}
			if tlsConf != nil {
				gRPCOpts = append(gRPCOpts, grpc.Creds(credentials.NewTLS(tlsConf)))
			}
			gRPCOpts = append(gRPCOpts,
				grpc.ChainUnaryInterceptor(
					grpclogging.UnaryServerInterceptor(grpcLogger),
				),
				grpc.ChainStreamInterceptor(
					grpclogging.StreamServerInterceptor(grpcLogger),
				),
			)

			grpcServer := grpccomm.NewServer(
				grpccomm.ServerConfig{
					ListenAddress: config.Server.ListenAddress,
					Logger:        grpcLogger,
				},
				gRPCOpts...,
			)

			logger.Debug("initializing database", zap.String("data_dir", config.Ledger.DataDir))
			db, err := levelDB(ctx, config.Ledger.DataDir)
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
