// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: orderer/atomic_broadcast_api.proto

package orderer

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// BroadcastRequest wraps a Payload with a signature so that the message may be authenticated.
type BroadcastRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A marshaled Payload.
	Payload []byte `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	// A signature by the creator specified in the Payload header.
	Signature []byte `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (x *BroadcastRequest) Reset() {
	*x = BroadcastRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_orderer_atomic_broadcast_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BroadcastRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BroadcastRequest) ProtoMessage() {}

func (x *BroadcastRequest) ProtoReflect() protoreflect.Message {
	mi := &file_orderer_atomic_broadcast_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BroadcastRequest.ProtoReflect.Descriptor instead.
func (*BroadcastRequest) Descriptor() ([]byte, []int) {
	return file_orderer_atomic_broadcast_api_proto_rawDescGZIP(), []int{0}
}

func (x *BroadcastRequest) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *BroadcastRequest) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

// BroadcastResponse contains information indicating whether the Broadcast was successful.
type BroadcastResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Status code, which may be used to programatically respond to success/failure.
	Status Status `protobuf:"varint,1,opt,name=status,proto3,enum=orderer.Status" json:"status,omitempty"`
	// Info string which may contain additional information about the status returned.
	Info string `protobuf:"bytes,2,opt,name=info,proto3" json:"info,omitempty"`
}

func (x *BroadcastResponse) Reset() {
	*x = BroadcastResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_orderer_atomic_broadcast_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BroadcastResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BroadcastResponse) ProtoMessage() {}

func (x *BroadcastResponse) ProtoReflect() protoreflect.Message {
	mi := &file_orderer_atomic_broadcast_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BroadcastResponse.ProtoReflect.Descriptor instead.
func (*BroadcastResponse) Descriptor() ([]byte, []int) {
	return file_orderer_atomic_broadcast_api_proto_rawDescGZIP(), []int{1}
}

func (x *BroadcastResponse) GetStatus() Status {
	if x != nil {
		return x.Status
	}
	return Status_STATUS_INVALID
}

func (x *BroadcastResponse) GetInfo() string {
	if x != nil {
		return x.Info
	}
	return ""
}

// DeliverRequest wraps a Payload with a signature so that the message may be authenticated.
type DeliverRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A marshaled Payload.
	Payload []byte `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	// A signature by the creator specified in the Payload header.
	Signature []byte `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (x *DeliverRequest) Reset() {
	*x = DeliverRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_orderer_atomic_broadcast_api_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeliverRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeliverRequest) ProtoMessage() {}

func (x *DeliverRequest) ProtoReflect() protoreflect.Message {
	mi := &file_orderer_atomic_broadcast_api_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeliverRequest.ProtoReflect.Descriptor instead.
func (*DeliverRequest) Descriptor() ([]byte, []int) {
	return file_orderer_atomic_broadcast_api_proto_rawDescGZIP(), []int{2}
}

func (x *DeliverRequest) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *DeliverRequest) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

// DeliverResponse contains either a status or a transaction hashes.
type DeliverResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Type:
	//	*DeliverResponse_Status
	//	*DeliverResponse_Txid
	Type isDeliverResponse_Type `protobuf_oneof:"type"`
}

func (x *DeliverResponse) Reset() {
	*x = DeliverResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_orderer_atomic_broadcast_api_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeliverResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeliverResponse) ProtoMessage() {}

func (x *DeliverResponse) ProtoReflect() protoreflect.Message {
	mi := &file_orderer_atomic_broadcast_api_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeliverResponse.ProtoReflect.Descriptor instead.
func (*DeliverResponse) Descriptor() ([]byte, []int) {
	return file_orderer_atomic_broadcast_api_proto_rawDescGZIP(), []int{3}
}

func (m *DeliverResponse) GetType() isDeliverResponse_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *DeliverResponse) GetStatus() Status {
	if x, ok := x.GetType().(*DeliverResponse_Status); ok {
		return x.Status
	}
	return Status_STATUS_INVALID
}

