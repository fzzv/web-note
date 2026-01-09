package main

import (
	"context"
	"go-grpc-example/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

// server 用于实现 Greeter 服务
type Server struct {
	proto.UnimplementedGreeterServiceServer
}

// 实现SayHello方法
func (s *Server) SayHello(ctx context.Context, in *proto.HelloRequest) (*proto.HelloResponse, error) {
	log.Printf("收到客户端请求: %v", in.GetName())
	return &proto.HelloResponse{Message: "Hello, " + in.Name}, nil
}

func main() {
	// 1. 监听端口
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("监听端口失败: %v", err)
	}
	// 2. 创建 grpc 服务器
	s := grpc.NewServer()
	// 3. 注册服务
	proto.RegisterGreeterServiceServer(s, &Server{})
	log.Printf("服务正在监听 :50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
