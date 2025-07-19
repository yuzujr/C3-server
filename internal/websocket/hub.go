package websocket

import (
	"github.com/yuzujr/C3/internal/logger"
	"github.com/yuzujr/C3/internal/models"
	"github.com/yuzujr/C3/internal/service"
)

type hub struct {
	agents     map[string]*client // 客户端连接
	users      map[*client]bool   // 前端用户连接
	register   chan *client
	unregister chan *client
}

var HubInstance = newHub()

func newHub() *hub {
	return &hub{
		agents:     make(map[string]*client),
		users:      make(map[*client]bool),
		register:   make(chan *client),
		unregister: make(chan *client),
	}
}

func (h *hub) Run() {
	logger.Infof("WebSocket hub started")
	for {
		select {
		case client := <-h.register:
			if client.Role == RoleAgent {
				h.agents[client.ID] = client
				service.UpsertClient(&models.Client{
					ClientID:     client.ID,
					OnlineStatus: true,
				})
				logger.Infof("Client %s connected", client.ID)
			} else {
				h.users[client] = true
				logger.Infof("User Websocket connected")
			}
		case client := <-h.unregister:
			if client.Role == RoleAgent {
				delete(h.agents, client.ID)
				service.UpsertClient(&models.Client{
					ClientID:     client.ID,
					OnlineStatus: false,
				})
				close(client.Send)
				logger.Infof("Client %s disconnected", client.ID)
			} else {
				delete(h.users, client)
				close(client.Send)
				logger.Infof("User Websocket disconnected")
			}
		}
	}
}
