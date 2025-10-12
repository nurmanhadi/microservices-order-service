package middleware

import (
	"fmt"
	"net/http"
	"order-service/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ErrorHandling() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		if len(ctx.Errors) > 0 {
			err := ctx.Errors.Last().Err

			if validationErr, ok := err.(validator.ValidationErrors); ok {
				var values []string
				for _, fieldErr := range validationErr {
					value := fmt.Sprintf("field %s is %s %s", fieldErr.Field(), fieldErr.Tag(), fieldErr.Param())
					values = append(values, value)
				}
				arguments := strings.Join(values, ", ")
				ctx.JSON(http.StatusBadRequest, map[string]any{
					"error": arguments,
					"path":  ctx.Request.URL.Path,
				})
				return
			} else if responseStatusException, ok := err.(*response.ErrorResponse); ok {
				ctx.JSON(responseStatusException.Code, map[string]any{
					"error": responseStatusException.Message,
					"path":  ctx.Request.URL.Path,
				})
				return
			} else {
				ctx.JSON(http.StatusInternalServerError, map[string]any{
					"error": err.Error(),
					"path":  ctx.Request.URL.Path,
				})
				return
			}
		}
	}
}
