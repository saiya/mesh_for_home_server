package testutil

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func NewHTTPStubServer(t *testing.T, patterns ...HttpStub) (*httptest.Server, int) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		matcher := requestToRequestMatcher(t, r)

		for _, p := range patterns {
			if matcher.Match(t, p) {
				p.handle(t, w, r)
				return
			}
		}

		var failMessage strings.Builder
		failMessage.WriteString("None of registered pattern found, unexpected request ")
		failMessage.WriteString(matcher.String())
		failMessage.WriteString("\n")
		failMessage.WriteString("Available stubs are...")
		for _, p := range patterns {
			failMessage.WriteString("\n")
			failMessage.WriteString("  ")
			failMessage.WriteString(p.String())
		}
		assert.Fail(t, failMessage.String())
	}))
	url, err := url.Parse(srv.URL)
	assert.NoError(t, err)
	port, err := strconv.Atoi(url.Port())
	assert.NoError(t, err)
	return srv, port
}

func (p *HttpStub) handle(t *testing.T, w http.ResponseWriter, r *http.Request) {
	if p.ResponseHeaders != nil {
		for k, v := range p.ResponseHeaders {
			w.Header()[k] = v
		}
	}

	w.WriteHeader(p.StatusCode)

	if p.ResponseBody != nil {
		respBody := *p.ResponseBody
		if p.Options.ChunkResponseBodyPer == 0 {
			_, err := w.Write(respBody)
			assert.NoError(t, err)
		} else {
			for i := 0; i < len(respBody); i += p.Options.ChunkResponseBodyPer {
				sliceEnd := i + p.Options.ChunkResponseBodyPer
				if sliceEnd > len(respBody) {
					sliceEnd = len(respBody)
				}

				_, err := w.Write(respBody[i:sliceEnd])
				assert.NoError(t, err)
				w.(http.Flusher).Flush()
			}
		}
	}
}
