package logging

import (
	"net/http"

	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
)

// RequestPayloadLogger handles logging of request payload information
type RequestPayloadLogger struct {
	cloudflareHelper CloudflareHelper
	correlationID    CorrelationID
	applicationContext ApplicationContext
}

// CloudflareHelper defines the interface for Cloudflare-related operations
type CloudflareHelper interface {
	StoreCloudflareHeaders(payload map[string]interface{}, r *http.Request)
}

// CorrelationID defines the interface for correlation ID operations
type CorrelationID interface {
	CurrentID() string
}

// ApplicationContext defines the interface for application context operations
type ApplicationContext interface {
	Current() map[string]interface{}
}

// NewRequestPayloadLogger creates a new instance of RequestPayloadLogger
func NewRequestPayloadLogger(
	cloudflareHelper CloudflareHelper,
	correlationID CorrelationID,
	applicationContext ApplicationContext,
) *RequestPayloadLogger {
	return &RequestPayloadLogger{
		cloudflareHelper: cloudflareHelper,
		correlationID: correlationID,
		applicationContext: applicationContext,
	}
}

// AppendInfoToPayload appends information to the payload for logging
func (r *RequestPayloadLogger) AppendInfoToPayload(payload map[string]interface{}, req *http.Request, authUser *model.User, urgency *model.RequestUrgency) map[string]interface{} {
	// Add user agent
	payload["ua"] = req.Header.Get("User-Agent")

	// Add remote IP
	payload["remote_ip"] = req.RemoteAddr

	// Add correlation ID
	payload["correlation_id"] = r.correlationID.CurrentID()

	// Add application context metadata
	payload["metadata"] = r.applicationContext.Current()

	// Add request urgency and target duration if defined
	if urgency != nil {
		payload["request_urgency"] = urgency.Name
		payload["target_duration_s"] = urgency.Duration
	}

	// Add user ID and username if a user is logged in
	if authUser != nil {
		payload["user_id"] = authUser.ID
		payload["username"] = authUser.Username
	}

	// Add queue duration
	payload["queue_duration_s"] = req.Header.Get("X-Gitlab-Rails-Queue-Duration")

	// Store Cloudflare headers
	r.cloudflareHelper.StoreCloudflareHeaders(payload, req)

	return payload
}
