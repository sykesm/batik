// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package grpc

import (
	"context"
	"encoding/hex"
	"fmt"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	tb "github.com/sykesm/batik/pkg/pb/transaction"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

var _ = Describe("Grpc", func() {
	var (
		cmd     *exec.Cmd
		session *gexec.Session
		address string
	)

	BeforeEach(func() {
		var err error

		address = fmt.Sprintf("127.0.0.1:%d", StartPort())

		cmd = exec.Command(batikPath, "start", "-a", address)
		session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).To(BeNil())
	})

	It("starts a grpc server", func() {
		Eventually(session, testTimeout).Should(gbytes.Say(fmt.Sprintf("Starting server at %s", address)))
	})

	Describe("Encode transaction api", func() {
		It("encodes a transaction", func() {
			testTx := newTestTransaction()

			clientConn, err := grpc.Dial(address, grpc.WithInsecure())
			Expect(err).NotTo(HaveOccurred())

			encodeTransactionClient := tb.NewEncodeTransactionAPIClient(clientConn)

			req := &tb.EncodeTransactionRequest{
				Transaction: testTx,
			}

			resp, err := encodeTransactionClient.EncodeTransaction(context.Background(), req)
			Expect(err).NotTo(HaveOccurred())

			expectedTxID, err := fromHex("53e33ae87fb6cf2e4aaaabcdae3a93d578d9b7366e905dfff0446356774f726f")
			Expect(err).NotTo(HaveOccurred())

			Expect(resp.Txid).To(Equal(expectedTxID))

			expectedEncoded, err := proto.MarshalOptions{Deterministic: true}.Marshal(testTx)
			Expect(resp.EncodedTransaction).To(Equal(expectedEncoded))
		})
	})
})

func newTestTransaction() *tb.Transaction {
	return &tb.Transaction{
		Inputs: []*tb.StateReference{
			{Txid: []byte("input-transaction-id-0"), OutputIndex: 1},
			{Txid: []byte("input-transaction-id-1"), OutputIndex: 0},
		},
		References: []*tb.StateReference{
			{Txid: []byte("ref-transaction-id-0"), OutputIndex: 1},
			{Txid: []byte("ref-transaction-id-1"), OutputIndex: 0},
		},
		Outputs: []*tb.State{
			{
				Info: &tb.StateInfo{
					Owners: []*tb.Party{{Credential: []byte("owner-1")}, {Credential: []byte("owner-2")}},
					Kind:   "state-kind-0",
				},
				State: []byte("state-0"),
			},
			{
				Info: &tb.StateInfo{
					Owners: []*tb.Party{{Credential: []byte("owner-1")}, {Credential: []byte("owner-2")}},
					Kind:   "state-kind-1",
				},
				State: []byte("state-1"),
			},
		},
		Parameters: []*tb.Parameter{
			{Name: "name-0", Value: []byte("value-0")},
			{Name: "name-1", Value: []byte("value-1")},
		},
		RequiredSigners: []*tb.Party{
			{Credential: []byte("observer-1")},
			{Credential: []byte("observer-2")},
		},
		Salt: []byte("NaCl"),
	}
}

func fromHex(s string) ([]byte, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("failed to decode %q as hex string", s)
	}

	return b, nil
}
