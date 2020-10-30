// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package encodeservice

import (
	"context"
	"crypto"
	"testing"

	. "github.com/onsi/gomega"

	txv1 "github.com/sykesm/batik/pkg/pb/tx/v1"
	"github.com/sykesm/batik/pkg/protomsg"
	"github.com/sykesm/batik/pkg/transaction"
)

func TestEncode(t *testing.T) {
	gt := NewGomegaWithT(t)

	testTx := newTestTransaction()
	itx, err := transaction.New(crypto.SHA256, testTx)
	gt.Expect(err).NotTo(HaveOccurred())
	req := &txv1.EncodeRequest{Transaction: testTx}

	encodeSvc := &EncodeService{}
	resp, err := encodeSvc.Encode(context.Background(), req)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(resp.Txid).To(Equal(itx.ID.Bytes()))

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
