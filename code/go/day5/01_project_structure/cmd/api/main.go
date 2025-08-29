// Package main 演示标准 Go 项目结构中的 API 服务器入口点
package main

import (
	"fmt"
	"log"
	"net/http"

	"day5/01_project_structure/internal/config"
	"day5/01_project_structure/internal/handler"
	"day5/01_project_structure/internal/service"
	"day5/01_project_structure/pkg/logger"
)

func main() {
	// 1. 初始化配置
	cfg := config.Load()
	fmt.Printf("启动 API 服务器，端口: %s\n", cfg.Port)

	// 2. 初始化日志器
	log := logger.New(cfg.LogLevel)
	log.Info("API 服务器启动中...")

	// 3. 初始化服务层
	userService := service.NewUserService()

	// 4. 初始化处理器层
	userHandler := handler.NewUserHandler(userService, log)

	// 5. 设置路由
	mux := http.NewServeMux()
	mux.HandleFunc("/users", userHandler.GetUsers)
	mux.HandleFunc("/users/", userHandler.GetUser)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// 6. 启动服务器
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	log.Info("API 服务器已启动", "port", cfg.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Error("服务器启动失败", "error", err)
	}
}
