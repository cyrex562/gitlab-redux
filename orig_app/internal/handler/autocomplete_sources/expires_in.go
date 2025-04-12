package autocomplete_sources

import (
	"time"

	"github.com/gin-gonic/gin"
)

const (
	// AutoCompleteExpiresIn is the duration for which autocomplete results should be cached
	AutoCompleteExpiresIn = 3 * time.Minute
)

// AutoCompleteCachedActions is a list of actions that should be cached
var AutoCompleteCachedActions = []string{"members", "labels"}

// ExpiresIn provides caching functionality for autocomplete sources
type ExpiresIn struct {
	// Add any dependencies here if needed
}

// NewExpiresIn creates a new instance of ExpiresIn
func NewExpiresIn() *ExpiresIn {
	return &ExpiresIn{}
}

// RegisterRoutes registers the routes for autocomplete sources
func (e *ExpiresIn) RegisterRoutes(r *gin.RouterGroup) {
	// Apply middleware for cached actions
	for _, action := range AutoCompleteCachedActions {
		r.Use(e.setExpiresIn)
	}
}

// setExpiresIn middleware sets the cache expiration time for autocomplete responses
func (e *ExpiresIn) setExpiresIn(ctx *gin.Context) {
	// Set cache control headers
	ctx.Header("Cache-Control", "public, max-age=180") // 3 minutes in seconds
	ctx.Header("Expires", time.Now().Add(AutoCompleteExpiresIn).Format(time.RFC1123))
}
