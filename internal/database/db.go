package database

import (
	"fmt"

	"github.com/yuzujr/C3/internal/config"
	"github.com/yuzujr/C3/internal/logger"
	"github.com/yuzujr/C3/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	cfg := config.Get()

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Name, cfg.DB.Password,
	)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // 使用 simple protocol（pgx 特性）
	}), &gorm.Config{
		// Logger: gormLogger.Default.LogMode(gormLogger.Silent), // 禁用 GORM 日志
	})

	if err != nil {
		logger.Fatalf("failed to connect to database: %v", err)
	}

	DB = db

	DB.AutoMigrate(
		&models.User{},
		&models.Client{},
		&models.CommandLog{},
		&models.ScreenshotLog{},
		&models.UserSession{},
		&models.ClientConfig{},
	)

	//使用配置文件创建默认用户
	if cfg.Auth.Enabled && cfg.Auth.Username != "" {
		var user models.User
		if err := DB.Where("username = ?", cfg.Auth.Username).First(&user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				bytes, err := bcrypt.GenerateFromPassword([]byte(cfg.Auth.Password), bcrypt.DefaultCost)
				if err != nil {
					logger.Fatalf("failed to hash password: %v", err)
				}
				user = models.User{
					Username:     cfg.Auth.Username,
					PasswordHash: string(bytes),
				}
				if err := DB.Create(&user).Error; err != nil {
					logger.Fatalf("failed to create default user: %v", err)
				}
				logger.Infof("Default user created: %s", user.Username)
			} else {
				logger.Fatalf("failed to check default user: %v", err)
			}
		}
	}

	logger.Infof("Database connected successfully")
}
