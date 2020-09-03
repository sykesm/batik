// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	sb "github.com/sykesm/batik/pkg/pb/store"
	tb "github.com/sykesm/batik/pkg/pb/transaction"
	"github.com/sykesm/batik/pkg/store"
	"github.com/sykesm/batik/pkg/transaction"
)

// BatikServer is a gRPC server that provides services for interacting
// with the store and encoding transactions.
type BatikServer struct {
	address string
	server  *grpc.Server

	logger *zap.Logger

	db *store.LevelDBKV
}

// TODO(mjs): Replace stdout and stderr with loggers.

func NewServer(config Config, logger *zap.Logger) (*BatikServer, error) {
	server := &BatikServer{
		address: config.Server.Address,
		server:  grpc.NewServer(),
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
	tb.RegisterEncodeTransactionAPIServer(s.server, encodeTxSvc)

	if s.db == nil {
		return errors.New("server db not initialized")
	}

	storeSvc := &store.StoreService{
		Db: s.db,
	}
	sb.RegisterStoreAPIServer(s.server, storeSvc)

	return nil
}

func (s *BatikServer) Start() error {
	s.logger.Info("Starting server")
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	serve := make(chan error)

	s.handleSignals(map[os.Signal]func(){
		syscall.SIGINT:  func() { s.server.GracefulStop(); serve <- nil },
		syscall.SIGTERM: func() { s.server.GracefulStop(); serve <- nil },
	})

	go func() {
		var grpcErr error
		if grpcErr = s.server.Serve(listener); grpcErr != nil {
			grpcErr = errors.Wrap(grpcErr, "grpc server exited with error")
		}
		serve <- grpcErr
	}()

	s.logger.Info("Server started")

	// Block until grpc server exits
	return <-serve
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
