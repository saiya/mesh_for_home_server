// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: peering/proto/peering.proto

package generated

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// PeeringClient is the client API for Peering service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PeeringClient interface {
	Peer(ctx context.Context, opts ...grpc.CallOption) (Peering_PeerClient, error)
}

type peeringClient struct {
	cc grpc.ClientConnInterface
}

func NewPeeringClient(cc grpc.ClientConnInterface) PeeringClient {
	return &peeringClient{cc}
}

func (c *peeringClient) Peer(ctx context.Context, opts ...grpc.CallOption) (Peering_PeerClient, error) {
	stream, err := c.cc.NewStream(ctx, &Peering_ServiceDesc.Streams[0], "/test.Peering/Peer", opts...)
	if err != nil {
		return nil, err
	}
	x := &peeringPeerClient{stream}
	return x, nil
}

type Peering_PeerClient interface {
	Send(*PeerClientMessage) error
	Recv() (*PeerServerMessage, error)
	grpc.ClientStream
}

type peeringPeerClient struct {
	grpc.ClientStream
}

func (x *peeringPeerClient) Send(m *PeerClientMessage) error {
	return x.ClientStream.SendMsg(m)
}

func (x *peeringPeerClient) Recv() (*PeerServerMessage, error) {
	m := new(PeerServerMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PeeringServer is the server API for Peering service.
// All implementations should embed UnimplementedPeeringServer
// for forward compatibility
type PeeringServer interface {
	Peer(Peering_PeerServer) error
}

// UnimplementedPeeringServer should be embedded to have forward compatible implementations.
type UnimplementedPeeringServer struct {
}

func (UnimplementedPeeringServer) Peer(Peering_PeerServer) error {
	return status.Errorf(codes.Unimplemented, "method Peer not implemented")
}

// UnsafePeeringServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PeeringServer will
// result in compilation errors.
type UnsafePeeringServer interface {
	mustEmbedUnimplementedPeeringServer()
}

func RegisterPeeringServer(s grpc.ServiceRegistrar, srv PeeringServer) {
	s.RegisterService(&Peering_ServiceDesc, srv)
}

func _Peering_Peer_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(PeeringServer).Peer(&peeringPeerServer{stream})
}

type Peering_PeerServer interface {
	Send(*PeerServerMessage) error
	Recv() (*PeerClientMessage, error)
	grpc.ServerStream
}

type peeringPeerServer struct {
	grpc.ServerStream
}

func (x *peeringPeerServer) Send(m *PeerServerMessage) error {
	return x.ServerStream.SendMsg(m)
}

func (x *peeringPeerServer) Recv() (*PeerClientMessage, error) {
	m := new(PeerClientMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Peering_ServiceDesc is the grpc.ServiceDesc for Peering service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Peering_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "test.Peering",
	HandlerType: (*PeeringServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Peer",
			Handler:       _Peering_Peer_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "peering/proto/peering.proto",
}
