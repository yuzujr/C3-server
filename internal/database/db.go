package database

import (
	"fmt"

	"github.com/yuzujr/C3/internal/config"
	"github.com/yuzujr/C3/internal/logger"
	"github.com/yuzujr/C3/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
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
		Logger: gormLogger.Default.LogMode(gormLogger.Silent), // 禁用 GORM 日志
	})

	if err != nil {
		logger.Fatalf("failed to connect to database: %v", err)
	}

	DB = db

	DB.AutoMigrate(
		&models.User{},
		&models.Client{},
		&models.CommandLog{},
		&models.Screenshot{},
	)

	logger.Infof("Database connected successfully")
}
