package service

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yuzujr/C3/internal/config"
	"github.com/yuzujr/C3/internal/logger"
	"github.com/yuzujr/C3/internal/models"
	"github.com/yuzujr/C3/internal/repository"
	"gorm.io/gorm"
)

// CreateClient 只做“新增”，客户端已存在将报错
func CreateClient(info *models.Client) error {
	// 检查是否已存在
	if _, err := repository.FindClientByID(info.ClientID); err == nil {
		return fmt.Errorf("client %s already exists", info.ClientID)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 填充默认字段
	now := time.Now()
	newClient := &models.Client{
		ClientID:     info.ClientID,
		Alias:        info.ClientID, // 默认别名为 client_id
		IPAddress:    info.IPAddress,
		OnlineStatus: info.OnlineStatus,
		LastSeen:     now,
	}
	return repository.UpsertClient(newClient)
}

// UpdateClient 只做“更新”，客户端不存在将报错
func UpdateClient(info *models.Client) error {
	client, err := repository.FindClientByID(info.ClientID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("client %s not found", info.ClientID)
		}
		return err
	}

	// 有值就更新
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

// UpsertClient 根据是否存在调用 CreateClient 或 UpdateClient
func UpsertClient(info *models.Client) error {
	if _, err := repository.FindClientByID(info.ClientID); errors.Is(err, gorm.ErrRecordNotFound) {
		return CreateClient(info)
	} else if err != nil {
		return err
	}
	return UpdateClient(info)
}

// 删除客户端
func DeleteClient(clientID string) error {
	client, err := repository.FindClientByID(clientID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil // 客户端不存在
		}
		logger.Errorf("Failed to find client: %v", err)
		return err
	}

	return repository.DeleteClient(client)
}

// 获取客户端信息
func GetClient(clientID string) (*models.Client, error) {
	client, err := repository.FindClientByID(clientID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 客户端不存在
		}
		logger.Errorf("Failed to find client: %v", err)
		return nil, err
	}
	return client, nil
}

// 获取所有客户端列表
func GetAllClients() ([]*models.Client, error) {
	clients, err := repository.GetAllClients()
	if err != nil {
		logger.Errorf("Failed to get clients: %v", err)
		return nil, err
	}
	return clients, nil
}

// 获取在线客户端ID列表
func GetOnlineClientsID() []string {
	onlineIDs := make([]string, 0)
	for _, client := range repository.GetOnlineClients() {
		onlineIDs = append(onlineIDs, client.ClientID)
	}
	return onlineIDs
}

// 查询别名
func GetAlias(clientID string) (string, error) {
	client, err := repository.FindClientByID(clientID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil // 客户端不存在
		}
		logger.Errorf("Failed to find client: %v", err)
		return "", err
	}
	return client.Alias, nil
}

// 更新别名
func UpdateClientAlias(clientID, alias string) error {
	client, err := repository.FindClientByID(clientID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil // 客户端不存在
		}
		logger.Errorf("Failed to find client: %v", err)
		return err
	}

	client.Alias = alias
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

// 保存客户端截图文件并记录日志
func SaveScreenshot(c *gin.Context, file *multipart.FileHeader, clientID string) error {
	// 生成带随机数的文件名，防止同一秒文件被覆盖
	ext := filepath.Ext(file.Filename)
	name := file.Filename[:len(file.Filename)-len(ext)]
	randomSuffix := fmt.Sprintf("_%d", time.Now().UnixNano()%1e6)
	newFilename := fmt.Sprintf("%s%s%s", name, randomSuffix, ext)

	dst := filepath.Join(config.Get().Upload.Directory, clientID, time.Now().Format("2006-01-02_15"), newFilename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		logger.Errorf("Failed to save file: %v", err)
		return err
	}

	// 记录截图日志
	// 如果记录失败，则删除已保存的文件
	client, err := repository.FindClientByID(clientID)
	if err != nil {
		logger.Errorf("Failed to find client: %v", err)
		_ = os.Remove(dst)
		return err
	}

	screenshot := models.ScreenshotLog{
		ClientID: client.ClientID,
		Filename: newFilename,
		FilePath: dst,
		FileSize: int(file.Size),
	}

	if err := repository.CreateScreenshotLog(&screenshot); err != nil {
		logger.Errorf("Failed to log screenshot: %v", err)
		_ = os.Remove(dst)
		return err
	}

	return nil
}
