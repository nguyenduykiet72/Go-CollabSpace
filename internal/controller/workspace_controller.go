package controller

import (
	"Go-CollabSpace/internal/common/token"
	"Go-CollabSpace/internal/dto"
	"Go-CollabSpace/internal/service"
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
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userClaims, ok := payload.(*token.UserClaims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token payload"})
		return
	}

	var req dto.CreateWorkspaceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	resp, err := c.service.Create(ctx.Request.Context(), req, userClaims.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Workspace created successfully", "data": resp})
}
