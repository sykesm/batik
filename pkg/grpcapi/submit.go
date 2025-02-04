// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package grpcapi

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

type SubmitterMap interface {
	Submitter(namespace string) Submitter
}

type Submitter interface {
	Submit(ctx context.Context, signedTx *transaction.Signed) error
}

// SubmitService implements the EncodeAPIServer gRPC interface.
type SubmitService struct {
	// Unnsafe has been chosed to ensure there's a compilation failure when the
	// implementation diverges from the gRPC service.
	txv1.UnsafeSubmitAPIServer
	// hasher implements the hash algorithm used to build and validate the
	// transaction ID.
	hasher merkle.Hasher
	// submitters are the set of domain specific transaction processors asociated
	// with each namespace
	submitters SubmitterMap
}

var _ txv1.SubmitAPIServer = (*SubmitService)(nil)

// NewSubmitService creates a new instance of the SubmitService.
func NewSubmitService(submitters SubmitterMap) *SubmitService {
	return &SubmitService{
		hasher:     crypto.SHA256,
		submitters: submitters,
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

	signed := &transaction.Signed{
		Transaction: itx,
		Signatures:  transaction.ToSignatures(signedTx.Signatures...),
	}
	err = s.submitters.Submitter(req.Namespace).Submit(ctx, signed)
	if err != nil {
		code := codes.Unknown
		switch {
		case store.IsAlreadyExists(err):
			code = codes.AlreadyExists
		case store.IsNotFound(err):
			code = codes.FailedPrecondition
		case isNamespaceNotFound(err):
			code = codes.InvalidArgument
		}
		return nil, status.Errorf(code, "storing transaction %s failed: %s", itx.ID, err)
	}

	return &txv1.SubmitResponse{Txid: itx.ID}, nil
}
