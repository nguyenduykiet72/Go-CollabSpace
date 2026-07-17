package router

import (
	"Go-CollabSpace/internal/constant"
	"Go-CollabSpace/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func registerWorkspaceRoutes(rg *gin.RouterGroup, ah AppHandlers, db *gorm.DB) {
	wp := rg.Group("/workspace")
	wp.POST("", ah.WorkspaceController.CreateWorkspace)

	wpScoped := wp.Group("/:workspaceId")

	viewerGroup := wpScoped.Group("/")
	viewerGroup.Use(middleware.RequireWorkSpaceRole(db, constant.RoleViewer))
	{
		viewerGroup.GET("", ah.WorkspaceController.GetWorkspaceByID)
	}

	adminGroup := wpScoped.Group("/")
	adminGroup.Use(middleware.RequireWorkSpaceRole(db, constant.RoleAdmin))
	{
		adminGroup.POST("/members", ah.WorkspaceController.AddMembers)
	}

	ownerGroup := wpScoped.Group("/")
	ownerGroup.Use(middleware.RequireWorkSpaceRole(db, constant.RoleOwner))
	{
		// TODO: Implement 2 functions below
		// ownerGroup.DELETE("", ah.WorkspaceController.DeleteWorkspace)
		// ownerGroup.PUT("", ah.WorkspaceController.UpdateWorkspace)
	}
}
