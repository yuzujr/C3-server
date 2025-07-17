package websocket

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/yuzujr/C3/internal/logger"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ServeWs(c *gin.Context) {
	connType := c.Query("type")
	switch connType {
	case "web":
		//not impl yet
		//handleWebConnection(c)
		return
	case "client":
		// 处理客户端连接
		handleClientConnection(c)
		return
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid connection type"})
		return
	}
}

func handleClientConnection(c *gin.Context) {
	clientID := c.Query("client_id")

	if clientID == "" {
		c.JSON(400, gin.H{"error": "Missing client_id"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		logger.Errorf("Failed to upgrade connection: %v", err)
		return
	}

	client := &client{
		Conn: conn,
		Send: make(chan []byte, 256),
		ID:   clientID,
	}

	HubInstance.Register <- client

	go client.writePump()
	go client.readPump()
}
