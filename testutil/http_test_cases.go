package testutil

import "math/rand"

var shortHTTPBody1a = RandomBytes(32)
var shortHTTPBody1b = RandomBytes(32)

var largeHTTPBody1a = RandomBytes(32*1024*1024 + rand.Intn(1024))
var largeHTTPBody1b = RandomBytes(32*1024*1024 + rand.Intn(1024))

var DefaultHTTPTestCases = []HTTPTestCase{
	{
		"GET", "/get?param1=a&param2=日本語",
		nil, nil,
		201, nil, nil,
		HTTPTestCaseOptions{},
	},
	{
		"POST", "/post",
		map[string][]string{"Test-Req-Header": {"test1", "test2"}, "Test-Req-Header2": {"foo; bar"}}, &shortHTTPBody1a,
		200, map[string][]string{"Test-Res-Header": {"test1", "test2"}, "Test-Res-Header2": {"foo; bar"}}, &shortHTTPBody1b,
		HTTPTestCaseOptions{},
	},
	{
		"POST", "/post/large",
		nil, &largeHTTPBody1a,
		200, nil, &largeHTTPBody1b,
		HTTPTestCaseOptions{},
	},
}

var DefaultHTTPStubs = []HTTPStub{
	{
		"GET", "/get", map[string]string{"param1": "a", "param2": "日本語"},
		nil, nil,
		201, nil, nil,
		HTTPStubOptions{},
	},
	{
		"POST", "/post", nil,
		map[string][]string{"Test-Req-Header": {"test1", "test2"}, "Test-Req-Header2": {"foo; bar"}}, &shortHTTPBody1a,
		200, map[string][]string{"Test-Res-Header": {"test1", "test2"}, "Test-Res-Header2": {"foo; bar"}}, &shortHTTPBody1b,
		HTTPStubOptions{},
	},
	{
		"POST", "/post/large", nil,
		nil, &largeHTTPBody1a,
		200, nil, &largeHTTPBody1b,
		HTTPStubOptions{ChunkResponseBodyPer: 1024},
	},
}
