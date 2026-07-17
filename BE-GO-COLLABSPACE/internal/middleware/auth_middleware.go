package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"Go-CollabSpace/internal/common/apperror"
	"Go-CollabSpace/internal/common/token"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "Bearer"
)

func AuthMiddleware(tokenProvider token.ITokenProvider) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authHeader) == 0 {
			_ = ctx.Error(apperror.ErrUnauthorized)
			ctx.Abort()
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			_ = ctx.Error(apperror.ErrUnauthorized.WithMessage("invalid authorization header format"))
			ctx.Abort()
			return
		}

		authType := fields[0]
		if authType != authorizationTypeBearer {
			_ = ctx.Error(apperror.ErrUnauthorized.WithMessage("invalid authorization header type"))
			ctx.Abort()
			return
		}

		accessToken := fields[1]

		payload, err := tokenProvider.ValidateAccessToken(accessToken)
		if err != nil {
			_ = ctx.Error(apperror.ErrUnauthorized.WithMessage("invalid credentials"))
			ctx.Abort()
			return
		}

		ctx.Set(ContextKeyAuthPayload, payload)
		ctx.Set(ContextKeyUserID, payload.UserID)
		ctx.Next()
	}
}
