package repository

import (
	"github.com/yuzujr/C3/internal/database"
	"github.com/yuzujr/C3/internal/models"
)

func UpsertUserSession(session *models.UserSession) error {
	return database.DB.Save(session).Error
}

func GetUserSessionByID(sessionID string) (*models.UserSession, error) {
	var session models.UserSession
	err := database.DB.Preload("User").Where("session_id = ?", sessionID).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func DeleteUserSession(sessionID string) error {
	return database.DB.Where("session_id = ?", sessionID).Delete(&models.UserSession{}).Error
}
