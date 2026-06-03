package service

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"logic-grpc-service/internal/model"
	"logic-grpc-service/proto/gen/user"
)

type UserService struct {
	user.UnimplementedUserServiceServer
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetProfile(ctx context.Context, req *user.GetProfileReq) (*user.GetProfileResp, error) {
	var userModel model.User
	if err := s.db.Where("id = ?", req.UserId).First(&userModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &user.GetProfileResp{}, nil
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	return &user.GetProfileResp{
		UserId:    userModel.ID,
		Username:  userModel.Username,
		Role:      userModel.Role,
		RealName:  getStringPtr(userModel.RealName),
		Phone:     getStringPtr(userModel.Phone),
		Education: getStringPtr(userModel.Education),
		School:    getStringPtr(userModel.School),
		Experience: getStringPtr(userModel.Experience),
		Skills:    getStringPtr(userModel.Skills),
	}, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, req *user.UpdateProfileReq) (*user.UpdateProfileResp, error) {
	updateData := map[string]interface{}{}
	if req.RealName != "" {
		updateData["real_name"] = req.RealName
	}
	if req.Phone != "" {
		updateData["phone"] = req.Phone
	}
	if req.Education != "" {
		updateData["education"] = req.Education
	}
	if req.School != "" {
		updateData["school"] = req.School
	}
	if req.Experience != "" {
		updateData["experience"] = req.Experience
	}
	if req.Skills != "" {
		updateData["skills"] = req.Skills
	}

	if len(updateData) == 0 {
		return &user.UpdateProfileResp{
			Success: false,
			Message: "没有需要更新的字段",
		}, nil
	}

	if err := s.db.Model(&model.User{}).Where("id = ?", req.UserId).Updates(updateData).Error; err != nil {
		return nil, fmt.Errorf("更新用户信息失败: %w", err)
	}

	return &user.UpdateProfileResp{
		Success: true,
		Message: "更新成功",
	}, nil
}

func getStringPtr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}