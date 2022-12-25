package testutil

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type HttpTestCase struct {
	Method         string
	Path           string
	RequestBody    *string
	ExpectedStatus int
	ExpectedBody   string
}

func (c *HttpTestCase) String() string {
	bodyBytes := -1
	if c.RequestBody != nil {
		bodyBytes = len(*c.RequestBody)
	}

	return fmt.Sprintf("%s %s (body: %d bytes)", c.Method, c.Path, bodyBytes)
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
		reqBody = strings.NewReader(*c.RequestBody)
	}
	req, err := http.NewRequestWithContext(
		Context(t),
		c.Method, fmt.Sprintf("http://localhost:%d%s", port, c.Path),
		reqBody)
	assert.NoError(t, err)
	return req
}

func (c *HttpTestCase) AssertResponse(t *testing.T, res *http.Response) {
	assert.Equal(t, c.ExpectedStatus, res.StatusCode)

	resBody, err := io.ReadAll(res.Body)
	assert.NoError(t, err)
	assert.Equal(t, c.ExpectedBody, string(resBody))
}
