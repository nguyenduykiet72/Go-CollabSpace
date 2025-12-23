package router

import (
	"Go-CollabSpace/internal/controller"

	"github.com/gin-gonic/gin"
)

func NewRouter(userController *controller.UserController) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1/user")
	{
		v1.POST("register", userController.Register)
	}
	return r
}
