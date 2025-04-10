package error

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
	"gitlab.com/gitlab-org/gitlab-redux/internal/util"
)

// GracefulTimeoutHandling handles graceful timeout handling for database queries
type GracefulTimeoutHandling struct {
	logger *util.Logger
}

// NewGracefulTimeoutHandling creates a new instance of GracefulTimeoutHandling
func NewGracefulTimeoutHandling(logger *util.Logger) *GracefulTimeoutHandling {
	return &GracefulTimeoutHandling{
		logger: logger,
	}
}

// HandleQueryTimeoutMiddleware creates a middleware that handles query timeouts
func (g *GracefulTimeoutHandling) HandleQueryTimeoutMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Use a defer-recover pattern to catch panics
		defer func() {
			if r := recover(); r != nil {
				// Check if the panic is a QueryTimeoutError
				if err, ok := r.(*service.QueryTimeoutError); ok {
					// Check if the request format is JSON
					if ctx.GetHeader("Accept") == "application/json" || ctx.GetHeader("Content-Type") == "application/json" {
						// Log the exception
						g.logger.Error("Query timeout error", "error", err)

						// Return a JSON error response
						ctx.JSON(http.StatusRequestTimeout, gin.H{
							"error": "There is too much data to calculate. Please change your selection.",
						})
						ctx.Abort()
						return
					}
				}

				// Re-panic if it's not a QueryTimeoutError or not a JSON request
				panic(r)
			}
		}()

		ctx.Next()
	}
}
