package controller

import (
	"Go-CollabSpace/internal/common/apperror"
	"Go-CollabSpace/internal/common/token"
	"Go-CollabSpace/internal/dto"
	"Go-CollabSpace/pkg/httpx"
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DocumentUseCase interface {
	CreateDoc(ctx context.Context, req dto.CreateDocRequest, userID uuid.UUID) (*dto.DocumentResponse, error)
	GetWorkspaceDocs(ctx context.Context, workspaceID uuid.UUID, userID uuid.UUID) ([]dto.DocumentResponse, error)
	GetDocDetail(ctx context.Context, docID uuid.UUID, userID uuid.UUID) (*dto.DocumentResponse, error)
}

type DocumentController struct {
	documentService DocumentUseCase
}

func NewDocumentController(documentService DocumentUseCase) *DocumentController {
	return &DocumentController{documentService: documentService}
}

func (c *DocumentController) CreateDoc(ctx *gin.Context) {
	userID, err := c.getUserIDFromContext(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	workspaceIDStr := ctx.Param("workspace_id")
	workspaceID, err := uuid.Parse(workspaceIDStr)
	if err != nil {
		_ = ctx.Error(apperror.ErrBadRequest)
		return
	}

	var req dto.CreateDocRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.Error(apperror.ErrBadRequest)
		return
	}

	req.WorkspaceID = workspaceID

	resp, err := c.documentService.CreateDoc(ctx.Request.Context(), req, userID)
	if err != nil {
		_ = ctx.Error(apperror.ErrInternal)
		return
	}

	httpx.WriteJSON(ctx, http.StatusCreated, resp, "Document created successfully")
}

func (c *DocumentController) GetWorkspaceDocs(ctx *gin.Context) {
	userID, err := c.getUserIDFromContext(ctx)
	if err != nil {
		_ = ctx.Error(apperror.ErrUnauthorized)
		return
	}

	workspaceIDStr := ctx.Param("workspace_id")
	workspaceID, err := uuid.Parse(workspaceIDStr)
	if err != nil {
		_ = ctx.Error(apperror.ErrBadRequest)
		return
	}

	docs, err := c.documentService.GetWorkspaceDocs(ctx.Request.Context(), workspaceID, userID)
	if err != nil {
		_ = ctx.Error(apperror.ErrInternal)
		return
	}

	httpx.WriteJSON(ctx, http.StatusOK, docs, "Documents retrieved successfully")
}

func (c *DocumentController) GetDocDetail(ctx *gin.Context) {
	userID, err := c.getUserIDFromContext(ctx)
	if err != nil {
		_ = ctx.Error(apperror.ErrUnauthorized)
		return
	}

	docIDStr := ctx.Param("doc_id")
	docID, err := uuid.Parse(docIDStr)
	if err != nil {
		_ = ctx.Error(apperror.ErrBadRequest)
		return
	}

	doc, err := c.documentService.GetDocDetail(ctx.Request.Context(), docID, userID)
	if err != nil {
		_ = ctx.Error(apperror.ErrInternal)
		return
	}

	httpx.WriteJSON(ctx, http.StatusOK, doc, "Document detail retrieved successfully")
}

func (c *DocumentController) getUserIDFromContext(ctx *gin.Context) (uuid.UUID, error) {
	payload, exists := ctx.Get("authorization_payload")
	if !exists {
		return uuid.Nil, errors.New("unauthorized: missing token payload")
	}

	claims, ok := payload.(*token.UserClaims)
	if !ok {
		return uuid.Nil, errors.New("internal error: invalid token payload type")
	}

	return claims.UserID, nil
}
