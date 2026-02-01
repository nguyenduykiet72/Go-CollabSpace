package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"Go-CollabSpace/internal/common/apperror"
	"Go-CollabSpace/internal/dto"
	"Go-CollabSpace/pkg/httpx"
)

type UserUseCase interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.UserResponse, error)
	Login(ctx context.Context, req dto.LoginRequest, userAgent string) (*dto.TokenResponse, error)
	GetAllUsers(ctx context.Context, req dto.PaginationReq) ([]*dto.UserResponse, error)
}

type UserController struct {
	userService UserUseCase
}

func NewUserController(userUseCase UserUseCase) *UserController {
	return &UserController{userService: userUseCase}
}

func (c *UserController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.Error(apperror.ErrBadRequest)
		return
	}

	resp, err := c.userService.Register(ctx.Request.Context(), req)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	httpx.WriteJSON(ctx, http.StatusCreated, resp, "User created successfully")
}

func (c *UserController) Login(ctx *gin.Context) {
	var req dto.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.Error(apperror.ErrBadRequest)
		return
	}

	userAgent := ctx.GetHeader("User-Agent")
	if userAgent == "" {
		userAgent = "unknown"
	}

	tokenResp, err := c.userService.Login(ctx.Request.Context(), req, userAgent)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	httpx.WriteJSON(ctx, http.StatusOK, tokenResp, "Login successful")
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
