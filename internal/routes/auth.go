package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yuzujr/C3/internal/handler"
)

// RegisterAuthRoutes 注册认证相关接口
func RegisterAuthRoutes(r *gin.RouterGroup) {
	// 前端登录
	r.POST("/login", handler.HandleLogin)

	// 前端登出
	r.POST("/logout", handler.HandleLogout)

	// 获取当前登录用户信息
	r.GET("/session", handler.HandleGetSessionInfo)
}
