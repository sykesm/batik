// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package grpccomm

import (
	"net"
	"os"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	*grpc.Server

	listenAddress string
	listener      net.Listener
	logger        *zap.Logger
}

type ServerConfig struct {
	ListenAddress string
	Listener      net.Listener
	Logger        *zap.Logger
}

func NewServer(config ServerConfig, grpcOptions ...grpc.ServerOption) *Server {
	return &Server{
		Server:        grpc.NewServer(grpcOptions...),
		listenAddress: config.ListenAddress,
		listener:      config.Listener,
		logger:        config.Logger,
	}
}

func (s *Server) Run(signals <-chan os.Signal, ready chan<- struct{}) error {
	if s.listener == nil {
		l, err := net.Listen("tcp", s.listenAddress)
		if err != nil {
			return errors.Wrap(err, "grpccomm: listen failed")
		}
		s.listener = l
	}

	errCh := make(chan error, 1)
	go func() {
		err := s.Server.Serve(s.listener)
		if err != nil {
			err = errors.Wrap(err, "grpccomm: server terminated")
		}
		errCh <- err
	}()

	close(ready)

	select {
	case <-signals:
	case err := <-errCh:
		return err
	}

	s.Server.GracefulStop()
	return nil
}
