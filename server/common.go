package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

var (
	ShanghaiLocation, _ = time.LoadLocation("Asia/Shanghai")
)

func InitLogger() error {
	// 创建日志目录
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("create log directory failed: %v", err)
	}

	// 使用当前日期作为日志文件名
	currentTime := time.Now().In(ShanghaiLocation)
	logFileName := fmt.Sprintf("%s.log", currentTime.Format("2006-01-02_15"))
	logFilePath := path.Join(logDir, logFileName)

	// 打开日志文件
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("open log file failed: %v", err)
	}

	// 配置 logrus
	logrus.SetOutput(logFile)
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})

	// 设置日志级别
	logrus.SetLevel(logrus.InfoLevel)

	// 每天凌晨自动切换日志文件
	go nextLogFile(logFile)

	return nil
}

// nextLogFile 用于每天凌晨更新日志文件
func nextLogFile(oldLogFile *os.File) {
	logDir := "logs"
	logFile := oldLogFile

	for {
		now := time.Now().In(ShanghaiLocation)
		// 计算下一个整点
		next := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, ShanghaiLocation)
		// 等待到下一个整点
		time.Sleep(next.Sub(now))

		// 创建新的日志文件
		newLogFileName := fmt.Sprintf("%s.log", next.Format("2006-01-02_15"))
		newLogFilePath := path.Join(logDir, newLogFileName)
		newLogFile, err := os.OpenFile(newLogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			logrus.Errorf("create new log file failed: %v", err)
			continue
		}

		// 切换到新的日志文件
		logrus.SetOutput(newLogFile)
		// 关闭旧的日志文件
		logFile.Close()
		logFile = newLogFile
	}
}
