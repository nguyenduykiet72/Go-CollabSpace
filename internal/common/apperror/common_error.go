package apperror

import "net/http"

var (
	ErrNotFound = &AppError{
		Code:       "NOT_FOUND",
		Message:    "Resource not found",
		HTTPStatus: http.StatusNotFound,
	}
	ErrUnauthorized = &AppError{
		Code:       "UNAUTHORIZED",
		Message:    "Unauthorized",
		HTTPStatus: http.StatusUnauthorized,
	}
	ErrBadRequest = &AppError{
		Code:       "BAD_REQUEST",
		Message:    "Bad request",
		HTTPStatus: http.StatusBadRequest,
	}
	ErrInternal = &AppError{
		Code:       "INTERNAL_ERROR",
		Message:    "Internal server error",
		HTTPStatus: http.StatusInternalServerError,
	}
)