func (x *DeliverResponse) GetTxid() []byte {
	if x, ok := x.GetType().(*DeliverResponse_Txid); ok {
		return x.Txid
	}
	return nil
}

type isDeliverResponse_Type interface {
	isDeliverResponse_Type()
}

type DeliverResponse_Status struct {
	Status Status `protobuf:"varint,1,opt,name=status,proto3,enum=orderer.Status,oneof"`
}

type DeliverResponse_Txid struct {
	Txid []byte `protobuf:"bytes,2,opt,name=txid,proto3,oneof"`
}

func (*DeliverResponse_Status) isDeliverResponse_Type() {}

func (*DeliverResponse_Txid) isDeliverResponse_Type() {}

var File_orderer_atomic_broadcast_api_proto protoreflect.FileDescriptor

var file_orderer_atomic_broadcast_api_proto_rawDesc = []byte{
	0x0a, 0x22, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x65, 0x72, 0x2f, 0x61, 0x74, 0x6f, 0x6d, 0x69, 0x63,
	0x5f, 0x62, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x5f, 0x61, 0x70, 0x69, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x65, 0x72, 0x1a, 0x14, 0x6f,
	0x72, 0x64, 0x65, 0x72, 0x65, 0x72, 0x2f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x4a, 0x0a, 0x10, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f,
	0x61, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61,
	0x64, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x22,
	0x50, 0x0a, 0x11, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x27, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x0f, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x65, 0x72, 0x2e, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x12, 0x0a,
	0x04, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x69, 0x6e, 0x66,
	0x6f, 0x22, 0x48, 0x0a, 0x0e, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x1c, 0x0a,
	0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x22, 0x5a, 0x0a, 0x0f, 0x44,
	0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x29,
	0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0f,
	0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x65, 0x72, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x48,
	0x00, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x14, 0x0a, 0x04, 0x74, 0x78, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x00, 0x52, 0x04, 0x74, 0x78, 0x69, 0x64, 0x42,
	0x06, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x32, 0x9e, 0x01, 0x0a, 0x12, 0x41, 0x74, 0x6f, 0x6d,
	0x69, 0x63, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x41, 0x50, 0x49, 0x12, 0x46,
	0x0a, 0x09, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x12, 0x19, 0x2e, 0x6f, 0x72,
	0x64, 0x65, 0x72, 0x65, 0x72, 0x2e, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x65, 0x72,
	0x2e, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x28, 0x01, 0x30, 0x01, 0x12, 0x40, 0x0a, 0x07, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65,
	0x72, 0x12, 0x17, 0x2e, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x65, 0x72, 0x2e, 0x44, 0x65, 0x6c, 0x69,
	0x76, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x6f, 0x72, 0x64,
	0x65, 0x72, 0x65, 0x72, 0x2e, 0x44, 0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x28, 0x01, 0x30, 0x01, 0x42, 0x2b, 0x5a, 0x29, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x69, 0x62, 0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x61, 0x74, 0x69, 0x6b,
	0x2f, 0x62, 0x61, 0x74, 0x69, 0x6b, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x62, 0x2f, 0x6f, 0x72,
	0x64, 0x65, 0x72, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_orderer_atomic_broadcast_api_proto_rawDescOnce sync.Once
	file_orderer_atomic_broadcast_api_proto_rawDescData = file_orderer_atomic_broadcast_api_proto_rawDesc
)

func file_orderer_atomic_broadcast_api_proto_rawDescGZIP() []byte {
	file_orderer_atomic_broadcast_api_proto_rawDescOnce.Do(func() {
		file_orderer_atomic_broadcast_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_orderer_atomic_broadcast_api_proto_rawDescData)
	})
	return file_orderer_atomic_broadcast_api_proto_rawDescData
}

