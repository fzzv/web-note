package main

import (
	"log"
	"time"

	"github.com/fzzv/go-gin-example/models"
	"github.com/robfig/cron"
)

// func main() {
func Cron() {
	log.Println("Starting...")
	// cron.New()会根据本地时间创建一个新（空白）的 Cron job runner
	c := cron.New()
	// cron.AddFunc 会向 Cron job runner 添加一个 func ，以按给定的时间表运行，
	// 第一个参数是 cron 表达式，第二个参数是任务函数
	c.AddFunc("* * * * * *", func() {
		log.Println("Run models.CleanAllTag...")
		models.CleanAllTag()
	})
	c.AddFunc("* * * * * *", func() {
		log.Println("Run models.CleanAllArticle...")
		models.CleanAllArticle()
	})
	// 在当前执行的程序中启动 Cron 调度程序
	c.Start()
	// 会创建一个新的定时器，持续你设定的时间 d 后发送一个 channel 消息
	t1 := time.NewTimer(time.Second * 10)
	// for + select 阻塞 select 等待 channel
	// for {
	// 	select {
	// 	case <-t1.C:
	// 		// 重置定时器，让它重新开始计时
	// 		t1.Reset(time.Second * 10)
	// 	}
	// }
	// 可以简写为
	for range t1.C {
		t1.Reset(time.Second * 10)
	}
}
