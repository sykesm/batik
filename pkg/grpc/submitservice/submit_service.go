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

// A Repository abstracts the data persistence layer for transactions and
// states.
type Repository interface {
	PutTransaction(*transaction.Transaction) error
	GetTransaction(transaction.ID) (*transaction.Transaction, error)
	PutState(*transaction.State) error
	GetState(transaction.StateID) (*transaction.State, error)
	ConsumeStates(...transaction.StateID) error
}

// SubmitService implements the EncodeAPIServer gRPC interface.
type SubmitService struct {
	// Unnsafe has been chosed to ensure there's a compilation failure when the
	// implementation diverges from the gRPC service.
	txv1.UnsafeSubmitAPIServer
	// hasher implements the hash algorithm used to build and validate the
	// transaction ID.
	hasher merkle.Hasher
	// repo is a reference to the transaction state repository.
	repo Repository
}

var _ txv1.SubmitAPIServer = (*SubmitService)(nil)

// NewSubmitService creates a new instance of the SubmitService.
func NewSubmitService(repo Repository) *SubmitService {
	return &SubmitService{
		hasher: crypto.SHA256,
		repo:   repo,
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
	itx, err := transaction.New(s.hasher, tx)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	// Check if the transaction already exists
	_, err = s.repo.GetTransaction(itx.ID)
	if !store.IsNotFound(err) {
		return nil, status.Errorf(codes.AlreadyExists, "transaction %s already exists", itx.ID)
	}

	_, err = resolve(s.repo, tx)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "state resolution for transaction %s failed: %s", itx.ID, err)
	}

	err = s.repo.PutTransaction(itx)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, "storing transaction %s failed: %s", itx.ID, err)
	}

	for i := range tx.Outputs {
		state := &transaction.State{
			ID:        transaction.StateID{TxID: itx.ID, OutputIndex: uint64(i)},
			StateInfo: tx.Outputs[i].Info,
			Data:      tx.Outputs[i].State,
		}

		err = s.repo.PutState(state)
		if err != nil {
			return nil, status.Errorf(codes.Unknown, "storing transaction %s failed: %s", itx.ID, err)
		}
	}

	var inputs []transaction.StateID
	for _, input := range tx.Inputs {
		inputs = append(inputs, transaction.StateID{TxID: input.Txid, OutputIndex: input.OutputIndex})
	}
	err = s.repo.ConsumeStates(inputs...)
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

func resolve(repo Repository, tx *txv1.Transaction) (*ResolvedTransaction, error) {
	resolved := &ResolvedTransaction{
		Salt:            tx.Salt,
		Outputs:         tx.Outputs,
		Parameters:      tx.Parameters,
		RequiredSigners: tx.RequiredSigners,
	}

	var inputs []*transaction.State
	for _, input := range tx.Inputs {
		stateID := transaction.StateID{TxID: input.Txid, OutputIndex: input.OutputIndex}
		state, err := repo.GetState(stateID)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, state)
	}
	var refs []*transaction.State
	for _, ref := range tx.References {
		stateID := transaction.StateID{TxID: ref.Txid, OutputIndex: ref.OutputIndex}
		state, err := repo.GetState(stateID)
		if err != nil {
			return nil, err
		}
		refs = append(refs, state)
	}

	for _, input := range inputs {
		resolved.Inputs = append(resolved.Inputs, &txv1.State{
			Info:  input.StateInfo,
			State: input.Data,
		})
	}
	for _, ref := range refs {
		resolved.References = append(resolved.References, &txv1.State{
			Info:  ref.StateInfo,
			State: ref.Data,
		})
	}

	return resolved, nil
}

func verifySignature(pk crypto.PublicKey, signature, hash []byte) bool {
	return false
}
