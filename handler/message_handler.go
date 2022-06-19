package handler

import "github.com/saiya/mesh_for_home_server/interfaces"

func NewMessageHandlers(router interfaces.Router) []interfaces.MessageHandler {
	return []interfaces.MessageHandler{
		NewPingHandler(router),
	}
}
