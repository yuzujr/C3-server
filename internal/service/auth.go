package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/yuzujr/C3/internal/config"
	"github.com/yuzujr/C3/internal/models"
	"github.com/yuzujr/C3/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type LoginResult struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	SessionID string `json:"session_id,omitempty"` // 登录成功时返回的会话ID
}

func Login(username, password string) (*LoginResult, error) {
	user, err := repository.GetUserByUsername(username)
	if err != nil {
		return &LoginResult{Success: false, Message: "User not found"}, err
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return &LoginResult{Success: false, Message: "Invalid password"}, err
	}

	// 创建会话
	sessionBytes := make([]byte, 32)
	rand.Read(sessionBytes)
	sessionID := hex.EncodeToString(sessionBytes)
	expiresAt := time.Now().Add(time.Duration(config.Get().Auth.SessionExpireHours) * time.Hour)
	session := &models.UserSession{
		SessionID: sessionID,
		UserID:    user.ID,
		ExpiresAt: expiresAt,
	}
	if err := repository.UpsertUserSession(session); err != nil {
		return &LoginResult{Success: false, Message: "Failed to create session"}, err
	}

	// 登录成功，返回结果
	return &LoginResult{Success: true, Message: "Login successful", SessionID: sessionID}, nil
}

func Logout(sessionId string) error {
	// 删除会话
	if err := repository.DeleteUserSession(sessionId); err != nil {
		return err
	}

	return nil
}

func ValidateSession(sessionID string) (*models.UserSession, error) {
	session, err := repository.GetUserSessionByID(sessionID)
	if err != nil {
		return nil, err
	}

	// 检查会话是否过期
	if session.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("session expired")
	}

	return session, nil
}
