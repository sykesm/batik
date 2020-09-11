// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/tedsuo/ifrit"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/sykesm/batik/pkg/grpccomm"
	sb "github.com/sykesm/batik/pkg/pb/store"
	tb "github.com/sykesm/batik/pkg/pb/transaction"
	"github.com/sykesm/batik/pkg/store"
	"github.com/sykesm/batik/pkg/transaction"
)

// BatikServer is a gRPC server that provides services for interacting
// with the store and encoding transactions.
type BatikServer struct {
	address string
	server  *grpccomm.Server

	logger *zap.Logger

	db store.KV
}

func NewServer(config Config, logger *zap.Logger) (*BatikServer, error) {
	grpcserver := grpccomm.NewServer(
		grpccomm.ServerConfig{
			ListenAddress: config.Server.Address,
			Logger:        logger,
		},
	)
	server := &BatikServer{
		address: config.Server.Address,
		server:  grpcserver,
		logger:  logger.With(zap.String("address", config.Server.Address)),
	}

	if err := server.initializeDB(config.DBPath); err != nil {
		return nil, err
	}

	if err := server.registerServices(); err != nil {
		return nil, err
	}

	return server, nil
}

// initializeDB initializes a new LevelDB instance at the dbPath or in memory
// if the path is empty.
func (s *BatikServer) initializeDB(dbPath string) error {
	db, err := store.NewLevelDB("")
	if err != nil {
		return err
	}

	s.db = db

	return nil
}

// registerServices registers the gRPC services supported by this server.
func (s *BatikServer) registerServices() error {
	encodeTxSvc := &transaction.EncodeService{}
	tb.RegisterEncodeTransactionAPIServer(s.server.Server, encodeTxSvc)

	if s.db == nil {
		return errors.New("server db not initialized")
	}

	storeSvc := store.NewStoreService(s.db)
	sb.RegisterStoreAPIServer(s.server.Server, storeSvc)

	return nil
}

func (s *BatikServer) Start() error {
	s.logger.Info("Starting server")

	process := ifrit.Invoke(s.server)
	s.handleSignals(map[os.Signal]func(){
		syscall.SIGINT:  func() { process.Signal(syscall.SIGINT) },
		syscall.SIGTERM: func() { process.Signal(syscall.SIGTERM) },
	})
	s.logger.Info("Server started")

	// Block until grpc server exits
	return <-process.Wait()
}

func (s *BatikServer) Stop() {
	s.logger.Info("Stopping server")
	s.server.GracefulStop()
}

func (s *BatikServer) Status() error {
	s.logger.Info("Checking status of server")

	// create GRPC client conn
	clientConn, err := grpc.Dial(s.address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer clientConn.Close()

	//TODO add client healthcheck to verify grpc server status

	return nil
}

func (s *BatikServer) handleSignals(handlers map[os.Signal]func()) {
	var signals []os.Signal
	for sig := range handlers {
		signals = append(signals, sig)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, signals...)

	go func() {
		for sig := range signalChan {
			s.logger.Warn("Received signal", zap.String("signal", sig.String()))
			handlers[sig]()
		}
	}()
}
