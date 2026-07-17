package controller

import (
	"Go-CollabSpace/internal/common/apperror"
	"Go-CollabSpace/internal/dto"
	"Go-CollabSpace/pkg/httpx"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StorageUseCase interface {
	GeneratePresignedURL(ctx context.Context, userID uuid.UUID, req dto.FileUploadRequest) (*dto.FileURLResponse, error)
}

type StorageController struct {
	storageService StorageUseCase
}

func NewStorageController(storageService StorageUseCase) *StorageController {
	return &StorageController{storageService: storageService}
}

func (c *StorageController) GetUploadURL(ctx *gin.Context) {
	userID := ctx.MustGet("userId").(uuid.UUID)

	var req dto.FileUploadRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.Error(apperror.ErrBadRequest)
		return
	}

	//result, err := c.storageService.GeneratePresignedURL(ctx.Request.Context(), userID, req)
	result, err := c.storageService.GeneratePresignedURL(ctx.Request.Context(), userID, req)
	if err != nil {
		_ = ctx.Error(apperror.ErrBadRequest)
		return
	}

	httpx.WriteJSON(ctx, http.StatusOK, result, "Presigned URL generated successfully")
}
