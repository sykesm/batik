// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package transaction

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tb "github.com/sykesm/batik/pkg/pb/transaction"
)

// SubmitService implements the EncodeTransactionAPIServer gRPC interface.
type SubmitService struct {
	tb.UnimplementedSubmitTransactionAPIServer
}

var _ tb.SubmitTransactionAPIServer = (*SubmitService)(nil)

func (s *SubmitService) SubmitTransaction(ctx context.Context, req *tb.SubmitTransactionRequest) (*tb.SubmitTransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "I am not a teapot")
}
