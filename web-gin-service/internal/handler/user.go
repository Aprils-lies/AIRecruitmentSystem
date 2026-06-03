package handler

import (
	"context"

	"github.com/gin-gonic/gin"

	"web-gin-service/internal/grpc"
	"web-gin-service/internal/middleware"
	"web-gin-service/pkg"
	"web-gin-service/proto/gen/user"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		pkg.Unauthorized(c, "用户未登录")
		return
	}

	grpcReq := &user.GetProfileReq{
		UserId: userID,
	}

	resp, err := grpc.GetUserClient().GetProfile(context.Background(), grpcReq)
	if err != nil {
		pkg.InternalError(c, "获取用户信息失败："+err.Error())
		return
	}

	pkg.Success(c, gin.H{
		"user_id":    resp.UserId,
		"username":   resp.Username,
		"role":       resp.Role,
		"real_name":  resp.RealName,
		"phone":      resp.Phone,
		"education":  resp.Education,
		"school":     resp.School,
		"experience": resp.Experience,
		"skills":     resp.Skills,
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if userID == 0 {
		pkg.Unauthorized(c, "用户未登录")
		return
	}

	var req struct {
		RealName  string `json:"real_name"`
		Phone     string `json:"phone"`
		Education string `json:"education"`
		School    string `json:"school"`
		Experience string `json:"experience"`
		Skills    string `json:"skills"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ParamError(c, "参数错误："+err.Error())
		return
	}

	grpcReq := &user.UpdateProfileReq{
		UserId:     userID,
		RealName:  req.RealName,
		Phone:     req.Phone,
		Education: req.Education,
		School:    req.School,
		Experience: req.Experience,
		Skills:    req.Skills,
	}

	resp, err := grpc.GetUserClient().UpdateProfile(context.Background(), grpcReq)
	if err != nil {
		pkg.InternalError(c, "更新用户信息失败："+err.Error())
		return
	}

	if !resp.Success {
		pkg.Error(c, pkg.CodeBadRequest, resp.Message)
		return
	}

	pkg.SuccessWithMessage(c, resp.Message, nil)
}
