//
// Use this to re-generate file: make grpc
//

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: peering/proto/peering.proto

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

type CloseByServerReason int32

const (
	// Connection down detected (client or network transit closed connection)
	CloseByServerReason_CONNECTION_DOWN   CloseByServerReason = 0
	CloseByServerReason_HANDSHAKE_FAILURE CloseByServerReason = 1
	CloseByServerReason_INVALID_MESSAGE   CloseByServerReason = 2
)

// Enum value maps for CloseByServerReason.
var (
	CloseByServerReason_name = map[int32]string{
		0: "CONNECTION_DOWN",
		1: "HANDSHAKE_FAILURE",
		2: "INVALID_MESSAGE",
	}
	CloseByServerReason_value = map[string]int32{
		"CONNECTION_DOWN":   0,
		"HANDSHAKE_FAILURE": 1,
		"INVALID_MESSAGE":   2,
	}
)

func (x CloseByServerReason) Enum() *CloseByServerReason {
	p := new(CloseByServerReason)
	*p = x
	return p
}

func (x CloseByServerReason) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CloseByServerReason) Descriptor() protoreflect.EnumDescriptor {
	return file_peering_proto_peering_proto_enumTypes[0].Descriptor()
}

func (CloseByServerReason) Type() protoreflect.EnumType {
	return &file_peering_proto_peering_proto_enumTypes[0]
}

func (x CloseByServerReason) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CloseByServerReason.Descriptor instead.
func (CloseByServerReason) EnumDescriptor() ([]byte, []int) {
	return file_peering_proto_peering_proto_rawDescGZIP(), []int{0}
}

type PeerClientMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Message:
	//
	//	*PeerClientMessage_PeerMessage
	//	*PeerClientMessage_ClientHello
	//	*PeerClientMessage_HandshakeDone
	Message isPeerClientMessage_Message `protobuf_oneof:"message"`
}

func (x *PeerClientMessage) Reset() {
	*x = PeerClientMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peering_proto_peering_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PeerClientMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PeerClientMessage) ProtoMessage() {}

func (x *PeerClientMessage) ProtoReflect() protoreflect.Message {
	mi := &file_peering_proto_peering_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PeerClientMessage.ProtoReflect.Descriptor instead.
func (*PeerClientMessage) Descriptor() ([]byte, []int) {
	return file_peering_proto_peering_proto_rawDescGZIP(), []int{0}
}

func (m *PeerClientMessage) GetMessage() isPeerClientMessage_Message {
	if m != nil {
		return m.Message
	}
	return nil
}

func (x *PeerClientMessage) GetPeerMessage() *PeerMessage {
	if x, ok := x.GetMessage().(*PeerClientMessage_PeerMessage); ok {
		return x.PeerMessage
	}
	return nil
}

func (x *PeerClientMessage) GetClientHello() *ClientHello {
	if x, ok := x.GetMessage().(*PeerClientMessage_ClientHello); ok {
		return x.ClientHello
	}
	return nil
}

func (x *PeerClientMessage) GetHandshakeDone() *HandShakeDone {
	if x, ok := x.GetMessage().(*PeerClientMessage_HandshakeDone); ok {
		return x.HandshakeDone
	}
	return nil
}

type isPeerClientMessage_Message interface {
	isPeerClientMessage_Message()
}

type PeerClientMessage_PeerMessage struct {
	PeerMessage *PeerMessage `protobuf:"bytes,1,opt,name=peer_message,json=peerMessage,proto3,oneof"`
}

type PeerClientMessage_ClientHello struct {
	ClientHello *ClientHello `protobuf:"bytes,2,opt,name=client_hello,json=clientHello,proto3,oneof"`
}

type PeerClientMessage_HandshakeDone struct {
	HandshakeDone *HandShakeDone `protobuf:"bytes,3,opt,name=handshake_done,json=handshakeDone,proto3,oneof"`
}

func (*PeerClientMessage_PeerMessage) isPeerClientMessage_Message() {}

func (*PeerClientMessage_ClientHello) isPeerClientMessage_Message() {}

func (*PeerClientMessage_HandshakeDone) isPeerClientMessage_Message() {}

type PeerServerMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Message:
	//
	//	*PeerServerMessage_PeerMessage
	//	*PeerServerMessage_ServerHello
	//	*PeerServerMessage_CloseByServer
	Message isPeerServerMessage_Message `protobuf_oneof:"message"`
}

