package testutil_test

import (
	"net/http"
	"testing"

	util "github.com/saiya/mesh_for_home_server/testutil"
)

func TestHTTPStubServer(t *testing.T) {
	httpServer, port := util.NewHTTPStubServer(t, util.DefaultHTTPStubs...)
	defer httpServer.Close()

	httpClient := &http.Client{}

	for _, c := range util.DefaultHTTPTestCases {
		t.Run(c.String(), func(t *testing.T) { c.Do(t, httpClient, port) })
	}
}
