package controller

import (
	"Go-CollabSpace/internal/common/apperror"
	"Go-CollabSpace/internal/dto"
	"Go-CollabSpace/pkg/httpx"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserUseCase interface {
	GetAllUsers(ctx context.Context, req dto.PaginationReq) ([]*dto.UserResponse, error)
}

type UserController struct {
	userService UserUseCase
}

func NewUserController(userUseCase UserUseCase) *UserController {
	return &UserController{userService: userUseCase}
}

func (c *UserController) GetAllUsers(ctx *gin.Context) {
	var req dto.PaginationReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		_ = ctx.Error(apperror.ErrBadRequest)
		return
	}

	users, err := c.userService.GetAllUsers(ctx.Request.Context(), req)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	httpx.WriteJSON(ctx, http.StatusOK, users, "Users retrieved successfully")
}
