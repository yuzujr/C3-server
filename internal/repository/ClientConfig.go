package repository

import (
	"github.com/yuzujr/C3/internal/database"
	"github.com/yuzujr/C3/internal/models"
)

func UpsertClientConfig(config *models.ClientConfig) error {
	return database.DB.Save(config).Error
}

func GetClientConfigByID(clientID string) (*models.ClientConfig, error) {
	var config models.ClientConfig
	err := database.DB.Where("api_client_id = ?", clientID).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}
