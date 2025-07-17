package models

import (
	"time"
)

type Client struct {
	ClientID     string `gorm:"primaryKey;size:255;not null"`
	Alias        string `gorm:"size:100"`
	IPAddress    string `gorm:"size:45"`
	OnlineStatus bool   `gorm:"default:false"`
	LastSeen     time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time

	CommandLogs []CommandLog    `gorm:"foreignKey:ClientID;references:ClientID"`
	ScreenshotLogs []ScreenshotLog `gorm:"foreignKey:ClientID;references:ClientID"`
}
