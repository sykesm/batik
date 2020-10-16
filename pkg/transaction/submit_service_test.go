// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package transaction

import (
	"context"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"

	tb "github.com/sykesm/batik/pkg/pb/transaction"
	. "github.com/sykesm/batik/pkg/tested/matcher"
)

func TestSubmitTransaction(t *testing.T) {
	tests := map[string]struct {
		req        *tb.SubmitTransactionRequest
		resp       *tb.SubmitTransactionResponse
		errMatcher types.GomegaMatcher
	}{
		"nil transaction": {
			req:        &tb.SubmitTransactionRequest{},
			resp:       nil,
			errMatcher: HaveOccurred(),
		},
		"valid transaction": {
			req: &tb.SubmitTransactionRequest{
				Transaction: &tb.Transaction{
					Outputs: []*tb.State{{
						Info:  &tb.StateInfo{Kind: "test-kind"},
						State: []byte("test-state-1"),
					}},
				},
			},
			resp: &tb.SubmitTransactionResponse{
				Txid: fromHex(t, "c6892f1044b2e7fe7731f47a58297925fed43bd797c878028392454690fa973e"),
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			ss := NewSubmitService()
			resp, err := ss.SubmitTransaction(context.Background(), tt.req)
			if tt.errMatcher != nil {
				gt.Expect(err).To(tt.errMatcher)
				return
			}
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(resp).To(ProtoEqual(tt.resp))
		})
	}
}
