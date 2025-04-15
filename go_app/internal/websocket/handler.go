package websocket

import (
	"net/http"

	"github.com/cyrex562/gitlab-redux/internal/websocket/connection"
	"github.com/cyrex562/gitlab-redux/internal/websocket/logging"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Implement proper origin checking
		return true
	},
}

// Handler handles WebSocket connections
type Handler struct {
	hub    *Hub
	logger logging.Logger
}

// NewHandler creates a new WebSocket handler
func NewHandler(hub *Hub) *Handler {
	return &Handler{
		hub:    hub,
		logger: logging.NewDefaultLogger(),
	}
}

// ServeHTTP implements the http.Handler interface
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error(logging.LogPayload{}, "Failed to upgrade connection", err)
		return
	}

	// Create a new connection
	wsConn := connection.NewConnection(conn, r)

	// Attempt to authenticate the connection
	if err := wsConn.Connect(); err != nil {
		h.logger.Error(logging.GetNotificationPayload(wsConn), "Authentication failed", err)
		wsConn.RejectConnection()
		return
	}

	// Create a new channel for this connection
	channel := NewChannel(conn, h.hub)
	channel.SetParams(wsConn.GetParams())

	// Register the channel with the hub
	h.hub.RegisterChannel(channel)

	h.logger.Info(logging.GetNotificationPayload(wsConn), "New WebSocket connection established")

	// Start handling messages
	go h.handleMessages(channel, wsConn)
}

// handleMessages handles incoming messages from a channel
func (h *Handler) handleMessages(channel *Channel, conn *connection.Connection) {
	defer func() {
		h.hub.UnregisterChannel(channel)
		conn.Close()
		h.logger.Info(logging.GetNotificationPayload(conn), "WebSocket connection closed")
	}()

	for {
		message, err := channel.Receive()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				h.logger.Error(logging.GetNotificationPayload(conn), "WebSocket error", err)
			}
			break
		}

		// TODO: Handle the message based on its type
		// This would involve parsing the message and routing it to the appropriate handler
		h.logger.Debug(logging.GetNotificationPayload(conn), "Received message: "+string(message))
	}
}
