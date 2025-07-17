package models

import (
	"time"
)

type CommandLog struct {
	ID       uint   `gorm:"primaryKey"`
	ClientID string `gorm:"not null"`
	Command  string `gorm:"type:text;not null"`
	Result   string `gorm:"type:text"`
	ExitCode int

	CreatedAt time.Time
	UpdatedAt time.Time
}
