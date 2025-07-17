package models

import (
	"time"
)

type ScreenshotLog struct {
	ID       uint   `gorm:"primaryKey"`
	ClientID string `gorm:"not null"`
	Filename string `gorm:"size:255;not null"`
	FilePath string `gorm:"size:500;not null"`
	FileSize int

	CreatedAt time.Time
	UpdatedAt time.Time
}
