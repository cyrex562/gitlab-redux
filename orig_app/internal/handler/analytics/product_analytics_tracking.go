package analytics

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// ProductAnalyticsTracking handles product analytics tracking functionality
type ProductAnalyticsTracking struct {
	trackingService *service.TrackingService
	redisService    *service.RedisService
	cookieService   *service.CookieService
	userService     *service.UserService
	projectService  *service.ProjectService
	namespaceService *service.NamespaceService
	logger          *service.Logger
}

// NewProductAnalyticsTracking creates a new instance of ProductAnalyticsTracking
func NewProductAnalyticsTracking(
	trackingService *service.TrackingService,
	redisService *service.RedisService,
	cookieService *service.CookieService,
	userService *service.UserService,
	projectService *service.ProjectService,
	namespaceService *service.NamespaceService,
	logger *service.Logger,
) *ProductAnalyticsTracking {
	return &ProductAnalyticsTracking{
		trackingService:  trackingService,
		redisService:     redisService,
		cookieService:    cookieService,
		userService:      userService,
		projectService:   projectService,
		namespaceService: namespaceService,
		logger:           logger,
	}
}

// TrackEvent sets up tracking for controller actions
func (p *ProductAnalyticsTracking) TrackEvent(c *gin.Context, controllerActions []string, name string, action string, label string, conditions []string, destinations []string, customIDFn func(c *gin.Context) string) {
	// Check if the current action is in the list of controller actions
	currentAction := c.GetString("action")
	actionIncluded := false
	for _, a := range controllerActions {
		if a == currentAction {
			actionIncluded = true
			break
		}
	}

	if !actionIncluded {
		return
	}

	// Check if the request is trackable
	if !p.isTrackableHTMLRequest(c) {
		return
	}

	// Check additional conditions
	for _, condition := range conditions {
		if !p.checkCondition(c, condition) {
			return
		}
	}

	// Route events to destinations
	p.routeEventsTo(c, destinations, name, action, label, customIDFn)
}

// TrackInternalEvent sets up internal event tracking for controller actions
func (p *ProductAnalyticsTracking) TrackInternalEvent(c *gin.Context, controllerActions []string, name string, conditions []string, eventArgs map[string]interface{}) {
	// Check if the current action is in the list of controller actions
	currentAction := c.GetString("action")
	actionIncluded := false
	for _, a := range controllerActions {
		if a == currentAction {
			actionIncluded = true
			break
		}
	}

	if !actionIncluded {
		return
	}

	// Check if the request is trackable
	if !p.isTrackableHTMLRequest(c) {
		return
	}

	// Check additional conditions
	for _, condition := range conditions {
		if !p.checkCondition(c, condition) {
			return
		}
	}

	// Get current user
	user, err := p.userService.GetCurrentUser(c)
	if err != nil {
		p.logger.Error("Failed to get current user for internal event tracking", err)
		return
	}

	// Get tracking project source
	project, err := p.getTrackingProjectSource(c)
	if err != nil {
		p.logger.Error("Failed to get tracking project source", err)
	}

	// Get tracking namespace source
	namespace, err := p.getTrackingNamespaceSource(c)
	if err != nil {
		p.logger.Error("Failed to get tracking namespace source", err)
	}

	// Add user, project, and namespace to event args
	eventArgs["user"] = user
	if project != nil {
		eventArgs["project"] = project
	}
	if namespace != nil {
		eventArgs["namespace"] = namespace
	}

	// Track internal event
	err = p.trackingService.TrackInternalEvent(name, eventArgs)
	if err != nil {
		p.logger.Error("Failed to track internal event", err)
	}
}

