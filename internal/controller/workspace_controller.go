package controller

import (
	"Go-CollabSpace/internal/common/apperror"
	"Go-CollabSpace/internal/common/token"
	"Go-CollabSpace/internal/dto"
	"Go-CollabSpace/internal/service"
	"Go-CollabSpace/pkg/httpx"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WorkspaceController struct {
	service service.IWorkspaceService
}

func NewWorkspaceController(service service.IWorkspaceService) *WorkspaceController {
	return &WorkspaceController{service: service}
}

func (c *WorkspaceController) Create(ctx *gin.Context) {
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

	resp, err := c.service.Create(ctx.Request.Context(), req, userClaims.UserID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	httpx.Created(ctx, resp)
}
