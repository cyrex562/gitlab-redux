package registry

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

var (
	// ErrInvalidRegistryPath represents an invalid registry path error
	ErrInvalidRegistryPath = errors.New("invalid registry path")
)

// ConnectionErrorsHandler provides functionality for handling registry connection errors
type ConnectionErrorsHandler struct {
	registryService *service.RegistryService
}

// NewConnectionErrorsHandler creates a new instance of ConnectionErrorsHandler
func NewConnectionErrorsHandler(registryService *service.RegistryService) *ConnectionErrorsHandler {
	return &ConnectionErrorsHandler{
		registryService: registryService,
	}
}

// RegisterMiddleware registers the middleware for handling connection errors
func (h *ConnectionErrorsHandler) RegisterMiddleware(router *gin.RouterGroup) {
	router.Use(h.pingContainerRegistry)
	router.Use(h.handleErrors)
}

// handleErrors middleware handles registry errors
func (h *ConnectionErrorsHandler) handleErrors(ctx *gin.Context) {
	// Store error handlers in context for use in error handling
	ctx.Set("handleInvalidPath", h.handleInvalidPath)
	ctx.Set("handleConnectionError", h.handleConnectionError)

	// Continue with request
	ctx.Next()

	// Check for errors
	if err := ctx.Errors.Last(); err != nil {
		switch {
		case errors.Is(err.Err, ErrInvalidRegistryPath):
			h.handleInvalidPath(ctx)
		case errors.Is(err.Err, service.ErrRegistryConnection):
			h.handleConnectionError(ctx)
		}
	}
}

// pingContainerRegistry middleware pings the container registry
func (h *ConnectionErrorsHandler) pingContainerRegistry(ctx *gin.Context) {
	err := h.registryService.PingRegistry(ctx)
	if err != nil {
		ctx.Error(err)
		ctx.Abort()
		return
	}
	ctx.Next()
}

// handleInvalidPath handles invalid registry path errors
func (h *ConnectionErrorsHandler) handleInvalidPath(ctx *gin.Context) {
	// Set error flag for frontend
	ctx.Set("invalid_path_error", true)

	// Render index template
	h.renderIndex(ctx)
}

// handleConnectionError handles registry connection errors
func (h *ConnectionErrorsHandler) handleConnectionError(ctx *gin.Context) {
	// Set error flag for frontend
	ctx.Set("connection_error", true)

	// Render index template
	h.renderIndex(ctx)
}

// renderIndex renders the index template with error flags
func (h *ConnectionErrorsHandler) renderIndex(ctx *gin.Context) {
	// Get error flags
	invalidPathError, _ := ctx.Get("invalid_path_error")
	connectionError, _ := ctx.Get("connection_error")

	// Render template with error flags
	ctx.HTML(http.StatusOK, "registry/repositories/index", gin.H{
		"invalid_path_error": invalidPathError,
		"connection_error":   connectionError,
	})
}