var file_orderer_atomic_broadcast_api_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_orderer_atomic_broadcast_api_proto_goTypes = []interface{}{
	(*BroadcastRequest)(nil),  // 0: orderer.BroadcastRequest
	(*BroadcastResponse)(nil), // 1: orderer.BroadcastResponse
	(*DeliverRequest)(nil),    // 2: orderer.DeliverRequest
	(*DeliverResponse)(nil),   // 3: orderer.DeliverResponse
	(Status)(0),               // 4: orderer.Status
}
var file_orderer_atomic_broadcast_api_proto_depIdxs = []int32{
	4, // 0: orderer.BroadcastResponse.status:type_name -> orderer.Status
	4, // 1: orderer.DeliverResponse.status:type_name -> orderer.Status
	0, // 2: orderer.AtomicBroadcastAPI.Broadcast:input_type -> orderer.BroadcastRequest
	2, // 3: orderer.AtomicBroadcastAPI.Deliver:input_type -> orderer.DeliverRequest
	1, // 4: orderer.AtomicBroadcastAPI.Broadcast:output_type -> orderer.BroadcastResponse
	3, // 5: orderer.AtomicBroadcastAPI.Deliver:output_type -> orderer.DeliverResponse
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_orderer_atomic_broadcast_api_proto_init() }
func file_orderer_atomic_broadcast_api_proto_init() {
	if File_orderer_atomic_broadcast_api_proto != nil {
		return
	}
	file_orderer_status_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_orderer_atomic_broadcast_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BroadcastRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_orderer_atomic_broadcast_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BroadcastResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_orderer_atomic_broadcast_api_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeliverRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_orderer_atomic_broadcast_api_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeliverResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_orderer_atomic_broadcast_api_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*DeliverResponse_Status)(nil),
		(*DeliverResponse_Txid)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_orderer_atomic_broadcast_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_orderer_atomic_broadcast_api_proto_goTypes,
		DependencyIndexes: file_orderer_atomic_broadcast_api_proto_depIdxs,
		MessageInfos:      file_orderer_atomic_broadcast_api_proto_msgTypes,
	}.Build()
	File_orderer_atomic_broadcast_api_proto = out.File
	file_orderer_atomic_broadcast_api_proto_rawDesc = nil
	file_orderer_atomic_broadcast_api_proto_goTypes = nil
	file_orderer_atomic_broadcast_api_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// AtomicBroadcastAPIClient is the client API for AtomicBroadcastAPI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AtomicBroadcastAPIClient interface {
	// Broadcast receives a reply of Acknowledgement for each common.Envelope in order, indicating success or type of failure.
	Broadcast(ctx context.Context, opts ...grpc.CallOption) (AtomicBroadcastAPI_BroadcastClient, error)
	// Deliver first requires an Envelope of type DELIVER_SEEK_INFO with Payload data as a mashaled SeekInfo message, then a stream of block replies is received.
	Deliver(ctx context.Context, opts ...grpc.CallOption) (AtomicBroadcastAPI_DeliverClient, error)
}

type atomicBroadcastAPIClient struct {
	cc grpc.ClientConnInterface
}

func NewAtomicBroadcastAPIClient(cc grpc.ClientConnInterface) AtomicBroadcastAPIClient {
	return &atomicBroadcastAPIClient{cc}
}

func (c *atomicBroadcastAPIClient) Broadcast(ctx context.Context, opts ...grpc.CallOption) (AtomicBroadcastAPI_BroadcastClient, error) {
	stream, err := c.cc.NewStream(ctx, &_AtomicBroadcastAPI_serviceDesc.Streams[0], "/orderer.AtomicBroadcastAPI/Broadcast", opts...)
	if err != nil {
		return nil, err
	}
	x := &atomicBroadcastAPIBroadcastClient{stream}
	return x, nil
}

type AtomicBroadcastAPI_BroadcastClient interface {
	Send(*BroadcastRequest) error
	Recv() (*BroadcastResponse, error)
	grpc.ClientStream
}

type atomicBroadcastAPIBroadcastClient struct {
	grpc.ClientStream
}

