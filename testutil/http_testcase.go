package testutil

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/saiya/mesh_for_home_server/logger"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

type HttpTestCase struct {
	Method         string
	Path           string
	RequestHeaders map[string][]string
	RequestBody    *[]byte

	ExpectedStatus  int
	ExpectedHeaders map[string][]string
	ExpectedBody    *[]byte

	Options HttpTestCaseOptions
}

// Options that usually not need to set
type HttpTestCaseOptions struct{}

func (c *HttpTestCase) String() string {
	bodyBytes := -1
	if c.RequestBody != nil {
		bodyBytes = len(*c.RequestBody)
	}

	return fmt.Sprintf("%s_%s_body%d", c.Method, strings.ReplaceAll(c.Path, "/", "_"), bodyBytes)
}

func (c *HttpTestCase) Do(t *testing.T, httpClient *http.Client, port int) {
	res, err := httpClient.Do(c.ToRequest(t, port))
	if !assert.NoError(t, err) {
		return
	}
	c.AssertResponse(t, res)
}

func (c *HttpTestCase) ToRequest(t *testing.T, port int) *http.Request {
	var reqBody io.Reader
	if c.RequestBody != nil {
		reqBody = bytes.NewReader(*c.RequestBody)
	}
	req, err := http.NewRequestWithContext(
		Context(t),
		c.Method, fmt.Sprintf("http://localhost:%d%s", port, c.Path),
		reqBody,
	)
	if c.RequestHeaders != nil {
		for k, v := range c.RequestHeaders {
			req.Header[k] = v
		}
	}
	assert.NoError(t, err)
	return req
}

func (c *HttpTestCase) AssertResponse(t *testing.T, res *http.Response) {
	if !assert.Equal(t, c.ExpectedStatus, res.StatusCode) {
		return
	}

	if c.ExpectedHeaders != nil {
		for k, v := range c.ExpectedHeaders {
			assert.Equal(t, v, res.Header[k])
		}
	}

	logger.Get().Debugw("Reading HTTP client's response body...")
	resBody, err := io.ReadAll(res.Body)
	logger.Get().Debugw("Read HTTP client's", "resBody", len(resBody), "err", err)
	assert.NoError(t, err)
	if c.ExpectedBody == nil {
		assert.Equal(t, 0, len(resBody))
	} else {
		assert.Equal(t, len(*c.ExpectedBody), len(resBody))

		// To avoid dumping body (it can be long), not pass []byte itself to assert
		if !slices.Equal(*c.ExpectedBody, resBody) {
			assert.Fail(t, "Response body not match")
		}
	}
}
