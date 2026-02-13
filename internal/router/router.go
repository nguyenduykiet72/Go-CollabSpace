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
	DocumentController  *controller.DocumentController
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
			user := protectedGroup.Group("/users")
			{
				user.GET("", ah.UserController.GetAllUsers)
			}
			wp := protectedGroup.Group("/workspace")
			{
				wp.POST("", ah.WorkspaceController.CreateWorkspace)
				wp.GET("/:workspaceId", ah.WorkspaceController.GetWorkspaceByID)
				wp.POST("/:workspaceId/members", ah.WorkspaceController.AddMembers)
			}
			doc := protectedGroup.Group("/document")
			{
				doc.POST("", ah.DocumentController.CreateDoc)
				doc.GET("/doc/:workspaceId", ah.DocumentController.GetWorkspaceDocs)
				doc.GET("/:docId", ah.DocumentController.GetDocDetail)
			}
		}
	}
}
