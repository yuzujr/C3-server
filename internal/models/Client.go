package models

import (
	"time"
)

type Client struct {
	ClientID  string `gorm:"primaryKey;size:255;not null"`
	Alias     string `gorm:"size:100"`
	IPAddress string `gorm:"size:45"`
	Online    bool   `gorm:"default:false"`
	LastSeen  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time

	ClientConfig   ClientConfig    `gorm:"foreignKey:ClientID;references:ClientID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CommandLogs    []CommandLog    `gorm:"foreignKey:ClientID;references:ClientID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ScreenshotLogs []ScreenshotLog `gorm:"foreignKey:ClientID;references:ClientID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
