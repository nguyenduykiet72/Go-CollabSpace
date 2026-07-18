package dto

import "github.com/google/uuid"

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"fullName" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"fullName"`
	Avatar    string    `json:"avatar"`
	CreatedAt string    `json:"createdAt"`
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type ForgotPasswordRequest struct {
	Email     string `json:"email" binding:"required,email"`
	ReturnURL string `json:"returnUrl" binding:"required,url"`
}

type ResetPasswordRequest struct {
	RawToken    string `json:"token" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

type SocialLoginReq struct {
	Code string `json:"code" binding:"required"`
}
