package apperror

import "net/http"

var (
	ErrNotFound = &AppError{
		Code:       "NOT_FOUND",
		Message:    "Resource not found",
		HTTPStatus: http.StatusNotFound,
	}
	ErrUnauthorized = &AppError{
		HTTPStatus: http.StatusUnauthorized,
		Message:    "Unauthorized",
		Code:       "UNAUTHORIZED",
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
	ErrEmailAlreadyExists = &AppError{
		Code:       "EMAIL_ALREADY_EXISTS",
		Message:    "Email already exists",
		HTTPStatus: http.StatusBadRequest,
	}
	ErrSlugAlreadyExists = &AppError{
		Code:       "SLUG_ALREADY_EXISTS",
		Message:    "Workspace with this slug already exists",
		HTTPStatus: http.StatusConflict,
	}
)
