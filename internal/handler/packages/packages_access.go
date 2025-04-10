package packages

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// PackagesAccess handles access control for packages
type PackagesAccess struct {
	configService *service.ConfigService
	authService   *service.AuthService
	projectService *service.ProjectService
	logger        *service.Logger
}

// NewPackagesAccess creates a new instance of PackagesAccess
func NewPackagesAccess(
	configService *service.ConfigService,
	authService *service.AuthService,
	projectService *service.ProjectService,
	logger *service.Logger,
) *PackagesAccess {
	return &PackagesAccess{
		configService:  configService,
		authService:    authService,
		projectService: projectService,
		logger:         logger,
	}
}

// SetupMiddleware sets up the middleware for packages access
func (p *PackagesAccess) SetupMiddleware(router *gin.Engine) {
	router.Use(func(c *gin.Context) {
		// Verify packages are enabled
		if err := p.verifyPackagesEnabled(c); err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		// Verify read package permission
		if err := p.verifyReadPackage(c); err != nil {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	})
}

// verifyPackagesEnabled checks if packages are enabled
func (p *PackagesAccess) verifyPackagesEnabled(c *gin.Context) error {
	// Get configuration
	config, err := p.configService.GetConfiguration()
	if err != nil {
		return err
	}

	// Check if packages are enabled
	if !config.Packages.Enabled {
		return ErrPackagesDisabled
	}

	return nil
}

// verifyReadPackage checks if the current user can read packages
func (p *PackagesAccess) verifyReadPackage(c *gin.Context) error {
	// Get current user
	user, err := p.authService.GetCurrentUser(c)
	if err != nil {
		return err
	}

	// Get project from context
	project, err := p.projectService.GetProjectFromContext(c)
	if err != nil {
		return err
	}

	// Check if user can read package
	canRead, err := p.authService.CanReadPackage(c, user, project)
	if err != nil {
		return err
	}

	if !canRead {
		return ErrAccessDenied
	}

	return nil
}

// Errors
var (
	ErrPackagesDisabled = &service.Error{
		Code:    "packages_disabled",
		Message: "Packages are not enabled",
	}
	ErrAccessDenied = &service.Error{
		Code:    "access_denied",
		Message: "Access denied",
	}
)
