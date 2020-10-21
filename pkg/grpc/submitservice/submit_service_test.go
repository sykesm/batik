// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package submitservice

import (
	"context"
	"encoding/hex"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"

	txv1 "github.com/sykesm/batik/pkg/pb/tx/v1"
	"github.com/sykesm/batik/pkg/store"
	"github.com/sykesm/batik/pkg/tested"
	. "github.com/sykesm/batik/pkg/tested/matcher"
)

func TestSubmit(t *testing.T) {
	tests := map[string]struct {
		req        *txv1.SubmitRequest
		resp       *txv1.SubmitResponse
		errMatcher types.GomegaMatcher
	}{
		"nil transaction": {
			req:        &txv1.SubmitRequest{},
			resp:       nil,
			errMatcher: HaveOccurred(),
		},
		"valid transaction": {
			req: &txv1.SubmitRequest{
				SignedTransaction: &txv1.SignedTransaction{
					Transaction: &txv1.Transaction{
						Salt: []byte("potassium permanganate (KMnO4) is a salt"),
						Outputs: []*txv1.State{{
							Info:  &txv1.StateInfo{Kind: "test-kind"},
							State: []byte("test-state-1"),
						}},
					},
				},
			},
			resp: &txv1.SubmitResponse{
				Txid: fromHex(t, "5cfb2ad672e2ac73ff7d8d008bf1e8bb32224279722a5ee562f3d3a8726f277e"),
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			kv, err := store.NewLevelDB("")
			gt.Expect(err).NotTo(HaveOccurred())
			defer tested.Close(t, kv)

			ss := NewSubmitService(kv)
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

func fromHex(t *testing.T, s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatalf("failed to decode %q as hex string", s)
	}
	return b
}
