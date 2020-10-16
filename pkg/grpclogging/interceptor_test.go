// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package grpclogging

import (
	"context"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"

	"github.com/sykesm/batik/pkg/grpclogging/internal/testprotos/echo"
	"github.com/sykesm/batik/pkg/tested"
	. "github.com/sykesm/batik/pkg/tested/matcher"
)

func TestUnaryServerInterceptor(t *testing.T) {
	gt := NewGomegaWithT(t)
	tc := newTestContext(t)

	server := grpc.NewServer(
		grpc.Creds(credentials.NewTLS(tc.serverKeyPair.ServerTLSConfig(t, nil))),
		grpc.UnaryInterceptor(UnaryServerInterceptor(tc.logger)),
	)
	echoServer := &echoServiceServer{}
	echo.RegisterEchoServiceServer(server, echoServer)
	serveCompleteCh := tc.serve(server)

	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(tc.serverCA.TLSConfig(t))),
		grpc.WithBlock(),
	}
	clientConn := tc.dial(t, dialOpts...)
	defer clientConn.Close()

	echoServiceClient := echo.NewEchoServiceClient(clientConn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	resp, err := echoServiceClient.Echo(ctx, &echo.Message{Message: "hi"})
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(resp).To(EqualProto(&echo.Message{Message: "hi", Sequence: 1}))

	tc.listener.Close()
	gt.Eventually(serveCompleteCh).Should(Receive())

	t.Run("DecoratedContext", func(t *testing.T) {
		testDecoratedContext(t, echoServer.context())
	})

	t.Run("MessagesAndFields", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		var logMessages []string
		for _, entry := range tc.observed.AllUntimed() {
			logMessages = append(logMessages, entry.Message)
			testEntryFields(t, "Echo", echoServer.context(), entry)
		}
		gt.Expect(logMessages).To(ConsistOf(
			"received unary request", // received payload
			"sending unary response", // sending payload
			"unary call completed",
		))
	})
}

func TestStreamServerInterceptor(t *testing.T) {
	gt := NewGomegaWithT(t)
	tc := newTestContext(t)

	server := grpc.NewServer(
		grpc.Creds(credentials.NewTLS(tc.serverKeyPair.ServerTLSConfig(t, &tc.clientCA.Certificate))),
		grpc.StreamInterceptor(StreamServerInterceptor(tc.logger)),
	)
	echoServer := &echoServiceServer{}
	echo.RegisterEchoServiceServer(server, echoServer)
	serveCompleteCh := tc.serve(server)

	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(tc.serverCA.TLSConfig(t))),
		grpc.WithBlock(),
	}
	clientConn := tc.dial(t, dialOpts...)
	defer clientConn.Close()

	echoServiceClient := echo.NewEchoServiceClient(clientConn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	streamClient, err := echoServiceClient.EchoStream(ctx)
	gt.Expect(err).NotTo(HaveOccurred())

	err = streamClient.Send(&echo.Message{Message: "hi", Sequence: 1})
	gt.Expect(err).NotTo(HaveOccurred())
	msg, err := streamClient.Recv()
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(msg).To(EqualProto(&echo.Message{Message: "hi", Sequence: 2}))

	err = streamClient.CloseSend()
	gt.Expect(err).NotTo(HaveOccurred())
	_, err = streamClient.Recv()
	gt.Expect(err).To(Equal(io.EOF))

	tc.listener.Close()
	gt.Eventually(serveCompleteCh).Should(Receive())

	t.Run("DecoratedContext", func(t *testing.T) {
		testDecoratedContext(t, echoServer.context())
	})

	t.Run("MessagesAndFields", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		var logMessages []string
		for _, entry := range tc.observed.AllUntimed() {
			logMessages = append(logMessages, entry.Message)
			testEntryFields(t, "EchoStream", echoServer.context(), entry)
		}
		gt.Expect(logMessages).To(ConsistOf(
			"received stream message", // received payload
			"sending stream message",  // sending payload
			"streaming call completed",
		))
	})
}

