package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yuzujr/C3/internal/handler"
)

// RegisterClientRoutes 注册客户端相关接口
func RegisterClientRoutes(r *gin.RouterGroup) {
	// 上传客户端截图
	r.POST("/screenshot", handler.HandleClientScreenshot)

	// 上传客户端配置
	r.POST("/client_config", handler.HandleClientConfig)
}
