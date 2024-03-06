// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.11.2
// source: db_list.proto

package db_list

import (
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

type RPUSHRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key    []byte   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Values [][]byte `protobuf:"bytes,2,rep,name=values,proto3" json:"values,omitempty"`
	Ttl    []uint64 `protobuf:"varint,3,rep,packed,name=ttl,proto3" json:"ttl,omitempty"`
}

func (x *RPUSHRequest) Reset() {
	*x = RPUSHRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_list_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RPUSHRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RPUSHRequest) ProtoMessage() {}

func (x *RPUSHRequest) ProtoReflect() protoreflect.Message {
	mi := &file_db_list_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RPUSHRequest.ProtoReflect.Descriptor instead.
func (*RPUSHRequest) Descriptor() ([]byte, []int) {
	return file_db_list_proto_rawDescGZIP(), []int{0}
}

func (x *RPUSHRequest) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *RPUSHRequest) GetValues() [][]byte {
	if x != nil {
		return x.Values
	}
	return nil
}

func (x *RPUSHRequest) GetTtl() []uint64 {
	if x != nil {
		return x.Ttl
	}
	return nil
}

type RPUSHResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *RPUSHResponse) Reset() {
	*x = RPUSHResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_list_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RPUSHResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RPUSHResponse) ProtoMessage() {}

func (x *RPUSHResponse) ProtoReflect() protoreflect.Message {
	mi := &file_db_list_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RPUSHResponse.ProtoReflect.Descriptor instead.
func (*RPUSHResponse) Descriptor() ([]byte, []int) {
	return file_db_list_proto_rawDescGZIP(), []int{1}
}

func (x *RPUSHResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type LPUSHRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key    []byte   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Values [][]byte `protobuf:"bytes,2,rep,name=values,proto3" json:"values,omitempty"`
	Ttl    []uint64 `protobuf:"varint,3,rep,packed,name=ttl,proto3" json:"ttl,omitempty"`
}

func (x *LPUSHRequest) Reset() {
	*x = LPUSHRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_list_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LPUSHRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LPUSHRequest) ProtoMessage() {}

func (x *LPUSHRequest) ProtoReflect() protoreflect.Message {
	mi := &file_db_list_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LPUSHRequest.ProtoReflect.Descriptor instead.
func (*LPUSHRequest) Descriptor() ([]byte, []int) {
	return file_db_list_proto_rawDescGZIP(), []int{2}
}

func (x *LPUSHRequest) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *LPUSHRequest) GetValues() [][]byte {
	if x != nil {
		return x.Values
	}
	return nil
}

func (x *LPUSHRequest) GetTtl() []uint64 {
	if x != nil {
		return x.Ttl
	}
	return nil
}

type LPUSHResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *LPUSHResponse) Reset() {
	*x = LPUSHResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_list_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LPUSHResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LPUSHResponse) ProtoMessage() {}

