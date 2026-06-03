package service

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"logic-grpc-service/internal/middleware"
	"logic-grpc-service/internal/model"
	"logic-grpc-service/proto/gen/auth"
)

type AuthService struct {
	auth.UnimplementedAuthServiceServer
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) Register(ctx context.Context, req *auth.RegisterReq) (*auth.RegisterResp, error) {
	var existingUser model.User
	if err := s.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return &auth.RegisterResp{
			UserId:  0,
			Message: "用户名已存在",
		}, nil
	}

	if req.Role != model.RoleHR && req.Role != model.RoleCandidate {
		return &auth.RegisterResp{
			UserId:  0,
			Message: "无效的角色类型",
		}, nil
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	user := model.User{
		Username:     req.Username,
		PasswordHash: hashedPassword,
		Role:         req.Role,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	return &auth.RegisterResp{
		UserId:  user.ID,
		Message: "注册成功",
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req *auth.LoginReq) (*auth.LoginResp, error) {
	var user model.User
	if err := s.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &auth.LoginResp{
				Token:    "",
				UserId:   0,
				Role:     "",
				Username: "",
			}, nil
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	if !CheckPasswordHash(req.Password, user.PasswordHash) {
		return &auth.LoginResp{
			Token:    "",
			UserId:   0,
			Role:     "",
			Username: "",
		}, nil
	}

	token, err := middleware.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("生成Token失败: %w", err)
	}

	return &auth.LoginResp{
		Token:    token,
		UserId:   user.ID,
		Role:     user.Role,
		Username: user.Username,
	}, nil
}

func (s *AuthService) VerifyToken(ctx context.Context, req *auth.VerifyTokenReq) (*auth.VerifyTokenResp, error) {
	claims, err := middleware.ParseToken(req.Token)
	if err != nil {
		return &auth.VerifyTokenResp{
			Valid:  false,
			UserId: 0,
			Role:   "",
		}, nil
	}

	fmt.Printf("VerifyToken: %v", claims)

	return &auth.VerifyTokenResp{
		Valid:  true,
		UserId: claims.UserID,
		Role:   claims.Role,
	}, nil
}
