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

// UserClaims is the JWT body. Workspace-level roles are NOT carried in the
// token: authorisation is resolved at the middleware layer by joining
// tbl_workspace_members on every protected request. The token only proves
// "this user is authenticated".
type UserClaims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

type ITokenProvider interface {
	GenerateAccessToken(userID uuid.UUID) (string, error)
	GenerateRefreshToken() (string, error)
	ValidateAccessToken(tokenString string) (*UserClaims, error)
}
