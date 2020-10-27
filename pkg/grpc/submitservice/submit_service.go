// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package submitservice

import (
	"context"
	"crypto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sykesm/batik/pkg/merkle"
	txv1 "github.com/sykesm/batik/pkg/pb/tx/v1"
	"github.com/sykesm/batik/pkg/store"
	"github.com/sykesm/batik/pkg/transaction"
)

// SubmitService implements the EncodeAPIServer gRPC interface.
type SubmitService struct {
	// Unnsafe has been chosed to ensure there's a compilation failure when the
	// implementation diverges from the gRPC service.
	txv1.UnsafeSubmitAPIServer
	// hasher implements the hash algorithm used to build and validate the
	// transaction ID.
	hasher merkle.Hasher
	// kv is a reference to the key value store backing this service
	kv store.KV
}

var _ txv1.SubmitAPIServer = (*SubmitService)(nil)

// NewSubmitService creates a new instance of the SubmitService.
func NewSubmitService(kv store.KV) *SubmitService {
	return &SubmitService{
		hasher: crypto.SHA256,
		kv:     kv,
	}
}

// Submit submits a transaction for validation and commit processing.
//
// NOTE: This is an implementation for prototyping.
func (s *SubmitService) Submit(ctx context.Context, req *txv1.SubmitRequest) (*txv1.SubmitResponse, error) {
	signedTx := req.GetSignedTransaction()
	if signedTx == nil {
		return nil, status.Errorf(codes.InvalidArgument, "signed transaction was not provided")
	}
	tx := signedTx.GetTransaction()
	if tx == nil {
		return nil, status.Errorf(codes.InvalidArgument, "transaction was not provided")
	}
	itx, err := transaction.Marshal(s.hasher, tx)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	// Check if the transaction already exists
	_, err = transaction.LoadTransaction(s.kv, itx.ID)
	if !store.IsNotFound(err) {
		return nil, status.Errorf(codes.AlreadyExists, "transaction %s already exists", itx.ID)
	}

	_, err = resolve(s.kv, tx)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "state resolution for transaction %s failed: %s", itx.ID, err)
	}

	// TODO: Store of transaction and states *must* be atomic.
	// TODO: Consumed outputs must be marked as consumed.
	// TODO: The data store should be using the intermediate tx with the marshaled state

	err = transaction.StoreTransactions(s.kv, []*txv1.Transaction{tx})
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "storing transaction %s failed: %s", itx.ID, err)
	}

	var resolvedOutputs []*txv1.ResolvedState
	for i := range tx.Outputs {
		resolvedOutputs = append(resolvedOutputs, &txv1.ResolvedState{
			Txid:        itx.ID,
			OutputIndex: uint64(i),
			Info:        tx.Outputs[i].Info,
			State:       tx.Outputs[i].State,
		})
	}

	err = transaction.StoreStates(s.kv, resolvedOutputs)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "storing transaction %s failed: %s", itx.ID, err)
	}

	err = transaction.ConsumeStates(s.kv, tx.Inputs)
	if err != nil {
		return nil, err
	}

	return &txv1.SubmitResponse{Txid: itx.ID}, nil
}

// A Resolved transaction is a Transaction where state references have been
// resolved and populated with data from the ledger.
type ResolvedTransaction struct {
	Salt            []byte
	Inputs          []*txv1.State
	References      []*txv1.State
	Outputs         []*txv1.State
	Parameters      []*txv1.Parameter
	RequiredSigners []*txv1.Party
}

func resolve(kv store.KV, tx *txv1.Transaction) (*ResolvedTransaction, error) {
	resolved := &ResolvedTransaction{
		Salt:            tx.Salt,
		Outputs:         tx.Outputs,
		Parameters:      tx.Parameters,
		RequiredSigners: tx.RequiredSigners,
	}

	inputs, err := transaction.LoadStates(kv, tx.Inputs)
	if err != nil {
		return nil, err
	}
	for _, input := range inputs {
		resolved.Inputs = append(resolved.Inputs, &txv1.State{
			Info:  input.Info,
			State: input.State,
		})
	}
	refs, err := transaction.LoadStates(kv, tx.References)
	if err != nil {
		return nil, err
	}
	for _, ref := range refs {
		resolved.References = append(resolved.References, &txv1.State{
			Info:  ref.Info,
			State: ref.State,
		})
	}

	return resolved, nil
}

func verifySignature(pk crypto.PublicKey, signature, hash []byte) bool {
	return false
}
