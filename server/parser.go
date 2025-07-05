package main

import (
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// LogRow 表示解析后的一行日志
type LogRow struct {
	Time       time.Time
	ClientIP   string
	ClientPort string
	XFF        string
	Host       string
	ServerPort string
	Method     string
	Path       string
	HTTPVer    string
	Status     int
	BodyBytes  int64
	ReqBytes   int64
	ReqTime    float64
	UpConnTime sql.NullFloat64 // 允许 NULL
	UpRespTime sql.NullFloat64
	UpStatus   sql.NullInt64
	UpAddr     sql.NullString
	Referer    string
	UA         string
	TLSProto   string
	TLSCipher  string
	ReqID      string
}

// 预编译的正则表达式，用于解析日志行
var logRegex *regexp.Regexp

func init() {
	// 初始化正则表达式
	// 初始化正则表达式 - 使用更灵活的模式匹配请求部分
	pattern := `^(\S+):(\d+) - (\S+):(\d+) \[([^\]]+)\] "(.*?)" (\d+) (\d+) (\d+) rt=([0-9.]+) uct=([0-9.-]+) urt=([0-9.-]+) ust=([0-9-]+) ua=([^ ]+) "([^"]*)" "([^"]*)" (\S+) (\S+) rid=(\S+)`
	logRegex = regexp.MustCompile(pattern)
}

// ParseLogLine 解析单行日志
func ParseLogLine(line string) (*LogRow, error) {
	// 使用正则表达式匹配日志行
	// 使用正则表达式匹配日志行
	matches := logRegex.FindStringSubmatch(line)
	if matches == nil || len(matches) < 19 {
		// 如果没有匹配到或匹配组数量不足，尝试使用字符串分割的方法
		return parseLogLineAlternative(line)
	}

	// 将匹配结果映射到LogRow结构体
	result := &LogRow{}

	// 解析时间
	timeStr := matches[5]
	parsedTime, err := time.Parse("02/Jan/2006:15:04:05 -0700", timeStr)
	if err != nil {
		return nil, fmt.Errorf("解析时间失败: %v", err)
	}
	result.Time = parsedTime

	// 设置基本字段
	// 设置基本字段
	result.ClientIP = matches[1]
	result.ClientPort = matches[2]
	result.XFF = matches[3]
	result.Host = matches[3]
	result.ServerPort = matches[4]

	// 处理请求部分
	requestPart := matches[6]
	// 尝试分割请求部分为METHOD PATH HTTP_VERSION
	requestParts := strings.SplitN(requestPart, " ", 3)
	if len(requestParts) == 3 {
		// 正常的HTTP请求
		result.Method = requestParts[0]
		result.Path = requestParts[1]
		result.HTTPVer = requestParts[2]
	} else {
		// 异常的请求，例如包含二进制数据
		result.Method = requestPart
		result.Path = ""
		result.HTTPVer = ""
	}

	// 解析数值字段
	// 解析数值字段
	status, _ := strconv.Atoi(matches[7])
	result.Status = status

	bodyBytes, _ := strconv.ParseInt(matches[8], 10, 64)
	result.BodyBytes = bodyBytes

	reqBytes, _ := strconv.ParseInt(matches[9], 10, 64)
	result.ReqBytes = reqBytes

	reqTime, _ := strconv.ParseFloat(matches[10], 64)
	result.ReqTime = reqTime

	// 处理可能为空的字段
	if matches[11] != "-" {
		upConnTime, err := strconv.ParseFloat(matches[11], 64)
		if err == nil {
			result.UpConnTime = sql.NullFloat64{Float64: upConnTime, Valid: true}
		}
	}

	if matches[12] != "-" {
		upRespTime, err := strconv.ParseFloat(matches[12], 64)
		if err == nil {
			result.UpRespTime = sql.NullFloat64{Float64: upRespTime, Valid: true}
		}
	}

	if matches[13] != "-" {
		upStatus, err := strconv.ParseInt(matches[13], 10, 64)
		if err == nil {
			result.UpStatus = sql.NullInt64{Int64: upStatus, Valid: true}
		}
	}

	if matches[14] != "-" {
		result.UpAddr = sql.NullString{String: matches[14], Valid: true}
	}

	// 设置其他字段
	result.Referer = strings.Trim(matches[15], "\"")
	result.UA = strings.Trim(matches[16], "\"")
	result.TLSProto = matches[17]
	result.TLSCipher = matches[18]
	result.ReqID = matches[19]

	return result, nil
}

// parseLogLineAlternative 使用字符串分割的方式解析日志行（作为备选方案）
func parseLogLineAlternative(line string) (*LogRow, error) {
	// 实现字符串分割的解析逻辑
	// 此处为简化实现，实际项目中应当完整实现
	return nil, fmt.Errorf("使用备选解析方法失败")
}
