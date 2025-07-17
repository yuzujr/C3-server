package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/yuzujr/C3/internal/handler"
)

// RegisterAuthRoutes 注册认证相关接口
func RegisterAuthRoutes(r *gin.RouterGroup) {
	r.POST("/login", handler.HandleLogin)
	r.POST("/logout", handler.HandleLogout)
	r.GET("/session", handler.HandleGetSessionInfo)
}
