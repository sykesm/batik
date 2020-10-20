// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package encodeservice

import (
	"context"
	"crypto"

	txv1 "github.com/sykesm/batik/pkg/pb/tx/v1"
	"github.com/sykesm/batik/pkg/transaction"
)

// EncodeService implements the EncodeAPIServer gRPC interface.
type EncodeService struct {
	// Unsafe has been chosen to ensure there's a compilation failure when the
	// implementation does not match the service interface.
	txv1.UnsafeEncodeAPIServer
}

var _ txv1.EncodeAPIServer = (*EncodeService)(nil)

// Encode encodes a transaction via deterministic marshal and returns the
// encoded bytes as well as a hash over the transaction represented as a merkle
// root and generated via SHA256 as the internal hashing function.
func (e *EncodeService) Encode(ctx context.Context, req *txv1.EncodeRequest) (*txv1.EncodeResponse, error) {
	tx := req.Transaction

	intTx, err := transaction.Marshal(crypto.SHA256, tx)
	if err != nil {
		return nil, err
	}

	return &txv1.EncodeResponse{
		Txid:               intTx.ID,
		EncodedTransaction: intTx.Encoded,
	}, nil
}
