package main

import (
	"google.golang.org/grpc"
	"log"
	"net"

	"user/config"
	"user/internal/handler"
	"user/internal/service/pb"
)

func main() {
	// 监听本地的8972端口
	lis, err := net.Listen("tcp", config.DefaultConfig.Grpc.Addr)
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer() // 创建gRPC服务器
	defer s.Stop()
	// 在gRPC服务端注册服务
	pb.RegisterUserServiceServer(s, &handler.UserServiceServer{})
	log.Printf("start user rpc server")
	// 启动服务
	err = s.Serve(lis)
	if err != nil {
		log.Printf("failed to serve: %v", err)
		return
	}
}
