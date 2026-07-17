package httpx

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"Go-CollabSpace/internal/common/apperror"
)

type Response struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	ErrorKey   string      `json:"errorKey,omitempty"`
}

func NewSuccess(data interface{}, msg string) Response {
	return Response{
		StatusCode: 200,
		Data:       data,
		Message:    msg,
	}
}

func WriteJSON(c *gin.Context, statusCode int, data interface{}, msg string) {
	c.JSON(statusCode, Response{
		StatusCode: statusCode,
		Message:    msg,
		Data:       data,
	})
}

func WriteError(c *gin.Context, err error) {
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.StatusCode, Response{
			StatusCode: appErr.StatusCode,
			Message:    appErr.Message,
			ErrorKey:   appErr.Key,
			Data:       nil,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, Response{
		StatusCode: http.StatusInternalServerError,
		Message:    "Internal server error",
		ErrorKey:   err.Error(),
		Data:       nil,
	})
}
