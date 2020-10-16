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
					Salt: []byte("potassium permanganate (KMnO4) is a salt"),
					Outputs: []*tb.State{{
						Info:  &tb.StateInfo{Kind: "test-kind"},
						State: []byte("test-state-1"),
					}},
				},
			},
			resp: &tb.SubmitTransactionResponse{
				Txid: fromHex(t, "5cfb2ad672e2ac73ff7d8d008bf1e8bb32224279722a5ee562f3d3a8726f277e"),
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
