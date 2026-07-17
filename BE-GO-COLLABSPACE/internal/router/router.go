package router

import (
	"Go-CollabSpace/internal/common/token"
	"Go-CollabSpace/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetUpRoutes(r *gin.Engine, ah AppHandlers, tokenProvider token.ITokenProvider, db *gorm.DB) {
	v1 := r.Group("/api/v1")

	publicRoutes := v1.Group("/")
	protectedRoutes := v1.Group("/")
	protectedRoutes.Use(middleware.AuthMiddleware(tokenProvider))

	registerAuthRoutes(publicRoutes, ah)
	registerUserRoutes(protectedRoutes, ah)
	registerWorkspaceRoutes(protectedRoutes, ah, db)
	registerDocumentRoutes(protectedRoutes, ah, db)
}
