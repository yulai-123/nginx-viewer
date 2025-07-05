package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
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

	// 注册路由
	http.HandleFunc("/api/logs", apiHandler.HandleLogs)

	// 启动服务
	addr := fmt.Sprintf(":%d", port)
	log.Printf("HTTP服务启动在 %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
