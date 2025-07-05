package main

import (
	"encoding/json"
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
func (h *APIHandler) HandleLogs(w http.ResponseWriter, r *http.Request) {
	// 解析过滤参数
	filter := h.parseFilter(r)

	// 获取日志数据
	logs, err := h.cache.GetLogs(filter)
	if err != nil {
		http.Error(w, "获取日志失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回JSON响应
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total": len(logs),
		"logs":  logs,
	})
}

// parseFilter 从请求参数中解析过滤条件
func (h *APIHandler) parseFilter(r *http.Request) LogFilter {
	query := r.URL.Query()

	filter := LogFilter{
		IP:     query.Get("ip"),
		Path:   query.Get("path"),
		Offset: 0,
		Limit:  200, // 默认每页200条
	}

	// 解析状态码
	if statusStr := query.Get("status"); statusStr != "" {
		if status, err := strconv.Atoi(statusStr); err == nil {
			filter.Status = status
		}
	}

	// 解析分页参数
	if limitStr := query.Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filter.Limit = limit
		}
	}

	if offsetStr := query.Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filter.Offset = offset
		}
	}

	// 解析时间范围
	if fromStr := query.Get("from"); fromStr != "" {
		if from, err := time.Parse(time.RFC3339, fromStr); err == nil {
			filter.From = from
		}
	}

	if toStr := query.Get("to"); toStr != "" {
		if to, err := time.Parse(time.RFC3339, toStr); err == nil {
			filter.To = to
		}
	}

	return filter
}
