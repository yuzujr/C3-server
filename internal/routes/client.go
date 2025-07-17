package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yuzujr/C3/internal/handler"
)

// RegisterClientRoutes 注册客户端相关接口
func RegisterClientRoutes(r *gin.RouterGroup) {
	r.POST("/screenshot", handler.HandleClientScreenshot)
	r.POST("/client_config", handler.HandleClientConfig)
	r.GET("/ws", handler.HandleClientWebSocketConnection)
}
