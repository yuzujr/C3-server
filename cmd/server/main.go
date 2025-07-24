package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/yuzujr/C3/internal/config"
	"github.com/yuzujr/C3/internal/database"
	"github.com/yuzujr/C3/internal/logger"
	"github.com/yuzujr/C3/internal/middleware"
	"github.com/yuzujr/C3/internal/routes"
	"github.com/yuzujr/C3/internal/service"
	"github.com/yuzujr/C3/internal/websocket"
)

func main() {
	// 加载环境变量
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	cfg := config.Get()

	// 初始化数据库
	database.InitDatabase()

	// 启动 HTTP 和 WebSocket 服务
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	go websocket.HubInstance.Run()

	// 中间件
	r.Use(gin.Recovery())

	// 静态文件
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("os.Executable: %v", err)
	}
	r.Static("/static", filepath.Join(cwd, "web"))
	r.Static("/uploads", config.Get().Upload.Directory)

	// 注册路由
	r.GET("/login", func(c *gin.Context) {
		c.File(filepath.Join("./web", "login.html"))
	})
	r.GET("/", middleware.RequireAuth(), func(c *gin.Context) {
		c.File(filepath.Join("./web", "index.html"))
	})
	clientGroup := r.Group(cfg.Server.BasePath + "client")
	routes.RegisterClientRoutes(clientGroup)
	authGroup := r.Group(cfg.Server.BasePath + "auth")
	routes.RegisterAuthRoutes(authGroup)
	webGroup := r.Group(cfg.Server.BasePath + "web")
	routes.RegisterWebRoutes(webGroup)
	wsGroup := r.Group(cfg.Server.BasePath + "ws")
	routes.RegisterWsRoutes(wsGroup)

	// 设置服务器监听地址
	addr := cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port)
	logger.Infof("Access the web interface at http://%s", addr)

	// 监听中断信号（ctrl+c）以执行清理操作
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 启动 HTTP 服务并等待中断信号
	go func() {
		if err := r.Run(addr); err != nil {
			logger.Errorf("Failed to start server: %v", err)
		}
	}()

	// 等待收到中断信号，执行清理操作
	<-sigChan
	// 执行清理操作，将所有客户端状态设置为下线
	logger.Infof("Shutting down gracefully...")
	clients, _ := service.GetAllClients()
	for _, c := range clients {
		c.Online = false
		service.UpdateClient(c)
	}
}