func (x *PeerServerMessage) Reset() {
	*x = PeerServerMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peering_proto_peering_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PeerServerMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PeerServerMessage) ProtoMessage() {}

func (x *PeerServerMessage) ProtoReflect() protoreflect.Message {
	mi := &file_peering_proto_peering_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PeerServerMessage.ProtoReflect.Descriptor instead.
func (*PeerServerMessage) Descriptor() ([]byte, []int) {
	return file_peering_proto_peering_proto_rawDescGZIP(), []int{1}
}

func (m *PeerServerMessage) GetMessage() isPeerServerMessage_Message {
	if m != nil {
		return m.Message
	}
	return nil
}

func (x *PeerServerMessage) GetPeerMessage() *PeerMessage {
	if x, ok := x.GetMessage().(*PeerServerMessage_PeerMessage); ok {
		return x.PeerMessage
	}
	return nil
}

func (x *PeerServerMessage) GetServerHello() *ServerHello {
	if x, ok := x.GetMessage().(*PeerServerMessage_ServerHello); ok {
		return x.ServerHello
	}
	return nil
}

func (x *PeerServerMessage) GetCloseByServer() *PeerCloseByServer {
	if x, ok := x.GetMessage().(*PeerServerMessage_CloseByServer); ok {
		return x.CloseByServer
	}
	return nil
}

type isPeerServerMessage_Message interface {
	isPeerServerMessage_Message()
}

type PeerServerMessage_PeerMessage struct {
	PeerMessage *PeerMessage `protobuf:"bytes,1,opt,name=peer_message,json=peerMessage,proto3,oneof"`
}

type PeerServerMessage_ServerHello struct {
	ServerHello *ServerHello `protobuf:"bytes,2,opt,name=server_hello,json=serverHello,proto3,oneof"`
}

type PeerServerMessage_CloseByServer struct {
	CloseByServer *PeerCloseByServer `protobuf:"bytes,3,opt,name=close_by_server,json=closeByServer,proto3,oneof"`
}

func (*PeerServerMessage_PeerMessage) isPeerServerMessage_Message() {}

func (*PeerServerMessage_ServerHello) isPeerServerMessage_Message() {}

func (*PeerServerMessage_CloseByServer) isPeerServerMessage_Message() {}

type PeerCloseByServer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Reason CloseByServerReason `protobuf:"varint,1,opt,name=reason,proto3,enum=test.CloseByServerReason" json:"reason,omitempty"`
}

func (x *PeerCloseByServer) Reset() {
	*x = PeerCloseByServer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peering_proto_peering_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PeerCloseByServer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PeerCloseByServer) ProtoMessage() {}

func (x *PeerCloseByServer) ProtoReflect() protoreflect.Message {
	mi := &file_peering_proto_peering_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PeerCloseByServer.ProtoReflect.Descriptor instead.
func (*PeerCloseByServer) Descriptor() ([]byte, []int) {
	return file_peering_proto_peering_proto_rawDescGZIP(), []int{2}
}

func (x *PeerCloseByServer) GetReason() CloseByServerReason {
	if x != nil {
		return x.Reason
	}
	return CloseByServerReason_CONNECTION_DOWN
}

type PeerMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Message:
	//
	//	*PeerMessage_Advertisement
	//	*PeerMessage_Ping
	//	*PeerMessage_Pong
	//	*PeerMessage_Http
	Message isPeerMessage_Message `protobuf_oneof:"message"`
}

func (x *PeerMessage) Reset() {
	*x = PeerMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peering_proto_peering_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PeerMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PeerMessage) ProtoMessage() {}

