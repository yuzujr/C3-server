package handler

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yuzujr/C3/internal/service"
)

// contains 判断字符串切片中是否包含指定项
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// HandleGetClients 返回客户端列表，包含在线状态
func HandleGetClients(c *gin.Context) {
	clients, err := service.GetAllClients()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	onlineIDs := service.GetOnlineClientsID()
	list := make([]gin.H, 0, len(clients))
	for _, cli := range clients {
		list = append(list, gin.H{
			"client_id": cli.ClientID,
			"alias":     cli.Alias,
			"online":    contains(onlineIDs, cli.ClientID),
		})
	}

	c.JSON(http.StatusOK, list)
}

// HandleGetClientConfig 获取指定客户端的配置
func HandleGetClientConfig(c *gin.Context) {
	clientID := c.Param("client_id")
	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing client_id"})
		return
	}

	config, err := service.GetConfig(clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if config == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Config not found"})
		return
	}

	c.JSON(http.StatusOK, config)
}

// HandleUpdateClientAlias 更新客户端别名
func HandleUpdateClientAlias(c *gin.Context) {
	clientID := c.Param("client_id")
	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Missing client_id"})
		return
	}

	var req struct {
		NewAlias string `json:"newAlias"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Missing newAlias"})
		return
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(req.NewAlias) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid alias format"})
		return
	}

	if err := service.UpdateClientAlias(clientID, req.NewAlias); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "newAlias": req.NewAlias})
}

// HandleDeleteClient 删除客户端
func HandleDeleteClient(c *gin.Context) {
	clientID := c.Param("client_id")
	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Missing client_id"})
		return
	}

	if err := service.DeleteClient(clientID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// HandleGetClientScreenshots 获取客户端截图列表
func HandleGetClientScreenshots(c *gin.Context) {
	clientID := c.Param("client_id")
	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing client_id"})
		return
	}

	since := int64(0)
	if s := c.Query("since"); s != "" {
		if v, err := strconv.ParseInt(s, 10, 64); err == nil {
			since = v
		}
	}

	shots, err := service.GetScreenshotsSince(clientID, since)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shots)
}

// HandleDeleteScreenshotsByTime 删除指定时间范围内截图
func HandleDeleteScreenshotsByTime(c *gin.Context) {
	clientID := c.Param("client_id")
	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Missing client_id"})
		return
	}

	var req struct {
		Hours int `json:"hours"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Hours <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid hours"})
		return
	}

	deletedCount, err := service.DeleteScreenshotsAfterHours(clientID, req.Hours)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"deletedCount": deletedCount,
	})
}

// HandleDeleteAllScreenshots 删除所有截图
func HandleDeleteAllScreenshots(c *gin.Context) {
	clientID := c.Param("client_id")
	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Missing client_id"})
		return
	}

	deletedCount, err := service.DeleteAllScreenshots(clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"deletedCount": deletedCount,
	})
}
