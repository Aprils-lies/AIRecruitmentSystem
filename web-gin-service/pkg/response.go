package pkg

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "操作成功",
		Data:    data,
	})
}

func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

func ErrorWithStatus(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

func ParamError(c *gin.Context, message string) {
	Error(c, CodeParamError, message)
}

func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "未授权，请先登录"
	}
	ErrorWithStatus(c, http.StatusUnauthorized, CodeUnauthorized, message)
}

func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = "无权访问"
	}
	ErrorWithStatus(c, http.StatusForbidden, CodeForbidden, message)
}

func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = "资源不存在"
	}
	Error(c, CodeNotFound, message)
}

func InternalError(c *gin.Context, message string) {
	if message == "" {
		message = "服务器内部错误"
	}
	Error(c, CodeInternalError, message)
}

func BadRequest(c *gin.Context, message string) {
	if message == "" {
		message = "请求参数错误"
	}
	Error(c, CodeBadRequest, message)
}

func SafeString(s string) string {
	if s == "" {
		return ""
	}
	return s
}
