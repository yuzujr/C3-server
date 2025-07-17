package repository

import (
	"github.com/yuzujr/C3/internal/database"
	"github.com/yuzujr/C3/internal/models"
)

// 根据 client_id 查找 Client
func FindClientByID(clientID string) (*models.Client, error) {
	var client models.Client
	err := database.DB.First(&client, clientID).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}

// 查别名
func FindAliasByID(clientID string) (string, error) {
	client, err := FindClientByID(clientID)
	if err != nil {
		return "", err
	}
	return client.Alias, nil
}

// 创建新 Client
func UpsertClient(client *models.Client) error {
	return database.DB.Create(client).Error
}
