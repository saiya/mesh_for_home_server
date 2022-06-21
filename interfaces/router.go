package interfaces

import (
	"context"
	"errors"
	"net/http"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/peering/proto/generated"
)

var ErrUnknownPeer = errors.New("no route to deliver")
var ErrPeerNoConnection = errors.New("no connection available for the peer")

type RouterShink = func(context.Context, Message) error
type RouterListener = func(ctx context.Context, from config.NodeID, msg Message) error
type RouterUnregister = func()

type Router interface {
	// NodeID returns this node itself's ID
	NodeID() config.NodeID

	Update(ctx context.Context, routing *generated.RoutingMessage)
	// If no route found, returns ErrUnknownPeer
	RouteHTTPRequest(req *http.Request) (config.NodeID, error)

	// Deliver transfer given message or handle that message in this node itself
	Deliver(ctx context.Context, from config.NodeID, dest config.NodeID, msg Message)

	// RegisterSink registers message deliverery route to another node.
	// Given destination must not equal to this node's ID.
	RegisterSink(dest config.NodeID, callback RouterShink) RouterUnregister

	// Listen regisgers callback for incoming messsage.
	// Lisnter will be called for all mesages that destination is this node.
	// All listeners called every time, should close unnesessary listener.
	Listen(callback RouterListener) RouterUnregister
}
