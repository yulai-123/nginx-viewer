package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// APIHandler 处理API请求
type APIHandler struct {
	cache  *CacheManager
	config *Config
}

// NewAPIHandler 创建一个新的API处理器
func NewAPIHandler(cache *CacheManager, config *Config) *APIHandler {
	return &APIHandler{
		cache:  cache,
		config: config,
	}
}

// HandleLogs 处理 /api/logs 请求
func (h *APIHandler) HandleLogs(c *gin.Context) {
	// 解析过滤参数
	filter := h.parseFilter(c)

	// 获取日志数据
	logs, err := h.cache.GetLogs(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取日志失败: " + err.Error()})
		return
	}

	// 返回JSON响应
	c.JSON(http.StatusOK, gin.H{
		"total": len(logs),
		"logs":  logs,
	})
}

// parseFilter 从请求参数中解析过滤条件
func (h *APIHandler) parseFilter(c *gin.Context) LogFilter {
	filter := LogFilter{
		IP:     c.Query("ip"),
		Path:   c.Query("path"),
		Offset: 0,
		Limit:  200, // 默认每页200条
	}

	// 解析状态码
	if statusStr := c.Query("status"); statusStr != "" {
		if status, err := strconv.Atoi(statusStr); err == nil {
			filter.Status = status
		}
	}

	// 解析分页参数
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filter.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filter.Offset = offset
		}
	}

	// 解析时间范围
	if fromStr := c.Query("from"); fromStr != "" {
		if from, err := time.Parse(time.RFC3339, fromStr); err == nil {
			filter.From = from
		}
	}

	if toStr := c.Query("to"); toStr != "" {
		if to, err := time.Parse(time.RFC3339, toStr); err == nil {
			filter.To = to
		}
	}

	return filter
}
