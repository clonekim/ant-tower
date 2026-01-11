package server

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// WsHub manages WebSocket connections and broadcasts
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WsHub struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan interface{}
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.RWMutex
}

func NewHub(broadcastChan chan interface{}) *WsHub {
	return &WsHub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  broadcastChan,
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

func (h *WsHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Info().Int("clients_count", len(h.clients)).Msg("WebSocket client connected")

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
			}
			h.mu.Unlock()
			log.Info().Int("clients_count", len(h.clients)).Msg("WebSocket client disconnected")

		case message := <-h.broadcast:
			jsonData, err := json.Marshal(message)
			if err != nil {
				log.Error().Err(err).Msg("Failed to marshal broadcast message")
				continue
			}

			h.mu.RLock()
			for client := range h.clients {
				err := client.WriteMessage(websocket.TextMessage, jsonData)
				if err != nil {
					client.Close()
					h.mu.RUnlock()
					h.mu.Lock()
					delete(h.clients, client)
					h.mu.Unlock()
					h.mu.RLock()
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *WsHub) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Error().Err(err).Msg("WebSocket upgrade error")
		return
	}

	h.register <- conn

	defer func() {
		h.unregister <- conn
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
