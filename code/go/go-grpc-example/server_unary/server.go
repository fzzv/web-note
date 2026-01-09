package main

import (
	"context"
	"go-grpc-example/proto"
	"log"
	"net"
	"time"

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

// 实现一个简单的日志拦截器
func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	res, err := handler(ctx, req)
	duration := time.Since(start)
	log.Printf("📢 [gRPC 日志] 方法: %s | 耗时: %v | 错误: %v", info.FullMethod, duration, err)
	return res, err
}

func main() {
	// 1. 监听端口
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("监听端口失败: %v", err)
	}
	// 2. 创建 grpc 服务器
	s := grpc.NewServer(
		// 注册拦截器
		grpc.UnaryInterceptor(loggingInterceptor),
	)
	// 3. 注册服务
	proto.RegisterGreeterServiceServer(s, &Server{})
	log.Printf("服务正在监听 :50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
