package middleware

import (
	"Go-CollabSpace/pkg/httpx"

	"github.com/gin-gonic/gin"
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
