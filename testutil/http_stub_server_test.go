package testutil_test

import (
	"net/http"
	"testing"

	"github.com/saiya/mesh_for_home_server/testutil"
)

func TestStubServer(t *testing.T) {
	httpServer, port := testutil.NewHTTPStubServer(t, []testutil.HttpStubPattern{
		{"GET", "/get", map[string]string{"param1": "a", "param2": "日本語"}, "", 200, nil},
	}...)
	defer httpServer.Close()

	httpClient := &http.Client{}

	for _, c := range []testutil.HttpTestCase{
		{"GET", "/get?param1=a&param2=日本語", nil, 200, ""},
	} {
		t.Run(c.String(), func(t *testing.T) { c.Do(t, httpClient, port) })
	}
}
