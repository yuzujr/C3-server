package config

import (
	"log"
	"os"
	"strconv"
	"sync"
)

type Config struct {
	Server ServerConfig
	DB     DatabaseConfig
	Auth   AuthConfig
	Upload UploadConfig
	Log    LogConfig
}

var (
	cfg  *Config
	once sync.Once
)

func Get() *Config {
	once.Do(func() {
		load()
	})
	return cfg
}

func load() {
	// PORT
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Invalid PORT value: %v", err)
	}
	// DB
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalf("Invalid DB_PORT value: %v", err)
	}
	// SessionExpireHours
	sessionExpireHours, err := strconv.Atoi(os.Getenv("SESSION_EXPIRE_HOURS"))
	if err != nil {
		log.Fatalf("Invalid SESSION_EXPIRE_HOURS value: %v", err)
	}
	if sessionExpireHours <= 0 {
		log.Fatalf("SESSION_EXPIRE_HOURS must be greater than 0")
	}

	cfg = &Config{
		Server: ServerConfig{
			BasePath: os.Getenv("BASE_PATH"),
			Host:     os.Getenv("HOST"),
			Port:     port,
			Env:      os.Getenv("ENV"),
		},
		DB: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     dbPort,
			Name:     os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
		},
		Auth: AuthConfig{
			Enabled:            os.Getenv("AUTH_ENABLED") == "true",
			Username:           os.Getenv("AUTH_USERNAME"),
			Password:           os.Getenv("AUTH_PASSWORD"),
			SessionExpireHours: int64(sessionExpireHours),
		},
		Upload: UploadConfig{
			Directory: os.Getenv("UPLOAD_DIR"),
		},
		Log: LogConfig{
			Directory: os.Getenv("LOG_DIR"),
			Level:     os.Getenv("LOG_LEVEL"),
		},
	}
}
