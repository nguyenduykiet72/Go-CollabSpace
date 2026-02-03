package controller

import (
	"Go-CollabSpace/internal/common/token"
	"Go-CollabSpace/internal/realtime"
	"Go-CollabSpace/pkg/httpx"
	"context"

	"github.com/coder/websocket"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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
	tokenStr := ctx.Query("token")

	if tokenStr == "" {
		httpx.WriteJSON(ctx, 401, nil, "Unauthorized: missing token")
		return
	}
	claims, err := c.TokenProvider.ValidateAccessToken(tokenStr)
	if err != nil {
		httpx.WriteJSON(ctx, 401, nil, "Unauthorized: invalid token")
		return
	}

	docIDStr := ctx.Query("doc_id")
	docID, err := uuid.Parse(docIDStr)
	if err != nil {
		httpx.WriteJSON(ctx, 400, nil, "Bad Request: invalid document ID")
		return
	}

	conn, err := websocket.Accept(ctx.Writer, ctx.Request, &websocket.AcceptOptions{
		OriginPatterns: []string{"*"},
	})

	if err != nil {
		return
	}

	clientCtx, cancel := context.WithCancel(context.Background())

	client := &realtime.Client{
		Hub:     c.Hub,
		Conn:    conn,
		UserID:  claims.UserID,
		DocID:   docID,
		Send:    make(chan []byte, 256),
		Context: clientCtx,
		Cancel:  cancel,
	}

	client.Hub.Register <- client

	go client.WriteLoop()
	go client.ReadLoop()

}
