package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yuzujr/C3/internal/logger"
	"github.com/yuzujr/C3/internal/models"
	"github.com/yuzujr/C3/internal/service"
)

// 上传客户端配置
// 若为新客户端，则创建新记录；若已存在，则更新记录
func HandleClientConfig(c *gin.Context) {
	var config models.ClientConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		logger.Errorf("Invalid client config: %v", err)
		c.String(http.StatusBadRequest, "Invalid client config")
		return
	}

	clientID := config.Api.ClientID
	if clientID == "" {
		logger.Errorf("Missing client_id in config")
		c.String(http.StatusBadRequest, "Missing client_id")
		return
	}

	// 创建或更新客户端信息
	service.SetClient(&models.Client{
		ClientID:     clientID,
		IPAddress:    c.ClientIP(),
		OnlineStatus: true,
	})

	// 配置存储到数据库
	config.LastUpload = time.Now().Format("2006-01-02 15:04:05")
	err := service.SetConfig(&config)
	if err != nil {
		logger.Errorf("Failed to save config: %v", err)
		c.String(http.StatusInternalServerError, "Failed to save config")
		return
	}

	logger.Infof("Client config received from %s", clientID)
	c.String(http.StatusOK, "Client config uploaded")
}

// 上传截图
// 要求客户端已存在
func HandleClientScreenshot(c *gin.Context) {
	clientID := c.Query("client_id")
	if clientID == "" {
		c.String(http.StatusBadRequest, "Missing client_id")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		logger.Errorf("File upload error: %v", err)
		c.String(http.StatusBadRequest, "File upload error")
		return
	}

	dst, err := service.SaveClientScreenshot(c, file, clientID)
	if err != nil {
		logger.Errorf("Failed to save screenshot: %v", err)
		c.String(http.StatusInternalServerError, "Failed to save screenshot")
		return
	}

	// 日志存储数据库
	err = service.LogScreenshot(clientID, file.Filename, dst, file.Size)
	if err != nil {
		logger.Errorf("Database error: %v", err)
		c.String(http.StatusInternalServerError, "Database error")
		return
	}

	logger.Infof("Screenshot received from %s", clientID)
	c.String(http.StatusCreated, "Screenshot uploaded")
}
