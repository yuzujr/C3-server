package repository

import (
	"github.com/yuzujr/C3/internal/database"
	"github.com/yuzujr/C3/internal/models"
)

// 创建或更新客户端信息
func UpsertClient(client *models.Client) error {
	return database.DB.Save(client).Error
}

// 删除客户端
func DeleteClient(client *models.Client) error {
	return database.DB.Delete(client).Error
}

// 根据 client_id 查找 Client
func FindClientByID(clientID string) (*models.Client, error) {
	var client models.Client
	err := database.DB.Where("client_id = ?", clientID).First(&client).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}

// 获取所有 Client
func GetAllClients() ([]*models.Client, error) {
	var clients []*models.Client
	err := database.DB.Find(&clients).Error
	if err != nil {
		return nil, err
	}
	return clients, nil
}

// 获取在线客户端列表
func GetOnlineClients() []*models.Client {
	var clients []*models.Client
	err := database.DB.Where("online_status = ?", true).Find(&clients).Error
	if err != nil {
		return nil
	}
	return clients
}
