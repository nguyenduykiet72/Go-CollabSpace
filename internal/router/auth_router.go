package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"

	"Go-CollabSpace/internal/middleware"
)

func registerAuthRoutes(rg *gin.RouterGroup, ah AppHandlers) {
	au := rg.Group("/auth")

	if ah.RateLimiter != nil {
		loginLimit := ah.RateLimiter.Limit("auth:login", middleware.LimitByIP, redis_rate.PerMinute(10))
		registerLimit := ah.RateLimiter.Limit("auth:register", middleware.LimitByIP, redis_rate.PerMinute(5))
		forgotLimit := ah.RateLimiter.Limit("auth:forgot", middleware.LimitByIP, redis_rate.PerHour(5))
		resetLimit := ah.RateLimiter.Limit("auth:reset", middleware.LimitByIP, redis_rate.PerHour(10))
		refreshLimit := ah.RateLimiter.Limit("auth:refresh", middleware.LimitByIP, redis_rate.PerMinute(30))
		oauthLimit := ah.RateLimiter.Limit("auth:oauth", middleware.LimitByIP, redis_rate.PerMinute(10))

		au.POST("/register", registerLimit, ah.AuthController.Register)
		au.POST("/login", loginLimit, ah.AuthController.Login)
		au.POST("/refresh", refreshLimit, ah.AuthController.Refresh)
		au.POST("/forgot-password", forgotLimit, ah.AuthController.ForgotPassword)
		au.POST("/reset-password", resetLimit, ah.AuthController.ResetPassword)
		au.POST("/oauth/:provider", oauthLimit, ah.AuthController.LoginWithSocial)
		return
	}

	au.POST("/register", ah.AuthController.Register)
	au.POST("/login", ah.AuthController.Login)
	au.POST("/refresh", ah.AuthController.Refresh)
	au.POST("/forgot-password", ah.AuthController.ForgotPassword)
	au.POST("/reset-password", ah.AuthController.ResetPassword)
	au.POST("/oauth/:provider", ah.AuthController.LoginWithSocial)
}
