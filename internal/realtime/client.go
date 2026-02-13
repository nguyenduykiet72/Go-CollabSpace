package realtime

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1048576 // 1 MB
)

type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	UserID uuid.UUID
	DocID  uuid.UUID
	Send   chan []byte // Outgoing message channel
}

func (c *Client) WriteLoop() {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case msg, ok := <-c.Send:
			if !ok {
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.BinaryMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) ReadLoop() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		msgType, msg, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		if msgType == websocket.BinaryMessage {
			log.Printf("Client %s raw msg (%d bytes): first bytes = %v", c.UserID, len(msg), firstN(msg, 6))

			// Pure relay: forward every binary message to all other clients in the same doc room.
			// y-websocket clients handle the Yjs protocol (SyncStep1/2, Updates, Awareness) themselves.
			c.Hub.Broadcast <- &BroadcastMessage{
				DocID:   c.DocID,
				Payload: msg,
				Sender:  c,
			}
		}
	}
}

// firstN returns the first n bytes of b, or all of b if len(b) < n.
func firstN(b []byte, n int) []byte {
	if len(b) < n {
		return b
	}
	return b[:n]
}
