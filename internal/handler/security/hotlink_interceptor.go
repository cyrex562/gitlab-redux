package security

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// HotlinkInterceptor handles intercepting hotlinking attempts
type HotlinkInterceptor struct {
	hotlinkingService *service.HotlinkingService
}

// NewHotlinkInterceptor creates a new instance of HotlinkInterceptor
func NewHotlinkInterceptor(hotlinkingService *service.HotlinkingService) *HotlinkInterceptor {
	return &HotlinkInterceptor{
		hotlinkingService: hotlinkingService,
	}
}

// InterceptHotlinkingMiddleware creates a middleware that intercepts hotlinking attempts
func (h *HotlinkInterceptor) InterceptHotlinkingMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Check if the request is a hotlinking attempt
		if h.hotlinkingService.IsHotlinking(ctx.Request) {
			// Render a 406 Not Acceptable response
			h.renderNotAcceptable(ctx)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// renderNotAcceptable renders a 406 Not Acceptable response
func (h *HotlinkInterceptor) renderNotAcceptable(ctx *gin.Context) {
	ctx.Status(http.StatusNotAcceptable)
}
