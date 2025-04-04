package harbor

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// Access provides authorization checks for Harbor registry access
type Access struct {
	authService *service.AuthService
}

// NewAccess creates a new instance of Access
func NewAccess(authService *service.AuthService) *Access {
	return &Access{
		authService: authService,
	}
}

// RegisterRoutes registers the routes for Harbor registry access
func (a *Access) RegisterRoutes(r *gin.RouterGroup) {
	// Add middleware for Harbor registry authorization
	r.Use(a.authorizeReadHarborRegistry)
}

// authorizeReadHarborRegistry middleware checks if the user has permission to read Harbor registry
func (a *Access) authorizeReadHarborRegistry(ctx *gin.Context) {
	// Get the project from context (should be set by previous middleware)
	project, exists := ctx.Get("project")
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		ctx.Abort()
		return
	}

	// Check if user has permission to read Harbor registry
	if err := a.authService.AuthorizeReadHarborRegistry(ctx, project); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to read Harbor registry"})
		ctx.Abort()
		return
	}
}
