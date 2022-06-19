package interfaces

import (
	"reflect"

	"github.com/saiya/mesh_for_home_server/peering/proto/generated"
)

type Message = *generated.PeerMessage

func MsgLogString(msg Message) string {
	return reflect.TypeOf(msg).String()
}