// RouteEventsTo routes events to the specified destinations
func (p *ProductAnalyticsTracking) routeEventsTo(c *gin.Context, destinations []string, name string, action string, label string, customIDFn func(c *gin.Context) string) {
	// Track unique Redis HLL event if destination includes redis_hll
	redisHLLIncluded := false
	for _, dest := range destinations {
		if dest == "redis_hll" {
			redisHLLIncluded = true
			break
		}
	}

	if redisHLLIncluded {
		p.trackUniqueRedisHLLEvent(c, name, customIDFn)
	}

	// Track Snowplow event if destination includes snowplow
	snowplowIncluded := false
	for _, dest := range destinations {
		if dest == "snowplow" {
			snowplowIncluded = true
			break
		}
	}

	if snowplowIncluded {
		// Validate required parameters for Snowplow
		if action == "" {
			p.logger.Error("Action is required when destination is snowplow", nil)
			return
		}
		if label == "" {
			p.logger.Error("Label is required when destination is snowplow", nil)
			return
		}

		// Get current user
		user, err := p.userService.GetCurrentUser(c)
		if err != nil {
			p.logger.Error("Failed to get current user for Snowplow tracking", err)
			return
		}

		// Get tracking project source
		project, err := p.getTrackingProjectSource(c)
		if err != nil {
			p.logger.Error("Failed to get tracking project source for Snowplow", err)
		}

		// Get tracking namespace source
		namespace, err := p.getTrackingNamespaceSource(c)
		if err != nil {
			p.logger.Error("Failed to get tracking namespace source for Snowplow", err)
		}

		// Create optional arguments
		optionalArgs := make(map[string]interface{})
		if namespace != nil {
			optionalArgs["namespace"] = namespace
		}
		if project != nil {
			optionalArgs["project"] = project
		}

		// Create service ping context
		servicePingContext := p.trackingService.NewServicePingContext("redis_hll", name)

		// Track event
		err = p.trackingService.TrackEvent(
			c.GetString("controller"),
			action,
			user,
			name,
			label,
			[]interface{}{servicePingContext},
			optionalArgs,
		)
		if err != nil {
			p.logger.Error("Failed to track Snowplow event", err)
		}
	}
}

// TrackUniqueRedisHLLEvent tracks a unique Redis HLL event
func (p *ProductAnalyticsTracking) trackUniqueRedisHLLEvent(c *gin.Context, eventName string, customIDFn func(c *gin.Context) string) {
	var uniqueID string

	// Get custom ID from function if provided
	if customIDFn != nil {
		uniqueID = customIDFn(c)
	}

	// If no custom ID, get visitor ID
	if uniqueID == "" {
		uniqueID = p.getVisitorID(c)
	}

	// Return if no unique ID
	if uniqueID == "" {
		return
	}

	// Track event in Redis HLL counter
	err := p.redisService.TrackHLLEvent(eventName, []string{uniqueID})
	if err != nil {
		p.logger.Error("Failed to track Redis HLL event", err)
	}
}

// GetVisitorID gets the visitor ID from cookies or creates a new one
func (p *ProductAnalyticsTracking) getVisitorID(c *gin.Context) string {
	// Check if visitor ID exists in cookies
	visitorID, err := p.cookieService.GetCookie(c, "visitor_id")
	if err == nil && visitorID != "" {
		return visitorID
	}

	// Get current user
	user, err := p.userService.GetCurrentUser(c)
	if err != nil {
		return ""
	}

	// Generate new UUID
	uuid := p.generateUUID()
	if uuid == "" {
		return ""
	}

	// Set cookie with 24 months expiration
	expiration := time.Now().Add(24 * 30 * 24 * time.Hour) // 24 months
	err = p.cookieService.SetCookie(c, "visitor_id", uuid, expiration)
	if err != nil {
		p.logger.Error("Failed to set visitor_id cookie", err)
		return ""
	}

	return uuid
}

// GenerateUUID generates a new UUID
func (p *ProductAnalyticsTracking) generateUUID() string {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		p.logger.Error("Failed to generate UUID", err)
		return ""
	}

	// Set version 4
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	// Set variant
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	return hex.EncodeToString(uuid)
}

// IsTrackableHTMLRequest checks if the request is trackable HTML request
func (p *ProductAnalyticsTracking) isTrackableHTMLRequest(c *gin.Context) bool {
	// Check if request is HTML
	acceptHeader := c.GetHeader("Accept")
	if acceptHeader == "" {
		return true
	}

	return acceptHeader == "*/*" || acceptHeader == "text/html" || acceptHeader == "application/xhtml+xml"
}

// CheckCondition checks a condition
func (p *ProductAnalyticsTracking) checkCondition(c *gin.Context, condition string) bool {
	switch condition {
	case "trackable_html_request?":
		return p.isTrackableHTMLRequest(c)
	default:
		// Handle other conditions as needed
		return true
	}
}

// GetTrackingProjectSource gets the tracking project source
func (p *ProductAnalyticsTracking) getTrackingProjectSource(c *gin.Context) (interface{}, error) {
	// This is a placeholder for actual implementation
	// In a real implementation, this would get the project from the context
	return nil, nil
}

// GetTrackingNamespaceSource gets the tracking namespace source
func (p *ProductAnalyticsTracking) getTrackingNamespaceSource(c *gin.Context) (interface{}, error) {
	// This is a placeholder for actual implementation
	// In a real implementation, this would get the namespace from the context
	return nil, nil
}
