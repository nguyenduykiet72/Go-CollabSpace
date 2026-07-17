package router

import (
	"Go-CollabSpace/internal/constant"
	"Go-CollabSpace/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func registerDocumentRoutes(rg *gin.RouterGroup, ah AppHandlers, db *gorm.DB) {
	docScoped := rg.Group("/workspace/:workspaceId/document")

	viewerGroup := docScoped.Group("/")
	viewerGroup.Use(middleware.RequireWorkSpaceRole(db, constant.RoleViewer))
	{
		viewerGroup.GET("", ah.DocumentController.GetWorkspaceDocs)
		viewerGroup.GET("/:docId", ah.DocumentController.GetDocDetail)
		viewerGroup.GET("/tree", ah.DocumentController.GetDocTree)
	}

	editorGroup := docScoped.Group("/")
	editorGroup.Use(middleware.RequireWorkSpaceRole(db, constant.RoleEditor))
	{
		editorGroup.POST("", ah.DocumentController.CreateDoc)
		editorGroup.PUT("/:docId/move", ah.DocumentController.MoveDoc)
		editorGroup.PUT("/:docId/snapshot", ah.DocumentController.SaveDocSnapshot)
	}
}
