// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: store/v1/store_api.proto

package storev1

import (
	proto "github.com/golang/protobuf/proto"
	v1 "github.com/sykesm/batik/pkg/pb/tx/v1"
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

// GetTransactionRequest contains a hashed transaction id.
type GetTransactionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Txid []byte `protobuf:"bytes,1,opt,name=txid,proto3" json:"txid,omitempty"`
}

func (x *GetTransactionRequest) Reset() {
	*x = GetTransactionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_v1_store_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetTransactionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTransactionRequest) ProtoMessage() {}

func (x *GetTransactionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_store_v1_store_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTransactionRequest.ProtoReflect.Descriptor instead.
func (*GetTransactionRequest) Descriptor() ([]byte, []int) {
	return file_store_v1_store_api_proto_rawDescGZIP(), []int{0}
}

func (x *GetTransactionRequest) GetTxid() []byte {
	if x != nil {
		return x.Txid
	}
	return nil
}

// GetTransactionResponse contains the transaction retrieved from the backing
// store indexed by the txid that it hashes to.
type GetTransactionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Transaction *v1.Transaction `protobuf:"bytes,1,opt,name=transaction,proto3" json:"transaction,omitempty"`
}

func (x *GetTransactionResponse) Reset() {
	*x = GetTransactionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_v1_store_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetTransactionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTransactionResponse) ProtoMessage() {}

func (x *GetTransactionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_store_v1_store_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTransactionResponse.ProtoReflect.Descriptor instead.
func (*GetTransactionResponse) Descriptor() ([]byte, []int) {
	return file_store_v1_store_api_proto_rawDescGZIP(), []int{1}
}

func (x *GetTransactionResponse) GetTransaction() *v1.Transaction {
	if x != nil {
		return x.Transaction
	}
	return nil
}

// PutTransactionRequest contains a tx.
type PutTransactionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Transaction *v1.Transaction `protobuf:"bytes,1,opt,name=transaction,proto3" json:"transaction,omitempty"`
}

func (x *PutTransactionRequest) Reset() {
	*x = PutTransactionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_v1_store_api_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PutTransactionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PutTransactionRequest) ProtoMessage() {}

func (x *PutTransactionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_store_v1_store_api_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PutTransactionRequest.ProtoReflect.Descriptor instead.
func (*PutTransactionRequest) Descriptor() ([]byte, []int) {
	return file_store_v1_store_api_proto_rawDescGZIP(), []int{2}
}

func (x *PutTransactionRequest) GetTransaction() *v1.Transaction {
	if x != nil {
		return x.Transaction
	}
	return nil
}

// PutTransactionResponse is an empty response returned on attempting to store
// a transaction in the backing store.
type PutTransactionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PutTransactionResponse) Reset() {
	*x = PutTransactionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_v1_store_api_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PutTransactionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PutTransactionResponse) ProtoMessage() {}

func (x *PutTransactionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_store_v1_store_api_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PutTransactionResponse.ProtoReflect.Descriptor instead.
func (*PutTransactionResponse) Descriptor() ([]byte, []int) {
	return file_store_v1_store_api_proto_rawDescGZIP(), []int{3}
}

// GetStateRequest provides a state reference to resolve in the backing store.
// Consumed indicates whether to fetch a consumed state.
type GetStateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StateRef *v1.StateReference `protobuf:"bytes,1,opt,name=state_ref,json=stateRef,proto3" json:"state_ref,omitempty"`
	Consumed bool               `protobuf:"varint,2,opt,name=consumed,proto3" json:"consumed,omitempty"`
}

func (x *GetStateRequest) Reset() {
	*x = GetStateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_v1_store_api_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetStateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetStateRequest) ProtoMessage() {}

func (x *GetStateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_store_v1_store_api_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetStateRequest.ProtoReflect.Descriptor instead.
func (*GetStateRequest) Descriptor() ([]byte, []int) {
	return file_store_v1_store_api_proto_rawDescGZIP(), []int{4}
}

func (x *GetStateRequest) GetStateRef() *v1.StateReference {
	if x != nil {
		return x.StateRef
	}
	return nil
}

func (x *GetStateRequest) GetConsumed() bool {
	if x != nil {
		return x.Consumed
	}
	return false
}

// GetStateResponse contains a transaction state from the backing store and its
// reference.
type GetStateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StateReference *v1.StateReference `protobuf:"bytes,1,opt,name=state_reference,json=stateReference,proto3" json:"state_reference,omitempty"`
	State          *v1.State          `protobuf:"bytes,2,opt,name=state,proto3" json:"state,omitempty"`
}

func (x *GetStateResponse) Reset() {
	*x = GetStateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_v1_store_api_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetStateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetStateResponse) ProtoMessage() {}

