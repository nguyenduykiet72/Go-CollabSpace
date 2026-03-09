package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func registerDocumentRoutes(rg *gin.RouterGroup, ah AppHandlers, db *gorm.DB) {
	doc := rg.Group("/document")
	doc.POST("", ah.DocumentController.CreateDoc)
	doc.GET("/doc/:workspaceId", ah.DocumentController.GetWorkspaceDocs)
	doc.GET("/:docId", ah.DocumentController.GetDocDetail)
}
