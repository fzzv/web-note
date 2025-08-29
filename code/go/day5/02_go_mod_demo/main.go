// Package main 演示 Go 模块管理的各种用法
package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Go 模块管理演示")
	
	// 演示使用外部依赖
	demonstrateGin()
	
	// 演示版本管理
	demonstrateVersions()
}

// demonstrateGin 演示使用 Gin 框架
func demonstrateGin() {
	fmt.Println("\n=== Gin 框架演示 ===")
	
	// 创建 Gin 路由器
	r := gin.New()
	
	// 添加中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	
	// 定义路由
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello from Gin!",
			"version": gin.Version,
		})
	})
	
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})
	
	fmt.Printf("Gin 版本: %s\n", gin.Version)
	fmt.Println("服务器配置完成，可以通过 r.Run() 启动")
	
	// 注意：这里不实际启动服务器，只是演示配置
	// r.Run(":8080")
}

// demonstrateVersions 演示版本管理概念
func demonstrateVersions() {
	fmt.Println("\n=== 版本管理演示 ===")
	
	// 演示语义化版本
	versions := []string{
		"v1.0.0",   // 主版本.次版本.修订版本
		"v1.1.0",   // 新功能，向后兼容
		"v1.1.1",   // 错误修复，向后兼容
		"v2.0.0",   // 重大更改，可能不向后兼容
	}
	
	fmt.Println("语义化版本示例:")
	for _, v := range versions {
		fmt.Printf("  %s\n", v)
	}
	
	// 演示版本约束
	constraints := map[string]string{
		"^1.2.3":  "兼容 1.2.3，允许 1.x.x（x >= 2.3）",
		"~1.2.3":  "兼容 1.2.3，允许 1.2.x（x >= 3）",
		">=1.2.3": "大于等于 1.2.3",
		"<2.0.0":  "小于 2.0.0",
		"latest":  "最新版本",
	}
	
	fmt.Println("\n版本约束示例:")
	for constraint, desc := range constraints {
		fmt.Printf("  %-8s: %s\n", constraint, desc)
	}
}

// ModuleInfo 模块信息结构
type ModuleInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Sum     string `json:"sum"`
}

// GetModuleInfo 获取模块信息（模拟）
func GetModuleInfo() []ModuleInfo {
	return []ModuleInfo{
		{
			Name:    "github.com/gin-gonic/gin",
			Version: "v1.9.1",
			Sum:     "h1:4idEAncQnU5cB7BeOkPtxjfCSye0AAm1R0RVIqJ+Jmg=",
		},
		{
			Name:    "github.com/stretchr/testify",
			Version: "v1.8.4",
			Sum:     "h1:CcVxjf3Q8BY+uce2d4aBtdWc+/Y5wvFMnOtHjKEHduE=",
		},
	}
}

// DemonstrateModuleCommands 演示模块命令（注释形式）
func DemonstrateModuleCommands() {
	fmt.Println("\n=== 常用模块命令 ===")
	
	commands := []struct {
		command string
		desc    string
	}{
		{"go mod init <module-name>", "初始化新模块"},
		{"go get <package>@<version>", "添加或更新依赖"},
		{"go get -u ./...", "更新所有依赖"},
		{"go mod tidy", "清理依赖"},
		{"go mod download", "下载依赖"},
		{"go list -m all", "列出所有依赖"},
		{"go mod graph", "显示依赖图"},
		{"go mod verify", "验证依赖"},
		{"go mod why <package>", "解释为什么需要某个依赖"},
	}
	
	for _, cmd := range commands {
		fmt.Printf("  %-25s: %s\n", cmd.command, cmd.desc)
	}
}
