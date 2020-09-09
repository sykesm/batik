// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package grpc

import (
	"context"
	"crypto"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	sb "github.com/sykesm/batik/pkg/pb/store"
	tb "github.com/sykesm/batik/pkg/pb/transaction"
	"github.com/sykesm/batik/pkg/tested"
	"github.com/sykesm/batik/pkg/transaction"
)

var _ = Describe("Grpc", func() {
	var (
		session *gexec.Session
		address string
		cleanup func()
	)

	BeforeEach(func() {
		address = fmt.Sprintf("127.0.0.1:%d", StartPort())

		var dbPath string
		dbPath, cleanup = tested.TempDir(GinkgoT(), "", "level")

		cmd := exec.Command(batikPath, "start", "-a", address)
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "DB_PATH="+dbPath)

		var err error
		session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).To(BeNil())
		Eventually(session.Err, testTimeout).Should(gbytes.Say("Starting server"))
		Eventually(session.Err, testTimeout).Should(gbytes.Say("Server started"))
	})

	AfterEach(func() {
		if session != nil {
			session.Kill().Wait(testTimeout)
		}

		cleanup()
	})

	Describe("Encode transaction api", func() {
		It("encodes a transaction", func() {
			testTx := newTestTransaction()

			clientConn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
			Expect(err).NotTo(HaveOccurred())

			encodeTransactionClient := tb.NewEncodeTransactionAPIClient(clientConn)
			resp, err := encodeTransactionClient.EncodeTransaction(context.Background(), &tb.EncodeTransactionRequest{
				Transaction: testTx,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(hex.EncodeToString(resp.Txid)).To(Equal("53e33ae87fb6cf2e4aaaabcdae3a93d578d9b7366e905dfff0446356774f726f"))

			expectedEncoded, err := proto.MarshalOptions{Deterministic: true}.Marshal(testTx)
			Expect(resp.EncodedTransaction).To(Equal(expectedEncoded))
		})
	})

	Describe("Store service api", func() {
		var (
			storeServiceClient sb.StoreAPIClient
			testTx             *tb.Transaction
			txid               []byte
		)

		BeforeEach(func() {
			clientConn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
			Expect(err).NotTo(HaveOccurred())

			storeServiceClient = sb.NewStoreAPIClient(clientConn)

			testTx = newTestTransaction()
			txid, err = transaction.ID(crypto.SHA256, testTx)
			Expect(err).NotTo(HaveOccurred())
		})

		Describe("GetTransaction", func() {
			var req *sb.GetTransactionRequest

			BeforeEach(func() {
				req = &sb.GetTransactionRequest{
					Txid: txid,
				}
			})

			When("the transaction does not exist", func() {
				It("returns an error", func() {
					resp, err := storeServiceClient.GetTransaction(context.Background(), req)
					Expect(err).To(MatchError(MatchRegexp("leveldb: not found")))
					Expect(resp).To(BeNil())
				})
			})

			When("the transaction exists", func() {
				BeforeEach(func() {
					putReq := &sb.PutTransactionRequest{
						Transaction: testTx,
					}

					_, err := storeServiceClient.PutTransaction(context.Background(), putReq)
					Expect(err).NotTo(HaveOccurred())
				})

				It("retrieves a transaction from the store", func() {
					resp, err := storeServiceClient.GetTransaction(context.Background(), req)
					Expect(err).NotTo(HaveOccurred())
					Expect(proto.Equal(resp.Transaction, testTx)).To(BeTrue())
				})
			})
		})

		Describe("PutTransaction", func() {
			var req *sb.PutTransactionRequest

			BeforeEach(func() {
				req = &sb.PutTransactionRequest{
					Transaction: testTx,
				}
			})

			// TODO: this test is too similar to the retrieval one, maybe rewrite this and the retrieval one
			// to store and retrieve from the db directly somehow
			It("stores a transaction in the store", func() {
				_, err := storeServiceClient.PutTransaction(context.Background(), req)
				Expect(err).NotTo(HaveOccurred())

				getReq := &sb.GetTransactionRequest{
					Txid: txid,
				}
				resp, err := storeServiceClient.GetTransaction(context.Background(), getReq)
				Expect(err).NotTo(HaveOccurred())
				Expect(proto.Equal(resp.Transaction, testTx)).To(BeTrue())
			})
		})

		Describe("GetState", func() {
			var (
				req       *sb.GetStateRequest
				testState *tb.ResolvedState
			)

			BeforeEach(func() {
				req = &sb.GetStateRequest{
					StateRef: &tb.StateReference{
						Txid:        txid,
						OutputIndex: 0,
					},
				}
				testState = &tb.ResolvedState{
					Txid:        txid,
					OutputIndex: 0,
					State:       testTx.Outputs[0].State,
					Info:        testTx.Outputs[0].Info,
				}
			})

			When("the state does not exist", func() {
				It("returns an error", func() {
					resp, err := storeServiceClient.GetState(context.Background(), req)
					Expect(err).To(MatchError(MatchRegexp("leveldb: not found")))
					Expect(resp).To(BeNil())
				})
			})

			When("the state exists", func() {
				BeforeEach(func() {
					putReq := &sb.PutStateRequest{
						State: testState,
					}

					_, err := storeServiceClient.PutState(context.Background(), putReq)
					Expect(err).NotTo(HaveOccurred())
				})

				It("retrieves a state from the store", func() {
					resp, err := storeServiceClient.GetState(context.Background(), req)
					Expect(err).NotTo(HaveOccurred())
					Expect(proto.Equal(resp.State, testState)).To(BeTrue())
				})
			})
		})

		Describe("PutState", func() {
			var (
				req       *sb.PutStateRequest
				testState *tb.ResolvedState
			)

			BeforeEach(func() {
				testState = &tb.ResolvedState{
					Txid:        txid,
					OutputIndex: 0,
					State:       testTx.Outputs[0].State,
					Info:        testTx.Outputs[0].Info,
				}

				req = &sb.PutStateRequest{
					State: testState,
				}
			})

			// TODO: this test is too similar to the retrieval one, maybe rewrite this and the retrieval one
			// to store and retrieve from the db directly somehow
			It("stores a state in the store", func() {
				_, err := storeServiceClient.PutState(context.Background(), req)
				Expect(err).NotTo(HaveOccurred())

				getReq := &sb.GetStateRequest{
					StateRef: &tb.StateReference{
						Txid:        txid,
						OutputIndex: 0,
					},
				}
				resp, err := storeServiceClient.GetState(context.Background(), getReq)
				Expect(err).NotTo(HaveOccurred())
				Expect(proto.Equal(resp.State, testState)).To(BeTrue())
			})
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
