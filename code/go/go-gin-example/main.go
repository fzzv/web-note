package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/fzzv/go-gin-example/docs"
	"github.com/fzzv/go-gin-example/pkg/setting"
	"github.com/fzzv/go-gin-example/routers"
)

func main() {
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
// 	endless.DefaultReadTimeOut = setting.ReadTimeout
// 	endless.DefaultWriteTimeOut = setting.WriteTimeout
// 	endless.DefaultMaxHeaderBytes = 1 << 20
// 	endPoint := fmt.Sprintf(":%d", setting.HTTPPort)

// 	server := endless.NewServer(endPoint, routers.InitRouter())
// 	server.BeforeBegin = func(add string) {
// 		log.Printf("Actual pid is %d", syscall.Getpid())
// 	}

// 	err := server.ListenAndServe()
// 	if err != nil {
// 		log.Printf("Server err: %v", err)
// 	}
// }
