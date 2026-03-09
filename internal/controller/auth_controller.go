package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"Go-CollabSpace/internal/common/apperror"
	"Go-CollabSpace/internal/dto"
	"Go-CollabSpace/pkg/httpx"
)

type AuthUseCase interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.UserResponse, error)
	Login(ctx context.Context, req dto.LoginRequest, userAgent string) (*dto.TokenResponse, error)
}

type AuthController struct {
	authService AuthUseCase
}

func NewAuthController(userUseCase AuthUseCase) *AuthController {
	return &AuthController{authService: userUseCase}
}

func (c *AuthController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.Error(apperror.ErrBadRequest)
		return
	}

	resp, err := c.authService.Register(ctx.Request.Context(), req)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	httpx.WriteJSON(ctx, http.StatusCreated, resp, "User created successfully")
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.Error(apperror.ErrBadRequest)
		return
	}

	userAgent := ctx.GetHeader("User-Agent")
	if userAgent == "" {
		userAgent = "unknown"
	}

	tokenResp, err := c.authService.Login(ctx.Request.Context(), req, userAgent)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	httpx.WriteJSON(ctx, http.StatusOK, tokenResp, "Login successful")
}
