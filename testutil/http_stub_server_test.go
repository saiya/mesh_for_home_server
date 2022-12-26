package testutil_test

import (
	"net/http"
	"testing"

	util "github.com/saiya/mesh_for_home_server/testutil"
)

func TestStubServer(t *testing.T) {
	httpServer, port := util.NewHTTPStubServer(t, util.DEFAULT_HTTP_STUBS...)
	defer httpServer.Close()

	httpClient := &http.Client{}

	for _, c := range util.DEFAULT_HTTP_TEST_CASES {
		t.Run(c.String(), func(t *testing.T) { c.Do(t, httpClient, port) })
	}
}
