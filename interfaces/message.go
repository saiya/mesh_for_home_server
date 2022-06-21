package interfaces

import (
	"reflect"

	"github.com/saiya/mesh_for_home_server/peering/proto/generated"
)

type Message = *generated.PeerMessage

func MsgLogString(msg Message) string {
	return reflect.TypeOf(msg.Message).String()
}

func RequestIDOf(packet Message) *generated.RequestID {
	switch msg := packet.Message.(type) {
	case *generated.PeerMessage_HttpRequestStart:
		return msg.HttpRequestStart.RequestId
	case *generated.PeerMessage_HttpRequestBody:
		return msg.HttpRequestBody.RequestId
	case *generated.PeerMessage_HttpRequestEnd:
		return msg.HttpRequestEnd.RequestId
	case *generated.PeerMessage_HttpResponseStart:
		return msg.HttpResponseStart.RequestId
	case *generated.PeerMessage_HttpResponseBody:
		return msg.HttpResponseBody.RequestId
	case *generated.PeerMessage_HttpResponseEnd:
		return msg.HttpResponseEnd.RequestId
	default:
		return nil
	}
}
