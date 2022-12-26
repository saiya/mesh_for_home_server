package testutil

import "math/rand"

var shortHttpBody1a = RandomBytes(32)
var shortHttpBody1b = RandomBytes(32)

var largeHttpBody1a = RandomBytes(32*1024*1024 + rand.Intn(1024))
var largeHttpBody1b = RandomBytes(32*1024*1024 + rand.Intn(1024))

var DEFAULT_HTTP_TEST_CASES = []HttpTestCase{
	{
		"GET", "/get?param1=a&param2=日本語",
		nil, nil,
		201, nil, nil,
		HttpTestCaseOptions{},
	},
	{
		"POST", "/post",
		map[string][]string{"Test-Req-Header": {"test1", "test2"}, "Test-Req-Header2": {"foo; bar"}}, &shortHttpBody1a,
		200, map[string][]string{"Test-Res-Header": {"test1", "test2"}, "Test-Res-Header2": {"foo; bar"}}, &shortHttpBody1b,
		HttpTestCaseOptions{},
	},
	{
		"POST", "/post/large",
		nil, &largeHttpBody1a,
		200, nil, &largeHttpBody1b,
		HttpTestCaseOptions{},
	},
}

var DEFAULT_HTTP_STUBS = []HttpStub{
	{
		"GET", "/get", map[string]string{"param1": "a", "param2": "日本語"},
		nil, nil,
		201, nil, nil,
		HttpStubOptions{},
	},
	{
		"POST", "/post", nil,
		map[string][]string{"Test-Req-Header": {"test1", "test2"}, "Test-Req-Header2": {"foo; bar"}}, &shortHttpBody1a,
		200, map[string][]string{"Test-Res-Header": {"test1", "test2"}, "Test-Res-Header2": {"foo; bar"}}, &shortHttpBody1b,
		HttpStubOptions{},
	},
	{
		"POST", "/post/large", nil,
		nil, &largeHttpBody1a,
		200, nil, &largeHttpBody1b,
		HttpStubOptions{ChunkResponseBodyPer: 1024},
	},
}
