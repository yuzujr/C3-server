package models

import (
	"time"
)

type UserSession struct {
	SessionID string    `gorm:"primaryKey;size:64;not null"`
	UserID    uint      `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	User      User `gorm:"foreignKey:UserID;references:ID"`
}