func (x *LPUSHResponse) ProtoReflect() protoreflect.Message {
	mi := &file_db_list_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LPUSHResponse.ProtoReflect.Descriptor instead.
func (*LPUSHResponse) Descriptor() ([]byte, []int) {
	return file_db_list_proto_rawDescGZIP(), []int{3}
}

func (x *LPUSHResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type LRANGERequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   []byte `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Left  int32  `protobuf:"varint,2,opt,name=left,proto3" json:"left,omitempty"`
	Right int32  `protobuf:"varint,3,opt,name=right,proto3" json:"right,omitempty"`
}

func (x *LRANGERequest) Reset() {
	*x = LRANGERequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_list_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LRANGERequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LRANGERequest) ProtoMessage() {}

func (x *LRANGERequest) ProtoReflect() protoreflect.Message {
	mi := &file_db_list_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LRANGERequest.ProtoReflect.Descriptor instead.
func (*LRANGERequest) Descriptor() ([]byte, []int) {
	return file_db_list_proto_rawDescGZIP(), []int{4}
}

func (x *LRANGERequest) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *LRANGERequest) GetLeft() int32 {
	if x != nil {
		return x.Left
	}
	return 0
}

func (x *LRANGERequest) GetRight() int32 {
	if x != nil {
		return x.Right
	}
	return 0
}

type LRANGEResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Values [][]byte `protobuf:"bytes,1,rep,name=values,proto3" json:"values,omitempty"`
}

func (x *LRANGEResponse) Reset() {
	*x = LRANGEResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_list_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LRANGEResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LRANGEResponse) ProtoMessage() {}

func (x *LRANGEResponse) ProtoReflect() protoreflect.Message {
	mi := &file_db_list_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LRANGEResponse.ProtoReflect.Descriptor instead.
func (*LRANGEResponse) Descriptor() ([]byte, []int) {
	return file_db_list_proto_rawDescGZIP(), []int{5}
}

func (x *LRANGEResponse) GetValues() [][]byte {
	if x != nil {
		return x.Values
	}
	return nil
}

type LINDEXRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   []byte `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Index int32  `protobuf:"varint,2,opt,name=index,proto3" json:"index,omitempty"`
}

func (x *LINDEXRequest) Reset() {
	*x = LINDEXRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_list_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LINDEXRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LINDEXRequest) ProtoMessage() {}

func (x *LINDEXRequest) ProtoReflect() protoreflect.Message {
	mi := &file_db_list_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LINDEXRequest.ProtoReflect.Descriptor instead.
func (*LINDEXRequest) Descriptor() ([]byte, []int) {
	return file_db_list_proto_rawDescGZIP(), []int{6}
}

func (x *LINDEXRequest) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *LINDEXRequest) GetIndex() int32 {
	if x != nil {
		return x.Index
	}
	return 0
}

type LINDEXResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value []byte `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *LINDEXResponse) Reset() {
	*x = LINDEXResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_list_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LINDEXResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LINDEXResponse) ProtoMessage() {}

func (x *LINDEXResponse) ProtoReflect() protoreflect.Message {
	mi := &file_db_list_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LINDEXResponse.ProtoReflect.Descriptor instead.
func (*LINDEXResponse) Descriptor() ([]byte, []int) {
	return file_db_list_proto_rawDescGZIP(), []int{7}
}

func (x *LINDEXResponse) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

type LPOPRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key []byte `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *LPOPRequest) Reset() {
	*x = LPOPRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_list_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LPOPRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LPOPRequest) ProtoMessage() {}

func (x *LPOPRequest) ProtoReflect() protoreflect.Message {
	mi := &file_db_list_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LPOPRequest.ProtoReflect.Descriptor instead.
func (*LPOPRequest) Descriptor() ([]byte, []int) {
	return file_db_list_proto_rawDescGZIP(), []int{8}
}

func (x *LPOPRequest) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

type LPOPResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value []byte `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *LPOPResponse) Reset() {
	*x = LPOPResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_list_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LPOPResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LPOPResponse) ProtoMessage() {}

func (x *LPOPResponse) ProtoReflect() protoreflect.Message {
	mi := &file_db_list_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LPOPResponse.ProtoReflect.Descriptor instead.
func (*LPOPResponse) Descriptor() ([]byte, []int) {
	return file_db_list_proto_rawDescGZIP(), []int{9}
}

func (x *LPOPResponse) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

type RPOPRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key []byte `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *RPOPRequest) Reset() {
	*x = RPOPRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_list_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RPOPRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RPOPRequest) ProtoMessage() {}

func (x *RPOPRequest) ProtoReflect() protoreflect.Message {
	mi := &file_db_list_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RPOPRequest.ProtoReflect.Descriptor instead.
func (*RPOPRequest) Descriptor() ([]byte, []int) {
	return file_db_list_proto_rawDescGZIP(), []int{10}
}

func (x *RPOPRequest) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

type RPOPResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value []byte `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *RPOPResponse) Reset() {
	*x = RPOPResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_list_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RPOPResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RPOPResponse) ProtoMessage() {}

func (x *RPOPResponse) ProtoReflect() protoreflect.Message {
	mi := &file_db_list_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RPOPResponse.ProtoReflect.Descriptor instead.
func (*RPOPResponse) Descriptor() ([]byte, []int) {
	return file_db_list_proto_rawDescGZIP(), []int{11}
}

func (x *RPOPResponse) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

type LLENRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key []byte `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *LLENRequest) Reset() {
	*x = LLENRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_list_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LLENRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LLENRequest) ProtoMessage() {}

