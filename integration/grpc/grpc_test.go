// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package grpc

import (
	"context"
	"crypto"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	storev1 "github.com/sykesm/batik/pkg/pb/store/v1"
	txv1 "github.com/sykesm/batik/pkg/pb/tx/v1"
	"github.com/sykesm/batik/pkg/tested"
	. "github.com/sykesm/batik/pkg/tested/matcher"
	"github.com/sykesm/batik/pkg/transaction"
)

var _ = Describe("gRPC", func() {
	var (
		session     *gexec.Session
		grpcAddress string
		httpAddress string
		dbPath      string
		cleanup     func()

		clientConn *grpc.ClientConn
	)

	BeforeEach(func() {
		grpcAddress = fmt.Sprintf("127.0.0.1:%d", StartPort())
		httpAddress = fmt.Sprintf("127.0.0.1:%d", StartPort()+1)
		dbPath, cleanup = tested.TempDir(GinkgoT(), "", "level")
		cmd := exec.Command(
			batikPath,
			"--data-dir", dbPath,
			"--color=yes",
			"start",
			"--grpc-listen-address", grpcAddress,
			"--http-listen-address", httpAddress,
		)

		var err error
		session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).To(BeNil())
		Eventually(session.Err, testTimeout).Should(gbytes.Say("Starting server"))
		Eventually(session.Err, testTimeout).Should(gbytes.Say("Server started"))

		clientConn, err = grpc.Dial(grpcAddress, grpc.WithInsecure(), grpc.WithBlock())
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		if session != nil {
			session.Kill().Wait(testTimeout)
		}
		if clientConn != nil {
			clientConn.Close()
		}
		if cleanup != nil {
			cleanup()
		}
	})

	Describe("Encode transaction api", func() {
		It("encodes a transaction", func() {
			testTx := newTestTransaction()

			encodeTransactionClient := txv1.NewEncodeAPIClient(clientConn)
			resp, err := encodeTransactionClient.Encode(context.Background(), &txv1.EncodeRequest{
				Transaction: testTx,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(hex.EncodeToString(resp.Txid)).To(Equal("d213b9cad8f9d0d8316da198d4dd2bee53359c3a6b56a8f4fe4411c504678fa4"))

			expectedEncoded, err := proto.MarshalOptions{Deterministic: true}.Marshal(testTx)
			Expect(resp.EncodedTransaction).To(Equal(expectedEncoded))
		})
	})

	Describe("Submit Transaction API", func() {
		var submitClient txv1.SubmitAPIClient
		var storeClient storev1.StoreAPIClient

		BeforeEach(func() {
			submitClient = txv1.NewSubmitAPIClient(clientConn)
			storeClient = storev1.NewStoreAPIClient(clientConn)
		})

		It("submits a dummy transaction for processing", func() {
			uuid := make([]byte, 16)
			_, err := io.ReadFull(rand.Reader, uuid)
			Expect(err).NotTo(HaveOccurred())

			salt := make([]byte, 32)
			_, err = io.ReadFull(rand.Reader, salt)
			Expect(err).NotTo(HaveOccurred())

			tx := &txv1.Transaction{
				Salt: salt,
				Outputs: []*txv1.State{{
					Info:  &txv1.StateInfo{Kind: "test-state"},
					State: uuid,
				}},
			}

			resp, err := submitClient.Submit(
				context.Background(),
				&txv1.SubmitRequest{
					SignedTransaction: &txv1.SignedTransaction{
						Transaction: tx,
					},
				},
			)
			Expect(err).NotTo(HaveOccurred())

			By("retrieving the transaciton")
			itx, err := transaction.Marshal(crypto.SHA256, tx)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Txid).To(Equal(itx.ID))

			result, err := storeClient.GetTransaction(
				context.Background(),
				&storev1.GetTransactionRequest{Txid: resp.Txid},
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ProtoEqual(&storev1.GetTransactionResponse{Transaction: tx}))

			By("resubmitting the transaciton")
			_, err = submitClient.Submit(
				context.Background(),
				&txv1.SubmitRequest{
					SignedTransaction: &txv1.SignedTransaction{Transaction: tx},
				},
			)
			Expect(err).To(HaveOccurred())
			st, ok := status.FromError(err)
			Expect(ok).To(BeTrue())
			Expect(st.Code()).To(Equal(codes.AlreadyExists))
			Expect(st.Message()).To(ContainSubstring(hex.EncodeToString(itx.ID)))

			By("consuming the output")
			salt = make([]byte, 32)
			_, err = io.ReadFull(rand.Reader, salt)
			Expect(err).NotTo(HaveOccurred())

			tx = &txv1.Transaction{
				Salt: salt,
				Inputs: []*txv1.StateReference{{
					Txid:        itx.ID,
					OutputIndex: 0,
				}},
			}
			resp, err = submitClient.Submit(
				context.Background(),
				&txv1.SubmitRequest{
					SignedTransaction: &txv1.SignedTransaction{Transaction: tx},
				},
			)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Store service api", func() {
		var (
			storeServiceClient storev1.StoreAPIClient
			testTx             *txv1.Transaction
			txid               []byte
		)

		BeforeEach(func() {
			storeServiceClient = storev1.NewStoreAPIClient(clientConn)

			testTx = newTestTransaction()
			intTx, err := transaction.Marshal(crypto.SHA256, testTx)
			Expect(err).NotTo(HaveOccurred())
			txid = intTx.ID
		})

		Describe("GetTransaction", func() {
			var req *storev1.GetTransactionRequest

			BeforeEach(func() {
				req = &storev1.GetTransactionRequest{
					Txid: txid,
				}
			})

			When("the transaction does not exist", func() {
				It("returns an error", func() {
					resp, err := storeServiceClient.GetTransaction(context.Background(), req)
					Expect(err).To(MatchError(ContainSubstring("leveldb: not found")))
					Expect(resp).To(BeNil())
				})
			})

			When("the transaction exists", func() {
				BeforeEach(func() {
					putReq := &storev1.PutTransactionRequest{
						Transaction: testTx,
					}

					_, err := storeServiceClient.PutTransaction(context.Background(), putReq)
					Expect(err).NotTo(HaveOccurred())
				})

				It("retrieves a transaction from the store", func() {
					resp, err := storeServiceClient.GetTransaction(context.Background(), req)
					Expect(err).NotTo(HaveOccurred())
					Expect(resp.Transaction).To(ProtoEqual(testTx))
				})

				// TODO: Reorganize the integration tests
				It("works through the REST gateway", func() {
					client := &http.Client{}
					url := "http://" + httpAddress + "/v1/store/tx/" + base64.URLEncoding.EncodeToString(txid)
					resp, err := client.Get(url)
					Expect(err).NotTo(HaveOccurred())
					defer resp.Body.Close()

					Expect(resp.StatusCode).To(Equal(http.StatusOK))
					body, err := ioutil.ReadAll(resp.Body)
					Expect(err).NotTo(HaveOccurred())

					testJSON, err := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(&storev1.GetTransactionResponse{Transaction: testTx})
					Expect(err).NotTo(HaveOccurred())
					Expect(body).To(MatchJSON(testJSON))
				})
			})
		})

		Describe("PutTransaction", func() {
			var req *storev1.PutTransactionRequest

			BeforeEach(func() {
				req = &storev1.PutTransactionRequest{
					Transaction: testTx,
				}
			})

			// TODO: this test is too similar to the retrieval one, maybe rewrite this and the retrieval one
			// to store and retrieve from the db directly somehow
			It("stores a transaction in the store", func() {
				_, err := storeServiceClient.PutTransaction(context.Background(), req)
				Expect(err).NotTo(HaveOccurred())

				getReq := &storev1.GetTransactionRequest{
					Txid: txid,
				}
				resp, err := storeServiceClient.GetTransaction(context.Background(), getReq)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Transaction).To(ProtoEqual(testTx))
			})
		})

		Describe("GetState", func() {
			var (
				req       *storev1.GetStateRequest
				testState *txv1.ResolvedState
			)

			BeforeEach(func() {
				req = &storev1.GetStateRequest{
					StateRef: &txv1.StateReference{
						Txid:        txid,
						OutputIndex: 0,
					},
				}
				testState = &txv1.ResolvedState{
					Txid:        txid,
					OutputIndex: 0,
					State:       testTx.Outputs[0].State,
					Info:        testTx.Outputs[0].Info,
				}
			})

			When("the state does not exist", func() {
				It("returns an error", func() {
					resp, err := storeServiceClient.GetState(context.Background(), req)
					Expect(err).To(MatchError(ContainSubstring("leveldb: not found")))
					Expect(resp).To(BeNil())
				})
			})

			When("the state exists", func() {
				BeforeEach(func() {
					putReq := &storev1.PutStateRequest{
						State: testState,
					}

					_, err := storeServiceClient.PutState(context.Background(), putReq)
					Expect(err).NotTo(HaveOccurred())
				})

				It("retrieves a state from the store", func() {
					resp, err := storeServiceClient.GetState(context.Background(), req)
					Expect(err).NotTo(HaveOccurred())
					Expect(resp.State).To(ProtoEqual(testState))
				})
			})
		})

		Describe("PutState", func() {
			var (
				req       *storev1.PutStateRequest
				testState *txv1.ResolvedState
			)

			BeforeEach(func() {
				testState = &txv1.ResolvedState{
					Txid:        txid,
					OutputIndex: 0,
					State:       testTx.Outputs[0].State,
					Info:        testTx.Outputs[0].Info,
				}

				req = &storev1.PutStateRequest{
					State: testState,
				}
			})

			// TODO: this test is too similar to the retrieval one, maybe rewrite this and the retrieval one
			// to store and retrieve from the db directly somehow
			It("stores a state in the store", func() {
				_, err := storeServiceClient.PutState(context.Background(), req)
				Expect(err).NotTo(HaveOccurred())

				getReq := &storev1.GetStateRequest{
					StateRef: &txv1.StateReference{
						Txid:        txid,
						OutputIndex: 0,
					},
				}
				resp, err := storeServiceClient.GetState(context.Background(), getReq)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.State).To(ProtoEqual(testState))
			})
		})
	})
})

func newTestTransaction() *txv1.Transaction {
	return &txv1.Transaction{
		Salt: []byte("0123456789abcdef0123456789abcdef"),
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

func fromHex(s string) ([]byte, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("failed to decode %q as hex string", s)
	}

	return b, nil
}
