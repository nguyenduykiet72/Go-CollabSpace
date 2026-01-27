package middleware

import (
	"github.com/gin-gonic/gin"

	"Go-CollabSpace/pkg/httpx"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		httpx.WriteError(c, err)

		c.Abort()
	}
}
