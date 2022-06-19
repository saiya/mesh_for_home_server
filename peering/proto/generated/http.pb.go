// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.6.1
// source: peering/proto/http.proto

package generated

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

type HttpRequestStart struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RequestId *RequestID    `protobuf:"bytes,1,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`
	Method    string        `protobuf:"bytes,2,opt,name=method,proto3" json:"method,omitempty"`
	Hostname  string        `protobuf:"bytes,3,opt,name=hostname,proto3" json:"hostname,omitempty"`
	Path      string        `protobuf:"bytes,4,opt,name=path,proto3" json:"path,omitempty"`
	Headers   []*HttpHeader `protobuf:"bytes,5,rep,name=headers,proto3" json:"headers,omitempty"`
}

func (x *HttpRequestStart) Reset() {
	*x = HttpRequestStart{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peering_proto_http_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HttpRequestStart) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HttpRequestStart) ProtoMessage() {}

func (x *HttpRequestStart) ProtoReflect() protoreflect.Message {
	mi := &file_peering_proto_http_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HttpRequestStart.ProtoReflect.Descriptor instead.
func (*HttpRequestStart) Descriptor() ([]byte, []int) {
	return file_peering_proto_http_proto_rawDescGZIP(), []int{0}
}

func (x *HttpRequestStart) GetRequestId() *RequestID {
	if x != nil {
		return x.RequestId
	}
	return nil
}

func (x *HttpRequestStart) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *HttpRequestStart) GetHostname() string {
	if x != nil {
		return x.Hostname
	}
	return ""
}

func (x *HttpRequestStart) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *HttpRequestStart) GetHeaders() []*HttpHeader {
	if x != nil {
		return x.Headers
	}
	return nil
}

type HttpRequestBody struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RequestId *RequestID `protobuf:"bytes,1,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`
	Data      []byte     `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *HttpRequestBody) Reset() {
	*x = HttpRequestBody{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peering_proto_http_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HttpRequestBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HttpRequestBody) ProtoMessage() {}

func (x *HttpRequestBody) ProtoReflect() protoreflect.Message {
	mi := &file_peering_proto_http_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HttpRequestBody.ProtoReflect.Descriptor instead.
func (*HttpRequestBody) Descriptor() ([]byte, []int) {
	return file_peering_proto_http_proto_rawDescGZIP(), []int{1}
}

func (x *HttpRequestBody) GetRequestId() *RequestID {
	if x != nil {
		return x.RequestId
	}
	return nil
}

func (x *HttpRequestBody) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type HttpRequestEnd struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RequestId *RequestID `protobuf:"bytes,1,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`
}

func (x *HttpRequestEnd) Reset() {
	*x = HttpRequestEnd{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peering_proto_http_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HttpRequestEnd) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HttpRequestEnd) ProtoMessage() {}

