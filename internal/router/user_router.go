package router

import "github.com/gin-gonic/gin"

func registerUserRoutes(rg *gin.RouterGroup, ah AppHandlers) {
	ur := rg.Group("/user")
	ur.GET("/all", ah.UserController.GetAllUsers)
}
