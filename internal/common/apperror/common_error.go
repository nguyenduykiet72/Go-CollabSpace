package apperror

import "net/http"

var (
	ErrInvalidInput      = NewAppError(http.StatusBadRequest, "Invalid input", "INVALID_INPUT")
	ErrUnauthorized      = NewAppError(http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
	ErrNotFound          = NewAppError(http.StatusNotFound, "Resource not found", "NOT_FOUND")
	ErrWorkspaceNotFound = NewAppError(http.StatusNotFound, "Workspace not found", "WORKSPACE_NOT_FOUND")
	ErrRoleNotFound      = NewAppError(http.StatusNotFound, "Role not found", "ROLE_NOT_FOUND")
	ErrBadRequest        = NewAppError(http.StatusBadRequest, "Bad request", "BAD_REQUEST")
	ErrInternal          = NewAppError(http.StatusInternalServerError, "Internal server error", "INTERNAL_ERROR")
	ErrEmailExists       = NewAppError(http.StatusConflict, "Email already exists", "EMAIL_EXISTS")
	ErrUserBlocked       = NewAppError(http.StatusForbidden, "User is blocked", "USER_BLOCKED")
	ErrSlugExists        = NewAppError(http.StatusConflict, "Slug already exists", "SLUG_EXISTS")
)
