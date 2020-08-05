package transaction

import (
	"context"
	"crypto"

	tb "github.com/sykesm/batik/pkg/pb/transaction"
	"github.com/sykesm/batik/pkg/protomsg"
)

type EncodeService struct{}

func (e *EncodeService) EncodedTransaction(ctx context.Context, req *tb.EncodedTransactionRequest) (*tb.EncodedTransactionResponse, error) {
	tx := req.Transaction

	id, err := ID(crypto.SHA256, tx)
	if err != nil {
		return nil, err
	}

	encodedTx, err := protomsg.MarshalDeterministic(tx)
	if err != nil {
		return nil, err
	}

	return &tb.EncodedTransactionResponse{
		Txid:               id,
		EncodedTransaction: encodedTx,
	}, nil
}
