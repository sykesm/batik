// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package orderer

import (
	"runtime/debug"

	"github.com/sykesm/batik/pkg/orderer/broadcast"
	"github.com/sykesm/batik/pkg/orderer/deliver"
	ab "github.com/sykesm/batik/pkg/pb/orderer"
	"go.uber.org/zap"
)

// AtomicBroadcastService implements the AtomicBroadcastAPIServer gRPC interface.
type AtomicBroadcastService struct {
	logger *zap.Logger

	bh *broadcast.Handler
	dh *deliver.Handler
}

var _ ab.AtomicBroadcastAPIServer = (*AtomicBroadcastService)(nil)

func NewAtomicBroadcastService(logger *zap.Logger) ab.AtomicBroadcastAPIServer {
	return &AtomicBroadcastService{
		logger: logger,
		dh:     &deliver.Handler{Logger: logger.Named("deliver")},
		bh:     &broadcast.Handler{Logger: logger.Named("broadcast")},
	}
}

// Broadcast receives a stream of messages from a client for ordering.
func (s *AtomicBroadcastService) Broadcast(srv ab.AtomicBroadcastAPI_BroadcastServer) error {
	s.logger.Debug("Starting new Broadcast handler")
	defer func() {
		if r := recover(); r != nil {
			r := r.(string)
			s.logger.Fatal("Broadcast client triggered panic", zap.String("panic", r), zap.String("stack", string(debug.Stack())))
		}
		s.logger.Debug("Closing Broadcast stream")
	}()

	return s.bh.Handle(srv)
}

// Deliver sends a stream of transaction ids to a client after ordering.
func (s *AtomicBroadcastService) Deliver(srv ab.AtomicBroadcastAPI_DeliverServer) error {
	s.logger.Debug("Starting new Deliver handler")
	defer func() {
		if r := recover(); r != nil {
			r := r.(string)
			s.logger.Fatal("Deliver client triggered panic", zap.String("panic", r), zap.String("stack", string(debug.Stack())))
		}
		s.logger.Debug("Closing Deliver stream")
	}()

	return s.dh.Handle(srv)
}
