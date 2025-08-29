// Package worker 提供后台任务处理功能
package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"day5/01_project_structure/pkg/logger"
)

// Task 表示一个后台任务
type Task struct {
	ID   string                 `json:"id"`
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

// Worker 后台任务处理器
type Worker struct {
	workerCount int
	taskQueue   chan Task
	quit        chan struct{}
	wg          sync.WaitGroup
	logger      logger.Logger
}

// New 创建新的工作器
func New(workerCount int, logger logger.Logger) *Worker {
	return &Worker{
		workerCount: workerCount,
		taskQueue:   make(chan Task, 100), // 缓冲队列，最多100个任务
		quit:        make(chan struct{}),
		logger:      logger,
	}
}

// Start 启动工作器
func (w *Worker) Start(ctx context.Context) error {
	w.logger.Info("启动工作器", "workerCount", w.workerCount)

	// 启动工作协程
	for i := 0; i < w.workerCount; i++ {
		w.wg.Add(1)
		go w.worker(ctx, i+1)
	}

	// 等待上下文取消
	<-ctx.Done()
	w.logger.Info("收到停止信号，正在停止工作器...")

	// 关闭任务队列
	close(w.quit)

	// 等待所有工作协程完成
	w.wg.Wait()
	w.logger.Info("所有工作器已停止")

	return nil
}

// AddTask 添加任务到队列
func (w *Worker) AddTask(task Task) {
	select {
	case w.taskQueue <- task:
		w.logger.Debug("任务已添加到队列", "taskID", task.ID)
	default:
		w.logger.Warn("任务队列已满，丢弃任务", "taskID", task.ID)
	}
}

// Shutdown 优雅关闭工作器
func (w *Worker) Shutdown(ctx context.Context) error {
	// 关闭任务队列，不再接受新任务
	close(w.taskQueue)

	// 等待所有任务完成或超时
	done := make(chan struct{})
	go func() {
		w.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		w.logger.Info("所有任务已完成")
		return nil
	case <-ctx.Done():
		w.logger.Warn("关闭超时，强制停止")
		return ctx.Err()
	}
}

// worker 工作协程
func (w *Worker) worker(ctx context.Context, workerID int) {
	defer w.wg.Done()

	w.logger.Info("工作器启动", "workerID", workerID)

	for {
		select {
		case task, ok := <-w.taskQueue:
			if !ok {
				w.logger.Info("任务队列已关闭", "workerID", workerID)
				return
			}

			w.processTask(workerID, task)

		case <-w.quit:
			w.logger.Info("工作器收到停止信号", "workerID", workerID)
			return

		case <-ctx.Done():
			w.logger.Info("工作器上下文已取消", "workerID", workerID)
			return
		}
	}
}

// processTask 处理单个任务
func (w *Worker) processTask(workerID int, task Task) {
	start := time.Now()
	w.logger.Info("开始处理任务",
		"workerID", workerID,
		"taskID", task.ID,
		"taskType", task.Type)

	// 模拟任务处理
	switch task.Type {
	case "email":
		w.processEmailTask(task)
	case "notification":
		w.processNotificationTask(task)
	case "data_processing":
		w.processDataTask(task)
	default:
		w.logger.Warn("未知任务类型", "taskType", task.Type, "taskID", task.ID)
		return
	}

	duration := time.Since(start)
	w.logger.Info("任务处理完成",
		"workerID", workerID,
		"taskID", task.ID,
		"duration", duration)
}

// processEmailTask 处理邮件任务
func (w *Worker) processEmailTask(task Task) {
	// 模拟邮件发送
	to, _ := task.Data["to"].(string)
	subject, _ := task.Data["subject"].(string)

	w.logger.Info("发送邮件",
		"taskID", task.ID,
		"to", to,
		"subject", subject)

	// 模拟处理时间
	time.Sleep(time.Duration(500+workerID*100) * time.Millisecond)

	w.logger.Info("邮件发送成功", "taskID", task.ID)
}

// processNotificationTask 处理通知任务
func (w *Worker) processNotificationTask(task Task) {
	// 模拟推送通知
	userID, _ := task.Data["user_id"].(string)
	message, _ := task.Data["message"].(string)

	w.logger.Info("发送通知",
		"taskID", task.ID,
		"userID", userID,
		"message", message)

	// 模拟处理时间
	time.Sleep(300 * time.Millisecond)

	w.logger.Info("通知发送成功", "taskID", task.ID)
}

// processDataTask 处理数据任务
func (w *Worker) processDataTask(task Task) {
	// 模拟数据处理
	dataSize, _ := task.Data["size"].(int)

	w.logger.Info("处理数据",
		"taskID", task.ID,
		"dataSize", dataSize)

	// 模拟处理时间（数据越大处理时间越长）
	processingTime := time.Duration(dataSize*10) * time.Millisecond
	if processingTime > 2*time.Second {
		processingTime = 2 * time.Second
	}
	time.Sleep(processingTime)

	w.logger.Info("数据处理完成", "taskID", task.ID)
}
