package interfaces

import (
	"reflect"

	"github.com/saiya/mesh_for_home_server/peering/proto/generated"
)

type Message = *generated.PeerMessage

func MsgLogString(msg Message) string {
	if http := msg.GetHttp(); http != nil {
		return reflect.TypeOf(msg.GetHttp().Message).String()
	} else {
		return reflect.TypeOf(msg.Message).String()
	}
}
