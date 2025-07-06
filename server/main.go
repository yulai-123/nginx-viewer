package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

var (
	configPath string
	port       int
)

func init() {
	flag.StringVar(&configPath, "config", "config.yaml", "配置文件路径")
	flag.IntVar(&port, "port", 8080, "HTTP 服务端口")
}

func main() {
	InitLogger()

	flag.Parse()

	// 加载配置
	cfg, err := LoadConfig(configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 初始化日志缓存管理器
	cache := NewCacheManager(cfg.LogPath)

	// 初始化API处理器
	apiHandler := NewAPIHandler(cache, cfg)

	// 创建Gin路由器
	r := gin.Default()

	// 添加CORS中间件
	r.Use(corsMiddleware())

	// 注册路由
	r.GET("/api/logs", loggingMiddleware(), apiHandler.HandleLogs)

	// 启动服务
	addr := fmt.Sprintf(":%d", port)
	log.Printf("HTTP服务启动在 %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
