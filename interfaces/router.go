package interfaces

import (
	"context"
	"errors"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/peering/proto/generated"
)

var ErrUnknownPeer = errors.New("no route to deliver")
var ErrPeerNoConnection = errors.New("no connection available for the peer")
var ErrNoRoute = errors.New("no route found to handle given request")

type RouterShink = func(context.Context, Message) error
type RouterListener = func(ctx context.Context, from config.NodeID, msg Message) error
type RouterUnregister = func()

type Advertisement *generated.Advertisement
type Advertiser = func(ctx context.Context) (Advertisement, error)

type Router interface {
	// NodeID returns this node itself's ID
	NodeID() config.NodeID

	// Advertisement returns this node's Advertisement (= this node's capability)
	GenerateAdvertisement(ctx context.Context) (Advertisement, error)
	// Update routing based on incoming advertisement
	Update(ctx context.Context, node config.NodeID, ad Advertisement)
	// Route decides destination route
	// Returns ErrNoRoute if no route found
	Route(ctx context.Context, request Message) (config.NodeID, error)

	// Deliver transfer given message or handle that message in this node itself
	Deliver(ctx context.Context, from config.NodeID, dest config.NodeID, msg Message)

	// RegisterSink registers message deliverery route to another node.
	// Given destination must not equal to this node's ID.
	RegisterSink(dest config.NodeID, callback RouterShink) RouterUnregister

	// Listen regisgers callback for incoming messsage.
	// Lisnter will be called for all mesages that destination is this node.
	// All listeners called every time, should close unnesessary listener.
	Listen(callback RouterListener) RouterUnregister

	Close(ctx context.Context) error
}
