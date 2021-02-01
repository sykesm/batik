// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package grpcapi

import (
	"context"
	"crypto"
	"errors"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	txv1 "github.com/sykesm/batik/pkg/pb/tx/v1"
	"github.com/sykesm/batik/pkg/store"
	. "github.com/sykesm/batik/pkg/tested/matcher"
	"github.com/sykesm/batik/pkg/transaction"
)

type submitMapAdapter map[string]submitterFunc

func (sma submitMapAdapter) Submitter(namespace string) Submitter {
	s, ok := sma[namespace]
	if !ok {
		return notFoundSubmitter(namespace)
	}

	return s
}

type submitterFunc func(context.Context, *transaction.Signed) error

func (s submitterFunc) Submit(ctx context.Context, tx *transaction.Signed) error {
	return s(ctx, tx)
}

func TestSubmit(t *testing.T) {
	gt := NewGomegaWithT(t)

	tx, err := transaction.New(crypto.SHA256, &txv1.Transaction{
		Salt: []byte("potassium permanganate (KMnO4) is a salt"),
		Outputs: []*txv1.State{{
			Info:  &txv1.StateInfo{Kind: "test-kind"},
			State: []byte("test-state-1"),
		}},
	})
	gt.Expect(err).NotTo(HaveOccurred())

	submitRequest := &txv1.SubmitRequest{
		Namespace: "namespace",
		SignedTransaction: &txv1.SignedTransaction{
			Transaction: tx.Tx,
			Signatures:  nil,
		},
	}

	tests := map[string]struct {
		setup      func(sr *txv1.SubmitRequest)
		submitErr  error
		resp       *txv1.SubmitResponse
		errMatcher types.GomegaMatcher
	}{
		"nil signed transaction": {
			setup:      func(sr *txv1.SubmitRequest) { sr.SignedTransaction = nil },
			errMatcher: MatchError(status.Errorf(codes.InvalidArgument, "signed transaction was not provided")),
		},
		"nil transaction": {
			setup:      func(sr *txv1.SubmitRequest) { sr.SignedTransaction.Transaction = nil },
			errMatcher: MatchError(status.Errorf(codes.InvalidArgument, "transaction was not provided")),
		},
		"invalid transaction": {
			setup:      func(sr *txv1.SubmitRequest) { sr.SignedTransaction.Transaction = &txv1.Transaction{} },
			errMatcher: MatchError(status.Errorf(codes.InvalidArgument, "transaction salt is missing or less than 32 bytes in length")),
		},
		"unknown namespace": {
			setup:      func(sr *txv1.SubmitRequest) { sr.Namespace = "missing" },
			errMatcher: MatchError(status.Errorf(codes.InvalidArgument, "storing transaction %s failed: bad namespace %q: namespace not found", tx.ID, "missing")),
		},
		"valid transaction": {
			resp: &txv1.SubmitResponse{Txid: tx.ID},
		},
		"already exists error": {
			submitErr:  &store.AlreadyExistsError{Err: errors.New("already-exists")},
			errMatcher: MatchError(status.Errorf(codes.AlreadyExists, "storing transaction %s failed: already-exists", tx.ID)),
		},
		"not found error": {
			submitErr:  &store.NotFoundError{Err: errors.New("not-found")},
			errMatcher: MatchError(status.Errorf(codes.FailedPrecondition, "storing transaction %s failed: not-found", tx.ID)),
		},
		"unknown error": {
			submitErr:  errors.New("woops"),
			errMatcher: MatchError(status.Errorf(codes.Unknown, "storing transaction %s failed: woops", tx.ID)),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)
			req := proto.Clone(submitRequest).(*txv1.SubmitRequest)
			if tt.setup != nil {
				tt.setup(req)
			}
			var submitter submitterFunc = func(ctx context.Context, tx *transaction.Signed) error {
				return tt.submitErr
			}
			ss := NewSubmitService(submitMapAdapter(map[string]submitterFunc{
				"namespace": submitter,
			}))

			resp, err := ss.Submit(context.Background(), req)
			if tt.errMatcher != nil {
				gt.Expect(err).To(tt.errMatcher)
				return
			}
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(resp).To(ProtoEqual(tt.resp))
		})
	}
}
