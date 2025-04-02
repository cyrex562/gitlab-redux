package logging

import (
	"net/http"

	"gitlab.com/gitlab-org/gitlab/internal/websocket/connection"
)

// LogKey represents the correlation ID key in logs
const LogKey = "X-Request-ID"

// LogPayload represents the structured log data
type LogPayload struct {
	CorrelationID string                 `json:"correlation_id"`
	UserID        int64                  `json:"user_id,omitempty"`
	Username      string                 `json:"username,omitempty"`
	RemoteIP      string                 `json:"remote_ip"`
	UserAgent     string                 `json:"user_agent"`
	Params        map[string]interface{} `json:"params,omitempty"`
}

// GetNotificationPayload creates a log payload from a connection
func GetNotificationPayload(conn *connection.Connection) LogPayload {
	user := conn.GetCurrentUser()
	req := conn.GetRequest()

	payload := LogPayload{
		CorrelationID: req.Header.Get(LogKey),
		RemoteIP:      req.RemoteAddr,
		UserAgent:     req.Header.Get("User-Agent"),
		Params:        conn.GetParams(),
	}

	if user != nil {
		payload.UserID = user.ID
		payload.Username = user.Username
	}

	return payload
}

// WithLogging wraps a handler with logging functionality
func WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ensure correlation ID exists
		if r.Header.Get(LogKey) == "" {
			// TODO: Generate a new correlation ID if none exists
			// This could be a UUID or other unique identifier
		}

		next.ServeHTTP(w, r)
	})
}
