// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package txv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// EncodeAPIClient is the client API for EncodeAPI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EncodeAPIClient interface {
	// Encode encodes a transaction via deterministic marshal and returns the
	// encoded bytes as well as a hash over the transaction represented as a
	// merkle root and generated via SHA256 as the internal hashing function.
	Encode(ctx context.Context, in *EncodeRequest, opts ...grpc.CallOption) (*EncodeResponse, error)
}

type encodeAPIClient struct {
	cc grpc.ClientConnInterface
}

func NewEncodeAPIClient(cc grpc.ClientConnInterface) EncodeAPIClient {
	return &encodeAPIClient{cc}
}

func (c *encodeAPIClient) Encode(ctx context.Context, in *EncodeRequest, opts ...grpc.CallOption) (*EncodeResponse, error) {
	out := new(EncodeResponse)
	err := c.cc.Invoke(ctx, "/tx.v1.EncodeAPI/Encode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EncodeAPIServer is the server API for EncodeAPI service.
// All implementations must embed UnimplementedEncodeAPIServer
// for forward compatibility
type EncodeAPIServer interface {
	// Encode encodes a transaction via deterministic marshal and returns the
	// encoded bytes as well as a hash over the transaction represented as a
	// merkle root and generated via SHA256 as the internal hashing function.
	Encode(context.Context, *EncodeRequest) (*EncodeResponse, error)
	mustEmbedUnimplementedEncodeAPIServer()
}

// UnimplementedEncodeAPIServer must be embedded to have forward compatible implementations.
type UnimplementedEncodeAPIServer struct {
}

func (UnimplementedEncodeAPIServer) Encode(context.Context, *EncodeRequest) (*EncodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Encode not implemented")
}
func (UnimplementedEncodeAPIServer) mustEmbedUnimplementedEncodeAPIServer() {}

// UnsafeEncodeAPIServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EncodeAPIServer will
// result in compilation errors.
type UnsafeEncodeAPIServer interface {
	mustEmbedUnimplementedEncodeAPIServer()
}

func RegisterEncodeAPIServer(s grpc.ServiceRegistrar, srv EncodeAPIServer) {
	s.RegisterService(&_EncodeAPI_serviceDesc, srv)
}

func _EncodeAPI_Encode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EncodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EncodeAPIServer).Encode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tx.v1.EncodeAPI/Encode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EncodeAPIServer).Encode(ctx, req.(*EncodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _EncodeAPI_serviceDesc = grpc.ServiceDesc{
	ServiceName: "tx.v1.EncodeAPI",
	HandlerType: (*EncodeAPIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Encode",
			Handler:    _EncodeAPI_Encode_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "tx/v1/encode_api.proto",
}
