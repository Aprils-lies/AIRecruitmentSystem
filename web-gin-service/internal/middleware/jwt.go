package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"

	"web-gin-service/internal/grpc"
	"web-gin-service/pkg"
	"web-gin-service/proto/gen/auth"
)

const (
	ContextUserID = "user_id"
	ContextRole   = "role"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			pkg.Unauthorized(c, "未提供认证令牌")
			c.Abort()
			return
		}

		token := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}

		if token == "" {
			pkg.Unauthorized(c, "无效的认证令牌格式")
			c.Abort()
			return
		}

		req := &auth.VerifyTokenReq{
			Token: token,
		}

		resp, err := grpc.GetAuthClient().VerifyToken(context.Background(), req)
		if err != nil {
			pkg.InternalError(c, "验证令牌失败")
			c.Abort()
			return
		}

		if !resp.Valid {
			pkg.Unauthorized(c, "令牌无效或已过期")
			c.Abort()
			return
		}

		c.Set(ContextUserID, resp.UserId)
		c.Set(ContextRole, resp.Role)

		c.Next()
	}
}

func GetUserID(c *gin.Context) int64 {
	userID, exists := c.Get(ContextUserID)
	if !exists {
		return 0
	}
	return userID.(int64)
}

func GetRole(c *gin.Context) string {
	role, exists := c.Get(ContextRole)
	if !exists {
		return ""
	}
	return role.(string)
}
