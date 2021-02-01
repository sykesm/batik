// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"context"
	"net"
	"net/http"
	"os"
	"path/filepath"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/sigmon"
	cli "github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/sykesm/batik/pkg/grpcapi"
	"github.com/sykesm/batik/pkg/grpccomm"
	"github.com/sykesm/batik/pkg/grpclogging"
	"github.com/sykesm/batik/pkg/namespace"
	"github.com/sykesm/batik/pkg/options"
	storev1 "github.com/sykesm/batik/pkg/pb/store/v1"
	txv1 "github.com/sykesm/batik/pkg/pb/tx/v1"
)

func startCommand(config *options.Batik, interactive bool) *cli.Command {
	return &cli.Command{
		Name:        "start",
		Description: "Establish network connections and begin processing.",
		Usage:       "start the server",
		Flags:       config.Server.Flags(),
		Action: func(ctx *cli.Context) error {
			return startAction(ctx, config, interactive)
		},
	}
}

func startAction(ctx *cli.Context, config *options.Batik, interactive bool) error {
	logger, err := GetLogger(ctx)
	if err != nil {
		return cli.Exit(err, exitServerStartFailed)
	}

	grpcLogger := logger.Named("grpc")
	grpcServerOptions := config.Server.GRPC.BuildServerOptions()

	grpcServerOptions = append(grpcServerOptions,
		grpc.ChainUnaryInterceptor(
			grpclogging.UnaryServerInterceptor(grpcLogger),
		),
		grpc.ChainStreamInterceptor(
			grpclogging.StreamServerInterceptor(grpcLogger),
		),
	)

	tlsConf, err := config.Server.TLS.TLSConfig()
	if errors.Is(err, options.ErrServerTLSNotBootstrapped) {
		logger.Debug("initializing server TLS with self signed certificate", zap.String("certs_dir", config.Server.TLS.CertsDir))
		err = config.Server.TLS.Bootstrap()
		if err != nil {
			return cli.Exit(errors.WithMessage(err, "failed to bootstrap server tls"), exitServerCreateFailed)
		}

		tlsConf, err = config.Server.TLS.TLSConfig()
	}
	if err != nil {
		return cli.Exit(errors.WithMessage(err, "failed to create server"), exitServerCreateFailed)
	}

	grpcServerOptions = append(grpcServerOptions, grpc.Creds(credentials.NewTLS(tlsConf)))

	grpcServer := grpccomm.NewServer(
		grpccomm.ServerConfig{
			ListenAddress: config.Server.GRPC.ListenAddress,
			Logger:        grpcLogger,
		},
		grpcServerOptions...,
	)

	encodeService := &grpcapi.EncodeService{}
	txv1.RegisterEncodeAPIServer(grpcServer.Server, encodeService)

	namespaceNames := []string{"namespace"}

	namespaces := map[string]*namespace.Namespace{}
	for _, namespaceName := range namespaceNames {
		namespaceLogger := logger.With(zap.String("namespace", namespaceName))

		dbPath := filepath.Join(config.Ledger.DataDir, namespaceName)
		namespaceLogger.Debug("initializing database", zap.String("data_dir", dbPath))
		db, err := levelDB(ctx, dbPath)
		if err != nil {
			return cli.Exit(errors.Wrap(err, "failed to create server"), exitServerCreateFailed)
		}

		namespaces[namespaceName] = namespace.New(namespaceLogger, db)
	}

	grpcapiAdapter := grpcapi.NamespaceMapAdapter(namespaces)

	submitService := grpcapi.NewSubmitService(grpcapiAdapter)
	txv1.RegisterSubmitAPIServer(grpcServer.Server, submitService)

	storeService := grpcapi.NewStoreService(grpcapiAdapter)
	storev1.RegisterStoreAPIServer(grpcServer.Server, storeService)

	mux := gwruntime.NewServeMux()
	storev1.RegisterStoreAPIHandlerServer(context.Background(), mux, storeService)

	httpServer := config.Server.HTTP.BuildServer(tlsConf)
	httpServer.Handler = mux

	httpRunner := func(signals <-chan os.Signal, ready chan<- struct{}) error {
		lis, err := net.Listen("tcp", config.Server.HTTP.ListenAddress)
		if err != nil {
			return err
		}

		errCh := make(chan error, 1)
		if httpServer.TLSConfig == nil {
			go func() { errCh <- httpServer.Serve(lis) }()
		} else {
			go func() { errCh <- httpServer.ServeTLS(lis, "", "") }()
		}

		close(ready)
		select {
		case <-signals:
			return nil
		case err := <-errCh:
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			return err
		}
	}

	logger.Info(
		"Starting server",
		zap.String("grpc-address", config.Server.GRPC.ListenAddress),
		zap.String("http-address", config.Server.HTTP.ListenAddress),
	)
	grpcProcess := ifrit.Invoke(sigmon.New(grpcServer))
	httpProcess := ifrit.Invoke(sigmon.New(ifrit.RunFunc(httpRunner)))
	logger.Info("Server started")
	if !interactive {
		select {
		case err := <-grpcProcess.Wait():
			return err
		case err := <-httpProcess.Wait():
			return err
		}
	}

	return nil
}
