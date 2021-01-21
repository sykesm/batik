// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: tx/v1/encode_api.proto

package txv1

import (
	proto "github.com/golang/protobuf/proto"
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

// EncodeRequest contains a Transaction.
type EncodeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Transaction *Transaction `protobuf:"bytes,1,opt,name=transaction,proto3" json:"transaction,omitempty"`
}

func (x *EncodeRequest) Reset() {
	*x = EncodeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tx_v1_encode_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EncodeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EncodeRequest) ProtoMessage() {}

func (x *EncodeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_tx_v1_encode_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EncodeRequest.ProtoReflect.Descriptor instead.
func (*EncodeRequest) Descriptor() ([]byte, []int) {
	return file_tx_v1_encode_api_proto_rawDescGZIP(), []int{0}
}

func (x *EncodeRequest) GetTransaction() *Transaction {
	if x != nil {
		return x.Transaction
	}
	return nil
}

// EncodeTransactionResponse contains the transaction ID and encoded bytes
// representing the transaction passed in the EncodeResponse.
type EncodeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Txid               []byte `protobuf:"bytes,1,opt,name=txid,proto3" json:"txid,omitempty"`
	EncodedTransaction []byte `protobuf:"bytes,2,opt,name=encoded_transaction,json=encodedTransaction,proto3" json:"encoded_transaction,omitempty"`
}

func (x *EncodeResponse) Reset() {
	*x = EncodeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tx_v1_encode_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EncodeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EncodeResponse) ProtoMessage() {}

func (x *EncodeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_tx_v1_encode_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EncodeResponse.ProtoReflect.Descriptor instead.
func (*EncodeResponse) Descriptor() ([]byte, []int) {
	return file_tx_v1_encode_api_proto_rawDescGZIP(), []int{1}
}

func (x *EncodeResponse) GetTxid() []byte {
	if x != nil {
		return x.Txid
	}
	return nil
}

func (x *EncodeResponse) GetEncodedTransaction() []byte {
	if x != nil {
		return x.EncodedTransaction
	}
	return nil
}

var File_tx_v1_encode_api_proto protoreflect.FileDescriptor

var file_tx_v1_encode_api_proto_rawDesc = []byte{
	0x0a, 0x16, 0x74, 0x78, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x5f, 0x61,
	0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x74, 0x78, 0x2e, 0x76, 0x31, 0x1a,
	0x17, 0x74, 0x78, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x45, 0x0a, 0x0d, 0x45, 0x6e, 0x63, 0x6f,
	0x64, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x34, 0x0a, 0x0b, 0x74, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12,
	0x2e, 0x74, 0x78, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x0b, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x22,
	0x55, 0x0a, 0x0e, 0x45, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x78, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x04, 0x74, 0x78, 0x69, 0x64, 0x12, 0x2f, 0x0a, 0x13, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x64,
	0x5f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x12, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x64, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x32, 0x42, 0x0a, 0x09, 0x45, 0x6e, 0x63, 0x6f, 0x64, 0x65,
	0x41, 0x50, 0x49, 0x12, 0x35, 0x0a, 0x06, 0x45, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x14, 0x2e,
	0x74, 0x78, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x74, 0x78, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x6e, 0x63, 0x6f,
	0x64, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x2e, 0x5a, 0x2c, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x69, 0x62, 0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x61, 0x74,
	0x69, 0x6b, 0x2f, 0x62, 0x61, 0x74, 0x69, 0x6b, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x62, 0x2f,
	0x74, 0x78, 0x2f, 0x76, 0x31, 0x3b, 0x74, 0x78, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_tx_v1_encode_api_proto_rawDescOnce sync.Once
	file_tx_v1_encode_api_proto_rawDescData = file_tx_v1_encode_api_proto_rawDesc
)

func file_tx_v1_encode_api_proto_rawDescGZIP() []byte {
	file_tx_v1_encode_api_proto_rawDescOnce.Do(func() {
		file_tx_v1_encode_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_tx_v1_encode_api_proto_rawDescData)
	})
	return file_tx_v1_encode_api_proto_rawDescData
}

var file_tx_v1_encode_api_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_tx_v1_encode_api_proto_goTypes = []interface{}{
	(*EncodeRequest)(nil),  // 0: tx.v1.EncodeRequest
	(*EncodeResponse)(nil), // 1: tx.v1.EncodeResponse
	(*Transaction)(nil),    // 2: tx.v1.Transaction
}
var file_tx_v1_encode_api_proto_depIdxs = []int32{
	2, // 0: tx.v1.EncodeRequest.transaction:type_name -> tx.v1.Transaction
	0, // 1: tx.v1.EncodeAPI.Encode:input_type -> tx.v1.EncodeRequest
	1, // 2: tx.v1.EncodeAPI.Encode:output_type -> tx.v1.EncodeResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_tx_v1_encode_api_proto_init() }
func file_tx_v1_encode_api_proto_init() {
	if File_tx_v1_encode_api_proto != nil {
		return
	}
	file_tx_v1_transaction_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_tx_v1_encode_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EncodeRequest); i {
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
		file_tx_v1_encode_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EncodeResponse); i {
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
			RawDescriptor: file_tx_v1_encode_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_tx_v1_encode_api_proto_goTypes,
		DependencyIndexes: file_tx_v1_encode_api_proto_depIdxs,
		MessageInfos:      file_tx_v1_encode_api_proto_msgTypes,
	}.Build()
	File_tx_v1_encode_api_proto = out.File
	file_tx_v1_encode_api_proto_rawDesc = nil
	file_tx_v1_encode_api_proto_goTypes = nil
	file_tx_v1_encode_api_proto_depIdxs = nil
}
