package httpx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Code       string      `json:"code,omitempty"`
}

func Success(data interface{}, message string) Response {
	return Response{
		StatusCode: http.StatusOK,
		Data:       data,
		Message:    message,
	}
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		StatusCode: http.StatusCreated,
		Message:    "Resource created successfully",
		Data:       data,
	})
}

func ErrorResponse(c *gin.Context, httpStatus int, code string, msg string) {
	c.JSON(httpStatus, Response{
		StatusCode: httpStatus,
		Code:       code,
		Message:    msg,
		Data:       nil,
	})
}

//func Fail(statusCode int, message string) Response {
//	return Response{
//		StatusCode: statusCode,
//		Data:       nil,
//		Message:    message,
//	}
//}
