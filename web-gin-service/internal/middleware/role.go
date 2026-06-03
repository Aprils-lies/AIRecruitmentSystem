package middleware

import (
	"github.com/gin-gonic/gin"

	"web-gin-service/pkg"
)

const (
	RoleHR        = "hr"
	RoleCandidate = "candidate"
)

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := GetRole(c)
		if role == "" {
			pkg.Unauthorized(c, "未获取到用户角色")
			c.Abort()
			return
		}

		isAllowed := false
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			pkg.Forbidden(c, "无权访问该资源")
			c.Abort()
			return
		}

		c.Next()
	}
}

func HRMiddleware() gin.HandlerFunc {
	return RoleMiddleware(RoleHR)
}

func CandidateMiddleware() gin.HandlerFunc {
	return RoleMiddleware(RoleCandidate)
}

func AnyRoleMiddleware() gin.HandlerFunc {
	return RoleMiddleware(RoleHR, RoleCandidate)
}
