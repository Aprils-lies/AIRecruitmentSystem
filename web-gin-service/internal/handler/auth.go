package handler

import (
	"context"

	"github.com/gin-gonic/gin"

	"web-gin-service/internal/grpc"
	"web-gin-service/pkg"
	"web-gin-service/proto/gen/auth"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ParamError(c, "参数错误："+err.Error())
		return
	}

	if req.Role != "hr" && req.Role != "candidate" {
		pkg.ParamError(c, "角色必须是 hr 或 candidate")
		return
	}

	grpcReq := &auth.RegisterReq{
		Username: req.Username,
		Password: req.Password,
		Role:     req.Role,
	}

	resp, err := grpc.GetAuthClient().Register(context.Background(), grpcReq)
	if err != nil {
		pkg.InternalError(c, "注册失败："+err.Error())
		return
	}

	if resp.UserId == 0 {
		pkg.Error(c, pkg.CodeBadRequest, resp.Message)
		return
	}

	pkg.SuccessWithMessage(c, resp.Message, gin.H{
		"user_id": resp.UserId,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.ParamError(c, "参数错误："+err.Error())
		return
	}

	grpcReq := &auth.LoginReq{
		Username: req.Username,
		Password: req.Password,
	}

	resp, err := grpc.GetAuthClient().Login(context.Background(), grpcReq)
	if err != nil {
		pkg.InternalError(c, "登录失败："+err.Error())
		return
	}

	if resp.Token == "" {
		pkg.Error(c, pkg.CodeUnauthorized, "用户名或密码错误")
		return
	}

	pkg.Success(c, gin.H{
		"token":    resp.Token,
		"user_id":  resp.UserId,
		"role":     resp.Role,
		"username": resp.Username,
	})
}
