package main

import (
	"context"
	"go-grpc-example/proto"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// 服务端流 模式通信，类似下载文件或订阅消息
// 客户端发送一个请求，服务端返回一个流，客户端接收流中的数据
// func main() {
// 	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		log.Fatalf("连接服务失败: %v", err)
// 	}
// 	defer conn.Close()
// 	c := proto.NewGreeterServiceClient(conn)
// 	stream, err := c.LotsOfReplies(context.Background(), &proto.HelloRequest{Name: "world"})
// 	if err != nil {
// 		log.Fatalf("调用 LotsOfReplies 失败: %v", err)
// 	}
// 	for {
// 		// 循环接收流中的数据
// 		res, err := stream.Recv()
// 		if err == io.EOF {
// 			log.Println("✅ 所有数据接收完毕")
// 			break
// 		}
// 		if err != nil {
// 			log.Fatalf("接收出错: %v", err)
// 		}
// 		log.Printf("收到消息: %s", res.GetMessage())
// 	}
// }

// 客户端流 模式通信，类似上传文件
// 客户端发送一个流，服务端返回一个响应
// func main() {
// 	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		log.Fatalf("连接服务失败: %v", err)
// 	}
// 	defer conn.Close()
// 	c := proto.NewGreeterServiceClient(conn)
// 	stream, err := c.CollectNames(context.Background())
// 	if err != nil {
// 		log.Fatalf("调用 CollectNames 失败: %v", err)
// 	}
// 	names := []string{"Fan", "Xyu", "Xxx"}
// 	for _, name := range names {
// 		if err := stream.Send(&proto.HelloRequest{Name: name}); err != nil {
// 			log.Fatalf("发送名字失败: %v", err)
// 		}
// 	}
// 	// 发送完毕并获取响应
// 	res, err := stream.CloseAndRecv()
// 	if err != nil {
// 		log.Fatalf("接收响应失败: %v", err)
// 	}
// 	log.Printf("收到响应: %s", res.GetMessage())
// }

// 双向流 模式通信，类似聊天
// 客户端发送一个流，服务端返回一个流
func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接服务失败: %v", err)
	}
	defer conn.Close()
	c := proto.NewGreeterServiceClient(conn)
	stream, err := c.Chat(context.Background())
	// 1. 开启一个协程专门发送消息
	go func() {
		messages := []string{"你好", "在吗", "gRPC 真好用"}
		for _, m := range messages {
			stream.Send(&proto.HelloRequest{Name: m})
			time.Sleep(time.Second)
		}
		stream.CloseSend() // 发送完毕后关闭发送端
	}()

	// 2. 主线程负责接收消息
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("接收失败: %v", err)
		}
		log.Printf("来自服务器的回应: %s", res.GetMessage())
	}
}
