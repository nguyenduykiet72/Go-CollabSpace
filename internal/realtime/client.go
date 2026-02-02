package realtime

import (
	"context"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

type Client struct {
	Hub     *Hub
	Conn    *websocket.Conn
	UserID  uuid.UUID
	DocID   uuid.UUID
	Send    chan []byte // Outgoing message channel
	Context context.Context
	Cancel  context.CancelFunc
}

func (c *Client) WriteLoop() {
	defer func() {
		c.Conn.Close(websocket.StatusNormalClosure, "") // Close the connection when the write loop ends
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			if !ok {
				return // Channel closed, exit the write loop
			}
			ctx, cancel := context.WithTimeout(c.Context, 5*time.Second)
			err := c.Conn.Write(ctx, websocket.MessageBinary, msg)
			cancel()

			if err != nil {
				return // Error writing to connection, exit the write loop
			}
		case <-c.Context.Done():
			return // Context cancelled, exit the write loop
		}
	}
}

func (c *Client) ReadLoop() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close(websocket.StatusNormalClosure, "")
		c.Cancel()
	}()

	c.Conn.SetReadLimit(1048576) // 1 MB

	for {
		typ, msg, err := c.Conn.Read(c.Context)
		if err != nil {
			break
		}
		if typ == websocket.MessageBinary {
			c.Hub.Broadcast <- &BroadcastMessage{
				DocID:   c.DocID,
				Payload: msg,
				Sender:  c.UserID,
			}
		}
	}
}
