package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"logic-grpc-service/config"
	"logic-grpc-service/internal/ai"
	"logic-grpc-service/internal/db"
	"logic-grpc-service/internal/oss"
	"logic-grpc-service/internal/service"
	"logic-grpc-service/proto/gen/ai_chat"
	"logic-grpc-service/proto/gen/application"
	"logic-grpc-service/proto/gen/auth"
	"logic-grpc-service/proto/gen/position"
	"logic-grpc-service/proto/gen/resume"
	"logic-grpc-service/proto/gen/user"
)

func main() {
	fmt.Println("正在启动智能招聘系统 gRPC 服务...")

	if err := config.Init(); err != nil {
		log.Fatalf("初始化配置失败: %v", err)
	}
	fmt.Println("配置初始化成功")

	if err := db.Init(); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	fmt.Println("数据库连接成功")

	if err := oss.Init(); err != nil {
		log.Printf("初始化OSS客户端失败: %v", err)
		fmt.Println("注意：OSS功能将不可用")
	} else {
		fmt.Println("OSS客户端初始化成功")
	}

	database := db.GetDB()

	if err := ai.Init(database); err != nil {
		log.Printf("初始化AI客户端失败: %v", err)
		fmt.Println("注意：AI聊天功能将不可用")
	} else {
		fmt.Println("AI客户端初始化成功")
	}

	serverConfig := config.GetServerConfig()
	lis, err := net.Listen("tcp", serverConfig.Port)
	if err != nil {
		log.Fatalf("监听端口失败: %v", err)
	}
	fmt.Printf("正在监听端口 %s\n", serverConfig.Port)

	s := grpc.NewServer()
	auth.RegisterAuthServiceServer(s, service.NewAuthService(database))
	user.RegisterUserServiceServer(s, service.NewUserService(database))
	position.RegisterPositionServiceServer(s, service.NewPositionService(database))
	application.RegisterApplicationServiceServer(s, service.NewApplicationService(database))
	resume.RegisterResumeServiceServer(s, service.NewResumeService(database))
	ai_chat.RegisterAIChatServiceServer(s, service.NewAIChatService(database))

	reflection.Register(s)

	fmt.Println("所有服务已注册，准备启动 gRPC 服务器...")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
