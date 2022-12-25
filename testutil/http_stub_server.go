package testutil

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type HttpStubPattern struct {
	Method      string
	Path        string
	QueryParams map[string]string
	RequestBody string

	StatusCode   int
	ResponseBody *string
}

func NewHTTPStubServer(t *testing.T, patterns ...HttpStubPattern) (*httptest.Server, int) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, p := range patterns {
			if p.handle(t, w, r) {
				return
			}
		}
		assert.Failf(t, "None of registered pattern found, unexpected request %s %s", r.Method, r.URL)
	}))
	url, err := url.Parse(srv.URL)
	assert.NoError(t, err)
	port, err := strconv.Atoi(url.Port())
	assert.NoError(t, err)
	return srv, port
}

func (p *HttpStubPattern) handle(t *testing.T, w http.ResponseWriter, r *http.Request) bool {
	if strings.ToUpper(p.Method) != strings.ToUpper(r.Method) {
		return false
	}
	if p.Path != r.URL.Path && !strings.HasPrefix(p.Path, r.URL.Path+"?") {
		return false
	}

	if p.QueryParams == nil {
		if len(r.URL.Query()) != 0 {
			return false
		}
	} else {
		for k, v := range p.QueryParams {
			if r.URL.Query().Get(k) != v {
				return false
			}
		}
	}

	bodyBytes, err := io.ReadAll(r.Body)
	assert.NoError(t, err)
	if p.RequestBody != string(bodyBytes) {
		return false
	}

	w.WriteHeader(p.StatusCode)
	if p.ResponseBody != nil {
		w.Write([]byte(*p.ResponseBody))
	}
	return true
}
