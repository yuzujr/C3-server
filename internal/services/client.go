package services

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/yuzujr/C3/internal/config"
	"github.com/yuzujr/C3/internal/database"
	"github.com/yuzujr/C3/internal/logger"
	"github.com/yuzujr/C3/internal/models"
	"gorm.io/gorm"
)

// 获取客户端信息
func GetClient(clientID string) (*models.Client, error) {
	var client models.Client
	err := database.DB.Where("client_id = ?", clientID).First(&client).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}

// 查询别名
func GetAlias(clientID string) (string, error) {
	var client models.Client
	err := database.DB.Where("client_id = ?", clientID).First(&client).Error
	if err != nil {
		return "", err
	}
	return client.Alias, nil
}

// 设置或更新客户端信息
func SetClient(clientID string, info *models.Client) error {
	var client models.Client
	err := database.DB.Where("client_id = ?", clientID).First(&client).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 不存在，新建
		client = models.Client{
			ClientID:     clientID,
			Alias:        clientID,
			IPAddress:    info.IPAddress,
			OnlineStatus: info.OnlineStatus,
			LastSeen:     time.Now(),
		}
		return database.DB.Create(&client).Error
	} else if err != nil {
		logger.Errorf("Failed to find client: %v", err)
		return err
	}

	// 更新
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
	return database.DB.Save(&client).Error
}

// 保存客户端配置
func SaveClientConfig(clientConfig models.ClientConfig) error {
	// 验证配置
	if err := validateApiConfig(clientConfig.Api); err != nil {
		logger.Errorf("Invalid client config: %v", err)
		return err
	}

	// 保存配置文件到UPLOAD_DIR/client_id/config.json
	dir := filepath.Join(config.Get().Upload.Directory, clientConfig.ClientID)
	if dir == "" {
		return errors.New("upload directory is not set")
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		logger.Errorf("Failed to create config directory: %v", err)
		return err
	}
	configPath := filepath.Join(dir, "config.json")
	file, err := os.Create(configPath)
	if err != nil {
		logger.Errorf("Failed to create config file: %v", err)
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(clientConfig); err != nil {
		logger.Errorf("Failed to write config to file: %v", err)
		return err
	}
	return nil
}

// 记录截图信息到数据库
func LogScreenshot(clientID string, filename, path string, size int64) error {
	// 查 client 的主键 ID
	var client models.Client
	if err := database.DB.Where("client_id = ?", clientID).First(&client).Error; err != nil {
		return err
	}

	screenshot := models.Screenshot{
		ClientID: client.ClientID,
		Filename: filename,
		FilePath: path,
		FileSize: int(size),
	}

	return database.DB.Create(&screenshot).Error
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
