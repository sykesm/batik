// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package grpc

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"

	"github.com/sykesm/batik/pkg/options"
	storev1 "github.com/sykesm/batik/pkg/pb/store/v1"
	txv1 "github.com/sykesm/batik/pkg/pb/tx/v1"
	"github.com/sykesm/batik/pkg/tested"
	. "github.com/sykesm/batik/pkg/tested/matcher"
	"github.com/sykesm/batik/pkg/transaction"
)

var _ = Describe("gRPC", func() {
	var (
		session        *gexec.Session
		grpcAddress    string
		httpAddress    string
		storagePath    string
		storageCleanup func()

		clientConn *grpc.ClientConn
	)

	BeforeEach(func() {
		grpcAddress = fmt.Sprintf("127.0.0.1:%d", StartPort())
		httpAddress = fmt.Sprintf("127.0.0.1:%d", StartPort()+1)

		storagePath, storageCleanup = tested.TempDir(GinkgoT(), "", "grpc-integration")

		confFilePath := filepath.Join(storagePath, "batik.yaml")
		err := writeNewConfig(confFilePath)
		Expect(err).NotTo(HaveOccurred())

		cmd := exec.Command(
			batikPath,
			"--config", confFilePath,
			"--color=yes",
			"start",
			"--grpc-listen-address", grpcAddress,
			"--http-listen-address", httpAddress,
		)

		session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).To(BeNil())
		Eventually(session.Err, testTimeout).Should(gbytes.Say("Starting server"))
		Eventually(session.Err, testTimeout).Should(gbytes.Say("Server started"))

		creds, err := credentials.NewClientTLSFromFile(filepath.Join(storagePath, "tls-certs", "server-cert.pem"), "")
		Expect(err).NotTo(HaveOccurred())
		clientConn, err = grpc.Dial(grpcAddress, grpc.WithTransportCredentials(creds), grpc.WithBlock())
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		if session != nil {
			session.Kill().Wait(testTimeout)
		}
		if clientConn != nil {
			clientConn.Close()
		}
		if storageCleanup != nil {
			storageCleanup()
		}
	})

	Describe("Encode transaction api", func() {
		It("encodes a transaction", func() {
			testTx := newTestTransaction()
			itx, err := transaction.New(crypto.SHA256, testTx)
			Expect(err).NotTo(HaveOccurred())

			encodeTransactionClient := txv1.NewEncodeAPIClient(clientConn)
			resp, err := encodeTransactionClient.Encode(context.Background(), &txv1.EncodeRequest{
				Transaction: testTx,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.Txid).To(Equal(itx.ID.Bytes()))

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

			for _, namespace := range []string{"ns2", "ns1"} {
				tx := proto.Clone(tx).(*txv1.Transaction)
				By("performing the submit process in " + namespace)
				resp, err := submitClient.Submit(
					context.Background(),
					&txv1.SubmitRequest{
						Namespace: namespace,
						SignedTransaction: &txv1.SignedTransaction{
							Transaction: tx,
						},
					},
				)
				Expect(err).NotTo(HaveOccurred())

				By("retrieving the transaction")
				itx, err := transaction.New(crypto.SHA256, tx)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Txid).To(Equal(itx.ID.Bytes()))

				result, err := storeClient.GetTransaction(
					context.Background(),
					&storev1.GetTransactionRequest{
						Namespace: namespace,
						Txid:      resp.Txid,
					},
				)
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(ProtoEqual(&storev1.GetTransactionResponse{Transaction: tx}))

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
				itx2, err := transaction.New(crypto.SHA256, tx)
				Expect(err).NotTo(HaveOccurred())
				resp, err = submitClient.Submit(
					context.Background(),
					&txv1.SubmitRequest{
						Namespace: namespace,
						SignedTransaction: &txv1.SignedTransaction{
							Transaction: tx,
						},
					},
				)
				Expect(resp.Txid).To(Equal(itx2.ID.Bytes()))
				Expect(err).NotTo(HaveOccurred())

				By("verifying the output is consumed")
				salt = make([]byte, 32)
				_, err = io.ReadFull(rand.Reader, salt)
				Expect(err).NotTo(HaveOccurred())

				tx2 := &txv1.Transaction{
					Salt: salt,
					Inputs: []*txv1.StateReference{{
						Txid:        itx.ID,
						OutputIndex: 0,
					}},
				}
				_, err = submitClient.Submit(
					context.Background(),
					&txv1.SubmitRequest{
						Namespace: namespace,
						SignedTransaction: &txv1.SignedTransaction{
							Transaction: tx2,
						},
					},
				)
				Expect(err).To(HaveOccurred())

				st, ok := status.FromError(err)
				Expect(ok).To(BeTrue())
				Expect(st.Code()).To(Equal(codes.FailedPrecondition))

				itx3, err := transaction.New(crypto.SHA256, tx2)
				Expect(err).NotTo(HaveOccurred())
				Expect(st.Message()).To(ContainSubstring(hex.EncodeToString(itx3.ID)))

				By("fetching the consumed state")
				_, err = storeClient.GetState(
					context.Background(),
					&storev1.GetStateRequest{
						Namespace: namespace,
						StateRef: &txv1.StateReference{
							Txid:        itx.ID,
							OutputIndex: 0,
						},
						Consumed: true,
					},
				)
				Expect(err).NotTo(HaveOccurred())
			}
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
			intTx, err := transaction.New(crypto.SHA256, testTx)
			Expect(err).NotTo(HaveOccurred())
			txid = intTx.ID
		})

		Describe("GetTransaction", func() {
			var req *storev1.GetTransactionRequest

			BeforeEach(func() {
				req = &storev1.GetTransactionRequest{
					Namespace: "ns1",
					Txid:      txid,
				}
			})

			When("the transaction does not exist", func() {
				It("returns an error", func() {
					resp, err := storeServiceClient.GetTransaction(context.Background(), req)
					Expect(err).To(MatchError(ContainSubstring("leveldb: not found")))
					Expect(resp).To(BeNil())
				})
			})

			When("the namespace does not exist", func() {
				BeforeEach(func() {
					req.Namespace = "missing"
				})

				It("returns an error", func() {
					resp, err := storeServiceClient.GetTransaction(context.Background(), req)
					Expect(err).To(MatchError(ContainSubstring("namespace not found")))
					Expect(resp).To(BeNil())
				})
			})

			When("the transaction exists", func() {
				BeforeEach(func() {
					putReq := &storev1.PutTransactionRequest{
						Namespace:   "ns1",
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
					caCertPEM, err := ioutil.ReadFile(filepath.Join(storagePath, "tls-certs", "server-cert.pem"))
					Expect(err).NotTo(HaveOccurred())

					caCertPool := x509.NewCertPool()
					caCertPool.AppendCertsFromPEM(caCertPEM)

					client := &http.Client{
						Transport: &http.Transport{
							TLSClientConfig: &tls.Config{
								RootCAs: caCertPool,
							},
						},
					}
					url := "https://" + httpAddress + "/v1/store/ns1/tx/" + base64.URLEncoding.EncodeToString(txid)
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
					Namespace:   "ns1",
					Transaction: testTx,
				}
			})

			// TODO: this test is too similar to the retrieval one, maybe rewrite this and the retrieval one
			// to store and retrieve from the db directly somehow
			It("stores a transaction in the store", func() {
				_, err := storeServiceClient.PutTransaction(context.Background(), req)
				Expect(err).NotTo(HaveOccurred())

				getReq := &storev1.GetTransactionRequest{
					Namespace: "ns1",
					Txid:      txid,
				}
				resp, err := storeServiceClient.GetTransaction(context.Background(), getReq)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Transaction).To(ProtoEqual(testTx))
			})
		})

		Describe("GetState", func() {
			var (
				req      *storev1.GetStateRequest
				stateRef *txv1.StateReference
				state    *txv1.State
			)

			BeforeEach(func() {
				state = testTx.Outputs[0]
				stateRef = &txv1.StateReference{
					Txid:        txid,
					OutputIndex: 0,
				}
				req = &storev1.GetStateRequest{
					Namespace: "ns1",
					StateRef:  stateRef,
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
						Namespace: "ns1",
						StateRef:  stateRef,
						State:     state,
					}

					_, err := storeServiceClient.PutState(context.Background(), putReq)
					Expect(err).NotTo(HaveOccurred())
				})

				It("retrieves a state from the store", func() {
					resp, err := storeServiceClient.GetState(context.Background(), req)
					Expect(err).NotTo(HaveOccurred())
					Expect(resp.State).To(ProtoEqual(state))
				})
			})
		})

		Describe("PutState", func() {
			var (
				req      *storev1.PutStateRequest
				stateRef *txv1.StateReference
				state    *txv1.State
			)

			BeforeEach(func() {
				state = testTx.Outputs[0]
				stateRef = &txv1.StateReference{
					Txid:        txid,
					OutputIndex: 0,
				}

				req = &storev1.PutStateRequest{
					Namespace: "ns1",
					StateRef:  stateRef,
					State:     state,
				}
			})

			// TODO: this test is too similar to the retrieval one, maybe rewrite this and the retrieval one
			// to store and retrieve from the db directly somehow
			It("stores a state in the store", func() {
				_, err := storeServiceClient.PutState(context.Background(), req)
				Expect(err).NotTo(HaveOccurred())

				getReq := &storev1.GetStateRequest{
					Namespace: "ns1",
					StateRef:  stateRef,
				}
				resp, err := storeServiceClient.GetState(context.Background(), getReq)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.State).To(ProtoEqual(state))
			})
		})
	})
})

// writeNewConfig creates a config file with a single statically
// defined namespace inside it.
func writeNewConfig(path string) error {
	config := options.Batik{
		Namespaces: []options.Namespace{
			{
				Name:      "ns1",
				Validator: "signature-builtin",
			},
			{
				Name:      "ns2",
				Validator: "signature-wasm",
			},
		},
		Validators: []options.Validator{
			{
				Name: "signature-builtin",
				Type: "builtin",
			},
			{
				Name: "signature-wasm",
				Type: "wasm",
				Path: wasmSigValidatorPath,
			},
		},
	}

	confFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer confFile.Close()

	encoder := yaml.NewEncoder(confFile)
	err = encoder.Encode(config)
	if err != nil {
		return err
	}
	encoder.Close()
	return nil
}

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

func fromHex(s string) ([]byte, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("failed to decode %q as hex string", s)
	}

	return b, nil
}