func (x *LLENRequest) ProtoReflect() protoreflect.Message {
	mi := &file_db_list_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LLENRequest.ProtoReflect.Descriptor instead.
func (*LLENRequest) Descriptor() ([]byte, []int) {
	return file_db_list_proto_rawDescGZIP(), []int{12}
}

func (x *LLENRequest) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

type LLENResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Length int32 `protobuf:"varint,1,opt,name=length,proto3" json:"length,omitempty"`
}

func (x *LLENResponse) Reset() {
	*x = LLENResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_db_list_proto_msgTypes[13]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LLENResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LLENResponse) ProtoMessage() {}

func (x *LLENResponse) ProtoReflect() protoreflect.Message {
	mi := &file_db_list_proto_msgTypes[13]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LLENResponse.ProtoReflect.Descriptor instead.
func (*LLENResponse) Descriptor() ([]byte, []int) {
	return file_db_list_proto_rawDescGZIP(), []int{13}
}

func (x *LLENResponse) GetLength() int32 {
	if x != nil {
		return x.Length
	}
	return 0
}

var File_db_list_proto protoreflect.FileDescriptor

var file_db_list_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x64, 0x62, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x0b, 0x61, 0x70, 0x69, 0x5f, 0x64, 0x62, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x22, 0x4a, 0x0a, 0x0c,
	0x52, 0x50, 0x55, 0x53, 0x48, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x16,
	0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x06,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x74, 0x74, 0x6c, 0x18, 0x03, 0x20,
	0x03, 0x28, 0x04, 0x52, 0x03, 0x74, 0x74, 0x6c, 0x22, 0x29, 0x0a, 0x0d, 0x52, 0x50, 0x55, 0x53,
	0x48, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x22, 0x4a, 0x0a, 0x0c, 0x4c, 0x50, 0x55, 0x53, 0x48, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x12, 0x10, 0x0a,
	0x03, 0x74, 0x74, 0x6c, 0x18, 0x03, 0x20, 0x03, 0x28, 0x04, 0x52, 0x03, 0x74, 0x74, 0x6c, 0x22,
	0x29, 0x0a, 0x0d, 0x4c, 0x50, 0x55, 0x53, 0x48, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x22, 0x4b, 0x0a, 0x0d, 0x4c, 0x52,
	0x41, 0x4e, 0x47, 0x45, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x12, 0x0a,
	0x04, 0x6c, 0x65, 0x66, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x6c, 0x65, 0x66,
	0x74, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x69, 0x67, 0x68, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x05, 0x72, 0x69, 0x67, 0x68, 0x74, 0x22, 0x28, 0x0a, 0x0e, 0x4c, 0x52, 0x41, 0x4e, 0x47,
	0x45, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x73, 0x22, 0x37, 0x0a, 0x0d, 0x4c, 0x49, 0x4e, 0x44, 0x45, 0x58, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x22, 0x26, 0x0a, 0x0e, 0x4c, 0x49,
	0x4e, 0x44, 0x45, 0x58, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x22, 0x1f, 0x0a, 0x0b, 0x4c, 0x50, 0x4f, 0x50, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x22, 0x24, 0x0a, 0x0c, 0x4c, 0x50, 0x4f, 0x50, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x1f, 0x0a, 0x0b, 0x52, 0x50, 0x4f,
	0x50, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x22, 0x24, 0x0a, 0x0c, 0x52, 0x50,
	0x4f, 0x50, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x22, 0x1f, 0x0a, 0x0b, 0x4c, 0x4c, 0x45, 0x4e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x22, 0x26, 0x0a, 0x0c, 0x4c, 0x4c, 0x45, 0x4e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x16, 0x0a, 0x06, 0x6c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x06, 0x6c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x32, 0xcb, 0x03, 0x0a, 0x0c, 0x48, 0x61,
	0x73, 0x68, 0x44, 0x61, 0x74, 0x61, 0x62, 0x61, 0x73, 0x65, 0x12, 0x3e, 0x0a, 0x05, 0x52, 0x50,
	0x55, 0x53, 0x48, 0x12, 0x19, 0x2e, 0x61, 0x70, 0x69, 0x5f, 0x64, 0x62, 0x5f, 0x6c, 0x69, 0x73,
	0x74, 0x2e, 0x52, 0x50, 0x55, 0x53, 0x48, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a,
	0x2e, 0x61, 0x70, 0x69, 0x5f, 0x64, 0x62, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x2e, 0x52, 0x50, 0x55,
	0x53, 0x48, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3e, 0x0a, 0x05, 0x4c, 0x50,
	0x55, 0x53, 0x48, 0x12, 0x19, 0x2e, 0x61, 0x70, 0x69, 0x5f, 0x64, 0x62, 0x5f, 0x6c, 0x69, 0x73,
	0x74, 0x2e, 0x4c, 0x50, 0x55, 0x53, 0x48, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a,
	0x2e, 0x61, 0x70, 0x69, 0x5f, 0x64, 0x62, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x2e, 0x4c, 0x50, 0x55,
	0x53, 0x48, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a, 0x06, 0x4c, 0x52,
	0x41, 0x4e, 0x47, 0x45, 0x12, 0x1a, 0x2e, 0x61, 0x70, 0x69, 0x5f, 0x64, 0x62, 0x5f, 0x6c, 0x69,
	0x73, 0x74, 0x2e, 0x4c, 0x52, 0x41, 0x4e, 0x47, 0x45, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1b, 0x2e, 0x61, 0x70, 0x69, 0x5f, 0x64, 0x62, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x2e, 0x4c,
	0x52, 0x41, 0x4e, 0x47, 0x45, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x41, 0x0a,
	0x06, 0x4c, 0x49, 0x4e, 0x44, 0x45, 0x58, 0x12, 0x1a, 0x2e, 0x61, 0x70, 0x69, 0x5f, 0x64, 0x62,
	0x5f, 0x6c, 0x69, 0x73, 0x74, 0x2e, 0x4c, 0x49, 0x4e, 0x44, 0x45, 0x58, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x61, 0x70, 0x69, 0x5f, 0x64, 0x62, 0x5f, 0x6c, 0x69, 0x73,
	0x74, 0x2e, 0x4c, 0x49, 0x4e, 0x44, 0x45, 0x58, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x3b, 0x0a, 0x04, 0x4c, 0x50, 0x4f, 0x50, 0x12, 0x18, 0x2e, 0x61, 0x70, 0x69, 0x5f, 0x64,
	0x62, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x2e, 0x4c, 0x50, 0x4f, 0x50, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x19, 0x2e, 0x61, 0x70, 0x69, 0x5f, 0x64, 0x62, 0x5f, 0x6c, 0x69, 0x73, 0x74,
	0x2e, 0x4c, 0x50, 0x4f, 0x50, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3b, 0x0a,
	0x04, 0x52, 0x50, 0x4f, 0x50, 0x12, 0x18, 0x2e, 0x61, 0x70, 0x69, 0x5f, 0x64, 0x62, 0x5f, 0x6c,
	0x69, 0x73, 0x74, 0x2e, 0x52, 0x50, 0x4f, 0x50, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x19, 0x2e, 0x61, 0x70, 0x69, 0x5f, 0x64, 0x62, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x2e, 0x52, 0x50,
	0x4f, 0x50, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3b, 0x0a, 0x04, 0x4c, 0x4c,
	0x45, 0x4e, 0x12, 0x18, 0x2e, 0x61, 0x70, 0x69, 0x5f, 0x64, 0x62, 0x5f, 0x6c, 0x69, 0x73, 0x74,
	0x2e, 0x4c, 0x4c, 0x45, 0x4e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x61,
	0x70, 0x69, 0x5f, 0x64, 0x62, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x2e, 0x4c, 0x4c, 0x45, 0x4e, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0d, 0x5a, 0x0b, 0x61, 0x70, 0x69, 0x2f, 0x64,
	0x62, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_db_list_proto_rawDescOnce sync.Once
	file_db_list_proto_rawDescData = file_db_list_proto_rawDesc
)

