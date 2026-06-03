package main

import (
	"fmt"
	"log"

	"web-gin-service/config"
	grpcclient "web-gin-service/internal/grpc"
	"web-gin-service/internal/router"
)

func main() {
	fmt.Println("正在启动智能招聘系统 HTTP 网关服务...")

	if err := config.Init(); err != nil {
		log.Fatalf("初始化配置失败: %v", err)
	}
	fmt.Println("配置初始化成功")

	if err := grpcclient.Init(); err != nil {
		log.Fatalf("初始化gRPC客户端失败: %v", err)
	}
	fmt.Println("gRPC客户端初始化成功")

	r := router.NewRouter()
	engine := r.SetupRoutes()

	serverConfig := config.GetServerConfig()
	addr := serverConfig.Port
	if addr == "" {
		addr = ":8080"
	}

	fmt.Printf("正在监听端口 %s\n", addr)
	fmt.Println("HTTP网关启动成功")

	if err := engine.Run(addr); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
