package realtime

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type eventMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type client struct {
	adminID int
	conn    *websocket.Conn
	send    chan []byte
}

type Hub struct {
	mu             sync.RWMutex
	clients        map[*client]struct{}
	clientsByAdmin map[int]*client
}

func NewHub() *Hub {
	return &Hub{
		clients:        make(map[*client]struct{}),
		clientsByAdmin: make(map[int]*client),
	}
}

var GlobalHub = NewHub()

func (h *Hub) RegisterConnection(conn *websocket.Conn, adminID int) {
	c := &client{
		adminID: adminID,
		conn:    conn,
		send:    make(chan []byte, 32),
	}

	h.mu.Lock()
	if existing, exists := h.clientsByAdmin[adminID]; exists {
		delete(h.clients, existing)
		delete(h.clientsByAdmin, adminID)
		_ = existing.conn.Close()
	}
	h.clientsByAdmin[adminID] = c
	h.clients[c] = struct{}{}
	h.mu.Unlock()

	go h.writePump(c)
	h.readPump(c)
}

func (h *Hub) Broadcast(event string, data interface{}) {
	payload, err := json.Marshal(eventMessage{Event: event, Data: data})
	if err != nil {
		log.Printf("realtime: failed to marshal event %s: %v", event, err)
		return
	}

	h.mu.RLock()
	clients := make([]*client, 0, len(h.clients))
	for c := range h.clients {
		clients = append(clients, c)
	}
	h.mu.RUnlock()

	for _, c := range clients {
		select {
		case c.send <- payload:
		default:
			h.unregister(c)
		}
	}
}

func (h *Hub) readPump(c *client) {
	defer h.unregister(c)

	_ = c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	})

	for {
		if _, _, err := c.conn.ReadMessage(); err != nil {
			return
		}
	}
}

func (h *Hub) writePump(c *client) {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			_ = c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (h *Hub) unregister(c *client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.clients[c]; !exists {
		return
	}

	delete(h.clients, c)
	if existing, exists := h.clientsByAdmin[c.adminID]; exists && existing == c {
		delete(h.clientsByAdmin, c.adminID)
	}
	close(c.send)
	_ = c.conn.Close()
}
