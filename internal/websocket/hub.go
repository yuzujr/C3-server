package websocket

import (
	"github.com/yuzujr/C3/internal/logger"
	"github.com/yuzujr/C3/internal/models"
	"github.com/yuzujr/C3/internal/service"
)

type hub struct {
	Clients    map[string]*client
	Register   chan *client
	Unregister chan *client
}

var HubInstance = newHub()

func newHub() *hub {
	return &hub{
		Clients:    make(map[string]*client),
		Register:   make(chan *client),
		Unregister: make(chan *client),
	}
}

func (h *hub) Run() {
	logger.Infof("WebSocket hub started")
	for {
		select {
		case client := <-h.Register:
			logger.Infof("Client %s registered", client.ClientID)
			h.Clients[client.ClientID] = client
			service.SetClient(&models.Client{
				ClientID:     client.ClientID,
				OnlineStatus: true,
			})
			logger.Infof("Client %s connected", client.ClientID)
		case client := <-h.Unregister:
			delete(h.Clients, client.ClientID)
			close(client.Send)
			service.SetClient(&models.Client{
				ClientID:     client.ClientID,
				OnlineStatus: false,
			})
			logger.Infof("Client %s disconnected", client.ClientID)
		}
	}
}
