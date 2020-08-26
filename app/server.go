// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	sb "github.com/sykesm/batik/pkg/pb/store"
	tb "github.com/sykesm/batik/pkg/pb/transaction"
	"github.com/sykesm/batik/pkg/store"
	"github.com/sykesm/batik/pkg/transaction"
	"google.golang.org/grpc"
)

// BatikServer is a gRPC server that provides services for interacting
// with the store and encoding transactions.
type BatikServer struct {
	address string
	server  *grpc.Server

	db *store.LevelDBKV
}

func NewServer(config Config) (*BatikServer, error) {
	server := &BatikServer{
		address: config.Server.Address,
		server:  grpc.NewServer(),
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
	fmt.Printf("Starting server at %s\n", s.address)
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	serve := make(chan error)

	handleSignals(map[os.Signal]func(){
		syscall.SIGINT:  func() { s.server.GracefulStop(); serve <- nil },
		syscall.SIGTERM: func() { s.server.GracefulStop(); serve <- nil },
	})

	go func() {
		var grpcErr error
		if grpcErr = s.server.Serve(listener); grpcErr != nil {
			grpcErr = fmt.Errorf("grpc server exited with error: %s", grpcErr)
		}
		serve <- grpcErr
	}()

	fmt.Println("Server started")

	// Block until grpc server exits
	return <-serve
}

func (s *BatikServer) Stop() {
	fmt.Println("Stopping server")
	s.server.GracefulStop()
}

func (s *BatikServer) Status() error {
	fmt.Printf("Checking status of server at %s\n", s.address)

	// create GRPC client conn
	clientConn, err := grpc.Dial(s.address, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer clientConn.Close()

	//TODO add client healthcheck to verify grpc server status

	return nil
}

func handleSignals(handlers map[os.Signal]func()) {
	var signals []os.Signal
	for sig := range handlers {
		signals = append(signals, sig)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, signals...)

	go func() {
		for sig := range signalChan {
			fmt.Printf("\nReceived signal: %d (%s)\n", sig, sig)
			handlers[sig]()
		}
	}()
}
