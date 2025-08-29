// Package main 演示后台任务处理器
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"day5/01_project_structure/internal/config"
	"day5/01_project_structure/internal/worker"
	"day5/01_project_structure/pkg/logger"
)

func main() {
	// 1. 初始化配置
	cfg := config.Load()
	fmt.Printf("启动后台任务处理器，工作线程数: %d\n", cfg.WorkerCount)

	// 2. 初始化日志器
	log := logger.New(cfg.LogLevel)
	log.Info("后台任务处理器启动中...")

	// 3. 创建上下文，用于优雅关闭
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 4. 初始化工作器
	w := worker.New(cfg.WorkerCount, log)

	// 5. 启动工作器
	go func() {
		if err := w.Start(ctx); err != nil {
			log.Error("工作器启动失败", "error", err)
			cancel()
		}
	}()

	// 6. 模拟添加任务
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		taskID := 1
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				task := worker.Task{
					ID:   fmt.Sprintf("task-%d", taskID),
					Type: "email",
					Data: map[string]interface{}{
						"to":      "user@example.com",
						"subject": fmt.Sprintf("任务通知 #%d", taskID),
						"body":    "这是一个测试任务",
					},
				}
				w.AddTask(task)
				log.Info("添加任务", "taskID", task.ID)
				taskID++
			}
		}
	}()

	// 7. 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Info("后台任务处理器已启动，按 Ctrl+C 停止")
	<-sigChan

	log.Info("收到停止信号，正在优雅关闭...")
	cancel()

	// 8. 等待工作器停止
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := w.Shutdown(shutdownCtx); err != nil {
		log.Error("工作器关闭失败", "error", err)
	} else {
		log.Info("后台任务处理器已停止")
	}
}
