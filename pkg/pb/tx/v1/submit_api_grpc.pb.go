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

// SubmitAPIClient is the client API for SubmitAPI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SubmitAPIClient interface {
	// Submit submits a transaction for validation and commit processing.
	// NOTE: This is an implementation for prototyping.
	Submit(ctx context.Context, in *SubmitRequest, opts ...grpc.CallOption) (*SubmitResponse, error)
}

type submitAPIClient struct {
	cc grpc.ClientConnInterface
}

func NewSubmitAPIClient(cc grpc.ClientConnInterface) SubmitAPIClient {
	return &submitAPIClient{cc}
}

func (c *submitAPIClient) Submit(ctx context.Context, in *SubmitRequest, opts ...grpc.CallOption) (*SubmitResponse, error) {
	out := new(SubmitResponse)
	err := c.cc.Invoke(ctx, "/tx.v1.SubmitAPI/Submit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SubmitAPIServer is the server API for SubmitAPI service.
// All implementations must embed UnimplementedSubmitAPIServer
// for forward compatibility
type SubmitAPIServer interface {
	// Submit submits a transaction for validation and commit processing.
	// NOTE: This is an implementation for prototyping.
	Submit(context.Context, *SubmitRequest) (*SubmitResponse, error)
	mustEmbedUnimplementedSubmitAPIServer()
}

// UnimplementedSubmitAPIServer must be embedded to have forward compatible implementations.
type UnimplementedSubmitAPIServer struct {
}

func (UnimplementedSubmitAPIServer) Submit(context.Context, *SubmitRequest) (*SubmitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Submit not implemented")
}
func (UnimplementedSubmitAPIServer) mustEmbedUnimplementedSubmitAPIServer() {}

// UnsafeSubmitAPIServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SubmitAPIServer will
// result in compilation errors.
type UnsafeSubmitAPIServer interface {
	mustEmbedUnimplementedSubmitAPIServer()
}

func RegisterSubmitAPIServer(s grpc.ServiceRegistrar, srv SubmitAPIServer) {
	s.RegisterService(&_SubmitAPI_serviceDesc, srv)
}

func _SubmitAPI_Submit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubmitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SubmitAPIServer).Submit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/tx.v1.SubmitAPI/Submit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SubmitAPIServer).Submit(ctx, req.(*SubmitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _SubmitAPI_serviceDesc = grpc.ServiceDesc{
	ServiceName: "tx.v1.SubmitAPI",
	HandlerType: (*SubmitAPIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Submit",
			Handler:    _SubmitAPI_Submit_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "tx/v1/submit_api.proto",
}