func (x *GetStateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_store_v1_store_api_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetStateResponse.ProtoReflect.Descriptor instead.
func (*GetStateResponse) Descriptor() ([]byte, []int) {
	return file_store_v1_store_api_proto_rawDescGZIP(), []int{5}
}

func (x *GetStateResponse) GetStateReference() *v1.StateReference {
	if x != nil {
		return x.StateReference
	}
	return nil
}

func (x *GetStateResponse) GetState() *v1.State {
	if x != nil {
		return x.State
	}
	return nil
}

// PutStateRequest contains the data necessary to populate the backing store
// with a state.
type PutStateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StateReference *v1.StateReference `protobuf:"bytes,1,opt,name=state_reference,json=stateReference,proto3" json:"state_reference,omitempty"`
	State          *v1.State          `protobuf:"bytes,2,opt,name=state,proto3" json:"state,omitempty"`
}

func (x *PutStateRequest) Reset() {
	*x = PutStateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_v1_store_api_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PutStateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PutStateRequest) ProtoMessage() {}

func (x *PutStateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_store_v1_store_api_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PutStateRequest.ProtoReflect.Descriptor instead.
func (*PutStateRequest) Descriptor() ([]byte, []int) {
	return file_store_v1_store_api_proto_rawDescGZIP(), []int{6}
}

func (x *PutStateRequest) GetStateReference() *v1.StateReference {
	if x != nil {
		return x.StateReference
	}
	return nil
}

func (x *PutStateRequest) GetState() *v1.State {
	if x != nil {
		return x.State
	}
	return nil
}

// PutStateResponse is an empty response returned on attempting to store
// a resolved state in the backing store.
type PutStateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PutStateResponse) Reset() {
	*x = PutStateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_v1_store_api_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PutStateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PutStateResponse) ProtoMessage() {}

func (x *PutStateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_store_v1_store_api_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PutStateResponse.ProtoReflect.Descriptor instead.
func (*PutStateResponse) Descriptor() ([]byte, []int) {
	return file_store_v1_store_api_proto_rawDescGZIP(), []int{7}
}

var File_store_v1_store_api_proto protoreflect.FileDescriptor

var file_store_v1_store_api_proto_rawDesc = []byte{
	0x0a, 0x18, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x65,
	0x5f, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x73, 0x74, 0x6f, 0x72,
	0x65, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x17, 0x74, 0x78, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2b, 0x0a, 0x15, 0x47,
	0x65, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x78, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x04, 0x74, 0x78, 0x69, 0x64, 0x22, 0x4e, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x54,
	0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x34, 0x0a, 0x0b, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x74, 0x78, 0x2e, 0x76, 0x31, 0x2e,
	0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0b, 0x74, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x4d, 0x0a, 0x15, 0x50, 0x75, 0x74, 0x54,
	0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x34, 0x0a, 0x0b, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x74, 0x78, 0x2e, 0x76, 0x31, 0x2e, 0x54,
	0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0b, 0x74, 0x72, 0x61, 0x6e,
	0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x18, 0x0a, 0x16, 0x50, 0x75, 0x74, 0x54, 0x72,
	0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x61, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x32, 0x0a, 0x09, 0x73, 0x74, 0x61, 0x74, 0x65, 0x5f, 0x72, 0x65,
	0x66, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x74, 0x78, 0x2e, 0x76, 0x31, 0x2e,
	0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x08,
	0x73, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x66, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x6f, 0x6e, 0x73,
	0x75, 0x6d, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x63, 0x6f, 0x6e, 0x73,
	0x75, 0x6d, 0x65, 0x64, 0x22, 0x76, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3e, 0x0a, 0x0f, 0x73, 0x74, 0x61, 0x74,
	0x65, 0x5f, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x15, 0x2e, 0x74, 0x78, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52,
	0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x52, 0x0e, 0x73, 0x74, 0x61, 0x74, 0x65, 0x52,
	0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x12, 0x22, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x74, 0x78, 0x2e, 0x76, 0x31, 0x2e,
	0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x22, 0x75, 0x0a, 0x0f,
	0x50, 0x75, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x3e, 0x0a, 0x0f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x5f, 0x72, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e,
	0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x74, 0x78, 0x2e, 0x76, 0x31,
	0x2e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x52,
	0x0e, 0x73, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63, 0x65, 0x12,
	0x22, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c,
	0x2e, 0x74, 0x78, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x05, 0x73, 0x74,
	0x61, 0x74, 0x65, 0x22, 0x12, 0x0a, 0x10, 0x50, 0x75, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0xd7, 0x02, 0x0a, 0x08, 0x53, 0x74, 0x6f, 0x72,
	0x65, 0x41, 0x50, 0x49, 0x12, 0x70, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1f, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x76,
	0x31, 0x2e, 0x47, 0x65, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1b, 0x82, 0xd3, 0xe4, 0x93, 0x02,
	0x15, 0x12, 0x13, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x74, 0x78, 0x2f,
	0x7b, 0x74, 0x78, 0x69, 0x64, 0x7d, 0x12, 0x53, 0x0a, 0x0e, 0x50, 0x75, 0x74, 0x54, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1f, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65,
	0x2e, 0x76, 0x31, 0x2e, 0x50, 0x75, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x73, 0x74, 0x6f, 0x72,
	0x65, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x75, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a, 0x08, 0x47,
	0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x19, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e,
	0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65,
	0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x41,
	0x0a, 0x08, 0x50, 0x75, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x19, 0x2e, 0x73, 0x74, 0x6f,
	0x72, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x75, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x76, 0x31,
	0x2e, 0x50, 0x75, 0x74, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x42, 0x34, 0x5a, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x69, 0x62, 0x6d, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x61, 0x74, 0x69, 0x6b, 0x2f, 0x62, 0x61, 0x74, 0x69, 0x6b, 0x2f,
	0x70, 0x6b, 0x67, 0x2f, 0x70, 0x62, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x76, 0x31, 0x3b,
	0x73, 0x74, 0x6f, 0x72, 0x65, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_store_v1_store_api_proto_rawDescOnce sync.Once
	file_store_v1_store_api_proto_rawDescData = file_store_v1_store_api_proto_rawDesc
)

func file_store_v1_store_api_proto_rawDescGZIP() []byte {
	file_store_v1_store_api_proto_rawDescOnce.Do(func() {
		file_store_v1_store_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_store_v1_store_api_proto_rawDescData)
	})
	return file_store_v1_store_api_proto_rawDescData
}

