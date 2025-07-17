package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yuzujr/C3/internal/service"
)

// 认证中间件
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionId, err := c.Cookie("sessionId")
		if err != nil || sessionId == "" {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		session, err := service.ValidateSession(sessionId)
		if err != nil || session == nil {
			// 清除无效cookie
			c.SetCookie("sessionId", "", -1, "/", "", false, true)
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// 将用户信息添加到 context
		c.Set("user", session.User)
		c.Set("sessionId", sessionId)
		c.Next()
	}
}
