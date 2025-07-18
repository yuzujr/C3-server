package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yuzujr/C3/internal/logger"
	"github.com/yuzujr/C3/internal/websocket"
)

// 处理 WebSocket 连接
func HandleWSConnection(c *gin.Context) {
	connType := c.Query("type")
	switch connType {
	case "web":
		if err := websocket.CreateConn(c, "", websocket.RoleUser); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create WebSocket connection"})
			logger.Errorf("Failed to create WebSocket connection: %v", err)
		}
	case "client":
		// 处理客户端连接
		clientID := c.Query("client_id")

		if clientID == "" {
			c.JSON(400, gin.H{"error": "Missing client_id"})
			return
		}

		if err := websocket.CreateConn(c, clientID, websocket.RoleAgent); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create WebSocket connection"})
			logger.Errorf("Failed to create WebSocket connection: %v", err)
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid connection type"})
	}
}