func (x *PeerMessage) ProtoReflect() protoreflect.Message {
	mi := &file_peering_proto_peering_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PeerMessage.ProtoReflect.Descriptor instead.
func (*PeerMessage) Descriptor() ([]byte, []int) {
	return file_peering_proto_peering_proto_rawDescGZIP(), []int{3}
}

func (m *PeerMessage) GetMessage() isPeerMessage_Message {
	if m != nil {
		return m.Message
	}
	return nil
}

func (x *PeerMessage) GetAdvertisement() *Advertisement {
	if x, ok := x.GetMessage().(*PeerMessage_Advertisement); ok {
		return x.Advertisement
	}
	return nil
}

func (x *PeerMessage) GetPing() *Ping {
	if x, ok := x.GetMessage().(*PeerMessage_Ping); ok {
		return x.Ping
	}
	return nil
}

func (x *PeerMessage) GetPong() *Pong {
	if x, ok := x.GetMessage().(*PeerMessage_Pong); ok {
		return x.Pong
	}
	return nil
}

func (x *PeerMessage) GetHttp() *HttpMessage {
	if x, ok := x.GetMessage().(*PeerMessage_Http); ok {
		return x.Http
	}
	return nil
}

type isPeerMessage_Message interface {
	isPeerMessage_Message()
}

type PeerMessage_Advertisement struct {
	Advertisement *Advertisement `protobuf:"bytes,1,opt,name=advertisement,proto3,oneof"`
}

type PeerMessage_Ping struct {
	Ping *Ping `protobuf:"bytes,2,opt,name=ping,proto3,oneof"`
}

type PeerMessage_Pong struct {
	Pong *Pong `protobuf:"bytes,3,opt,name=pong,proto3,oneof"`
}

type PeerMessage_Http struct {
	Http *HttpMessage `protobuf:"bytes,4,opt,name=http,proto3,oneof"`
}

func (*PeerMessage_Advertisement) isPeerMessage_Message() {}

func (*PeerMessage_Ping) isPeerMessage_Message() {}

func (*PeerMessage_Pong) isPeerMessage_Message() {}

func (*PeerMessage_Http) isPeerMessage_Message() {}

var File_peering_proto_peering_proto protoreflect.FileDescriptor

var file_peering_proto_peering_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x70, 0x65, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x70, 0x65, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x74,
	0x65, 0x73, 0x74, 0x1a, 0x21, 0x70, 0x65, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2f, 0x61, 0x64, 0x76, 0x65, 0x72, 0x74, 0x69, 0x73, 0x65, 0x6d, 0x65, 0x6e, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1d, 0x70, 0x65, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x68, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x18, 0x70, 0x65, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x68, 0x74, 0x74, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x18, 0x70, 0x65, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70,
	0x69, 0x6e, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc2, 0x01, 0x0a, 0x11, 0x50, 0x65,
	0x65, 0x72, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x36, 0x0a, 0x0c, 0x70, 0x65, 0x65, 0x72, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x50, 0x65, 0x65,
	0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x48, 0x00, 0x52, 0x0b, 0x70, 0x65, 0x65, 0x72,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x31, 0x0a, 0x0c, 0x63, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x5f, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e,
	0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x48, 0x00, 0x52, 0x0b, 0x63,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x12, 0x37, 0x0a, 0x0e, 0x68, 0x61,
	0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x5f, 0x64, 0x6f, 0x6e, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x48, 0x61, 0x6e, 0x64, 0x53, 0x68, 0x61, 0x6b, 0x65, 0x44, 0x6f,
	0x6e, 0x65, 0x48, 0x00, 0x52, 0x0d, 0x68, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x44,
	0x6f, 0x6e, 0x65, 0x42, 0x09, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0xcc,
	0x01, 0x0a, 0x11, 0x50, 0x65, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x12, 0x36, 0x0a, 0x0c, 0x70, 0x65, 0x65, 0x72, 0x5f, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x74, 0x65, 0x73,
	0x74, 0x2e, 0x50, 0x65, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x48, 0x00, 0x52,
	0x0b, 0x70, 0x65, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x31, 0x0a, 0x0c,
	0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x5f, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x48, 0x65, 0x6c, 0x6c, 0x6f,
	0x48, 0x00, 0x52, 0x0b, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x12,
	0x41, 0x0a, 0x0f, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x5f, 0x62, 0x79, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e,
	0x50, 0x65, 0x65, 0x72, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x42, 0x79, 0x53, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x48, 0x00, 0x52, 0x0d, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x42, 0x79, 0x53, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x42, 0x09, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x46, 0x0a,
	0x11, 0x50, 0x65, 0x65, 0x72, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x42, 0x79, 0x53, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x12, 0x31, 0x0a, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x19, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x42,
	0x79, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x52, 0x06, 0x72,
	0x65, 0x61, 0x73, 0x6f, 0x6e, 0x22, 0xae, 0x01, 0x0a, 0x0b, 0x50, 0x65, 0x65, 0x72, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x36, 0x0a, 0x0d, 0x61, 0x64, 0x76, 0x65, 0x72, 0x74, 0x69,
	0x73, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x41,
	0x64, 0x76, 0x65, 0x72, 0x74, 0x69, 0x73, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x48, 0x00, 0x52, 0x0d,
	0x61, 0x64, 0x76, 0x65, 0x72, 0x74, 0x69, 0x73, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x1b, 0x0a,
	0x04, 0x70, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x50, 0x69,
	0x6e, 0x67, 0x48, 0x00, 0x52, 0x04, 0x70, 0x69, 0x6e, 0x67, 0x12, 0x1b, 0x0a, 0x04, 0x70, 0x6f,
	0x6e, 0x67, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x50, 0x6f, 0x6e, 0x67, 0x48,
	0x00, 0x52, 0x04, 0x70, 0x6f, 0x6e, 0x67, 0x12, 0x22, 0x0a, 0x04, 0x68, 0x74, 0x74, 0x70, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x48, 0x74, 0x74, 0x70, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x48, 0x00, 0x52, 0x04, 0x68, 0x74, 0x74, 0x70, 0x42, 0x09, 0x0a, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2a, 0x56, 0x0a, 0x13, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x42,
	0x79, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12, 0x13, 0x0a,
	0x0f, 0x43, 0x4f, 0x4e, 0x4e, 0x45, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x44, 0x4f, 0x57, 0x4e,
	0x10, 0x00, 0x12, 0x15, 0x0a, 0x11, 0x48, 0x41, 0x4e, 0x44, 0x53, 0x48, 0x41, 0x4b, 0x45, 0x5f,
	0x46, 0x41, 0x49, 0x4c, 0x55, 0x52, 0x45, 0x10, 0x01, 0x12, 0x13, 0x0a, 0x0f, 0x49, 0x4e, 0x56,
	0x41, 0x4c, 0x49, 0x44, 0x5f, 0x4d, 0x45, 0x53, 0x53, 0x41, 0x47, 0x45, 0x10, 0x02, 0x32, 0x47,
	0x0a, 0x07, 0x50, 0x65, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x12, 0x3c, 0x0a, 0x04, 0x50, 0x65, 0x65,
	0x72, 0x12, 0x17, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x50, 0x65, 0x65, 0x72, 0x43, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1a, 0x17, 0x2e, 0x74, 0x65, 0x73,
	0x74, 0x2e, 0x50, 0x65, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x28, 0x01, 0x30, 0x01, 0x42, 0x19, 0x5a, 0x17, 0x70, 0x65, 0x65, 0x72, 0x69,
	0x6e, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74,
	0x65, 0x64, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_peering_proto_peering_proto_rawDescOnce sync.Once
	file_peering_proto_peering_proto_rawDescData = file_peering_proto_peering_proto_rawDesc
)

