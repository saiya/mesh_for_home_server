package testutil

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

type HttpStub struct {
	Method         string
	Path           string
	QueryParams    map[string]string
	RequestHeaders map[string][]string
	RequestBody    *[]byte

	StatusCode      int
	ResponseHeaders map[string][]string
	ResponseBody    *[]byte

	Options HttpStubOptions
}

// Options that usually not need to set
type HttpStubOptions struct {
	// If not zero, chunk response body per this bytes
	ChunkResponseBodyPer int
}

func (p *HttpStub) String() string {
	return fmt.Sprintf("%s %s ? %v", p.Method, p.Path, p.QueryParams)
}

type requestMatcher struct {
	r    *http.Request
	body []byte
}

func (req *requestMatcher) String() string {
	return fmt.Sprintf("%s %s", req.r.Method, req.r.RequestURI)
}

func requestToRequestMatcher(t *testing.T, r *http.Request) requestMatcher {
	bodyBytes, err := io.ReadAll(r.Body)
	assert.NoError(t, err)

	return requestMatcher{r, bodyBytes}
}

func (req *requestMatcher) Match(t *testing.T, p HttpStub) bool {
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
