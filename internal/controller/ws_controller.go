package controller

import (
	"Go-CollabSpace/internal/common/token"
	"Go-CollabSpace/internal/realtime"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (tighten in production)
	},
}

type WsController struct {
	Hub           *realtime.Hub
	TokenProvider token.ITokenProvider
	DB            *gorm.DB
}

func NewWsController(hub *realtime.Hub, tokenProvider token.ITokenProvider, db *gorm.DB) *WsController {
	return &WsController{
		Hub:           hub,
		TokenProvider: tokenProvider,
		DB:            db,
	}
}

func (c *WsController) HandleWS(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic in WebSocket handler: %v", r)
		}
	}()

	if c.Hub == nil {
		log.Println("CRITICAL ERROR: c.Hub is NIL inside HandleWS!")
		ctx.JSON(500, gin.H{"error": "Internal Server Error: Hub missing"})
		return
	}

	// --- Validate BEFORE upgrading to WebSocket ---
	tokenStr := ctx.Query("token")
	if tokenStr == "" {
		ctx.JSON(401, gin.H{"error": "Unauthorized: missing token"})
		return
	}

	claims, err := c.TokenProvider.ValidateAccessToken(tokenStr)
	if err != nil {
		ctx.JSON(401, gin.H{"error": "Unauthorized: invalid token"})
		return
	}

	docIDStr := ctx.Query("doc_id")
	docID, err := uuid.Parse(docIDStr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "Bad Request: invalid document ID"})
		return
	}

	var roleID uint
	result := c.DB.Table("tbl_workspace_members").
		Select("tbl_workspace_members.wpm_role_id").
		Joins("JOIN tbl_documents ON tbl_documents.doc_workspace_id = tbl_workspace_members.wpm_workspace_id").
		Where("tbl_documents.doc_id = ? AND tbl_workspace_members.wpm_user_id", docID, claims.UserID).
		Take(&roleID)

	if result.Error != nil {
		log.Printf("Error in WebSocket handler: %v", result.Error)
		ctx.JSON(403, gin.H{"error": "Access Denied: You do not have permission to access this document"})
		return
	}

	// --- Upgrade to WebSocket ---
	// gorilla/websocket hijacks BEFORE writing headers, bypassing Gin entirely.
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("WebSocket Upgrade failed: %v", err)
		return
	}

	log.Printf("WebSocket connected: user=%s doc=%s", claims.UserID, docID)

	client := &realtime.Client{
		Hub:    c.Hub,
		Conn:   conn,
		UserID: claims.UserID,
		DocID:  docID,
		Send:   make(chan []byte, 256),
	}

	client.Hub.Register <- client

	go client.WriteLoop()
	client.ReadLoop()
}
