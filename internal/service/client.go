package service

import (
	"errors"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yuzujr/C3/internal/config"
	"github.com/yuzujr/C3/internal/logger"
	"github.com/yuzujr/C3/internal/models"
	"github.com/yuzujr/C3/internal/repository"
	"gorm.io/gorm"
)

// 获取客户端信息
func GetClient(clientID string) (*models.Client, error) {
	return repository.FindClientByID(clientID)
}

// 查询别名
func GetAlias(clientID string) (string, error) {
	return repository.FindAliasByID(clientID)
}

// 创建或更新客户端信息
func SetClient(info *models.Client) error {
	client, err := repository.FindClientByID(info.ClientID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果不存在，则创建新客户端
		client = &models.Client{
			ClientID:     info.ClientID,
			Alias:        info.ClientID,
			IPAddress:    info.IPAddress,
			OnlineStatus: info.OnlineStatus,
			LastSeen:     time.Now(),
		}
		return repository.UpsertClient(client)
	} else if err != nil {
		logger.Errorf("Failed to find client: %v", err)
		return err
	}

	if info.ClientID != "" {
		client.ClientID = info.ClientID
	}
	if info.Alias != "" {
		client.Alias = info.Alias
	}
	if info.IPAddress != "" {
		client.IPAddress = info.IPAddress
	}
	client.OnlineStatus = info.OnlineStatus
	client.LastSeen = time.Now()

	return repository.UpsertClient(client)
}

// 保存或更新客户端配置
func SetConfig(config *models.ClientConfig) error {
	if err := validateApiConfig(config.Api); err != nil {
		logger.Errorf("Invalid API config: %v", err)
		return err
	}
	return repository.UpsertClientConfig(config)
}

func GetConfig(clientID string) (*models.ClientConfig, error) {
	return repository.GetClientConfigByID(clientID)
}

func validateApiConfig(config models.ApiConfig) error {
	if config.Hostname == "" {
		return errors.New("hostname is required")
	}
	if config.Port <= 0 || config.Port > 65535 {
		return errors.New("port must be between 1 and 65535")
	}
	if config.IntervalSeconds <= 0 {
		return errors.New("interval_seconds must be greater than 0")
	}
	if config.RetryDelayMilliseconds <= 0 {
		return errors.New("retry_delay_ms must be greater than 0")
	}
	return nil
}

func LogScreenshot(clientID string, filename, path string, size int64) error {
	client, err := repository.FindClientByID(clientID)
	if err != nil {
		return err
	}

	screenshot := models.ScreenshotLog{
		ClientID: client.ClientID,
		Filename: filename,
		FilePath: path,
		FileSize: int(size),
	}
	return repository.CreateScreenshotLog(&screenshot)
}

// 保存客户端截图文件并记录日志
func SaveClientScreenshot(c *gin.Context, file *multipart.FileHeader, clientID string) (dst string, err error) {
	dst = filepath.Join(config.Get().Upload.Directory, clientID, time.Now().Format("2006-01-02_15"), file.Filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		logger.Errorf("Failed to save file: %v", err)
		c.String(http.StatusInternalServerError, "Failed to save file")
		return "", err
	}

	return dst, nil
}
