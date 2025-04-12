package connection

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// Connection represents a WebSocket connection with authentication
type Connection struct {
	conn       *websocket.Conn
	request    *http.Request
	currentUser *User
	params     map[string]interface{}
}

// NewConnection creates a new WebSocket connection
func NewConnection(conn *websocket.Conn, r *http.Request) *Connection {
	return &Connection{
		conn:    conn,
		request: r,
		params:  make(map[string]interface{}),
	}
}

// Connect establishes the connection and authenticates the user
func (c *Connection) Connect() error {
	// Try to find user from bearer token first
	if user, err := c.findUserFromBearerToken(); err == nil && user != nil {
		c.currentUser = user
		return nil
	}

	// If no bearer token, try to find user from session
	if user, err := c.findUserFromSession(); err == nil && user != nil {
		c.currentUser = user
		return nil
	}

	return ErrUnauthorized
}

// findUserFromBearerToken attempts to find a user from the Authorization header
func (c *Connection) findUserFromBearerToken() (*User, error) {
	authHeader := c.request.Header.Get("Authorization")
	if authHeader == "" {
		return nil, nil
	}

	// Check if it's a bearer token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, nil
	}

	token := parts[1]
	// TODO: Implement actual token validation and user lookup
	// This should validate the token and return the associated user
	return nil, nil
}

// findUserFromSession attempts to find a user from the session cookie
func (c *Connection) findUserFromSession() (*User, error) {
	// Get the session cookie
	sessionCookie, err := c.request.Cookie("_gitlab_session")
	if err != nil {
		return nil, nil
	}

	// TODO: Implement session validation and user lookup
	// This should validate the session and return the associated user
	return nil, nil
}

// GetCurrentUser returns the current authenticated user
func (c *Connection) GetCurrentUser() *User {
	return c.currentUser
}

// SetParams sets the connection parameters
func (c *Connection) SetParams(params map[string]interface{}) {
	c.params = params
}

// GetParams returns the connection parameters
func (c *Connection) GetParams() map[string]interface{} {
	return c.params
}

// GetRequest returns the underlying HTTP request
func (c *Connection) GetRequest() *http.Request {
	return c.request
}

// RejectConnection closes the connection with an unauthorized error
func (c *Connection) RejectConnection() {
	c.conn.WriteControl(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Unauthorized"),
		time.Now().Add(time.Second))
	c.conn.Close()
}

// Send sends a message to the client
func (c *Connection) Send(message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return c.conn.WriteMessage(websocket.TextMessage, data)
}

// Receive receives a message from the client
func (c *Connection) Receive() ([]byte, error) {
	_, message, err := c.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	return message, nil
}

// Close closes the connection
func (c *Connection) Close() error {
	return c.conn.Close()
}
