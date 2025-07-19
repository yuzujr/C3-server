package service

import (
	"errors"

	"github.com/yuzujr/C3/internal/logger"
	"github.com/yuzujr/C3/internal/models"
	"github.com/yuzujr/C3/internal/repository"
)

// 保存或更新客户端配置
func UpsertConfig(config *models.ClientConfig) error {
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
