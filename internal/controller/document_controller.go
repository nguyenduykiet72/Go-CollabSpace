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
	GetDocTree(ctx context.Context, workspaceID uuid.UUID, userID uuid.UUID) ([]dto.DocTreeItem, error)
	MoveDoc(ctx context.Context, docID uuid.UUID, req dto.MoveDocRequest, userID uuid.UUID) error
	SaveDocSnapshot(ctx context.Context, docID uuid.UUID, req dto.SaveSnapshotRequest) error
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

	var req dto.CreateDocRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.Error(apperror.ErrBadRequest)
		return
	}

	// req.WorkspaceID = workspaceID

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

	workspaceIDStr := ctx.Param("workspaceId")
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

	docIDStr := ctx.Param("docId")
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

func (c *DocumentController) GetDocTree(ctx *gin.Context) {
	workspaceID, err := uuid.Parse(ctx.Param("workspaceId"))
	if err != nil {
		_ = ctx.Error(apperror.ErrBadRequest)
		return
	}

	userID, err := c.getUserIDFromContext(ctx)
	if err != nil {
		_ = ctx.Error(apperror.ErrUnauthorized)
		return
	}

	tree, err := c.documentService.GetDocTree(ctx.Request.Context(), workspaceID, userID)
	if err != nil {
		_ = ctx.Error(apperror.ErrBadRequest)
		return
	}

	httpx.WriteJSON(ctx, http.StatusOK, tree, "Document tree retrieved successfully")
}

func (c *DocumentController) MoveDoc(ctx *gin.Context) {
	docID, err := uuid.Parse(ctx.Param("docId"))
	if err != nil {
		_ = ctx.Error(apperror.ErrBadRequest)
		return
	}

	var req dto.MoveDocRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.Error(apperror.ErrBadRequest)
		return
	}

	userID, err := c.getUserIDFromContext(ctx)
	if err != nil {
		_ = ctx.Error(apperror.ErrUnauthorized)
		return
	}

	if err := c.documentService.MoveDoc(ctx.Request.Context(), docID, req, userID); err != nil {
		_ = ctx.Error(apperror.ErrBadRequest)
		return
	}

	httpx.WriteJSON(ctx, http.StatusOK, nil, "Document moved successfully")
}

func (c *DocumentController) SaveDocSnapshot(ctx *gin.Context) {
	docIDStr := ctx.Param("docId")
	docID, err := uuid.Parse(docIDStr)
	if err != nil {
		_ = ctx.Error(apperror.ErrBadRequest)
		return
	}

	var req dto.SaveSnapshotRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.Error(apperror.ErrBadRequest)
		return
	}

	err = c.documentService.SaveDocSnapshot(ctx.Request.Context(), docID, req)
	if err != nil {
		_ = ctx.Error(apperror.ErrInternal)
		return
	}

	httpx.WriteJSON(ctx, http.StatusOK, nil, "Document snapshot saved successfully")
}
