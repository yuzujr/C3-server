package websocket

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/yuzujr/C3/internal/logger"
)

type client struct {
	Conn *websocket.Conn
	Send chan []byte
	ID   string
}

type Command struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

func SendCommandToClient(clientID string, message Command) {
	// 将命令转换为JSON格式
	messageBytes, err := json.Marshal(message)
	if err != nil {
		logger.Errorf("Failed to marshal message for client %s: %v", clientID, err)
		return
	}

	// 查找客户端并发送消息
	client, ok := HubInstance.Clients[clientID]
	if !ok {
		logger.Errorf("Client %s not found", clientID)
		return
	}
	select {
	case client.Send <- messageBytes:
	default:
		logger.Errorf("Failed to send message to client %s: channel full", clientID)
	}
}

func (c *client) readPump() {
	defer func() {
		HubInstance.Unregister <- c
		c.Conn.Close()
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		handleClientResponse(c.ID, message)
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

func handleClientResponse(clientID string, message []byte) {
	// 处理客户端响应逻辑
	// 此处应该转发给web websocket，根据requestId找到web的websocket连接
	logger.Infof("Received message from client %s: %s", clientID, message)
}
