package apperror

import "net/http"

var (
	ErrInvalidInput        = NewAppError(http.StatusBadRequest, "Invalid input", "INVALID_INPUT")
	ErrUnauthorized        = NewAppError(http.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED")
	ErrNotFound            = NewAppError(http.StatusNotFound, "Resource not found", "NOT_FOUND")
	ErrForbidden           = NewAppError(http.StatusForbidden, "Forbidden", "FORBIDDEN")
	ErrWorkspaceNotFound   = NewAppError(http.StatusNotFound, "Workspace not found", "WORKSPACE_NOT_FOUND")
	ErrRoleNotFound        = NewAppError(http.StatusNotFound, "Role not found", "ROLE_NOT_FOUND")
	ErrBadRequest          = NewAppError(http.StatusBadRequest, "Bad request", "BAD_REQUEST")
	ErrInternal            = NewAppError(http.StatusInternalServerError, "Internal server error", "INTERNAL_ERROR")
	ErrEmailExists         = NewAppError(http.StatusConflict, "Email already exists", "EMAIL_EXISTS")
	ErrUserBlocked         = NewAppError(http.StatusForbidden, "User is blocked", "USER_BLOCKED")
	ErrSlugExists          = NewAppError(http.StatusConflict, "Slug already exists", "SLUG_EXISTS")
	ErrMemberExists        = NewAppError(http.StatusConflict, "User is already a member of this workspace", "MEMBER_EXISTS")
	ErrUserNotFound        = NewAppError(http.StatusNotFound, "User not found", "USER_NOT_FOUND")
	ErrCircularReference   = NewAppError(http.StatusBadRequest, "Circular reference detected", "CIRCULAR_REFERENCE")
	ErrCannotMoveToSelf    = NewAppError(http.StatusBadRequest, "Cannot move a document into itself", "CANNOT_MOVE_TO_SELF")
	ErrUnsupportedFileType = NewAppError(http.StatusBadRequest, "Unsupported file type", "UNSUPPORTED_FILE_TYPE")
	ErrFileTooLarge        = NewAppError(http.StatusBadRequest, "File size exceeds the limit (5MB)", "FILE_TOO_LARGE")
	ErrInvalidToken        = NewAppError(http.StatusUnauthorized, "Invalid token", "INVALID_TOKEN")
)
