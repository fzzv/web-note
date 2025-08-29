// Package main 演示基础命令行参数解析
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// 命令行参数变量
var (
	name     = flag.String("name", "World", "要问候的名字")
	age      = flag.Int("age", 0, "年龄")
	verbose  = flag.Bool("verbose", false, "详细输出模式")
	output   = flag.String("output", "", "输出文件路径")
	logLevel = flag.String("log-level", "info", "日志级别 (debug|info|warn|error)")
)

// LogLevel 自定义日志级别类型
type LogLevel string

const (
	DEBUG LogLevel = "debug"
	INFO  LogLevel = "info"
	WARN  LogLevel = "warn"
	ERROR LogLevel = "error"
)

// String 实现 Stringer 接口
func (l LogLevel) String() string {
	return string(l)
}

// IsValid 检查日志级别是否有效
func (l LogLevel) IsValid() bool {
	switch l {
	case DEBUG, INFO, WARN, ERROR:
		return true
	default:
		return false
	}
}

// Config 应用程序配置
type Config struct {
	Name     string
	Age      int
	Verbose  bool
	Output   string
	LogLevel LogLevel
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("名字不能为空")
	}
	
	if c.Age < 0 {
		return fmt.Errorf("年龄不能为负数")
	}
	
	if !c.LogLevel.IsValid() {
		return fmt.Errorf("无效的日志级别: %s", c.LogLevel)
	}
	
	return nil
}

func main() {
	// 自定义用法信息
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "用法: %s [选项] [文件...]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "这是一个演示基础命令行参数解析的程序。\n\n")
		fmt.Fprintf(os.Stderr, "选项:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n示例:\n")
		fmt.Fprintf(os.Stderr, "  %s -name=张三 -age=25 -verbose file1.txt file2.txt\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -output=result.txt -log-level=debug\n", os.Args[0])
	}
	
	// 解析命令行参数
	flag.Parse()
	
	// 创建配置
	config := &Config{
		Name:     *name,
		Age:      *age,
		Verbose:  *verbose,
		Output:   *output,
		LogLevel: LogLevel(*logLevel),
	}
	
	// 验证配置
	if err := config.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "配置错误: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}
	
	// 显示配置信息
	if config.Verbose {
		fmt.Printf("配置信息:\n")
		fmt.Printf("  名字: %s\n", config.Name)
		fmt.Printf("  年龄: %d\n", config.Age)
		fmt.Printf("  详细模式: %t\n", config.Verbose)
		fmt.Printf("  输出文件: %s\n", config.Output)
		fmt.Printf("  日志级别: %s\n", config.LogLevel)
		fmt.Printf("\n")
	}
	
	// 处理位置参数（文件列表）
	files := flag.Args()
	if len(files) == 0 {
		// 没有提供文件，执行默认操作
		executeDefaultAction(config)
	} else {
		// 处理指定的文件
		processFiles(config, files)
	}
}

// executeDefaultAction 执行默认操作
func executeDefaultAction(config *Config) {
	message := fmt.Sprintf("Hello, %s!", config.Name)
	
	if config.Age > 0 {
		message += fmt.Sprintf(" 你今年 %d 岁。", config.Age)
	}
	
	if config.Verbose {
		fmt.Printf("[%s] %s\n", config.LogLevel, message)
	} else {
		fmt.Println(message)
	}
	
	// 如果指定了输出文件，写入文件
	if config.Output != "" {
		writeToFile(config.Output, message)
	}
}

// processFiles 处理文件列表
func processFiles(config *Config, files []string) {
	fmt.Printf("处理 %d 个文件:\n", len(files))
	
	var results []string
	
	for i, filename := range files {
		if config.Verbose {
			fmt.Printf("[%d/%d] 处理文件: %s\n", i+1, len(files), filename)
		}
		
		result, err := processFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "处理文件 %s 失败: %v\n", filename, err)
			continue
		}
		
		results = append(results, result)
		
		if config.Verbose {
			fmt.Printf("  结果: %s\n", result)
		}
	}
	
	// 输出汇总信息
	fmt.Printf("\n处理完成，成功处理 %d/%d 个文件\n", len(results), len(files))
	
	// 如果指定了输出文件，写入结果
	if config.Output != "" {
		summary := fmt.Sprintf("处理结果汇总:\n%s", strings.Join(results, "\n"))
		writeToFile(config.Output, summary)
	}
}

// processFile 处理单个文件
func processFile(filename string) (string, error) {
	// 检查文件是否存在
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return "", fmt.Errorf("文件不存在: %s", filename)
	}
	
	// 获取文件信息
	info, err := os.Stat(filename)
	if err != nil {
		return "", fmt.Errorf("获取文件信息失败: %w", err)
	}
	
	// 读取文件内容
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %w", err)
	}
	
	// 统计信息
	lines := strings.Count(string(content), "\n") + 1
	words := len(strings.Fields(string(content)))
	chars := len(content)
	
	result := fmt.Sprintf("%s: %d 行, %d 词, %d 字符, %d 字节",
		filepath.Base(filename), lines, words, chars, info.Size())
	
	return result, nil
}

// writeToFile 写入文件
func writeToFile(filename, content string) {
	// 创建目录（如果不存在）
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "创建目录失败: %v\n", err)
		return
	}
	
	// 写入文件
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "写入文件失败: %v\n", err)
		return
	}
	
	fmt.Printf("结果已写入文件: %s\n", filename)
}

// 演示不同的参数解析方式
func demonstrateAdvancedParsing() {
	// 自定义 FlagSet
	fs := flag.NewFlagSet("advanced", flag.ExitOnError)
	
	var (
		host = fs.String("host", "localhost", "服务器主机")
		port = fs.Int("port", 8080, "服务器端口")
		ssl  = fs.Bool("ssl", false, "启用 SSL")
	)
	
	// 解析特定的参数
	args := []string{"-host=example.com", "-port=443", "-ssl"}
	fs.Parse(args)
	
	fmt.Printf("高级解析结果:\n")
	fmt.Printf("  主机: %s\n", *host)
	fmt.Printf("  端口: %d\n", *port)
	fmt.Printf("  SSL: %t\n", *ssl)
}

// 演示环境变量回退
func demonstrateEnvFallback() {
	// 从环境变量获取默认值
	defaultName := os.Getenv("USER_NAME")
	if defaultName == "" {
		defaultName = "Anonymous"
	}
	
	defaultPort := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		defaultPort = envPort
	}
	
	fmt.Printf("环境变量回退:\n")
	fmt.Printf("  默认用户名: %s\n", defaultName)
	fmt.Printf("  默认端口: %s\n", defaultPort)
}