func (x *HttpRequestEnd) ProtoReflect() protoreflect.Message {
	mi := &file_peering_proto_http_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HttpRequestEnd.ProtoReflect.Descriptor instead.
func (*HttpRequestEnd) Descriptor() ([]byte, []int) {
	return file_peering_proto_http_proto_rawDescGZIP(), []int{2}
}

func (x *HttpRequestEnd) GetRequestId() *RequestID {
	if x != nil {
		return x.RequestId
	}
	return nil
}

type HttpResponseStart struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RequestId *RequestID `protobuf:"bytes,1,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`
	// e.g. "200 OK"
	Status  string        `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	Headers []*HttpHeader `protobuf:"bytes,3,rep,name=headers,proto3" json:"headers,omitempty"`
}

func (x *HttpResponseStart) Reset() {
	*x = HttpResponseStart{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peering_proto_http_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HttpResponseStart) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HttpResponseStart) ProtoMessage() {}

func (x *HttpResponseStart) ProtoReflect() protoreflect.Message {
	mi := &file_peering_proto_http_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HttpResponseStart.ProtoReflect.Descriptor instead.
func (*HttpResponseStart) Descriptor() ([]byte, []int) {
	return file_peering_proto_http_proto_rawDescGZIP(), []int{3}
}

func (x *HttpResponseStart) GetRequestId() *RequestID {
	if x != nil {
		return x.RequestId
	}
	return nil
}

func (x *HttpResponseStart) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *HttpResponseStart) GetHeaders() []*HttpHeader {
	if x != nil {
		return x.Headers
	}
	return nil
}

type HttpResponseBody struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RequestId *RequestID `protobuf:"bytes,1,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`
	Data      []byte     `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *HttpResponseBody) Reset() {
	*x = HttpResponseBody{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peering_proto_http_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HttpResponseBody) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HttpResponseBody) ProtoMessage() {}

func (x *HttpResponseBody) ProtoReflect() protoreflect.Message {
	mi := &file_peering_proto_http_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HttpResponseBody.ProtoReflect.Descriptor instead.
func (*HttpResponseBody) Descriptor() ([]byte, []int) {
	return file_peering_proto_http_proto_rawDescGZIP(), []int{4}
}

func (x *HttpResponseBody) GetRequestId() *RequestID {
	if x != nil {
		return x.RequestId
	}
	return nil
}

func (x *HttpResponseBody) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

type HttpResponseEnd struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RequestId *RequestID `protobuf:"bytes,1,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`
}

func (x *HttpResponseEnd) Reset() {
	*x = HttpResponseEnd{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peering_proto_http_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HttpResponseEnd) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HttpResponseEnd) ProtoMessage() {}

func (x *HttpResponseEnd) ProtoReflect() protoreflect.Message {
	mi := &file_peering_proto_http_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HttpResponseEnd.ProtoReflect.Descriptor instead.
func (*HttpResponseEnd) Descriptor() ([]byte, []int) {
	return file_peering_proto_http_proto_rawDescGZIP(), []int{5}
}

func (x *HttpResponseEnd) GetRequestId() *RequestID {
	if x != nil {
		return x.RequestId
	}
	return nil
}

type HttpHeader struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name   string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Values []string `protobuf:"bytes,2,rep,name=values,proto3" json:"values,omitempty"`
}

func (x *HttpHeader) Reset() {
	*x = HttpHeader{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peering_proto_http_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HttpHeader) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HttpHeader) ProtoMessage() {}

func (x *HttpHeader) ProtoReflect() protoreflect.Message {
	mi := &file_peering_proto_http_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HttpHeader.ProtoReflect.Descriptor instead.
func (*HttpHeader) Descriptor() ([]byte, []int) {
	return file_peering_proto_http_proto_rawDescGZIP(), []int{6}
}

func (x *HttpHeader) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *HttpHeader) GetValues() []string {
	if x != nil {
		return x.Values
	}
	return nil
}

var File_peering_proto_http_proto protoreflect.FileDescriptor

var file_peering_proto_http_proto_rawDesc = []byte{
	0x0a, 0x18, 0x70, 0x65, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x68, 0x74, 0x74, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1a, 0x70, 0x65, 0x65, 0x72,
	0x69, 0x6e, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xac, 0x01, 0x0a, 0x10, 0x48, 0x74, 0x74, 0x70, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x53, 0x74, 0x61, 0x72, 0x74, 0x12, 0x29, 0x0a, 0x0a, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0a, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x44, 0x52, 0x09, 0x72, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x1a,
	0x0a, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61,
	0x74, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x12, 0x25,
	0x0a, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x0b, 0x2e, 0x48, 0x74, 0x74, 0x70, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x52, 0x07, 0x68, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x73, 0x22, 0x50, 0x0a, 0x0f, 0x48, 0x74, 0x74, 0x70, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x29, 0x0a, 0x0a, 0x72, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x44, 0x52, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x3b, 0x0a, 0x0e, 0x48, 0x74, 0x74, 0x70, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x45, 0x6e, 0x64, 0x12, 0x29, 0x0a, 0x0a, 0x72, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x44, 0x52, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x49, 0x64, 0x22, 0x7d, 0x0a, 0x11, 0x48, 0x74, 0x74, 0x70, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x53, 0x74, 0x61, 0x72, 0x74, 0x12, 0x29, 0x0a, 0x0a, 0x72, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x44, 0x52, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x25, 0x0a, 0x07,
	0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e,
	0x48, 0x74, 0x74, 0x70, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x52, 0x07, 0x68, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x73, 0x22, 0x51, 0x0a, 0x10, 0x48, 0x74, 0x74, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x42, 0x6f, 0x64, 0x79, 0x12, 0x29, 0x0a, 0x0a, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x44, 0x52, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x3c, 0x0a, 0x0f, 0x48, 0x74, 0x74, 0x70, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x45, 0x6e, 0x64, 0x12, 0x29, 0x0a, 0x0a, 0x72, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x44, 0x52, 0x09, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x49, 0x64, 0x22, 0x38, 0x0a, 0x0a, 0x48, 0x74, 0x74, 0x70, 0x48, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x42, 0x19,
	0x5a, 0x17, 0x70, 0x65, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_peering_proto_http_proto_rawDescOnce sync.Once
	file_peering_proto_http_proto_rawDescData = file_peering_proto_http_proto_rawDesc
)