func TestWithClientAuth(t *testing.T) {
	gt := NewGomegaWithT(t)
	tc := newTestContext(t)

	server := grpc.NewServer(
		grpc.Creds(credentials.NewTLS(tc.serverKeyPair.ServerTLSConfig(t, &tc.clientCA.Certificate))),
		grpc.UnaryInterceptor(UnaryServerInterceptor(tc.logger)),
		grpc.StreamInterceptor(StreamServerInterceptor(tc.logger)),
	)
	echoServer := &echoServiceServer{}
	echo.RegisterEchoServiceServer(server, echoServer)
	serveCompleteCh := tc.serve(server)

	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(tc.clientKeyPair.ClientTLSConfig(t, &tc.serverCA.Certificate))),
		grpc.WithBlock(),
	}
	clientConn := tc.dial(t, dialOpts...)
	defer clientConn.Close()

	echoServiceClient := echo.NewEchoServiceClient(clientConn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	resp, err := echoServiceClient.Echo(ctx, &echo.Message{Message: "hi"})
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(resp).To(EqualProto(&echo.Message{Message: "hi", Sequence: 1}))

	t.Run("DecoratedContext", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		gt.Expect(Fields(echoServer.context())).To(ContainElement(
			zap.String("grpc.peer_subject", "CN=client"),
		))
	})

	t.Run("Fields", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		gt.Expect(tc.observed.AllUntimed()).NotTo(HaveLen(0))
		for _, entry := range tc.observed.AllUntimed() {
			gt.Expect(entry.ContextMap()).To(HaveKey("grpc.peer_subject"))
			gt.Expect(entry.ContextMap()["grpc.peer_subject"]).To(HavePrefix("CN=client"))
		}
	})

	tc.listener.Close()
	gt.Eventually(serveCompleteCh).Should(Receive())
}

func TestOptions(t *testing.T) {
	gt := NewGomegaWithT(t)
	tc := newTestContext(t)

	var (
		unaryLeveler LevelerFunc = func(ctx context.Context, fullMethod string) zapcore.Level {
			if fullMethod == "/echo.EchoService/Echo" {
				return zapcore.Level(10)
			}
			return zapcore.ErrorLevel
		}
		unaryPayloadLeveler LevelerFunc = func(ctx context.Context, fullMethod string) zapcore.Level {
			if fullMethod == "/echo.EchoService/Echo" {
				return zapcore.Level(20)
			}
			return zapcore.ErrorLevel
		}
		streamLeveler LevelerFunc = func(ctx context.Context, fullMethod string) zapcore.Level {
			if fullMethod == "/echo.EchoService/EchoStream" {
				return zapcore.Level(30)
			}
			return zapcore.ErrorLevel
		}
		streamPayloadLeveler LevelerFunc = func(ctx context.Context, fullMethod string) zapcore.Level {
			if fullMethod == "/echo.EchoService/EchoStream" {
				return zapcore.Level(40)
			}
			return zapcore.ErrorLevel
		}
	)

	server := grpc.NewServer(
		grpc.Creds(credentials.NewTLS(tc.serverKeyPair.ServerTLSConfig(t, nil))),
		grpc.UnaryInterceptor(
			UnaryServerInterceptor(tc.logger, WithLeveler(unaryLeveler), WithPayloadLeveler(unaryPayloadLeveler)),
		),
		grpc.StreamInterceptor(
			StreamServerInterceptor(tc.logger, WithLeveler(streamLeveler), WithPayloadLeveler(streamPayloadLeveler)),
		),
	)
	echoServer := &echoServiceServer{}
	echo.RegisterEchoServiceServer(server, echoServer)
	serveCompleteCh := tc.serve(server)

	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(tc.serverCA.TLSConfig(t))),
		grpc.WithBlock(),
	}
	clientConn := tc.dial(t, dialOpts...)
	defer clientConn.Close()

	echoServiceClient := echo.NewEchoServiceClient(clientConn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	resp, err := echoServiceClient.Echo(ctx, &echo.Message{Message: "hi"})
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(resp).To(EqualProto(&echo.Message{Message: "hi", Sequence: 1}))

	streamClient, err := echoServiceClient.EchoStream(ctx)
	gt.Expect(err).NotTo(HaveOccurred())

	err = streamClient.Send(&echo.Message{Message: "hi", Sequence: 1})
	gt.Expect(err).NotTo(HaveOccurred())
	msg, err := streamClient.Recv()
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(msg).To(EqualProto(&echo.Message{Message: "hi", Sequence: 2}))

	err = streamClient.CloseSend()
	gt.Expect(err).NotTo(HaveOccurred())
	_, err = streamClient.Recv()
	gt.Expect(err).To(Equal(io.EOF))

	gt.Expect(tc.observed.AllUntimed()).To(HaveLen(6))
	for _, entry := range tc.observed.AllUntimed() {
		contextMap := entry.ContextMap()
		gt.Expect(contextMap).To(HaveKey("grpc.method"))
		method := contextMap["grpc.method"].(string)
		switch {
		case entry.LoggerName == "test-logger" && method == "Echo":
			gt.Expect(entry.Level).To(Equal(zapcore.Level(10)))
		case entry.LoggerName == "test-logger.payload" && method == "Echo":
			gt.Expect(entry.Level).To(Equal(zapcore.Level(20)))
		case entry.LoggerName == "test-logger" && method == "EchoStream":
			gt.Expect(entry.Level).To(Equal(zapcore.Level(30)))
		case entry.LoggerName == "test-logger.payload" && method == "EchoStream":
			gt.Expect(entry.Level).To(Equal(zapcore.Level(40)))
		default:
			t.Fatalf("unexpected log entry, name: %s, method: %s, level: %d", entry.LoggerName, method, entry.Level)
		}
	}

	tc.listener.Close()
	gt.Eventually(serveCompleteCh).Should(Receive())
}

