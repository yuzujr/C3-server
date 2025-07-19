package main

import (
	"log"
	"os"
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
	// 加载环境变量
	// 这里假设 .env 文件在项目根目录
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
	// 获取cwd
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("os.Executable: %v", err)
	}
	// web目录必须在cwd下
	// UPLOAD_DIR可以是绝对路径或相对路径
	r.Static("/static", filepath.Join(cwd, "web"))
	r.Static("/uploads", config.Get().Upload.Directory)

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
	// 前端api
	webGroup := r.Group(cfg.Server.BasePath + "web")
	routes.RegisterWebRoutes(webGroup)
	// WebSocket连接api
	wsGroup := r.Group(cfg.Server.BasePath + "ws")
	routes.RegisterWsRoutes(wsGroup)

	// 启动
	addr := cfg.Server.Host + ":" + strconv.Itoa(cfg.Server.Port)
	logger.Infof("Access the web interface at http://%s", addr)
	if err := r.Run(addr); err != nil {
		logger.Errorf("Failed to start server: %v", err)
	}
}
