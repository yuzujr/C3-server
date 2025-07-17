package models

import (
	"time"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"size:50;unique;not null"`
	PasswordHash string `gorm:"size:255;not null"`
	Role         string `gorm:"size:20;default:admin"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Client struct {
	ClientID     string `gorm:"primaryKey;size:255;not null;unique"` // 显式客户端唯一标识符
	Alias        string `gorm:"size:100"`
	IPAddress    string `gorm:"size:45"`
	OnlineStatus bool   `gorm:"default:false"`
	LastSeen     time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time

	CommandLogs []CommandLog `gorm:"foreignKey:ClientID"`
	Screenshots []Screenshot `gorm:"foreignKey:ClientID"`
}

type CommandLog struct {
	ID         uint   `gorm:"primaryKey"`
	ClientID   string `gorm:"not null"`
	Command    string `gorm:"type:text;not null"`
	Result     string `gorm:"type:text"`
	ExitCode   int
	
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Screenshot struct {
	ID         uint   `gorm:"primaryKey"`
	ClientID   string `gorm:"not null"`
	Filename   string `gorm:"size:255;not null"`
	FilePath   string `gorm:"size:500;not null"`
	FileSize   int

	CreatedAt  time.Time
	UpdatedAt  time.Time
}
