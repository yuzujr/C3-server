package repository

import (
	"github.com/yuzujr/C3/internal/database"
	"github.com/yuzujr/C3/internal/models"
)

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
