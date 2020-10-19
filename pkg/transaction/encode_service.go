// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package transaction

import (
	"context"
	"crypto"

	txv1 "github.com/sykesm/batik/pkg/pb/tx/v1"
)

// EncodeService implements the EncodeTransactionAPIServer gRPC interface.
type EncodeService struct {
	// Unsafe has been chosen to ensure there's a compilation failure when the
	// implementation does not match the service interface.
	txv1.UnsafeEncodeTransactionAPIServer
}

var _ txv1.EncodeTransactionAPIServer = (*EncodeService)(nil)

// EncodeTransaction encodes a transaction via deterministic marshal and returns
// the encoded bytes as well as a hash over the transaction represented as a merkle
// root and generated via SHA256 as the internal hashing function.
func (e *EncodeService) EncodeTransaction(ctx context.Context, req *txv1.EncodeTransactionRequest) (*txv1.EncodeTransactionResponse, error) {
	tx := req.Transaction

	intTx, err := Marshal(crypto.SHA256, tx)
	if err != nil {
		return nil, err
	}

	return &txv1.EncodeTransactionResponse{
		Txid:               intTx.ID,
		EncodedTransaction: intTx.Encoded,
	}, nil
}
