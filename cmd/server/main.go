package main

import (
	"log"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/yuzujr/C3/internal/config"
	"github.com/yuzujr/C3/internal/database"
	"github.com/yuzujr/C3/internal/logger"
	"github.com/yuzujr/C3/internal/middleware"
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

	// websocket
	go websocket.HubInstance.Run()

	// 中间件
	r.Use(gin.Recovery())

	// 静态文件
	r.Static("/static", "./web")

	// 注册路由
	// 登录页面路由
	r.GET("/login", func(c *gin.Context) {
		c.File(filepath.Join("./web", "login.html"))
	})
	// 主页重定向到登录检查
	r.GET("/", middleware.RequireAuth(), func(c *gin.Context) {
		c.File(filepath.Join("./web", "index.html"))
	})
	// 客户端api
	clientGroup := r.Group(cfg.Server.BasePath + "client")
	routes.RegisterClientRoutes(clientGroup)
	// 登录api
	authGroup := r.Group(cfg.Server.BasePath + "auth")
	routes.RegisterAuthRoutes(authGroup)

	// 启动
	addr := cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port)
	logger.Infof("Listening on %s", addr)
	if err := r.Run(addr); err != nil {
		logger.Errorf("Failed to start server: %v", err)
	}
}
