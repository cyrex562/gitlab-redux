package websocket

import (
	"sync"
)

// Hub maintains the set of active channels and broadcasts messages to the channels
type Hub struct {
	// Registered channels
	channels map[*Channel]bool

	// Inbound messages from the channels
	broadcast chan []byte

	// Register requests from the channels
	register chan *Channel

	// Unregister requests from channels
	unregister chan *Channel

	// Mutex for thread-safe operations
	mu sync.RWMutex
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Channel),
		unregister: make(chan *Channel),
		channels:   make(map[*Channel]bool),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case channel := <-h.register:
			h.mu.Lock()
			h.channels[channel] = true
			h.mu.Unlock()

		case channel := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.channels[channel]; ok {
				delete(h.channels, channel)
			}
			h.mu.Unlock()
			channel.Unsubscribe()

		case message := <-h.broadcast:
			h.mu.RLock()
			for channel := range h.channels {
				if err := channel.Send(message); err != nil {
					channel.Unsubscribe()
					delete(h.channels, channel)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Broadcast sends a message to all registered channels
func (h *Hub) Broadcast(message []byte) {
	h.broadcast <- message
}

// RegisterChannel registers a new channel
func (h *Hub) RegisterChannel(channel *Channel) {
	h.register <- channel
}

// UnregisterChannel unregisters a channel
func (h *Hub) UnregisterChannel(channel *Channel) {
	h.unregister <- channel
}

// GetChannelCount returns the number of active channels
func (h *Hub) GetChannelCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.channels)
}