func file_peering_proto_peering_proto_rawDescGZIP() []byte {
	file_peering_proto_peering_proto_rawDescOnce.Do(func() {
		file_peering_proto_peering_proto_rawDescData = protoimpl.X.CompressGZIP(file_peering_proto_peering_proto_rawDescData)
	})
	return file_peering_proto_peering_proto_rawDescData
}

var file_peering_proto_peering_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_peering_proto_peering_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_peering_proto_peering_proto_goTypes = []interface{}{
	(CloseByServerReason)(0),  // 0: test.CloseByServerReason
	(*PeerClientMessage)(nil), // 1: test.PeerClientMessage
	(*PeerServerMessage)(nil), // 2: test.PeerServerMessage
	(*PeerCloseByServer)(nil), // 3: test.PeerCloseByServer
	(*PeerMessage)(nil),       // 4: test.PeerMessage
	(*ClientHello)(nil),       // 5: ClientHello
	(*HandShakeDone)(nil),     // 6: HandShakeDone
	(*ServerHello)(nil),       // 7: ServerHello
	(*Advertisement)(nil),     // 8: Advertisement
	(*Ping)(nil),              // 9: Ping
	(*Pong)(nil),              // 10: Pong
	(*HttpMessage)(nil),       // 11: HttpMessage
}
var file_peering_proto_peering_proto_depIdxs = []int32{
	4,  // 0: test.PeerClientMessage.peer_message:type_name -> test.PeerMessage
	5,  // 1: test.PeerClientMessage.client_hello:type_name -> ClientHello
	6,  // 2: test.PeerClientMessage.handshake_done:type_name -> HandShakeDone
	4,  // 3: test.PeerServerMessage.peer_message:type_name -> test.PeerMessage
	7,  // 4: test.PeerServerMessage.server_hello:type_name -> ServerHello
	3,  // 5: test.PeerServerMessage.close_by_server:type_name -> test.PeerCloseByServer
	0,  // 6: test.PeerCloseByServer.reason:type_name -> test.CloseByServerReason
	8,  // 7: test.PeerMessage.advertisement:type_name -> Advertisement
	9,  // 8: test.PeerMessage.ping:type_name -> Ping
	10, // 9: test.PeerMessage.pong:type_name -> Pong
	11, // 10: test.PeerMessage.http:type_name -> HttpMessage
	1,  // 11: test.Peering.Peer:input_type -> test.PeerClientMessage
	2,  // 12: test.Peering.Peer:output_type -> test.PeerServerMessage
	12, // [12:13] is the sub-list for method output_type
	11, // [11:12] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_peering_proto_peering_proto_init() }
func file_peering_proto_peering_proto_init() {
	if File_peering_proto_peering_proto != nil {
		return
	}
	file_peering_proto_advertisement_proto_init()
	file_peering_proto_handshake_proto_init()
	file_peering_proto_http_proto_init()
	file_peering_proto_ping_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_peering_proto_peering_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PeerClientMessage); i {
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
		file_peering_proto_peering_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PeerServerMessage); i {
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
		file_peering_proto_peering_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PeerCloseByServer); i {
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
		file_peering_proto_peering_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PeerMessage); i {
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
	file_peering_proto_peering_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*PeerClientMessage_PeerMessage)(nil),
		(*PeerClientMessage_ClientHello)(nil),
		(*PeerClientMessage_HandshakeDone)(nil),
	}
	file_peering_proto_peering_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*PeerServerMessage_PeerMessage)(nil),
		(*PeerServerMessage_ServerHello)(nil),
		(*PeerServerMessage_CloseByServer)(nil),
	}
	file_peering_proto_peering_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*PeerMessage_Advertisement)(nil),
		(*PeerMessage_Ping)(nil),
		(*PeerMessage_Pong)(nil),
		(*PeerMessage_Http)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_peering_proto_peering_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_peering_proto_peering_proto_goTypes,
		DependencyIndexes: file_peering_proto_peering_proto_depIdxs,
		EnumInfos:         file_peering_proto_peering_proto_enumTypes,
		MessageInfos:      file_peering_proto_peering_proto_msgTypes,
	}.Build()
	File_peering_proto_peering_proto = out.File
	file_peering_proto_peering_proto_rawDesc = nil
	file_peering_proto_peering_proto_goTypes = nil
	file_peering_proto_peering_proto_depIdxs = nil
}
