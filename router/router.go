package router

import (
	"net/http"

	"github.com/saiya/mesh_for_home_server/interfaces"
)

type router struct {
	http http.RoundTripper
}

func NewRouter() interfaces.Router {
	return &router{
		// TODO: Implement HTTP RoundTripper
	}
}

func (r *router) HTTP() http.RoundTripper {
	return r.http
}
