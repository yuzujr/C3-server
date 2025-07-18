package websocket

import (
	"github.com/yuzujr/C3/internal/logger"
)

func broadcastMsgToUsers(message []byte) {
	for client := range HubInstance.users {
		select {
		case client.Send <- message:
		default:
			logger.Errorf("Failed to broadcast message: channel full")
		}
	}
}
