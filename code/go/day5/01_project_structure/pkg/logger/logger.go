// Package logger 提供结构化日志功能
package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// LogLevel 日志级别
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// String 返回日志级别的字符串表示
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger 日志器接口
type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
}

// SimpleLogger 简单的日志器实现
type SimpleLogger struct {
	level  LogLevel
	logger *log.Logger
}

// New 创建新的日志器
func New(levelStr string) Logger {
	level := parseLogLevel(levelStr)
	
	logger := log.New(os.Stdout, "", 0) // 不使用默认前缀，我们自己格式化

	return &SimpleLogger{
		level:  level,
		logger: logger,
	}
}

// parseLogLevel 解析日志级别字符串
func parseLogLevel(levelStr string) LogLevel {
	switch strings.ToUpper(levelStr) {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARN", "WARNING":
		return WARN
	case "ERROR":
		return ERROR
	default:
		return INFO
	}
}

// Debug 记录调试日志
func (l *SimpleLogger) Debug(msg string, fields ...interface{}) {
	if l.level <= DEBUG {
		l.log(DEBUG, msg, fields...)
	}
}

// Info 记录信息日志
func (l *SimpleLogger) Info(msg string, fields ...interface{}) {
	if l.level <= INFO {
		l.log(INFO, msg, fields...)
	}
}

// Warn 记录警告日志
func (l *SimpleLogger) Warn(msg string, fields ...interface{}) {
	if l.level <= WARN {
		l.log(WARN, msg, fields...)
	}
}

// Error 记录错误日志
func (l *SimpleLogger) Error(msg string, fields ...interface{}) {
	if l.level <= ERROR {
		l.log(ERROR, msg, fields...)
	}
}

// log 内部日志记录方法
func (l *SimpleLogger) log(level LogLevel, msg string, fields ...interface{}) {
	// 格式化时间戳
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	
	// 构建日志消息
	logMsg := fmt.Sprintf("[%s] %s %s", timestamp, level.String(), msg)
	
	// 添加字段
	if len(fields) > 0 {
		fieldStr := formatFields(fields...)
		if fieldStr != "" {
			logMsg += " " + fieldStr
		}
	}
	
	l.logger.Println(logMsg)
}

// formatFields 格式化字段
func formatFields(fields ...interface{}) string {
	if len(fields) == 0 {
		return ""
	}
	
	var parts []string
	
	// 字段应该成对出现：key, value, key, value...
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			key := fmt.Sprintf("%v", fields[i])
			value := fmt.Sprintf("%v", fields[i+1])
			parts = append(parts, fmt.Sprintf("%s=%s", key, value))
		} else {
			// 奇数个字段，最后一个作为单独的值
			parts = append(parts, fmt.Sprintf("%v", fields[i]))
		}
	}
	
	return strings.Join(parts, " ")
}

// NullLogger 空日志器，用于测试
type NullLogger struct{}

// NewNull 创建空日志器
func NewNull() Logger {
	return &NullLogger{}
}

func (n *NullLogger) Debug(msg string, fields ...interface{}) {}
func (n *NullLogger) Info(msg string, fields ...interface{})  {}
func (n *NullLogger) Warn(msg string, fields ...interface{})  {}
func (n *NullLogger) Error(msg string, fields ...interface{}) {}
