package grpc

import (
	"fmt"
	"log"

	"web-gin-service/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"web-gin-service/proto/gen/ai_chat"
	"web-gin-service/proto/gen/application"
	"web-gin-service/proto/gen/auth"
	"web-gin-service/proto/gen/position"
	"web-gin-service/proto/gen/resume"
	"web-gin-service/proto/gen/user"
)

var (
	authClient        auth.AuthServiceClient
	userClient        user.UserServiceClient
	positionClient    position.PositionServiceClient
	applicationClient application.ApplicationServiceClient
	resumeClient      resume.ResumeServiceClient
	aiChatClient      ai_chat.AIChatServiceClient
)

func Init() error {
	addr := config.GetGRPCAddress()

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("连接gRPC服务失败: %w", err)
	}

	log.Printf("gRPC连接成功: %s", addr)

	authClient = auth.NewAuthServiceClient(conn)
	userClient = user.NewUserServiceClient(conn)
	positionClient = position.NewPositionServiceClient(conn)
	applicationClient = application.NewApplicationServiceClient(conn)
	resumeClient = resume.NewResumeServiceClient(conn)
	aiChatClient = ai_chat.NewAIChatServiceClient(conn)

	return nil
}

func GetAuthClient() auth.AuthServiceClient {
	return authClient
}

func GetUserClient() user.UserServiceClient {
	return userClient
}

func GetPositionClient() position.PositionServiceClient {
	return positionClient
}

func GetApplicationClient() application.ApplicationServiceClient {
	return applicationClient
}

func GetResumeClient() resume.ResumeServiceClient {
	return resumeClient
}

func GetAIChatClient() ai_chat.AIChatServiceClient {
	return aiChatClient
}
