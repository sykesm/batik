// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package store

import (
	"encoding/hex"
	"testing"

	. "github.com/onsi/gomega"

	txv1 "github.com/sykesm/batik/pkg/pb/tx/v1"
)

func TestPending(t *testing.T) {
	gt := NewGomegaWithT(t)
	gt.Expect(true).To(BeTrue())
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
					Owners: []*txv1.Party{{Credential: []byte("owner-1")}, {Credential: []byte("owner-2")}},
					Kind:   "state-kind-0",
				},
				State: []byte("state-0"),
			},
			{
				Info: &txv1.StateInfo{
					Owners: []*txv1.Party{{Credential: []byte("owner-1")}, {Credential: []byte("owner-2")}},
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
			{Credential: []byte("observer-1")},
			{Credential: []byte("observer-2")},
		},
	}
}

func fromHex(t *testing.T, s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatalf("failed to decode %q as hex string", s)
	}
	return b
}
