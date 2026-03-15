package controller

import (
	"context"
	"fmt"
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
	AddMembers(ctx context.Context, workspaceID uuid.UUID, req dto.AddMembersRequest, callerID uuid.UUID) (int, error)
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

func (c *WorkspaceController) GetWorkspaceByID(ctx *gin.Context) {
	var param dto.GetWorkspaceParams
	if err := ctx.ShouldBindUri(&param); err != nil {
		fmt.Println("Error binding URI:", err)
		_ = ctx.Error(apperror.ErrBadRequest)
		return
	}

	mapUUID, err := uuid.Parse(param.WorkspaceID)
	if err != nil {
		_ = ctx.Error(apperror.ErrBadRequest)
		return
	}

	resp, err := c.workspaceService.GetWorkspaceByID(ctx.Request.Context(), mapUUID)
	if err != nil {
		_ = ctx.Error(apperror.ErrNotFound)
		return
	}
	httpx.WriteJSON(ctx, http.StatusOK, resp, "Workspace found successfully")
}

func (c *WorkspaceController) AddMembers(ctx *gin.Context) {
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

	workspaceIDStr := ctx.Param("workspaceId")
	workspaceID, err := uuid.Parse(workspaceIDStr)
	if err != nil {
		_ = ctx.Error(apperror.ErrBadRequest)
		return
	}
	fmt.Print("Workspace ID:", workspaceID)
	var req dto.AddMembersRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fmt.Println("Error binding JSON:", err)
		_ = ctx.Error(apperror.ErrBadRequest)
		return
	}
	fmt.Print("something here......")
	added, err := c.workspaceService.AddMembers(ctx.Request.Context(), workspaceID, req, userClaims.UserID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	httpx.WriteJSON(ctx, http.StatusOK, gin.H{"added": added}, fmt.Sprintf("%d member(s) added successfully", added))
}
