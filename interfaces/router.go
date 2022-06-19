package interfaces

import "net/http"

type Router interface {
	HTTP() http.RoundTripper
}
