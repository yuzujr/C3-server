package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yuzujr/C3/internal/handler"
)

// RegisterWsRoutes 注册客户端相关接口
func RegisterWsRoutes(r *gin.RouterGroup) {
	// 建立WebSocket连接
	r.GET("", handler.HandleWSConnection)
}
