// Package config 提供应用程序配置管理
package config

import (
	"os"
	"strconv"
)

// Config 应用程序配置结构
type Config struct {
	Port        string // API 服务器端口
	LogLevel    string // 日志级别
	WorkerCount int    // 后台工作线程数量
	DatabaseURL string // 数据库连接字符串
}

// Load 从环境变量加载配置，提供默认值
func Load() *Config {
	cfg := &Config{
		Port:        getEnv("PORT", "8080"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		WorkerCount: getEnvAsInt("WORKER_COUNT", 3),
		DatabaseURL: getEnv("DATABASE_URL", "sqlite://./app.db"),
	}

	return cfg
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 获取环境变量并转换为整数，如果不存在或转换失败则返回默认值
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// Validate 验证配置的有效性
func (c *Config) Validate() error {
	// 这里可以添加配置验证逻辑
	// 例如：检查端口范围、日志级别有效性等
	return nil
}
