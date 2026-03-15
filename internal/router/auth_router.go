package router

import (
	"github.com/gin-gonic/gin"
)

func registerAuthRoutes(rg *gin.RouterGroup, ah AppHandlers) {
	au := rg.Group("/auth")
	au.POST("/register", ah.AuthController.Register)
	au.POST("/login", ah.AuthController.Login)
	au.POST("/forgot-password", ah.AuthController.ForgotPassword)
	au.POST("/reset-password", ah.AuthController.ResetPassword)
}
