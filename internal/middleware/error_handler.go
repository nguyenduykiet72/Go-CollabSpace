package middleware

import (
	"Go-CollabSpace/internal/common/apperror"
	"Go-CollabSpace/pkg/httpx"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			httpx.ErrorResponse(c, appErr.HTTPStatus, appErr.Code, appErr.Message)
			return
		}
		httpx.ErrorResponse(c, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), "Internal server error")
	}
}
