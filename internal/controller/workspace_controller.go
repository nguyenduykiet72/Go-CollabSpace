package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"Go-CollabSpace/internal/common/apperror"
	"Go-CollabSpace/internal/common/token"
	"Go-CollabSpace/internal/dto"
	"Go-CollabSpace/pkg/httpx"
)

type WorkspaceUseCase interface {
	CreateWorkspace(ctx context.Context, req dto.CreateWorkspaceRequest, userID uuid.UUID) (*dto.WorkSpaceResponse, error)
	GetWorkspaceByID(ctx context.Context, id uuid.UUID) (*dto.WorkSpaceResponse, error)
}

type WorkspaceController struct {
	workspaceService WorkspaceUseCase
}

func NewWorkspaceController(workspaceService WorkspaceUseCase) *WorkspaceController {
	return &WorkspaceController{workspaceService: workspaceService}
}

func (c *WorkspaceController) CreateWorkspace(ctx *gin.Context) {
	payload, exists := ctx.Get("authorization_payload")
	if !exists {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userClaims, ok := payload.(*token.UserClaims)
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var req dto.CreateWorkspaceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		_ = ctx.Error(apperror.ErrBadRequest)
		return
	}

	resp, err := c.workspaceService.CreateWorkspace(ctx.Request.Context(), req, userClaims.UserID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	httpx.WriteJSON(ctx, http.StatusCreated, resp, "Workspace created successfully")
}
