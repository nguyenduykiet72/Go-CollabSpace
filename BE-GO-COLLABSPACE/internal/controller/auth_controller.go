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
	Refresh(ctx context.Context, refreshToken string, userAgent string) (*dto.TokenResponse, error)
	ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error
	LoginWithSocial(ctx context.Context, providerName string, code string, userAgent string) (*dto.TokenResponse, error)
}

type AuthController struct {
	authService AuthUseCase
}

func NewAuthController(userUseCase AuthUseCase) *AuthController {
	return &AuthController{authService: userUseCase}
}

func (a *AuthController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.Error(apperror.ErrBadRequest)
		return
	}

	resp, err := a.authService.Register(ctx.Request.Context(), req)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	httpx.WriteJSON(ctx, http.StatusCreated, resp, "User created successfully")
}

func (a *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.Error(apperror.ErrBadRequest)
		return
	}

	userAgent := ctx.GetHeader("User-Agent")
	if userAgent == "" {
		userAgent = "unknown"
	}

	tokenResp, err := a.authService.Login(ctx.Request.Context(), req, userAgent)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	httpx.WriteJSON(ctx, http.StatusOK, tokenResp, "Login successful")
}

func (a *AuthController) Refresh(ctx *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.Error(apperror.ErrBadRequest)
		return
	}

	userAgent := ctx.GetHeader("User-Agent")
	if userAgent == "" {
		userAgent = "unknown"
	}

	tokenResp, err := a.authService.Refresh(ctx.Request.Context(), req.RefreshToken, userAgent)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	httpx.WriteJSON(ctx, http.StatusOK, tokenResp, "Token refreshed successfully")
}

func (a *AuthController) ForgotPassword(ctx *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.Error(apperror.ErrBadRequest)
		return
	}

	err := a.authService.ForgotPassword(ctx.Request.Context(), req)
	if err != nil {
		_ = ctx.Error(apperror.ErrInternal)
		return
	}

	httpx.WriteJSON(ctx, http.StatusOK, nil, "If your email is registered, you will receive a password reset link shortly.")
}

func (a *AuthController) ResetPassword(ctx *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.Error(apperror.ErrBadRequest)
		return
	}

	err := a.authService.ResetPassword(ctx.Request.Context(), req)
	if err != nil {
		_ = ctx.Error(apperror.ErrInternal)
		return
	}

	httpx.WriteJSON(ctx, http.StatusOK, nil, "Password has been reset successfully.")
}

func (a *AuthController) LoginWithSocial(ctx *gin.Context) {
	providerName := ctx.Param("provider")

	var req dto.SocialLoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		_ = ctx.Error(apperror.ErrBadRequest)
		return
	}

	result, err := a.authService.LoginWithSocial(ctx.Request.Context(), providerName, req.Code, ctx.GetHeader("User-Agent"))
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	httpx.WriteJSON(ctx, http.StatusOK, result, "Login successful")
}
