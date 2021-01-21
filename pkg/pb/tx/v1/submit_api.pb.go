// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: tx/v1/submit_api.proto

package txv1

import (
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

// SubmitRequest contains a Transaction.
type SubmitRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Namespace         string             `protobuf:"bytes,1,opt,name=namespace,proto3" json:"namespace,omitempty"`
	SignedTransaction *SignedTransaction `protobuf:"bytes,2,opt,name=signed_transaction,json=signedTransaction,proto3" json:"signed_transaction,omitempty"`
}

func (x *SubmitRequest) Reset() {
	*x = SubmitRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tx_v1_submit_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubmitRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitRequest) ProtoMessage() {}

func (x *SubmitRequest) ProtoReflect() protoreflect.Message {
	mi := &file_tx_v1_submit_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmitRequest.ProtoReflect.Descriptor instead.
func (*SubmitRequest) Descriptor() ([]byte, []int) {
	return file_tx_v1_submit_api_proto_rawDescGZIP(), []int{0}
}

func (x *SubmitRequest) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

func (x *SubmitRequest) GetSignedTransaction() *SignedTransaction {
	if x != nil {
		return x.SignedTransaction
	}
	return nil
}

// SubmitResponse returns the unique identifier for the transaction that was
// submitted.
type SubmitResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Txid []byte `protobuf:"bytes,1,opt,name=txid,proto3" json:"txid,omitempty"`
}

func (x *SubmitResponse) Reset() {
	*x = SubmitResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tx_v1_submit_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubmitResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitResponse) ProtoMessage() {}

func (x *SubmitResponse) ProtoReflect() protoreflect.Message {
	mi := &file_tx_v1_submit_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmitResponse.ProtoReflect.Descriptor instead.
func (*SubmitResponse) Descriptor() ([]byte, []int) {
	return file_tx_v1_submit_api_proto_rawDescGZIP(), []int{1}
}

func (x *SubmitResponse) GetTxid() []byte {
	if x != nil {
		return x.Txid
	}
	return nil
}

var File_tx_v1_submit_api_proto protoreflect.FileDescriptor

var file_tx_v1_submit_api_proto_rawDesc = []byte{
	0x0a, 0x16, 0x74, 0x78, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x5f, 0x61,
	0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x74, 0x78, 0x2e, 0x76, 0x31, 0x1a,
	0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x74,
	0x78, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x76, 0x0a, 0x0d, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x6e, 0x61, 0x6d, 0x65, 0x73,
	0x70, 0x61, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6e, 0x61, 0x6d, 0x65,
	0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x47, 0x0a, 0x12, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64, 0x5f,
	0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x18, 0x2e, 0x74, 0x78, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x69, 0x67, 0x6e, 0x65, 0x64,
	0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x11, 0x73, 0x69, 0x67,
	0x6e, 0x65, 0x64, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x24,
	0x0a, 0x0e, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x74, 0x78, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04,
	0x74, 0x78, 0x69, 0x64, 0x32, 0x76, 0x0a, 0x09, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x41, 0x50,
	0x49, 0x12, 0x69, 0x0a, 0x06, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x12, 0x14, 0x2e, 0x74, 0x78,
	0x2e, 0x76, 0x31, 0x2e, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x15, 0x2e, 0x74, 0x78, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x32, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x2c,
	0x22, 0x16, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x2f, 0x7b, 0x6e, 0x61,
	0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x7d, 0x3a, 0x12, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x64,
	0x5f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x2e, 0x5a, 0x2c,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x69, 0x62, 0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62,
	0x61, 0x74, 0x69, 0x6b, 0x2f, 0x62, 0x61, 0x74, 0x69, 0x6b, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70,
	0x62, 0x2f, 0x74, 0x78, 0x2f, 0x76, 0x31, 0x3b, 0x74, 0x78, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_tx_v1_submit_api_proto_rawDescOnce sync.Once
	file_tx_v1_submit_api_proto_rawDescData = file_tx_v1_submit_api_proto_rawDesc
)

func file_tx_v1_submit_api_proto_rawDescGZIP() []byte {
	file_tx_v1_submit_api_proto_rawDescOnce.Do(func() {
		file_tx_v1_submit_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_tx_v1_submit_api_proto_rawDescData)
	})
	return file_tx_v1_submit_api_proto_rawDescData
}

var file_tx_v1_submit_api_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_tx_v1_submit_api_proto_goTypes = []interface{}{
	(*SubmitRequest)(nil),     // 0: tx.v1.SubmitRequest
	(*SubmitResponse)(nil),    // 1: tx.v1.SubmitResponse
	(*SignedTransaction)(nil), // 2: tx.v1.SignedTransaction
}
var file_tx_v1_submit_api_proto_depIdxs = []int32{
	2, // 0: tx.v1.SubmitRequest.signed_transaction:type_name -> tx.v1.SignedTransaction
	0, // 1: tx.v1.SubmitAPI.Submit:input_type -> tx.v1.SubmitRequest
	1, // 2: tx.v1.SubmitAPI.Submit:output_type -> tx.v1.SubmitResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_tx_v1_submit_api_proto_init() }
func file_tx_v1_submit_api_proto_init() {
	if File_tx_v1_submit_api_proto != nil {
		return
	}
	file_tx_v1_transaction_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_tx_v1_submit_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubmitRequest); i {
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
		file_tx_v1_submit_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubmitResponse); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_tx_v1_submit_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_tx_v1_submit_api_proto_goTypes,
		DependencyIndexes: file_tx_v1_submit_api_proto_depIdxs,
		MessageInfos:      file_tx_v1_submit_api_proto_msgTypes,
	}.Build()
	File_tx_v1_submit_api_proto = out.File
	file_tx_v1_submit_api_proto_rawDesc = nil
	file_tx_v1_submit_api_proto_goTypes = nil
	file_tx_v1_submit_api_proto_depIdxs = nil
}
