// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package transaction

import (
	"context"
	"crypto"

	tb "github.com/sykesm/batik/pkg/pb/transaction"
	"github.com/sykesm/batik/pkg/protomsg"
)

// EncodeService implements the EncodeTransactionAPI gRPC interface.
type EncodeService struct{}

// EncodeTransaction encodes a transaction via deterministic marshal and returns
// the encoded bytes as well as a hash over the transaction represented as a merkle
// root and generated via SHA256 as the internal hashing function.
func (e *EncodeService) EncodeTransaction(ctx context.Context, req *tb.EncodeTransactionRequest) (*tb.EncodeTransactionResponse, error) {
	tx := req.Transaction

	id, err := ID(crypto.SHA256, tx)
	if err != nil {
		return nil, err
	}

	encodedTx, err := protomsg.MarshalDeterministic(tx)
	if err != nil {
		return nil, err
	}

	return &tb.EncodeTransactionResponse{
		Txid:               id,
		EncodedTransaction: encodedTx,
	}, nil
}
