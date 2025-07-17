package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yuzujr/C3/internal/config"
	"github.com/yuzujr/C3/internal/logger"
	"github.com/yuzujr/C3/internal/models"
	"github.com/yuzujr/C3/internal/service"
)

// 登录接口
func HandleLogin(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password required"})
		return
	}

	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password cannot be empty"})
		return
	}

	// 登录验证，创建会话
	loginResult, err := service.Login(req.Username, req.Password)
	if err != nil || !loginResult.Success {
		logger.Errorf("User: " + req.Username + " ip: " + c.ClientIP() + " login failed - " + loginResult.Message)
		c.JSON(http.StatusUnauthorized, gin.H{"error": loginResult.Message})
		return
	}

	// 设置cookie
	c.SetCookie(
		"sessionId",           // 名称
		loginResult.SessionID, // 值
		int(config.Get().Auth.SessionExpireHours)*60*60, // 有效期（秒），这里是配置的小时数转换为秒
		"/",   // 路径
		"",    // 域名，留空为当前域
		false, // secure，开发环境为false
		true,  // httpOnly
	)

	logger.Infof("User: " + req.Username + " ip: " + c.ClientIP() + " login")
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Login successful"})
}

// 登出接口
func HandleLogout(c *gin.Context) {
	// 获取会话ID
	sessionId, success := getSessionId(c)
	if !success {
		return
	}

	// 获取用户名
	username := "unknown"
	user, success := getUser(c)
	if success {
		username = user.Username
	}
	if username == "unknown" {
		// 如果没有从 context 获取到用户信息，则尝试从 session 中获取
		session, err := service.ValidateSession(sessionId)
		if err == nil && session != nil {
			username = session.User.Username
		}
	}

	// 注销会话
	if err := service.Logout(sessionId); err != nil {
		logger.Errorf("Session destruction failed: %v", err)
	}

	// 清除cookie
	c.SetCookie(
		"sessionId",
		"",
		-1, // 立即过期
		"/",
		"",
		false,
		true,
	)

	logger.Infof("User: %s ip: %s logout", username, c.ClientIP())
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Logout successful"})
}

// 获取会话信息
func HandleGetSessionInfo(c *gin.Context) {
	cfg := config.Get()
	if !cfg.Auth.Enabled {
		c.JSON(http.StatusOK, gin.H{
			"authenticated": false,
			"authEnabled":   false,
		})
		return
	}

	sessionId, success := getSessionId(c)
	if !success || sessionId == "" {
		c.JSON(http.StatusOK, gin.H{
			"authenticated": false,
			"authEnabled":   true,
		})
		return
	}

	session, err := service.ValidateSession(sessionId)
	if err != nil || session == nil {
		c.JSON(http.StatusOK, gin.H{
			"authenticated": false,
			"authEnabled":   true,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"authenticated": true,
		"authEnabled":   true,
		"username":      session.User.Username,
		"role":          session.User.Role,
		"userId":        session.User.ID,
	})
}

func getSessionId(c *gin.Context) (string, bool) {
	sessionId, exists := c.Get("sessionId")
	if !exists || sessionId == "" {
		// 如果从 context 没有获取到 sessionId，则尝试从 cookie 获取
		var err error
		sessionId, err = c.Cookie("sessionId")
		if err != nil || sessionId == "" {
			return "", false
		}
	}
	return sessionId.(string), true
}

func getUser(c *gin.Context) (*models.User, bool) {
	user, exists := c.Get("user")
	if !exists {
		return nil, false
	}

	if u, ok := user.(*models.User); ok {
		return u, true
	}

	return nil, false
}
