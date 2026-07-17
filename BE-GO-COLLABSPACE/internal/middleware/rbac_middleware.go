package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"Go-CollabSpace/internal/model"
	"Go-CollabSpace/pkg/httpx"
)

func RequireWorkSpaceRole(db *gorm.DB, requiredRole uint) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, err := GetUserID(ctx)
		if err != nil {
			httpx.WriteJSON(ctx, http.StatusUnauthorized, nil, "Unauthorized")
			ctx.Abort()
			return
		}

		workspaceIDStr := ctx.Param("workspaceId")
		workspaceID, err := uuid.Parse(workspaceIDStr)
		if err != nil {
			httpx.WriteJSON(ctx, http.StatusBadRequest, nil, "Invalid workspace ID")
			ctx.Abort()
			return
		}

		var member model.WorkspaceMember
		err = db.WithContext(ctx.Request.Context()).
			Select("wpm_role_id").
			Where("wpm_workspace_id = ? AND wpm_user_id = ?", workspaceID, userID).
			First(&member).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				httpx.WriteJSON(ctx, http.StatusForbidden, nil, "Forbidden: You are not a member of this workspace")
				ctx.Abort()
				return
			}
			httpx.WriteJSON(ctx, http.StatusInternalServerError, nil, "Internal server error")
			ctx.Abort()
			return
		}

		if member.WpmRoleID > requiredRole {
			httpx.WriteJSON(ctx, http.StatusForbidden, nil, "Forbidden: Insufficient permissions")
			ctx.Abort()
			return
		}

		ctx.Set("workspace_role", member.WpmRoleID)
		ctx.Next()
	}
}