func file_peering_proto_http_proto_rawDescGZIP() []byte {
	file_peering_proto_http_proto_rawDescOnce.Do(func() {
		file_peering_proto_http_proto_rawDescData = protoimpl.X.CompressGZIP(file_peering_proto_http_proto_rawDescData)
	})
	return file_peering_proto_http_proto_rawDescData
}

var file_peering_proto_http_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_peering_proto_http_proto_goTypes = []interface{}{
	(*HttpRequestStart)(nil),  // 0: HttpRequestStart
	(*HttpRequestBody)(nil),   // 1: HttpRequestBody
	(*HttpRequestEnd)(nil),    // 2: HttpRequestEnd
	(*HttpResponseStart)(nil), // 3: HttpResponseStart
	(*HttpResponseBody)(nil),  // 4: HttpResponseBody
	(*HttpResponseEnd)(nil),   // 5: HttpResponseEnd
	(*HttpHeader)(nil),        // 6: HttpHeader
	(*RequestID)(nil),         // 7: RequestID
}
var file_peering_proto_http_proto_depIdxs = []int32{
	7, // 0: HttpRequestStart.request_id:type_name -> RequestID
	6, // 1: HttpRequestStart.headers:type_name -> HttpHeader
	7, // 2: HttpRequestBody.request_id:type_name -> RequestID
	7, // 3: HttpRequestEnd.request_id:type_name -> RequestID
	7, // 4: HttpResponseStart.request_id:type_name -> RequestID
	6, // 5: HttpResponseStart.headers:type_name -> HttpHeader
	7, // 6: HttpResponseBody.request_id:type_name -> RequestID
	7, // 7: HttpResponseEnd.request_id:type_name -> RequestID
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_peering_proto_http_proto_init() }
func file_peering_proto_http_proto_init() {
	if File_peering_proto_http_proto != nil {
		return
	}
	file_peering_proto_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_peering_proto_http_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HttpRequestStart); i {
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
		file_peering_proto_http_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HttpRequestBody); i {
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
		file_peering_proto_http_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HttpRequestEnd); i {
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
		file_peering_proto_http_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HttpResponseStart); i {
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
		file_peering_proto_http_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HttpResponseBody); i {
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
		file_peering_proto_http_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HttpResponseEnd); i {
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
		file_peering_proto_http_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HttpHeader); i {
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
			RawDescriptor: file_peering_proto_http_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_peering_proto_http_proto_goTypes,
		DependencyIndexes: file_peering_proto_http_proto_depIdxs,
		MessageInfos:      file_peering_proto_http_proto_msgTypes,
	}.Build()
	File_peering_proto_http_proto = out.File
	file_peering_proto_http_proto_rawDesc = nil
	file_peering_proto_http_proto_goTypes = nil
	file_peering_proto_http_proto_depIdxs = nil
}