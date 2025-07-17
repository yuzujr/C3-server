package repository

import (
	"github.com/yuzujr/C3/internal/database"
	"github.com/yuzujr/C3/internal/models"
)

func CreateScreenshotLog(log *models.ScreenshotLog) error {
	return database.DB.Create(log).Error
}
