package interfaces

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/saiya/mesh_for_home_server/config"
)

var ErrUnknownPeer = errors.New("Unknown peer, don't know where to deliver")

type Router interface {
	// NodeID returns this node itself's ID
	NodeID() config.NodeID

	HTTP() http.RoundTripper

	// Deliver transfer given message or handle that message in this node itself
	Deliver(ctx context.Context, from config.NodeID, dest config.NodeID, msg Message) error

	// Register message deliverer.
	RegisterSink(dest config.NodeID, callback func(context.Context, Message) error) io.Closer

	// Listen regisger callback for incoming messsage.
	// Catches all mesages that destination is this node.
	// All listeners called every time, should close unnesessary listener.
	Listen(callback func(context.Context, config.NodeID, Message) error) io.Closer
}
