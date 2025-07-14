package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/websocket/v2"
	"github.com/hafiddna/infrastructure-be-helper/dto"
	"log"
	"sync"
)

type Hub struct {
	clients map[string][]*websocket.Conn
	mu      sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{clients: make(map[string][]*websocket.Conn)}
}

func (h *Hub) HandleUserConnection(c *websocket.Conn) {
	user := c.Locals("user").(map[string]interface{})
	userID := fmt.Sprintf("%v", user["sub"])

	h.mu.Lock()
	h.clients[userID] = append(h.clients[userID], c)
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		conns := h.clients[userID]
		for i, conn := range conns {
			if conn == c {
				h.clients[userID] = append(conns[:i], conns[i+1:]...)
				break
			}
		}
		// Clean up if no more connections for user
		if len(h.clients[userID]) == 0 {
			delete(h.clients, userID)
		}
		h.mu.Unlock()
		c.Close()
	}()

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			break
		}

		var wsMsg dto.WebSocketMessage
		if err := json.Unmarshal(msg, &wsMsg); err != nil {
			log.Println("Invalid WS message:", err)
			continue
		}

		// Handle message types here (ping/chat/notification/etc)
	}
}

func (h *Hub) SendToUser(userID string, data dto.WebSocketMessage) error {
	h.mu.RLock()
	conns, ok := h.clients[userID]
	h.mu.RUnlock()

	if !ok || len(conns) == 0 {
		return fmt.Errorf("no WebSocket connections for user %s", userID)
	}

	var failed []int
	for i, conn := range conns {
		if err := conn.WriteJSON(data); err != nil {
			log.Printf("failed to send to conn %d for user %s: %v", i, userID, err)
			failed = append(failed, i)
		}
	}

	if len(failed) > 0 {
		// Clean up failed connections
		h.mu.Lock()
		var newConns []*websocket.Conn
		for i, conn := range h.clients[userID] {
			skip := false
			for _, f := range failed {
				if f == i {
					conn.Close()
					skip = true
					break
				}
			}
			if !skip {
				newConns = append(newConns, conn)
			}
		}
		h.clients[userID] = newConns
		if len(newConns) == 0 {
			delete(h.clients, userID)
		}
		h.mu.Unlock()

		return fmt.Errorf("some messages failed to send to user %s", userID)
	}

	return nil
}
