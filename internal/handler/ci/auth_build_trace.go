package ci

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// AuthBuildTrace provides authorization checks for build trace access
type AuthBuildTrace struct {
	authService  *service.AuthService
	buildService *service.BuildService
}

// NewAuthBuildTrace creates a new instance of AuthBuildTrace
func NewAuthBuildTrace(
	authService *service.AuthService,
	buildService *service.BuildService,
) *AuthBuildTrace {
	return &AuthBuildTrace{
		authService:  authService,
		buildService: buildService,
	}
}

// RegisterRoutes registers the routes for build trace authorization
func (a *AuthBuildTrace) RegisterRoutes(r *gin.RouterGroup) {
	// Add middleware for build trace authorization
	r.Use(a.authorizeReadBuildTrace)
}

// authorizeReadBuildTrace middleware checks if the user has permission to read build trace
func (a *AuthBuildTrace) authorizeReadBuildTrace(ctx *gin.Context) {
	// Get the build from context (should be set by previous middleware)
	build, exists := ctx.Get("build")
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Build not found"})
		ctx.Abort()
		return
	}

	buildObj, ok := build.(*model.Build)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid build object"})
		ctx.Abort()
		return
	}

	// Check if user has permission to read build trace
	hasPermission, err := a.authService.CanReadBuildTrace(ctx, buildObj)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permissions"})
		ctx.Abort()
		return
	}

	if hasPermission {
		return
	}

	// If debug mode is enabled, provide more detailed error message
	if buildObj.DebugMode {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "You must have developer or higher permissions in the associated project to view job logs when debug " +
				"trace is enabled. To disable debug trace, set the 'CI_DEBUG_TRACE' and 'CI_DEBUG_SERVICES' variables to " +
				"'false' in your pipeline configuration or CI/CD settings. If you must view this job log, " +
				"a project maintainer or owner must add you to the project with developer permissions or higher.",
		})
	} else {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "The current user is not authorized to access the job log."})
	}

	ctx.Abort()
}
