package proto

import (
	"net/http"

	"github.com/saiya/mesh_for_home_server/peering/proto/generated"
)

func FromHTTPHeaders(headers http.Header) []*generated.HttpHeader {
	result := make([]*generated.HttpHeader, 0, len(headers))
	for k, v := range headers {
		result = append(result, &generated.HttpHeader{Name: k, Values: v})
	}
	return result
}

func ToHTTPHeaders(headers []*generated.HttpHeader) http.Header {
	result := make(http.Header, len(headers))
	for _, h := range headers {
		result[h.Name] = h.Values
	}
	return result
}

func FromHTTPResponse(r *http.Response) *generated.HttpResponseStart {
	return &generated.HttpResponseStart{
		Status:            r.Status,
		StatusCode:        int32(r.StatusCode),
		Proto:             r.Proto,
		ProtoMajor:        int32(r.ProtoMajor),
		ProtoMinor:        int32(r.ProtoMinor),
		Headers:           FromHTTPHeaders(r.Header),
		ContentLength:     r.ContentLength,
		TransferEncodings: r.TransferEncoding,
	}
}

func ToHTTPResponse(r *generated.HttpResponseStart) *http.Response {
	return &http.Response{
		Status:           r.Status,
		StatusCode:       int(r.StatusCode),
		Proto:            r.Proto,
		ProtoMajor:       int(r.ProtoMajor),
		ProtoMinor:       int(r.ProtoMinor),
		Header:           ToHTTPHeaders(r.Headers),
		Body:             nil, // Should be set by caller
		ContentLength:    r.ContentLength,
		TransferEncoding: r.TransferEncodings,
		Close:            false,
		Uncompressed:     false,
		Trailer:          nil, // Should be set by caller
		Request:          nil,
		TLS:              nil,
	}
}