var file_store_v1_store_api_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_store_v1_store_api_proto_goTypes = []interface{}{
	(*GetTransactionRequest)(nil),  // 0: store.v1.GetTransactionRequest
	(*GetTransactionResponse)(nil), // 1: store.v1.GetTransactionResponse
	(*PutTransactionRequest)(nil),  // 2: store.v1.PutTransactionRequest
	(*PutTransactionResponse)(nil), // 3: store.v1.PutTransactionResponse
	(*GetStateRequest)(nil),        // 4: store.v1.GetStateRequest
	(*GetStateResponse)(nil),       // 5: store.v1.GetStateResponse
	(*PutStateRequest)(nil),        // 6: store.v1.PutStateRequest
	(*PutStateResponse)(nil),       // 7: store.v1.PutStateResponse
	(*v1.Transaction)(nil),         // 8: tx.v1.Transaction
	(*v1.StateReference)(nil),      // 9: tx.v1.StateReference
	(*v1.State)(nil),               // 10: tx.v1.State
}
var file_store_v1_store_api_proto_depIdxs = []int32{
	8,  // 0: store.v1.GetTransactionResponse.transaction:type_name -> tx.v1.Transaction
	8,  // 1: store.v1.PutTransactionRequest.transaction:type_name -> tx.v1.Transaction
	9,  // 2: store.v1.GetStateRequest.state_ref:type_name -> tx.v1.StateReference
	9,  // 3: store.v1.GetStateResponse.state_reference:type_name -> tx.v1.StateReference
	10, // 4: store.v1.GetStateResponse.state:type_name -> tx.v1.State
	9,  // 5: store.v1.PutStateRequest.state_reference:type_name -> tx.v1.StateReference
	10, // 6: store.v1.PutStateRequest.state:type_name -> tx.v1.State
	0,  // 7: store.v1.StoreAPI.GetTransaction:input_type -> store.v1.GetTransactionRequest
	2,  // 8: store.v1.StoreAPI.PutTransaction:input_type -> store.v1.PutTransactionRequest
	4,  // 9: store.v1.StoreAPI.GetState:input_type -> store.v1.GetStateRequest
	6,  // 10: store.v1.StoreAPI.PutState:input_type -> store.v1.PutStateRequest
	1,  // 11: store.v1.StoreAPI.GetTransaction:output_type -> store.v1.GetTransactionResponse
	3,  // 12: store.v1.StoreAPI.PutTransaction:output_type -> store.v1.PutTransactionResponse
	5,  // 13: store.v1.StoreAPI.GetState:output_type -> store.v1.GetStateResponse
	7,  // 14: store.v1.StoreAPI.PutState:output_type -> store.v1.PutStateResponse
	11, // [11:15] is the sub-list for method output_type
	7,  // [7:11] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_store_v1_store_api_proto_init() }
func file_store_v1_store_api_proto_init() {
	if File_store_v1_store_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_store_v1_store_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetTransactionRequest); i {
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
		file_store_v1_store_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetTransactionResponse); i {
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
		file_store_v1_store_api_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PutTransactionRequest); i {
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
		file_store_v1_store_api_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PutTransactionResponse); i {
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
		file_store_v1_store_api_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetStateRequest); i {
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
		file_store_v1_store_api_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetStateResponse); i {
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
		file_store_v1_store_api_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PutStateRequest); i {
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
		file_store_v1_store_api_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PutStateResponse); i {
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
			RawDescriptor: file_store_v1_store_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_store_v1_store_api_proto_goTypes,
		DependencyIndexes: file_store_v1_store_api_proto_depIdxs,
		MessageInfos:      file_store_v1_store_api_proto_msgTypes,
	}.Build()
	File_store_v1_store_api_proto = out.File
	file_store_v1_store_api_proto_rawDesc = nil
	file_store_v1_store_api_proto_goTypes = nil
	file_store_v1_store_api_proto_depIdxs = nil
}
