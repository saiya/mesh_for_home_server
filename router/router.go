package router

import (
	"net/http"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/interfaces"
)

type router struct {
	nodeID config.NodeID

	http http.RoundTripper
}

func NewRouter() interfaces.Router {
	return &router{
		// FIXME: Generate NodeID

		// TODO: Implement HTTP RoundTripper
	}
}

func (r *router) NodeID() config.NodeID {
	return r.nodeID
}

func (r *router) HTTP() http.RoundTripper {
	return r.http
}
