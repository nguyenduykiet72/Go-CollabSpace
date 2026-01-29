package router

import (
	"Go-CollabSpace/internal/common/token"
	"Go-CollabSpace/internal/controller"
	"Go-CollabSpace/internal/middleware"

	"github.com/gin-gonic/gin"
)

type AppHandlers struct {
	UserController      *controller.UserController
	WorkspaceController *controller.WorkspaceController
}

func SetUpRoutes(r *gin.Engine, ah AppHandlers, tokenProvider token.ITokenProvider) {
	v1 := r.Group("/api/v1")
	{
		userGroup := v1.Group("/auth")
		{
			userGroup.POST("/register", ah.UserController.Register)
			userGroup.POST("/login", ah.UserController.Login)
		}

		protectedGroup := v1.Group("/")
		protectedGroup.Use(middleware.AuthMiddleware(tokenProvider))
		{
			wp := protectedGroup.Group("/workspace")
			{
				wp.POST("", ah.WorkspaceController.CreateWorkspace)
			}
		}
	}
}
