package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// AuthorizationScopes represents the required scopes for WebSocket connections
var AuthorizationScopes = []string{"api", "read_api"}

// Channel represents a WebSocket channel with authentication and validation
type Channel struct {
	conn     *websocket.Conn
	hub      *Hub
	params   map[string]interface{}
	ctx      context.Context
	cancel   context.CancelFunc
	subscribed bool
}

// NewChannel creates a new WebSocket channel
func NewChannel(conn *websocket.Conn, hub *Hub) *Channel {
	ctx, cancel := context.WithCancel(context.Background())
	return &Channel{
		conn:     conn,
		hub:      hub,
		params:   make(map[string]interface{}),
		ctx:      ctx,
		cancel:   cancel,
		subscribed: false,
	}
}

// Subscribe handles the channel subscription
func (c *Channel) Subscribe() error {
	// Validate token scope before allowing subscription
	if err := c.validateTokenScope(); err != nil {
		return err
	}

	// Set up periodic token validation
	go c.periodicTokenValidation()

	c.subscribed = true
	return nil
}

// Unsubscribe handles the channel unsubscription
func (c *Channel) Unsubscribe() {
	if c.subscribed {
		c.subscribed = false
		c.cancel()
		c.conn.Close()
	}
}

// validateTokenScope validates the access token and its scopes
func (c *Channel) validateTokenScope() error {
	// TODO: Implement actual token validation logic
	// This should check if the token has the required scopes (api, read_api)
	return nil
}

// periodicTokenValidation periodically validates the token
func (c *Channel) periodicTokenValidation() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			if err := c.validateTokenScope(); err != nil {
				log.Printf("Token validation failed: %v", err)
				c.Unsubscribe()
				return
			}
		}
	}
}

// Send sends a message to the client
func (c *Channel) Send(message interface{}) error {
	if !c.subscribed {
		return errors.New("channel not subscribed")
	}

	// Add params to the message
	payload := map[string]interface{}{
		"params": c.params,
	}

	if msg, ok := message.(map[string]interface{}); ok {
		for k, v := range msg {
			payload[k] = v
		}
	} else {
		payload["message"] = message
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return c.conn.WriteMessage(websocket.TextMessage, data)
}

// Receive handles incoming messages
func (c *Channel) Receive() ([]byte, error) {
	_, message, err := c.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	return message, nil
}

// SetParams sets the channel parameters
func (c *Channel) SetParams(params map[string]interface{}) {
	c.params = params
}

// GetParams returns the channel parameters
func (c *Channel) GetParams() map[string]interface{} {
	return c.params
}