func (x *atomicBroadcastAPIBroadcastClient) Send(m *BroadcastRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *atomicBroadcastAPIBroadcastClient) Recv() (*BroadcastResponse, error) {
	m := new(BroadcastResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *atomicBroadcastAPIClient) Deliver(ctx context.Context, opts ...grpc.CallOption) (AtomicBroadcastAPI_DeliverClient, error) {
	stream, err := c.cc.NewStream(ctx, &_AtomicBroadcastAPI_serviceDesc.Streams[1], "/orderer.AtomicBroadcastAPI/Deliver", opts...)
	if err != nil {
		return nil, err
	}
	x := &atomicBroadcastAPIDeliverClient{stream}
	return x, nil
}

type AtomicBroadcastAPI_DeliverClient interface {
	Send(*DeliverRequest) error
	Recv() (*DeliverResponse, error)
	grpc.ClientStream
}

type atomicBroadcastAPIDeliverClient struct {
	grpc.ClientStream
}

func (x *atomicBroadcastAPIDeliverClient) Send(m *DeliverRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *atomicBroadcastAPIDeliverClient) Recv() (*DeliverResponse, error) {
	m := new(DeliverResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// AtomicBroadcastAPIServer is the server API for AtomicBroadcastAPI service.
type AtomicBroadcastAPIServer interface {
	// Broadcast receives a reply of Acknowledgement for each common.Envelope in order, indicating success or type of failure.
	Broadcast(AtomicBroadcastAPI_BroadcastServer) error
	// Deliver first requires an Envelope of type DELIVER_SEEK_INFO with Payload data as a mashaled SeekInfo message, then a stream of block replies is received.
	Deliver(AtomicBroadcastAPI_DeliverServer) error
}

// UnimplementedAtomicBroadcastAPIServer can be embedded to have forward compatible implementations.
type UnimplementedAtomicBroadcastAPIServer struct {
}

func (*UnimplementedAtomicBroadcastAPIServer) Broadcast(AtomicBroadcastAPI_BroadcastServer) error {
	return status.Errorf(codes.Unimplemented, "method Broadcast not implemented")
}
func (*UnimplementedAtomicBroadcastAPIServer) Deliver(AtomicBroadcastAPI_DeliverServer) error {
	return status.Errorf(codes.Unimplemented, "method Deliver not implemented")
}

func RegisterAtomicBroadcastAPIServer(s *grpc.Server, srv AtomicBroadcastAPIServer) {
	s.RegisterService(&_AtomicBroadcastAPI_serviceDesc, srv)
}

func _AtomicBroadcastAPI_Broadcast_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(AtomicBroadcastAPIServer).Broadcast(&atomicBroadcastAPIBroadcastServer{stream})
}

type AtomicBroadcastAPI_BroadcastServer interface {
	Send(*BroadcastResponse) error
	Recv() (*BroadcastRequest, error)
	grpc.ServerStream
}

type atomicBroadcastAPIBroadcastServer struct {
	grpc.ServerStream
}

func (x *atomicBroadcastAPIBroadcastServer) Send(m *BroadcastResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *atomicBroadcastAPIBroadcastServer) Recv() (*BroadcastRequest, error) {
	m := new(BroadcastRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _AtomicBroadcastAPI_Deliver_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(AtomicBroadcastAPIServer).Deliver(&atomicBroadcastAPIDeliverServer{stream})
}

type AtomicBroadcastAPI_DeliverServer interface {
	Send(*DeliverResponse) error
	Recv() (*DeliverRequest, error)
	grpc.ServerStream
}

type atomicBroadcastAPIDeliverServer struct {
	grpc.ServerStream
}

func (x *atomicBroadcastAPIDeliverServer) Send(m *DeliverResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *atomicBroadcastAPIDeliverServer) Recv() (*DeliverRequest, error) {
	m := new(DeliverRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _AtomicBroadcastAPI_serviceDesc = grpc.ServiceDesc{
	ServiceName: "orderer.AtomicBroadcastAPI",
	HandlerType: (*AtomicBroadcastAPIServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Broadcast",
			Handler:       _AtomicBroadcastAPI_Broadcast_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "Deliver",
			Handler:       _AtomicBroadcastAPI_Deliver_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "orderer/atomic_broadcast_api.proto",
}
