// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package deliver

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

// Handle reads requests from a Deliver stream, processes them, and returns the responses to the stream.
func (dh *Handler) Handle(srv ab.AtomicBroadcastAPI_DeliverServer) error {
	addr := extractRemoteAddress(srv.Context())
	logger := dh.Logger.With(zap.String("addr", addr))
	logger.Debug("Starting new deliver loop")

	for {
		msg, err := srv.Recv()
		if err == io.EOF {
			logger.Debug("Received EOF, hangup")
			return nil
		}
		if err != nil {
			logger.Warn("Error reading", zap.String("err", err.Error()))
			return err
		}

		resp := dh.ProcessMessage(msg, addr)
		err = srv.Send(resp)
		if r, ok := resp.Type.(*ab.DeliverResponse_Status); ok && r.Status != ab.Status_STATUS_SUCCESS {
			return err
		}

		if err != nil {
			logger.Warn("Error sending", zap.String("err", err.Error()))
			return err
		}
	}
}

func (h *Handler) ProcessMessage(msg *ab.DeliverRequest, addr string) *ab.DeliverResponse {
	//TODO
	// 1. Parse request for payload and headers
	// 2. Get "chain" for channel id
	// 3. Check for client authorization
	// 4. Get seek info from payload
	// 5. Iterate over chain of txs using seek info, for each tx:
	//		1. Check if client authorized (is this necessary to do again?)
	//		2. Send deliver response with tx or txid
	// 6. Return deliver response with successful status
	// On any failures, return deliver response with failure status
	return &ab.DeliverResponse{Type: &ab.DeliverResponse_Status{Status: ab.Status_STATUS_SUCCESS}}
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
