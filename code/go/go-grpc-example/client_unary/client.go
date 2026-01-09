package main

import (
	"context"
	"go-grpc-example/proto"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 1. 建立与服务器的连接
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接服务器失败: %v", err)
	}
	defer conn.Close()
	// 2. 创建 Greeter 客户端
	c := proto.NewGreeterServiceClient(conn)

	// 3. 准备请求上下文（设置 1 秒超时）
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// 4. 发起 RPC 调用
	r, err := c.SayHello(ctx, &proto.HelloRequest{Name: "world"})
	if err != nil {
		log.Fatalf("调用 SayHello 失败: %v", err)
	}
	log.Printf("收到回复: %s", r.GetMessage())
}
