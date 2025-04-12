package channel

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Channel represents a WebSocket channel
type Channel struct {
	conn   *websocket.Conn
	hub    *Hub
	params map[string]interface{}
	mu     sync.RWMutex
}

// Hub manages WebSocket channels
type Hub struct {
	channels   map[*Channel]bool
	broadcast  chan []byte
	register   chan *Channel
	unregister chan *Channel
	mu         sync.RWMutex
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		channels:   make(map[*Channel]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Channel),
		unregister: make(chan *Channel),
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
			channel.conn.Close()

		case message := <-h.broadcast:
			h.mu.RLock()
			for channel := range h.channels {
				err := channel.conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					channel.conn.Close()
					delete(h.channels, channel)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// NewChannel creates a new Channel instance
func NewChannel(conn *websocket.Conn, hub *Hub) *Channel {
	return &Channel{
		conn:   conn,
		hub:    hub,
		params: make(map[string]interface{}),
	}
}

// SetParams sets the channel parameters
func (c *Channel) SetParams(params map[string]interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.params = params
}

// GetParams returns the channel parameters
func (c *Channel) GetParams() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.params
}

// Send sends a message through the channel
func (c *Channel) Send(message []byte) error {
	return c.conn.WriteMessage(websocket.TextMessage, message)
}

// Receive receives a message from the channel
func (c *Channel) Receive() ([]byte, error) {
	_, message, err := c.conn.ReadMessage()
	return message, err
}

// Close closes the channel connection
func (c *Channel) Close() error {
	return c.conn.Close()
}
