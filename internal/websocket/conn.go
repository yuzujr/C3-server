package websocket

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func CreateConn(c *gin.Context, id string, role Role) error {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return err
	}

	client := &client{
		Conn: conn,
		Send: make(chan []byte, 256),
		ID:   id,
		Role: role,
	}

	HubInstance.register <- client

	go client.writePump()
	go client.readPump()

	return nil
}

func (c *client) readPump() {
	defer func() {
		HubInstance.unregister <- c
		c.Conn.Close()
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		if c.Role == RoleAgent {
			handleAgentMessage(c.ID, message)
		} else {
			//目前用户不会通过WebSocket发送消息
		}
	}
}

func (c *client) writePump() {
	for msg := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			break
		}
	}
}
