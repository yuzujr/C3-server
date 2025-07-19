package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yuzujr/C3/internal/handler"
)

// RegisterWebRoutes 注册所有面向 Web 的路由
func RegisterWebRoutes(r *gin.RouterGroup) {
	// 客户端列表
	r.GET("/clients", handler.HandleGetClients)

	// 客户端配置
	r.GET("/config/:client_id", handler.HandleGetClientConfig)

	// 更新客户端别名
	r.PUT("/clients/:client_id/alias", handler.HandleUpdateClientAlias)

	// 删除客户端
	r.DELETE("/clients/:client_id", handler.HandleDeleteClient)

	// 获取截图
	r.GET("/screenshots/:client_id", handler.HandleGetClientScreenshots)

	// 删除指定时间范围内的截图
	r.DELETE("/screenshots/:client_id", handler.HandleDeleteScreenshotsByTime)

	// 删除客户端所有截图
	r.DELETE("/all-screenshots/:client_id", handler.HandleDeleteAllScreenshots)
}
