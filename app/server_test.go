// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"context"
	"io/ioutil"
	"net"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"google.golang.org/grpc"

	"github.com/sykesm/batik/app/internal/testprotos"
	"github.com/sykesm/batik/pkg/log"
)

type emptyServiceServer struct{}

func (ess *emptyServiceServer) EmptyCall(context.Context, *testprotos.Empty) (*testprotos.Empty, error) {
	return new(testprotos.Empty), nil
}

// invoke the EmptyCall RPC
func invokeEmptyCall(address string, dialOptions ...grpc.DialOption) (*testprotos.Empty, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//create GRPC client conn
	clientConn, err := grpc.DialContext(ctx, address, dialOptions...)
	if err != nil {
		return nil, err
	}
	defer clientConn.Close()

	//create GRPC client
	client := testprotos.NewTestServiceClient(clientConn)

	//invoke service
	empty, err := client.EmptyCall(context.Background(), new(testprotos.Empty))
	if err != nil {
		return nil, err
	}

	return empty, nil
}

func TestNewServer(t *testing.T) {
	gt := NewGomegaWithT(t)

	logger, err := log.NewLogger(log.Config{Writer: ioutil.Discard})
	gt.Expect(err).NotTo(HaveOccurred())

	testAddress := "127.0.0.1:9053"
	srv, err := NewServer(
		Config{
			Server: Server{Address: testAddress},
		},
		logger,
	)
	gt.Expect(err).NotTo(HaveOccurred())

	// resolve the address
	addr, err := net.ResolveTCPAddr("tcp", testAddress)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(addr.String()).To(Equal(srv.address))

	gt.Expect(srv.db).NotTo(BeNil())

	// register the GRPC test server
	testprotos.RegisterTestServiceServer(srv.server.Server, &emptyServiceServer{})

	go srv.Start()
	defer srv.Stop()

	// invoke the EmptyCall service
	_, err = invokeEmptyCall(testAddress, grpc.WithInsecure(), grpc.WithBlock())
	gt.Expect(err).NotTo(HaveOccurred())
}
