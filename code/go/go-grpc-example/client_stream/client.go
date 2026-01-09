package main

import (
	"context"
	"go-grpc-example/proto"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接服务失败: %v", err)
	}
	defer conn.Close()
	c := proto.NewGreeterServiceClient(conn)
	stream, err := c.LotsOfReplies(context.Background(), &proto.HelloRequest{Name: "world"})
	if err != nil {
		log.Fatalf("调用 LotsOfReplies 失败: %v", err)
	}
	for {
		// 循环接收流中的数据
		res, err := stream.Recv()
		if err == io.EOF {
			log.Println("✅ 所有数据接收完毕")
			break
		}
		if err != nil {
			log.Fatalf("接收出错: %v", err)
		}
		log.Printf("收到消息: %s", res.GetMessage())
	}
}