func file_db_list_proto_rawDescGZIP() []byte {
	file_db_list_proto_rawDescOnce.Do(func() {
		file_db_list_proto_rawDescData = protoimpl.X.CompressGZIP(file_db_list_proto_rawDescData)
	})
	return file_db_list_proto_rawDescData
}

var file_db_list_proto_msgTypes = make([]protoimpl.MessageInfo, 14)
var file_db_list_proto_goTypes = []interface{}{
	(*RPUSHRequest)(nil),   // 0: api_db_list.RPUSHRequest
	(*RPUSHResponse)(nil),  // 1: api_db_list.RPUSHResponse
	(*LPUSHRequest)(nil),   // 2: api_db_list.LPUSHRequest
	(*LPUSHResponse)(nil),  // 3: api_db_list.LPUSHResponse
	(*LRANGERequest)(nil),  // 4: api_db_list.LRANGERequest
	(*LRANGEResponse)(nil), // 5: api_db_list.LRANGEResponse
	(*LINDEXRequest)(nil),  // 6: api_db_list.LINDEXRequest
	(*LINDEXResponse)(nil), // 7: api_db_list.LINDEXResponse
	(*LPOPRequest)(nil),    // 8: api_db_list.LPOPRequest
	(*LPOPResponse)(nil),   // 9: api_db_list.LPOPResponse
	(*RPOPRequest)(nil),    // 10: api_db_list.RPOPRequest
	(*RPOPResponse)(nil),   // 11: api_db_list.RPOPResponse
	(*LLENRequest)(nil),    // 12: api_db_list.LLENRequest
	(*LLENResponse)(nil),   // 13: api_db_list.LLENResponse
}
var file_db_list_proto_depIdxs = []int32{
	0,  // 0: api_db_list.HashDatabase.RPUSH:input_type -> api_db_list.RPUSHRequest
	2,  // 1: api_db_list.HashDatabase.LPUSH:input_type -> api_db_list.LPUSHRequest
	4,  // 2: api_db_list.HashDatabase.LRANGE:input_type -> api_db_list.LRANGERequest
	6,  // 3: api_db_list.HashDatabase.LINDEX:input_type -> api_db_list.LINDEXRequest
	8,  // 4: api_db_list.HashDatabase.LPOP:input_type -> api_db_list.LPOPRequest
	10, // 5: api_db_list.HashDatabase.RPOP:input_type -> api_db_list.RPOPRequest
	12, // 6: api_db_list.HashDatabase.LLEN:input_type -> api_db_list.LLENRequest
	1,  // 7: api_db_list.HashDatabase.RPUSH:output_type -> api_db_list.RPUSHResponse
	3,  // 8: api_db_list.HashDatabase.LPUSH:output_type -> api_db_list.LPUSHResponse
	5,  // 9: api_db_list.HashDatabase.LRANGE:output_type -> api_db_list.LRANGEResponse
	7,  // 10: api_db_list.HashDatabase.LINDEX:output_type -> api_db_list.LINDEXResponse
	9,  // 11: api_db_list.HashDatabase.LPOP:output_type -> api_db_list.LPOPResponse
	11, // 12: api_db_list.HashDatabase.RPOP:output_type -> api_db_list.RPOPResponse
	13, // 13: api_db_list.HashDatabase.LLEN:output_type -> api_db_list.LLENResponse
	7,  // [7:14] is the sub-list for method output_type
	0,  // [0:7] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_db_list_proto_init() }
func file_db_list_proto_init() {
	if File_db_list_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_db_list_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RPUSHRequest); i {
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
		file_db_list_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RPUSHResponse); i {
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
		file_db_list_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LPUSHRequest); i {
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
		file_db_list_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LPUSHResponse); i {
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
		file_db_list_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LRANGERequest); i {
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
		file_db_list_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LRANGEResponse); i {
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
		file_db_list_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LINDEXRequest); i {
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
		file_db_list_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LINDEXResponse); i {
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
		file_db_list_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LPOPRequest); i {
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
		file_db_list_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LPOPResponse); i {
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
		file_db_list_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RPOPRequest); i {
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
		file_db_list_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RPOPResponse); i {
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
		file_db_list_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LLENRequest); i {
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
		file_db_list_proto_msgTypes[13].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LLENResponse); i {
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
			RawDescriptor: file_db_list_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   14,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_db_list_proto_goTypes,
		DependencyIndexes: file_db_list_proto_depIdxs,
		MessageInfos:      file_db_list_proto_msgTypes,
	}.Build()
	File_db_list_proto = out.File
	file_db_list_proto_rawDesc = nil
	file_db_list_proto_goTypes = nil
	file_db_list_proto_depIdxs = nil
}