func testDecoratedContext(t *testing.T, ctx context.Context) {
	gt := NewGomegaWithT(t)

	zapFields := Fields(ctx)
	keyNames := []string{}
	for _, field := range zapFields {
		keyNames = append(keyNames, field.Key)
	}
	gt.Expect(keyNames).To(ContainElements(
		"grpc.service",
		"grpc.method",
		"grpc.request_deadline",
		"grpc.peer_address",
	))
}

func testEntryFields(t *testing.T, method string, ctx context.Context, entry observer.LoggedEntry) {
	gt := NewGomegaWithT(t)

	fieldKeys := map[string]struct{}{}
	for _, field := range entry.Context {
		fieldKeys[field.Key] = struct{}{}
	}

	switch entry.LoggerName {
	case "test-logger":
		gt.Expect(entry.Level).To(Equal(zapcore.InfoLevel))
		gt.Expect(fieldKeys).To(HaveKey("grpc.call_duration"))
		gt.Expect(fieldKeys).To(HaveKey("grpc.code"))
		gt.Expect(fieldKeys).To(HaveKey("grpc.method"))
		gt.Expect(fieldKeys).To(HaveKey("grpc.peer_address"))
		gt.Expect(fieldKeys).To(HaveKey("grpc.request_deadline"))
		gt.Expect(fieldKeys).To(HaveKey("grpc.service"))
	case "test-logger.payload":
		gt.Expect(entry.Level).To(Equal(zapcore.DebugLevel - 1))
		gt.Expect(fieldKeys).To(HaveKey("grpc.method"))
		gt.Expect(fieldKeys).To(HaveKey("grpc.peer_address"))
		gt.Expect(fieldKeys).To(HaveKey("grpc.request_deadline"))
		gt.Expect(fieldKeys).To(HaveKey("grpc.service"))
		gt.Expect(fieldKeys).To(HaveKey("message"))
	default:
		t.Fatalf("unexpected logger name: %s", entry.LoggerName)
	}
	gt.Expect(entry.Caller.String()).To(ContainSubstring("grpclogging/interceptor.go"))

	for _, field := range entry.Context {
		switch field.Key {
		case "grpc.code":
			gt.Expect(field.Type).To(Equal(zapcore.StringerType))
			gt.Expect(field.Interface).To(Equal(codes.OK))
		case "grpc.call_duration":
			gt.Expect(field.Type).To(Equal(zapcore.DurationType))
			gt.Expect(field.Integer).NotTo(BeZero())
		case "grpc.service":
			gt.Expect(field.Type).To(Equal(zapcore.StringType))
			gt.Expect(field.String).To(Equal("echo.EchoService"))
		case "grpc.method":
			gt.Expect(field.Type).To(Equal(zapcore.StringType))
			gt.Expect(field.String).To(Equal(method))
		case "grpc.request_deadline":
			deadline, ok := ctx.Deadline()
			gt.Expect(ok).To(BeTrue())
			gt.Expect(field.Type).To(Equal(zapcore.TimeType))
			gt.Expect(field.Integer).NotTo(BeZero())
			gt.Expect(time.Unix(0, field.Integer)).To(BeTemporally("==", deadline))
		case "grpc.peer_address":
			gt.Expect(field.Type).To(Equal(zapcore.StringType))
			gt.Expect(field.String).To(HavePrefix("127.0.0.1"))
		case "message":
			gt.Expect(field.Type).To(Equal(zapcore.ReflectType))
		case "error":
			gt.Expect(field.Type).To(Equal(zapcore.ErrorType))
		case "":
			gt.Expect(field.Type).To(Equal(zapcore.SkipType))
		default:
			t.Fatalf("unexpected context field: %s", field.Key)
		}
	}
}

