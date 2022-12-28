package testutil

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/saiya/mesh_for_home_server/logger"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

type HTTPStub struct {
	Method         string
	Path           string
	QueryParams    map[string]string
	RequestHeaders map[string][]string
	RequestBody    *[]byte

	StatusCode      int
	ResponseHeaders map[string][]string
	ResponseBody    *[]byte

	Options HTTPStubOptions
}

// HTTPStubOptions is an options that usually not need to set
type HTTPStubOptions struct {
	// If not zero, chunk response body per this bytes
	ChunkResponseBodyPer int
}

func (p *HTTPStub) String() string {
	reqBodyLen := 0
	if p.RequestBody != nil {
		reqBodyLen = len(*p.RequestBody)
	}
	return fmt.Sprintf("%s %s ? %v (reqBody: %d bytes)", p.Method, p.Path, p.QueryParams, reqBodyLen)
}

type requestMatcher struct {
	r    *http.Request
	body []byte
}

func (req *requestMatcher) String() string {
	return fmt.Sprintf("%s %s (body: %d bytes)", req.r.Method, req.r.RequestURI, len(req.body))
}

func requestToRequestMatcher(t *testing.T, r *http.Request) requestMatcher {
	logger.Get().Debugw("Stub HTTP server reading request body...")
	bodyBytes, err := io.ReadAll(r.Body)
	logger.Get().Debugw("Stub HTTP server read request body", "bodyBytes", len(bodyBytes), "err", err)
	assert.NoError(t, err)

	return requestMatcher{r, bodyBytes}
}

func (req *requestMatcher) Match(t *testing.T, p HTTPStub) bool {
	r := req.r

	if !strings.EqualFold(p.Method, r.Method) {
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

	if p.RequestHeaders != nil {
		for k, v := range p.RequestHeaders {
			if !slices.Equal(r.Header[k], v) {
				return false
			}
		}
	}

	if p.RequestBody == nil {
		if len(req.body) != 0 {
			return false
		}
	} else {
		if !slices.Equal(*p.RequestBody, req.body) {
			return false
		}
	}

	return true
}
