package realtime

import (
	"Go-CollabSpace/internal/repository"
	"log"
	"sync"

	"github.com/google/uuid"
)

type BroadcastMessage struct {
	DocID   uuid.UUID
	Payload []byte
	Sender  uuid.UUID
}

type Hub struct {
	Rooms      map[uuid.UUID]map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *BroadcastMessage

	DocRepo repository.DocumentRepository
	Mutex   sync.RWMutex
}

func NewHub(docRepo repository.DocumentRepository) *Hub {
	return &Hub{
		Rooms:      make(map[uuid.UUID]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *BroadcastMessage),
		DocRepo:    docRepo,
	}
}

func (h *Hub) Run() {
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

		case message := <-h.Broadcast:
			// TODO: Implement save logic
			h.Mutex.RLock()
			clients := h.Rooms[message.DocID]
			for client := range clients {
				if client.UserID != message.Sender {
					select {
					case client.Send <- message.Payload:
					default:
						close(client.Send)
						delete(clients, client)
					}
				}
			}
			h.Mutex.RUnlock()
		}
	}
}
