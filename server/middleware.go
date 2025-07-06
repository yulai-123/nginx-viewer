package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

// 创建日志中间件
func loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		startTime := time.Now()

		// 调用下一个处理器
		c.Next()

		// 计算请求处理时间
		duration := time.Since(startTime)

		// 记录详细的请求信息
		logrus.Infof(
			"[%s] %s | Path: %s | Status: %d | Time: %v | IP: %s | Params: %v | UserAgent: %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.Request.URL.String(),
			c.Writer.Status(),
			duration,
			c.ClientIP(),
			c.Request.Form,
			c.Request.UserAgent(),
		)
	}
}

// corsMiddleware 处理跨域资源共享
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // 允许所有来源访问，生产环境可改为特定域名
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
