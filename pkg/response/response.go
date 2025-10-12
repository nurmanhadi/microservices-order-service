package response

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Success[T any](ctx *gin.Context, code int, data T) {
	resp := map[string]any{
		"data": data,
		"path": ctx.Request.URL.Path,
	}
	ctx.JSON(code, resp)
}

type ErrorResponse struct {
	Code    int
	Message string
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("code: %d, error: %s", e.Code, e.Message)
}
func Except(code int, message string) error {
	return &ErrorResponse{
		Code:    code,
		Message: message,
	}
}
