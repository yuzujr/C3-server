package service

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/yuzujr/C3/internal/config"
	"github.com/yuzujr/C3/internal/eventbus"
	"github.com/yuzujr/C3/internal/logger"
	"github.com/yuzujr/C3/internal/models"
	"github.com/yuzujr/C3/internal/repository"
	"gorm.io/gorm"
)

// createClient 只做新增
func createClient(info *models.Client) error {
	// 填充默认字段
	now := time.Now()
	newClient := &models.Client{
		ClientID:  info.ClientID,
		Alias:     info.ClientID, // 默认别名为 client_id
		IPAddress: info.IPAddress,
		Online:    info.Online,
		LastSeen:  now,
	}
	return repository.UpsertClient(newClient)
}

// updateClient 只做更新
func updateClient(old, info *models.Client) error {
	if info.Alias != "" {
		old.Alias = info.Alias
	}
	if info.IPAddress != "" {
		old.IPAddress = info.IPAddress
	}
	old.Online = info.Online
	old.LastSeen = time.Now()

	return repository.UpsertClient(old)
}

func UpdateClient(info *models.Client) error {
	old, err := repository.FindClientByID(info.ClientID)
	if err != nil {
		return err
	}
	return updateClient(old, info)
}

// UpsertClient 根据是否存在调用 CreateClient 或 UpdateClient
func UpsertClient(info *models.Client) error {
	old, err := repository.FindClientByID(info.ClientID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return createClient(info)
	} else if err != nil {
		return err
	}
	return updateClient(old, info)
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

	if client.Online {
		sendOfflineCommand(clientID, "Client deleted")
		time.Sleep(100 * time.Millisecond) // 等待客户端处理下线命令
	}

	cleanupClientFiles(clientID)

	logger.Infof("Client %s deleted", clientID)
	return repository.DeleteClient(clientID)
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

// 发送下线命令
func sendOfflineCommand(clientID, reason string) {
	msg := eventbus.Command{
		Type: "offline",
		Data: map[string]any{"reason": reason},
	}

	eventbus.Global.SendCommand(clientID, msg)
	logger.Infof("Sent offline command to client %s", clientID)
}

// 清理客户端相关文件
func cleanupClientFiles(clientID string) {
	cfg := config.Get()
	// 删除客户端目录
	clientDir := filepath.Join(cfg.Upload.Directory, clientID)
	if err := os.RemoveAll(clientDir); err != nil {
		logger.Errorf("Failed to remove client directory %s: %v", clientDir, err)
	}
}
