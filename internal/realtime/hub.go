package realtime

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

type DocRepoForHub interface {
	AppendYjsUpdate(ctx context.Context, docID uuid.UUID, update []byte) error
	GetYjsState(ctx context.Context, docID uuid.UUID) ([]byte, error)
}

type BroadcastMessage struct {
	DocID   uuid.UUID
	Payload []byte
	Sender  *Client
}

type Hub struct {
	Rooms      map[uuid.UUID]map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *BroadcastMessage

	saveQueue chan *BroadcastMessage
	DocRepo   DocRepoForHub
	Mutex     sync.RWMutex
}

func NewHub(docRepo DocRepoForHub) *Hub {
	return &Hub{
		Rooms:      make(map[uuid.UUID]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *BroadcastMessage),
		saveQueue:  make(chan *BroadcastMessage, 1000),
		DocRepo:    docRepo,
	}
}

func (h *Hub) runSaver() {
	for msg := range h.saveQueue {
		// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

		// err := h.DocRepo.AppendYjsUpdate(ctx, msg.DocID, msg.Payload)
		// if err != nil {
		// 	log.Printf("Error saving Yjs update for doc %s: %v", msg.DocID, err)
		// }
		// cancel()
		msgType, payload := ParseYjsMessage(msg.Payload)

		if msgType == MessageSync {
			syncType, data := ParseYjsSyncMessage(payload)

			if syncType == SyncUpdate {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

				err := h.DocRepo.AppendYjsUpdate(ctx, msg.DocID, data)
				if err != nil {
					log.Printf("Error saving update: %v", err)
				} else {
					log.Printf("Saved update for doc %s (%d bytes)", msg.DocID, len(data))
				}
				cancel()
			}
		}
	}
}

func (h *Hub) Run() {
	go h.runSaver()

	log.Println("Hub is running...")

	for {
		select {
		case client := <-h.Register:
			h.Mutex.Lock()
			if h.Rooms[client.DocID] == nil {
				h.Rooms[client.DocID] = make(map[*Client]bool)
			}
			h.Rooms[client.DocID][client] = true
			h.Mutex.Unlock()
			log.Printf("Client %s joined doc %s", client.UserID, client.DocID)

		case message := <-h.Broadcast:
			// Save sync updates to DB for persistence
			select {
			case h.saveQueue <- message:
			default:
				log.Println("Save queue is full, dropping message")
			}

			h.Mutex.RLock()
			clients := h.Rooms[message.DocID]
			count := 0

			for client := range clients {
				if client != message.Sender {
					select {
					case client.Send <- message.Payload:
						count++
					default:
						close(client.Send)
						delete(clients, client)
					}
				}
			}
			h.Mutex.RUnlock()
			if count > 0 {
				log.Printf("Relayed message to %d clients for doc %s", count, message.DocID)
			}

		case client := <-h.Unregister:
			h.Mutex.Lock()
			if clients, ok := h.Rooms[client.DocID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.Send)
					if len(clients) == 0 {
						delete(h.Rooms, client.DocID)
					}
				}
			}
			h.Mutex.Unlock()
		}
	}
}
