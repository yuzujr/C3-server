package main

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/yuzujr/C3/internal/config"
	"github.com/yuzujr/C3/internal/database"
	"github.com/yuzujr/C3/internal/logger"
	"github.com/yuzujr/C3/internal/routes"
	"github.com/yuzujr/C3/internal/websocket"
)

func main() {
	// config
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	cfg := config.Get()

	// database
	database.InitDatabase()

	// http server
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// websocket hub
	go websocket.HubInstance.Run()

	// 中间件
	r.Use(gin.Recovery())

	// 注册路由
	clientGroup := r.Group(cfg.Server.BasePath + "client") // 例如 /c3
	routes.RegisterClientRoutes(clientGroup)

	// 启动
	addr := cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port)
	logger.Infof("Listening on %s", addr)
	if err := r.Run(addr); err != nil {
		logger.Errorf("Failed to start server: %v", err)
	}
}
