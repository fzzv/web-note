package main

import (
	"fmt"
	"go-grpc-example/proto"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type Server struct {
	proto.UnimplementedGreeterServiceServer
}

// 服务端流 模式通信，类似下载文件或订阅消息
// 客户端发送一个请求，服务端返回一个流，客户端接收流中的数据
func (s *Server) LotsOfReplies(in *proto.HelloRequest, stream proto.GreeterService_LotsOfRepliesServer) error {
	log.Printf("收到流式请求，来自: %v", in.GetName())
	for i := range 10 {
		res := &proto.HelloResponse{
			Message: fmt.Sprintf("你好 %s, 这是第 %d 份数据包 📦", in.GetName(), i),
		}
		if err := stream.Send(res); err != nil {
			return err
		}
		time.Sleep(1 * time.Second) // 模拟耗时操作
	}
	log.Println("数据发送完毕，关闭流")
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("监听端口失败: %v", err)
	}
	defer lis.Close()
	s := grpc.NewServer()
	proto.RegisterGreeterServiceServer(s, &Server{})
	log.Printf("服务正在监听 :50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("启动服务失败: %v", err)
	}
}
