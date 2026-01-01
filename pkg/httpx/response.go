package httpx

import "net/http"

type Response struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
}

func Success(data interface{}, message string) Response {
	return Response{
		StatusCode: http.StatusOK,
		Data:       data,
		Message:    message,
	}
}

func Fail(statusCode int, message string) Response {
	return Response{
		StatusCode: statusCode,
		Data:       nil,
		Message:    message,
	}
}
