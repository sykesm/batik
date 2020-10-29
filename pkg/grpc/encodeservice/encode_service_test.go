// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package encodeservice

import (
	"context"
	"fmt"
	"testing"

	. "github.com/onsi/gomega"

	txv1 "github.com/sykesm/batik/pkg/pb/tx/v1"
	"github.com/sykesm/batik/pkg/protomsg"
)

func TestEncode(t *testing.T) {
	gt := NewGomegaWithT(t)

	testTx := newTestTransaction()
	req := &txv1.EncodeRequest{Transaction: testTx}

	encodeSvc := &EncodeService{}
	resp, err := encodeSvc.Encode(context.Background(), req)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(fmt.Sprintf("%x", resp.Txid)).To(Equal("77dc6e1729583cf7f1db9863b34a8951a3bb9369ab4cf0a86340ea92a8514cf5"))

	expectedEncoded, err := protomsg.MarshalDeterministic(testTx)
	gt.Expect(resp.EncodedTransaction).To(Equal(expectedEncoded))
}

func newTestTransaction() *txv1.Transaction {
	return &txv1.Transaction{
		Salt: []byte("NaCl - abcdefghijklmnopqrstuvwxyz"),
		Inputs: []*txv1.StateReference{
			{Txid: []byte("input-transaction-id-0"), OutputIndex: 1},
			{Txid: []byte("input-transaction-id-1"), OutputIndex: 0},
		},
		References: []*txv1.StateReference{
			{Txid: []byte("ref-transaction-id-0"), OutputIndex: 1},
			{Txid: []byte("ref-transaction-id-1"), OutputIndex: 0},
		},
		Outputs: []*txv1.State{
			{
				Info: &txv1.StateInfo{
					Owners: []*txv1.Party{{PublicKey: []byte("owner-1")}, {PublicKey: []byte("owner-2")}},
					Kind:   "state-kind-0",
				},
				State: []byte("state-0"),
			},
			{
				Info: &txv1.StateInfo{
					Owners: []*txv1.Party{{PublicKey: []byte("owner-1")}, {PublicKey: []byte("owner-2")}},
					Kind:   "state-kind-1",
				},
				State: []byte("state-1"),
			},
		},
		Parameters: []*txv1.Parameter{
			{Name: "name-0", Value: []byte("value-0")},
			{Name: "name-1", Value: []byte("value-1")},
		},
		RequiredSigners: []*txv1.Party{
			{PublicKey: []byte("observer-1")},
			{PublicKey: []byte("observer-2")},
		},
	}
}
