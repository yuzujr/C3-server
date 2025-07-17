package routes

import (
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yuzujr/C3/internal/config"
	"github.com/yuzujr/C3/internal/logger"
	"github.com/yuzujr/C3/internal/models"
	"github.com/yuzujr/C3/internal/services"
	"github.com/yuzujr/C3/internal/utils"
	"github.com/yuzujr/C3/internal/websocket"
)

// RegisterClientRoutes 注册客户端相关接口
func RegisterClientRoutes(r *gin.RouterGroup) {
	r.POST("/screenshot", handleClientScreenshot)
	r.POST("/client_config", handleClientConfig)
	r.GET("/ws", websocket.ServeWs)
}

// 上传截图
// 要求客户端已存在
func handleClientScreenshot(c *gin.Context) {
	clientID := c.Query("client_id")
	file, err := c.FormFile("file")
	if err != nil {
		logger.Errorf("File upload error: %v", err)
		c.String(http.StatusBadRequest, "File upload error")
		return
	}

	if clientID == "" {
		c.String(http.StatusBadRequest, "Missing client_id")
		return
	}

	// 保存文件
	dst := filepath.Join(config.Get().Upload.Directory, clientID, utils.GetHourlyTime(), file.Filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		logger.Errorf("Failed to save file: %v", err)
		c.String(http.StatusInternalServerError, "Failed to save file")
		return
	}

	// 存储数据库
	err = services.LogScreenshot(clientID, file.Filename, dst, file.Size)
	if err != nil {
		logger.Errorf("Database error: %v", err)
		c.String(http.StatusInternalServerError, "Database error")
		return
	}

	logger.Infof("Screenshot received from %s", clientID)
	c.String(http.StatusCreated, "Screenshot uploaded")
}

// 上传客户端配置
// 若为新客户端，则创建新记录；若已存在，则更新记录
func handleClientConfig(c *gin.Context) {
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
	services.SetClient(clientID, &models.Client{
		IPAddress:    c.ClientIP(),
		OnlineStatus: true,
	})

	// 保存配置文件
	config.ClientID = clientID
	config.LastUpload = time.Now().Format("2006-01-02 15:04:05")
	err := services.SaveClientConfig(config)
	if err != nil {
		logger.Errorf("Failed to handle config update: %v", err)
		c.String(http.StatusInternalServerError, "Failed to handle config update")
		return
	}

	logger.Infof("Client config received from %s", clientID)
	c.String(http.StatusOK, "Client config uploaded")
}