// testContext aggregates the assets necessary to configure and run a TLS
// secured gRPC interaction in test.
type testContext struct {
	listener      net.Listener
	serverCA      tested.CA
	serverKeyPair tested.CertKeyPair
	clientCA      tested.CA
	clientKeyPair tested.CertKeyPair

	observed *observer.ObservedLogs
	logger   *zap.Logger
}

func newTestContext(t *testing.T) *testContext {
	gt := NewGomegaWithT(t)

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	gt.Expect(err).NotTo(HaveOccurred())

	sca := tested.NewCA(t, "test-server-ca")
	skp := sca.IssueServerCertificate(t, "server", "127.0.0.1")
	cca := tested.NewCA(t, "test-client-ca")
	ckp := cca.IssueClientCertificate(t, "client", "127.0.0.1")

	core, observed := observer.New(zapcore.Level(zapcore.DebugLevel - 1))
	logger := zap.New(core, zap.AddCaller()).Named("test-logger")

	return &testContext{
		listener:      lis,
		serverCA:      sca,
		serverKeyPair: skp,
		clientCA:      cca,
		clientKeyPair: ckp,

		observed: observed,
		logger:   logger,
	}
}

func (tc *testContext) serve(s *grpc.Server) <-chan error {
	serveCompleteCh := make(chan error, 1)
	go func() { serveCompleteCh <- s.Serve(tc.listener) }()
	return serveCompleteCh
}

func (tc *testContext) dial(t *testing.T, dialOpts ...grpc.DialOption) *grpc.ClientConn {
	gt := NewGomegaWithT(t)
	dialCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	clientConn, err := grpc.DialContext(dialCtx, tc.listener.Addr().String(), dialOpts...)
	gt.Expect(err).NotTo(HaveOccurred())
	return clientConn
}

// Trivial implementation of the echo.EchoServiceServer
type echoServiceServer struct {
	echo.UnimplementedEchoServiceServer

	lock sync.Mutex
	ctx  context.Context
}

func (e *echoServiceServer) setContext(ctx context.Context) {
	e.lock.Lock()
	defer e.lock.Unlock()
	e.ctx = ctx
}

func (e *echoServiceServer) context() context.Context {
	e.lock.Lock()
	defer e.lock.Unlock()
	return e.ctx
}

func (e *echoServiceServer) Echo(ctx context.Context, msg *echo.Message) (*echo.Message, error) {
	e.setContext(ctx)
	msg.Sequence++
	return msg, nil
}

func (e *echoServiceServer) EchoStream(stream echo.EchoService_EchoStreamServer) error {
	e.setContext(stream.Context())

	msg, err := stream.Recv()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}

	msg.Sequence++
	return stream.Send(msg)
}
