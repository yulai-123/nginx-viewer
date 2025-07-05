package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

// CacheKey 表示缓存的键
type CacheKey struct {
	FilePath string
	ModTime  time.Time
}

// CacheManager 管理日志文件的缓存
type CacheManager struct {
	LogPath string                // Nginx 日志目录
	cache   map[CacheKey][]LogRow // 缓存的日志数据
	mu      sync.RWMutex          // 读写锁
}

// NewCacheManager 创建一个新的缓存管理器
func NewCacheManager(logPath string) *CacheManager {
	return &CacheManager{
		LogPath: logPath,
		cache:   make(map[CacheKey][]LogRow),
	}
}

// GetLogs 获取日志数据，根据过滤条件筛选
func (cm *CacheManager) GetLogs(filter LogFilter) ([]LogRow, error) {
	// 获取所有日志文件
	files, err := cm.GetLogFiles()
	if err != nil {
		return nil, fmt.Errorf("获取日志文件列表失败: %v", err)
	}

	// 合并所有日志数据
	var allLogs []LogRow
	for _, file := range files {
		logs, err := cm.GetLogsFromFile(file)
		if err != nil {
			return nil, fmt.Errorf("从文件读取日志失败 %s: %v", file, err)
		}
		allLogs = append(allLogs, logs...)
	}

	// 应用过滤条件
	filteredLogs := cm.ApplyFilter(allLogs, filter)

	// 按照时间排序
	sort.Slice(filteredLogs, func(i, j int) bool {
		return filteredLogs[i].Time.After(filteredLogs[j].Time) // 倒序，最新的在前
	})

	// 应用分页
	if filter.Limit > 0 {
		start := filter.Offset
		end := filter.Offset + filter.Limit
		if start >= len(filteredLogs) {
			return []LogRow{}, nil
		}
		if end > len(filteredLogs) {
			end = len(filteredLogs)
		}
		filteredLogs = filteredLogs[start:end]
	}

	return filteredLogs, nil
}

// GetLogFiles 获取所有需要处理的日志文件列表
func (cm *CacheManager) GetLogFiles() ([]string, error) {
	// 实现获取日志文件列表的逻辑
	// 包括当前的 access_security.log 和历史的 access_security.log.*.gz 文件
	// 打开日志目录
	dir, err := os.Open(cm.LogPath)
	if err != nil {
		return nil, fmt.Errorf("打开日志目录失败: %v", err)
	}
	defer dir.Close()

	// 读取目录中的所有文件
	files, err := dir.Readdir(-1)
	if err != nil {
		return nil, fmt.Errorf("读取日志目录失败: %v", err)
	}

	var logFiles []string

	// 过滤文件
	for _, file := range files {
		if file.IsDir() {
			continue // 跳过目录
		}

		filename := file.Name()

		// 检查是否是当前日志文件或历史压缩日志
		if filename == "access_security.log" ||
			(strings.HasPrefix(filename, "access_security.log.") && strings.HasSuffix(filename, ".gz")) {

			// 添加文件到结果中
			fullPath := cm.LogPath
			if !strings.HasSuffix(fullPath, "/") {
				fullPath += "/"
			}
			logFiles = append(logFiles, fullPath+filename)
		}
	}

	return logFiles, nil
}

// GetLogsFromFile 从单个文件中读取日志，优先使用缓存
func (cm *CacheManager) GetLogsFromFile(filePath string) ([]LogRow, error) {
	// 获取文件信息
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %v", err)
	}

	// 构造缓存键
	cacheKey := CacheKey{
		FilePath: filePath,
		ModTime:  fileInfo.ModTime(),
	}

	// 检查缓存
	cm.mu.RLock()
	cachedLogs, exists := cm.cache[cacheKey]
	cm.mu.RUnlock()

	if exists {
		return cachedLogs, nil
	}

	// 缓存未命中，解析文件
	logs, err := cm.ParseLogFile(filePath)
	if err != nil {
		return nil, err
	}

	// 更新缓存
	cm.mu.Lock()
	cm.cache[cacheKey] = logs
	cm.mu.Unlock()

	return logs, nil
}

// ParseLogFile 解析单个日志文件
func (cm *CacheManager) ParseLogFile(filePath string) ([]LogRow, error) {
	// 根据文件扩展名决定是普通文件还是gzip文件
	var reader io.Reader
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开日志文件: %v", err)
	}
	defer file.Close()

	if strings.HasSuffix(filePath, ".gz") {
		// gzip压缩文件
		gzReader, err := gzip.NewReader(file)
		if err != nil {
			return nil, fmt.Errorf("无法创建gzip解压器: %v", err)
		}
		defer gzReader.Close()
		reader = gzReader
	} else {
		// 普通文件
		reader = file
	}

	return cm.parseLogReader(reader, filePath)
}

// parseLogReader 从io.Reader中解析日志内容
func (cm *CacheManager) parseLogReader(reader io.Reader, filePath string) ([]LogRow, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 0, 1024*1024), 10*1024*1024) // 设置较大的缓冲区以应对长日志行

	var logs []LogRow
	var parseErrors int
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// 跳过空行
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		// 解析日志行
		logRow, err := ParseLogLine(line)
		if err != nil {
			// 记录错误但继续处理
			parseErrors++
			// 如果错误太多，可以考虑提前返回
			if parseErrors > 1000 {
				logrus.Errorf("文件 %s 解析错误过多，已解析 %d 行，错误 %d 个", filePath, lineNum, parseErrors)
				return nil, fmt.Errorf("文件 %s 解析错误过多，已解析 %d 行，错误 %d 个",
					filePath, lineNum, parseErrors)
			}
			logrus.Warnf("解析文件 %s 的第 %d 行失败: %v", filePath, lineNum, err)
			continue
		}

		if logRow != nil {
			logs = append(logs, *logRow)
		}
	}

	if err := scanner.Err(); err != nil {
		logrus.Errorf("扫描日志文件 %s 时发生错误: %v", filePath, err)
		return nil, fmt.Errorf("扫描日志文件错误: %v", err)
	}

	if parseErrors > 0 {
		logrus.Warnf("文件 %s 中有 %d 行解析失败", filePath, parseErrors)
		fmt.Printf("警告: 文件 %s 中有 %d 行解析失败\n", filePath, parseErrors)
	}

	return logs, nil
}

// ApplyFilter 应用过滤条件
func (cm *CacheManager) ApplyFilter(logs []LogRow, filter LogFilter) []LogRow {
	var result []LogRow

	for _, log := range logs {
		if cm.MatchFilter(log, filter) {
			result = append(result, log)
		}
	}

	return result
}

// MatchFilter 判断日志条目是否匹配过滤条件
func (cm *CacheManager) MatchFilter(log LogRow, filter LogFilter) bool {
	// IP 过滤
	if filter.IP != "" && !strings.Contains(log.ClientIP, filter.IP) {
		return false
	}

	// 状态码过滤
	if filter.Status > 0 && log.Status != filter.Status {
		return false
	}

	// 路径过滤
	if filter.Path != "" && !strings.Contains(log.Path, filter.Path) {
		return false
	}

	// 时间范围过滤
	if !filter.From.IsZero() && log.Time.Before(filter.From) {
		return false
	}
	if !filter.To.IsZero() && log.Time.After(filter.To) {
		return false
	}

	return true
}

// LogFilter 表示日志过滤条件
type LogFilter struct {
	IP     string
	Status int
	Path   string
	From   time.Time
	To     time.Time
	Limit  int
	Offset int
}
