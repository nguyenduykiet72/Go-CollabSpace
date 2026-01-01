package router

import (
	"Go-CollabSpace/internal/controller"

	"github.com/gin-gonic/gin"
)

type AppHandlers struct {
	UserController *controller.UserController
}

func SetUpRoutes(r *gin.Engine, ah AppHandlers) {
	v1 := r.Group("/api/v1")
	{
		userGroup := v1.Group("/users")
		{
			userGroup.POST("/register", ah.UserController.Register)
			userGroup.POST("/login", ah.UserController.Login)
		}
	}
}
