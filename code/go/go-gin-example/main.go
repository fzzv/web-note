package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/fzzv/go-gin-example/docs"
	"github.com/fzzv/go-gin-example/models"
	"github.com/fzzv/go-gin-example/pkg/gredis"
	"github.com/fzzv/go-gin-example/pkg/logging"
	"github.com/fzzv/go-gin-example/pkg/setting"
	"github.com/fzzv/go-gin-example/routers"
)

func main() {
	setting.Setup()
	models.Setup()
	logging.Setup()
	gredis.Setup()
	router := routers.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.ServerSetting.HttpPort),
		Handler:        router,
		ReadTimeout:    setting.ServerSetting.ReadTimeout,
		WriteTimeout:   setting.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	err := s.ListenAndServe()
	if err != nil {
		log.Printf("Server err: %v", err)
	}
}

// endless 实现优雅重启
// func main() {
// 	setting.Setup()
// 	models.Setup()
// 	logging.Setup()

// 	endless.DefaultReadTimeOut = setting.ServerSetting.ReadTimeout
// 	endless.DefaultWriteTimeOut = setting.ServerSetting.WriteTimeout
// 	endless.DefaultMaxHeaderBytes = 1 << 20
// 	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)

// 	server := endless.NewServer(endPoint, routers.InitRouter())
// 	server.BeforeBegin = func(add string) {
// 		log.Printf("Actual pid is %d", syscall.Getpid())
// 	}

// 	err := server.ListenAndServe()
// 	if err != nil {
// 		log.Printf("Server err: %v", err)
// 	}
// }
