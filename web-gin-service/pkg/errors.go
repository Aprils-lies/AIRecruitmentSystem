package pkg

const (
	CodeSuccess       = 0
	CodeParamError    = 400
	CodeUnauthorized  = 401
	CodeForbidden     = 403
	CodeNotFound      = 404
	CodeBadRequest    = 400
	CodeInternalError = 500
	CodeServiceError  = 500
)

var CodeMessages = map[int]string{
	CodeSuccess:       "操作成功",
	CodeParamError:    "参数错误",
	CodeUnauthorized:  "未授权，请先登录",
	CodeForbidden:     "无权访问",
	CodeNotFound:      "资源不存在",
	CodeInternalError: "服务器内部错误",
}

func GetMessage(code int) string {
	if msg, ok := CodeMessages[code]; ok {
		return msg
	}
	return "未知错误"
}

type BusinessError struct {
	Code    int
	Message string
}

func (e *BusinessError) Error() string {
	return e.Message
}

func NewBusinessError(code int, message string) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: message,
	}
}

func IsBusinessError(err error) bool {
	_, ok := err.(*BusinessError)
	return ok
}

func GetBusinessCode(err error) int {
	if be, ok := err.(*BusinessError); ok {
		return be.Code
	}
	return CodeInternalError
}
