// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package transaction

import (
	"encoding/json"
	"testing"

	. "github.com/onsi/gomega"
)

func TestResolvedMarshaling(t *testing.T) {
	resolved := Resolved{
		ID: NewID([]byte("transaction-id-100")),
		Inputs: []*State{
			{
				ID: StateID{TxID: ID([]byte("transaction-id-1")), OutputIndex: 1},
				StateInfo: &StateInfo{
					Kind: "dummy-state",
					Owners: []*Party{
						{PublicKey: []byte("owner-1-public-key")},
					},
				},
				Data: []byte("input-0-data"),
			},
			{
				ID: StateID{TxID: ID([]byte("transaction-id-2")), OutputIndex: 2},
				StateInfo: &StateInfo{
					Kind: "dummy-state",
					Owners: []*Party{
						{PublicKey: []byte("owner-2-public-key")},
					},
				},
				Data: []byte("input-1-data"),
			},
		},
		References: []*State{
			{
				ID: StateID{TxID: ID([]byte("transaction-id-3")), OutputIndex: 3},
				StateInfo: &StateInfo{
					Kind: "dummy-reference-state",
					Owners: []*Party{
						{PublicKey: []byte("owner-3-public-key")},
						{PublicKey: []byte("owner-4-public-key")},
					},
				},
				Data: []byte("reference-0-data"),
			},
			{
				ID: StateID{TxID: ID([]byte("transaction-id-4")), OutputIndex: 4},
				StateInfo: &StateInfo{
					Kind: "dummy-reference-state",
					Owners: []*Party{
						{PublicKey: []byte("owner-5-public-key")},
						{PublicKey: []byte("owner-6-public-key")},
					},
				},
				Data: []byte("reference-1-data"),
			},
		},
		Outputs: []*State{
			{
				ID: StateID{TxID: ID([]byte("transaction-id-100")), OutputIndex: 0},
				StateInfo: &StateInfo{
					Kind: "currency-kind",
					Owners: []*Party{
						{PublicKey: []byte("owner-100-public-key")},
					},
				},
				Data: []byte("output-data-0"),
			},
			{
				ID: StateID{TxID: ID([]byte("transaction-id-100")), OutputIndex: 1},
				StateInfo: &StateInfo{
					Kind: "currency-kind",
					Owners: []*Party{
						{PublicKey: []byte("owner-100-public-key")},
					},
				},
				Data: []byte("output-data-1"),
			},
		},
		Parameters: []*Parameter{
			{Name: "operation", Value: []byte("generate-some-cash")},
		},
		RequiredSigners: []*Party{
			{PublicKey: []byte("owner-5-public-key")},
		},
		Signatures: []*Signature{
			{PublicKey: []byte("owner-1-public-key"), Signature: []byte("owner-1-signature")},
			{PublicKey: []byte("owner-2-public-key"), Signature: []byte("owner-2-signature")},
			{PublicKey: []byte("owner-5-public-key"), Signature: []byte("owner-5-signature")},
		},
	}

	gt := NewGomegaWithT(t)
	out, err := json.Marshal(resolved)
	gt.Expect(err).NotTo(HaveOccurred())

	gt.Expect(out).To(MatchJSON(`{
			"id": "7472616e73616374696f6e2d69642d313030",
			"inputs": [
				{
					"id": {
						"txid": "7472616e73616374696f6e2d69642d31",
						"output_index": 1
					},
					"info": {
						"kind": "dummy-state",
						"owners": [
							{
								"public_key": "b3duZXItMS1wdWJsaWMta2V5"
							}
						]
					},
					"data": "aW5wdXQtMC1kYXRh"
				},
				{
					"id": {
						"txid": "7472616e73616374696f6e2d69642d32",
						"output_index": 2
					},
					"info": {
						"kind": "dummy-state",
						"owners": [
							{
								"public_key": "b3duZXItMi1wdWJsaWMta2V5"
							}
						]
					},
					"data": "aW5wdXQtMS1kYXRh"
				}
			],
			"references": [
				{
					"id": {
						"txid": "7472616e73616374696f6e2d69642d33",
						"output_index": 3
					},
					"info": {
						"kind": "dummy-reference-state",
						"owners": [
							{
								"public_key": "b3duZXItMy1wdWJsaWMta2V5"
							},
							{
								"public_key": "b3duZXItNC1wdWJsaWMta2V5"
							}
						]
					},
					"data": "cmVmZXJlbmNlLTAtZGF0YQ=="
				},
				{
					"id": {
						"txid": "7472616e73616374696f6e2d69642d34",
						"output_index": 4
					},
					"info": {
						"kind": "dummy-reference-state",
						"owners": [
							{
								"public_key": "b3duZXItNS1wdWJsaWMta2V5"
							},
							{
								"public_key": "b3duZXItNi1wdWJsaWMta2V5"
							}
						]
					},
					"data": "cmVmZXJlbmNlLTEtZGF0YQ=="
				}
			],
			"outputs": [
				{
					"id": {
						"txid": "7472616e73616374696f6e2d69642d313030",
						"output_index": 0
					},
					"info": {
						"kind": "currency-kind",
						"owners": [
							{
								"public_key": "b3duZXItMTAwLXB1YmxpYy1rZXk="
							}
						]
					},
					"data": "b3V0cHV0LWRhdGEtMA=="
				},
				{
					"id": {
						"txid": "7472616e73616374696f6e2d69642d313030",
						"output_index": 1
					},
					"info": {
						"kind": "currency-kind",
						"owners": [
							{
								"public_key": "b3duZXItMTAwLXB1YmxpYy1rZXk="
							}
						]
					},
					"data": "b3V0cHV0LWRhdGEtMQ=="
				}
			],
			"parameters": [
				{
					"name": "operation",
					"value": "Z2VuZXJhdGUtc29tZS1jYXNo"
				}
			],
			"required_signers": [
				{
					"public_key": "b3duZXItNS1wdWJsaWMta2V5"
				}
			],
			"signatures": [
				{
					"public_key": "b3duZXItMS1wdWJsaWMta2V5",
					"signature": "b3duZXItMS1zaWduYXR1cmU="
				},
				{
					"public_key": "b3duZXItMi1wdWJsaWMta2V5",
					"signature": "b3duZXItMi1zaWduYXR1cmU="
				},
				{
					"public_key": "b3duZXItNS1wdWJsaWMta2V5",
					"signature": "b3duZXItNS1zaWduYXR1cmU="
				}
			]
		}`))

	var r Resolved
	err = json.Unmarshal(out, &r)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(r).To(Equal(resolved))
}
