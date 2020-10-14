// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package broadcast

import (
	"context"
	"io"

	ab "github.com/sykesm/batik/pkg/pb/orderer"
	"go.uber.org/zap"
	"google.golang.org/grpc/peer"
)

// Handler is designed to handle connections from an AtomicBroadcastAPI gRPC service.
type Handler struct {
	Logger *zap.Logger
}

// Handle reads requests from a Broadcast stream, processes them, and returns the responses to the stream.
func (bh *Handler) Handle(srv ab.AtomicBroadcastAPI_BroadcastServer) error {
	addr := extractRemoteAddress(srv.Context())
	logger := bh.Logger.With(zap.String("addr", addr))
	logger.Debug("Starting new broadcast loop")

	for {
		msg, err := srv.Recv()
		if err == io.EOF {
			logger.Debug("Received EOF, hangup")
			return nil
		}
		if err != nil {
			logger.Warn("Error receiving", zap.String("err", err.Error()))
			return err
		}

		resp := bh.ProcessMessage(msg, addr)
		err = srv.Send(resp)
		if resp.Status != ab.Status_STATUS_SUCCESS {
			return err
		}

		if err != nil {
			logger.Warn("Error sending", zap.String("err", err.Error()))
			return err
		}
	}
}

// ProcessMessage validates and enqueues a single message.
func (bh *Handler) ProcessMessage(msg *ab.BroadcastRequest, addr string) *ab.BroadcastResponse {
	// TODO
	// 1. Parse request for payload and headers
	// 2. Get chain processor for channelid based on type of transaction (ie config, etc)
	// 3. Get sequence for tx via chain processor
	// 4. Wait for consenter to be ready to accept next tx
	// 5. Consenter orders the tx and appends to chain, or reconfigures if config tx
	return &ab.BroadcastResponse{Status: ab.Status_STATUS_SUCCESS}
}

func extractRemoteAddress(ctx context.Context) string {
	var remoteAddress string
	p, ok := peer.FromContext(ctx)
	if !ok {
		return ""
	}
	if address := p.Addr; address != nil {
		remoteAddress = address.String()
	}
	return remoteAddress
}
