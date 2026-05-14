package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"Go-CollabSpace/internal/common/token"
)

const (
	ContextKeyAuthPayload = "authorization_payload"
	ContextKeyUserID      = "user_id"
)

var ErrNoAuthPayload = errors.New("no authorization payload in context")

func GetUserID(ctx *gin.Context) (uuid.UUID, error) {
	if v, ok := ctx.Get(ContextKeyUserID); ok {
		if id, ok := v.(uuid.UUID); ok {
			return id, nil
		}
	}

	if v, ok := ctx.Get(ContextKeyAuthPayload); ok {
		if claims, ok := v.(*token.UserClaims); ok {
			return claims.UserID, nil
		}
	}

	return uuid.Nil, ErrNoAuthPayload
}

func GetUserClaims(ctx *gin.Context) (*token.UserClaims, error) {
	v, ok := ctx.Get(ContextKeyAuthPayload)
	if !ok {
		return nil, ErrNoAuthPayload
	}
	claims, ok := v.(*token.UserClaims)
	if !ok {
		return nil, ErrNoAuthPayload
	}
	return claims, nil
}
