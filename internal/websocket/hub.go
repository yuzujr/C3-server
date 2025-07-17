package websocket

import (
	"github.com/yuzujr/C3/internal/logger"
	"github.com/yuzujr/C3/internal/models"
	"github.com/yuzujr/C3/internal/services"
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
			logger.Infof("Client %s registered", client.ID)
			h.Clients[client.ID] = client
			services.SetClient(client.ID, &models.Client{
				OnlineStatus: true,
			})
			logger.Infof("Client %s connected", client.ID)
		case client := <-h.Unregister:
			delete(h.Clients, client.ID)
			close(client.Send)
			services.SetClient(client.ID, &models.Client{
				OnlineStatus: false,
			})
			logger.Infof("Client %s disconnected", client.ID)
		}
	}
}
