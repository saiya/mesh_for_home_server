package httpingress

import "net/http"

func logAttributesOfRequest(r *http.Request) []interface{} {
	return []interface{}{
		"method", r.Method, "url", r.URL,
		"remote-addr", r.RemoteAddr, "user-agent", r.UserAgent(),
	}
}
