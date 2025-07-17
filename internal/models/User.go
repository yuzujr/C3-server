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
