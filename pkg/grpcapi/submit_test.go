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

	txv1 "github.com/sykesm/batik/pkg/pb/tx/v1"
	"github.com/sykesm/batik/pkg/store"
	. "github.com/sykesm/batik/pkg/tested/matcher"
	"github.com/sykesm/batik/pkg/transaction"
)

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

	tests := map[string]struct {
		req        *txv1.SubmitRequest
		submitErr  error
		resp       *txv1.SubmitResponse
		errMatcher types.GomegaMatcher
	}{
		"nil signed transaction": {
			req:        &txv1.SubmitRequest{},
			errMatcher: MatchError(status.Errorf(codes.InvalidArgument, "signed transaction was not provided")),
		},
		"nil transaction": {
			req:        &txv1.SubmitRequest{SignedTransaction: &txv1.SignedTransaction{}},
			errMatcher: MatchError(status.Errorf(codes.InvalidArgument, "transaction was not provided")),
		},
		"invalid transaction": {
			req:        &txv1.SubmitRequest{SignedTransaction: &txv1.SignedTransaction{Transaction: &txv1.Transaction{}}},
			errMatcher: MatchError(status.Errorf(codes.InvalidArgument, "transaction salt is missing or less than 32 bytes in length")),
		},
		"valid transaction": {
			req:  &txv1.SubmitRequest{SignedTransaction: &txv1.SignedTransaction{Transaction: tx.Tx}},
			resp: &txv1.SubmitResponse{Txid: tx.ID},
		},
		"already exists error": {
			req:        &txv1.SubmitRequest{SignedTransaction: &txv1.SignedTransaction{Transaction: tx.Tx}},
			submitErr:  &store.AlreadyExistsError{Err: errors.New("already-exists")},
			errMatcher: MatchError(status.Errorf(codes.AlreadyExists, "storing transaction %s failed: already-exists", tx.ID)),
		},
		"not found error": {
			req:        &txv1.SubmitRequest{SignedTransaction: &txv1.SignedTransaction{Transaction: tx.Tx}},
			submitErr:  &store.NotFoundError{Err: errors.New("not-found")},
			errMatcher: MatchError(status.Errorf(codes.FailedPrecondition, "storing transaction %s failed: not-found", tx.ID)),
		},
		"unknown error": {
			req:        &txv1.SubmitRequest{SignedTransaction: &txv1.SignedTransaction{Transaction: tx.Tx}},
			submitErr:  errors.New("woops"),
			errMatcher: MatchError(status.Errorf(codes.Unknown, "storing transaction %s failed: woops", tx.ID)),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)
			var submitter submitterFunc = func(ctx context.Context, tx *transaction.Signed) error {
				return tt.submitErr
			}
			ss := NewSubmitService(submitter)
			resp, err := ss.Submit(context.Background(), tt.req)
			if tt.errMatcher != nil {
				gt.Expect(err).To(tt.errMatcher)
				return
			}
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(resp).To(ProtoEqual(tt.resp))
		})
	}
}
