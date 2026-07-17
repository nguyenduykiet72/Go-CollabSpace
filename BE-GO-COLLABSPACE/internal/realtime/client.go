package realtime

import (
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"Go-CollabSpace/internal/constant"
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
	RoleID uint
	Send   chan []byte // Outgoing message channel

	closeOnce sync.Once
}

// CloseSend closes the Send channel exactly once. Safe to call from multiple
// goroutines (hub eviction, shutdown, etc.).
func (c *Client) CloseSend() {
	c.closeOnce.Do(func() {
		close(c.Send)
	})
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
		// Best-effort unregister; if the hub is shutting down its run loop will no
		// longer drain Unregister, so don't block on it.
		select {
		case c.Hub.Unregister <- c:
		case <-c.Hub.quit:
		}
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
		if msgType == websocket.BinaryMessage && len(msg) > 0 {
			yjsMsgType := msg[0]

			if c.RoleID >= constant.RoleViewer && yjsMsgType == MessageSync {
				log.Printf("SECURITY ALERT: Viewer %s attempted to send a Sync message to Doc %s. Packet dropped.", c.UserID, c.DocID)
				continue
			}

			// Persist only real Sync Update payloads. Awareness/SyncStep messages are relay-only.
			saveToDB := false
			if yjsMsgType == MessageSync && len(msg) >= 2 && msg[1] == SyncUpdate {
				saveToDB = true
			}

			// Pure relay: forward every binary message to all other clients in the same doc room.
			// y-websocket clients handle the Yjs protocol (SyncStep1/2, Updates, Awareness) themselves.
			c.Hub.Broadcast <- &BroadcastMessage{
				DocID:    c.DocID,
				Payload:  msg,
				SenderId: c.UserID,
				SaveToDB: saveToDB,
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
