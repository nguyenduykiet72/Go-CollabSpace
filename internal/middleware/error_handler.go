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
			c.JSON(appErr.HTTPStatus, httpx.Fail(appErr.HTTPStatus, appErr.Message))
			return
		}

		c.JSON(http.StatusInternalServerError, httpx.Fail(http.StatusInternalServerError, "Internal Server Error"))
	}
}
