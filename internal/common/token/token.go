package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type ConfigToken struct {
	AccessTokenSecret    string
	RefreshTokenSecret   string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

type UserClaims struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

type ITokenProvider interface {
	GenerateAccessToken(userID uuid.UUID, role string) (string, error)
	GenerateRefreshToken() (string, error)
	ValidateAccessToken(tokenString string) (*UserClaims, error)
}
