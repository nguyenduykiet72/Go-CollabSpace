package middleware

import (
	"Go-CollabSpace/internal/common/token"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func AuthMiddleware(tokenProvider token.ITokenProvider) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authHeader) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is missing"})
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is invalid"})
			return
		}

		authType := fields[0]
		if authType != authorizationTypeBearer {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization type is invalid"})
			return
		}

		accessToken := fields[1]

		payload, err := tokenProvider.ValidateAccessToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid access token"})
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
