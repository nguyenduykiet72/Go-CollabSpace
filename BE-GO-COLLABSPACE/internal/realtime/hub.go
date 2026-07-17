package realtime

import (
	"Go-CollabSpace/internal/common/infrastructure"
	"Go-CollabSpace/internal/telemetry"
	"context"
	"encoding/json"
	"fmt"
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
	// Sender   *Client
	SenderId uuid.UUID `json:"sender_id"`
	SaveToDB bool      `json:"-"`
}

type Hub struct {
	Rooms      map[uuid.UUID]map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *BroadcastMessage

	saveQueue   chan *BroadcastMessage
	DocRepo     DocRepoForHub
	Mutex       sync.RWMutex
	RedisClient *infrastructure.RedisClient

	quit     chan struct{}
	saverWg  sync.WaitGroup
	stopOnce sync.Once
}

func NewHub(docRepo DocRepoForHub, redisClient *infrastructure.RedisClient) *Hub {
	return &Hub{
		Rooms:       make(map[uuid.UUID]map[*Client]bool),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Broadcast:   make(chan *BroadcastMessage),
		saveQueue:   make(chan *BroadcastMessage, 1000),
		DocRepo:     docRepo,
		RedisClient: redisClient,
		quit:        make(chan struct{}),
	}
}

func (h *Hub) runSaver() {
	defer h.saverWg.Done()
	for msg := range h.saveQueue {
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

func (h *Hub) subscribeToRedisChannel(docID uuid.UUID) {
	ctx := context.Background()
	channelName := fmt.Sprintf("ws:doc:%s", docID.String())

	pubsub := h.RedisClient.Client.Subscribe(ctx, channelName)

	go func() {
		defer pubsub.Close()
		ch := pubsub.Channel()

		for msg := range ch {
			var boardcastMsg BroadcastMessage
			if err := json.Unmarshal([]byte(msg.Payload), &boardcastMsg); err != nil {
				log.Printf("Error unmarshaling Redis message: %v", err)
				continue
			}

			h.broadcastToLocalClients(&boardcastMsg)
		}
	}()
}

func (h *Hub) broadcastToLocalClients(msg *BroadcastMessage) {
	// Phase 1: under RLock, snapshot recipients and collect slow clients.
	h.Mutex.RLock()
	clients := h.Rooms[msg.DocID]
	var dead []*Client
	for client := range clients {
		if client.UserID == msg.SenderId {
			continue
		}
		select {
		case client.Send <- msg.Payload:
		default:
			dead = append(dead, client)
		}
	}
	h.Mutex.RUnlock()

	if len(dead) == 0 {
		return
	}

	// Phase 2: upgrade to write lock to evict slow clients.
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	room, ok := h.Rooms[msg.DocID]
	if !ok {
		return
	}
	for _, client := range dead {
		if _, stillThere := room[client]; stillThere {
			delete(room, client)
			client.CloseSend()
		}
	}
	if len(room) == 0 {
		delete(h.Rooms, msg.DocID)
	}
}

func (h *Hub) Run() {
	h.saverWg.Add(1)
	go h.runSaver()

	log.Println("Hub is running...")

	for {
		select {
		case <-h.quit:
			h.shutdown()
			return

		case client := <-h.Register:
			h.Mutex.Lock()
			if h.Rooms[client.DocID] == nil {
				h.Rooms[client.DocID] = make(map[*Client]bool)
				h.subscribeToRedisChannel(client.DocID)
			}
			h.Rooms[client.DocID][client] = true
			h.Mutex.Unlock()
			log.Printf("Client %s joined doc %s", client.UserID, client.DocID)
			telemetry.ActiveConnections.WithLabelValues(client.DocID.String()).Inc()

		case message := <-h.Broadcast:
			if message.SaveToDB {
				select {
				case h.saveQueue <- message:
				default:
					log.Println("Save queue full, dropping save request")
				}
			}

			payload, err := json.Marshal(message)
			if err != nil {
				log.Printf("Error marshaling broadcast message: %v", err)
				continue
			}

			channelName := fmt.Sprintf("ws:doc:%s", message.DocID.String())

			err = h.RedisClient.Client.Publish(context.Background(), channelName, payload).Err()
			if err != nil {
				log.Printf("Error publishing to Redis channel: %v", err)
			}
			telemetry.ProcessedMessages.WithLabelValues("update").Inc()

		case client := <-h.Unregister:
			h.Mutex.Lock()
			if clients, ok := h.Rooms[client.DocID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					client.CloseSend()
					if len(clients) == 0 {
						delete(h.Rooms, client.DocID)
					}
				}
			}
			h.Mutex.Unlock()
			telemetry.ActiveConnections.WithLabelValues(client.DocID.String()).Dec()
		}
	}
}

// Stop signals the hub to stop. Safe to call multiple times. It blocks until
// the run loop exits and the save queue has been drained.
func (h *Hub) Stop() {
	h.stopOnce.Do(func() {
		close(h.quit)
	})
	h.saverWg.Wait()
}

// shutdown is invoked by the run loop after it receives the quit signal. It
// closes all live client sends and lets the saver drain the in-flight queue.
func (h *Hub) shutdown() {
	log.Println("Hub shutting down...")

	h.Mutex.Lock()
	for docID, clients := range h.Rooms {
		for client := range clients {
			client.CloseSend()
		}
		delete(h.Rooms, docID)
	}
	h.Mutex.Unlock()

	close(h.saveQueue)
}
