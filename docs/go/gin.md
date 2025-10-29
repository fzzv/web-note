# Gin

## 初始化目录

- conf：用于存储配置文件
- middleware：应用中间件
- models：应用数据库模型
- pkg：第三方包
- routers 路由逻辑处理
- runtime：应用运行时数据

## Go Modules Replace

打开 `go.mod` 文件，新增 `replace` 配置项，使用的是完整的外部模块引用路径（`github.com/fzzv/go-gin-example/xxx`），而这个模块还没推送到远程，是没有办法下载下来的，因此需要用 `replace` 将其指定读取本地的模块路径，这样子就可以解决本地模块读取的问题

```go
module github.com/fzzv/go-gin-example

go 1.24.6

require (
	github.com/gin-gonic/gin v1.11.0
	github.com/go-ini/ini v1.67.0
)

require (...)

replace (
	github.com/fzzv/go-gin-example/conf => E:/other-project/web-note/code/go/go-gin-example/pkg/conf
	github.com/fzzv/go-gin-example/middleware => E:/other-project/web-note/code/go/go-gin-example/middleware
	github.com/fzzv/go-gin-example/models => E:/other-project/web-note/code/go/go-gin-example/models
	github.com/fzzv/go-gin-example/pkg/e => E:/other-project/web-note/code/go/go-gin-example/pkg/e
	github.com/fzzv/go-gin-example/pkg/setting => E:/other-project/web-note/code/go/go-gin-example/pkg/setting
	github.com/fzzv/go-gin-example/routers => E:/other-project/web-note/code/go/go-gin-example/routers
)
```

## 一些第三方包

- github.com/go-ini/ini：用于编写应用配置文件
- github.com/unknwon/com：提供一组**实用的通用函数（helper functions）**，用于字符串、文件、类型转换、切片、结构体反射等常见任务
- github.com/jinzhu/gorm：Go 对象关系映射框架

## 简单的demo

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fzzv/go-gin-example/pkg/setting"
)

func main() {
	// 返回Gin的type Engine struct{...}，里面包含RouterGroup，
	// 相当于创建一个路由Handlers，可以后期绑定各类的路由规则和函数、中间件等
	router := gin.Default()
	// 创建不同的HTTP方法绑定到Handlers中，
	// 也支持POST、PUT、DELETE、PATCH、OPTIONS、HEAD 等常用的Restful方法
	// Context是gin中的上下文，它允许我们在中间件之间传递变量、管理流、验证JSON请求、响应JSON请求等，
	// 在gin中包含大量Context的方法，例如我们常用的DefaultQuery、Query、DefaultPostForm、PostForm等等
	router.GET("/test", func(c *gin.Context) {
		// gin.H 是 map[string]interface{} 的缩写
		c.JSON(200, gin.H{
			"message": "test",
		})
	})

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort), // 监听的TCP地址，格式为:8000
		Handler:        router,                               // http句柄，实质为ServeHTTP，用于处理程序响应HTTP请求
		ReadTimeout:    setting.ReadTimeout,                  // 允许读取的最大时间
		WriteTimeout:   setting.WriteTimeout,                 // 允许写入的最大时间
		MaxHeaderBytes: 1 << 20,                              // 请求头的最大字节数
	}

	s.ListenAndServe() // 开始监听服务，监听TCP网络地址，Addr和调用应用程序处理连接上的请求
}
```
