package controller

import (
	"Go-CollabSpace/internal/common/token"
	"Go-CollabSpace/internal/realtime"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// upgrader hijacks the connection BEFORE writing the 101 response,
// completely bypassing Gin's ResponseWriter. This avoids the
// "gin: response already written" error from coder/websocket.
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
}

func NewWsController(hub *realtime.Hub, tokenProvider token.ITokenProvider) *WsController {
	return &WsController{
		Hub:           hub,
		TokenProvider: tokenProvider,
	}
}

func (c *WsController) HandleWS(ctx *gin.Context) {
	// Custom panic recovery — this route has NO Gin Recovery middleware
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
