// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.12
// source: peering/proto/handshake.proto

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

type ClientHello struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeId        string         `protobuf:"bytes,1,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
	Advertisement *Advertisement `protobuf:"bytes,10,opt,name=advertisement,proto3" json:"advertisement,omitempty"`
}

func (x *ClientHello) Reset() {
	*x = ClientHello{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peering_proto_handshake_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientHello) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientHello) ProtoMessage() {}

func (x *ClientHello) ProtoReflect() protoreflect.Message {
	mi := &file_peering_proto_handshake_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientHello.ProtoReflect.Descriptor instead.
func (*ClientHello) Descriptor() ([]byte, []int) {
	return file_peering_proto_handshake_proto_rawDescGZIP(), []int{0}
}

func (x *ClientHello) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

func (x *ClientHello) GetAdvertisement() *Advertisement {
	if x != nil {
		return x.Advertisement
	}
	return nil
}

type ServerHello struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NodeId        string         `protobuf:"bytes,1,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
	Advertisement *Advertisement `protobuf:"bytes,10,opt,name=advertisement,proto3" json:"advertisement,omitempty"`
}

func (x *ServerHello) Reset() {
	*x = ServerHello{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peering_proto_handshake_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerHello) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerHello) ProtoMessage() {}

func (x *ServerHello) ProtoReflect() protoreflect.Message {
	mi := &file_peering_proto_handshake_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerHello.ProtoReflect.Descriptor instead.
func (*ServerHello) Descriptor() ([]byte, []int) {
	return file_peering_proto_handshake_proto_rawDescGZIP(), []int{1}
}

func (x *ServerHello) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

func (x *ServerHello) GetAdvertisement() *Advertisement {
	if x != nil {
		return x.Advertisement
	}
	return nil
}

type HandShakeDone struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *HandShakeDone) Reset() {
	*x = HandShakeDone{}
	if protoimpl.UnsafeEnabled {
		mi := &file_peering_proto_handshake_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HandShakeDone) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HandShakeDone) ProtoMessage() {}

func (x *HandShakeDone) ProtoReflect() protoreflect.Message {
	mi := &file_peering_proto_handshake_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HandShakeDone.ProtoReflect.Descriptor instead.
func (*HandShakeDone) Descriptor() ([]byte, []int) {
	return file_peering_proto_handshake_proto_rawDescGZIP(), []int{2}
}

var File_peering_proto_handshake_proto protoreflect.FileDescriptor

var file_peering_proto_handshake_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x70, 0x65, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x68, 0x61, 0x6e, 0x64, 0x73, 0x68, 0x61, 0x6b, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x21, 0x70, 0x65, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61,
	0x64, 0x76, 0x65, 0x72, 0x74, 0x69, 0x73, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x5c, 0x0a, 0x0b, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x48, 0x65, 0x6c, 0x6c,
	0x6f, 0x12, 0x17, 0x0a, 0x07, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x34, 0x0a, 0x0d, 0x61, 0x64,
	0x76, 0x65, 0x72, 0x74, 0x69, 0x73, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0e, 0x2e, 0x41, 0x64, 0x76, 0x65, 0x72, 0x74, 0x69, 0x73, 0x65, 0x6d, 0x65, 0x6e,
	0x74, 0x52, 0x0d, 0x61, 0x64, 0x76, 0x65, 0x72, 0x74, 0x69, 0x73, 0x65, 0x6d, 0x65, 0x6e, 0x74,
	0x22, 0x5c, 0x0a, 0x0b, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x12,
	0x17, 0x0a, 0x07, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x34, 0x0a, 0x0d, 0x61, 0x64, 0x76, 0x65,
	0x72, 0x74, 0x69, 0x73, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0e, 0x2e, 0x41, 0x64, 0x76, 0x65, 0x72, 0x74, 0x69, 0x73, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x52,
	0x0d, 0x61, 0x64, 0x76, 0x65, 0x72, 0x74, 0x69, 0x73, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x22, 0x0f,
	0x0a, 0x0d, 0x48, 0x61, 0x6e, 0x64, 0x53, 0x68, 0x61, 0x6b, 0x65, 0x44, 0x6f, 0x6e, 0x65, 0x42,
	0x19, 0x5a, 0x17, 0x70, 0x65, 0x65, 0x72, 0x69, 0x6e, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_peering_proto_handshake_proto_rawDescOnce sync.Once
	file_peering_proto_handshake_proto_rawDescData = file_peering_proto_handshake_proto_rawDesc
)

func file_peering_proto_handshake_proto_rawDescGZIP() []byte {
	file_peering_proto_handshake_proto_rawDescOnce.Do(func() {
		file_peering_proto_handshake_proto_rawDescData = protoimpl.X.CompressGZIP(file_peering_proto_handshake_proto_rawDescData)
	})
	return file_peering_proto_handshake_proto_rawDescData
}

var file_peering_proto_handshake_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_peering_proto_handshake_proto_goTypes = []interface{}{
	(*ClientHello)(nil),   // 0: ClientHello
	(*ServerHello)(nil),   // 1: ServerHello
	(*HandShakeDone)(nil), // 2: HandShakeDone
	(*Advertisement)(nil), // 3: Advertisement
}
var file_peering_proto_handshake_proto_depIdxs = []int32{
	3, // 0: ClientHello.advertisement:type_name -> Advertisement
	3, // 1: ServerHello.advertisement:type_name -> Advertisement
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_peering_proto_handshake_proto_init() }
func file_peering_proto_handshake_proto_init() {
	if File_peering_proto_handshake_proto != nil {
		return
	}
	file_peering_proto_advertisement_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_peering_proto_handshake_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientHello); i {
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
		file_peering_proto_handshake_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerHello); i {
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
		file_peering_proto_handshake_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HandShakeDone); i {
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
			RawDescriptor: file_peering_proto_handshake_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_peering_proto_handshake_proto_goTypes,
		DependencyIndexes: file_peering_proto_handshake_proto_depIdxs,
		MessageInfos:      file_peering_proto_handshake_proto_msgTypes,
	}.Build()
	File_peering_proto_handshake_proto = out.File
	file_peering_proto_handshake_proto_rawDesc = nil
	file_peering_proto_handshake_proto_goTypes = nil
	file_peering_proto_handshake_proto_depIdxs = nil
}
