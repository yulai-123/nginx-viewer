package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

// Config 应用配置结构体
type Config struct {
	LogPath string `yaml:"log_path"` // Nginx日志路径
}

// LoadConfig 从YAML文件中加载配置
func LoadConfig(configPath string) (*Config, error) {
	// 检查文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		logrus.Warnf("配置文件不存在: %s, 将使用默认配置", configPath)
		return &Config{LogPath: "/var/log/nginx"}, nil
	}

	// 读取配置文件
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析YAML
	config := &Config{
		LogPath: "/var/log/nginx", // 默认值
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 验证配置
	if config.LogPath == "" {
		return nil, fmt.Errorf("log_path 不能为空")
	}

	// 确保日志路径存在
	if _, err := os.Stat(config.LogPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("日志路径不存在: %s", config.LogPath)
	}

	return config, nil
}
