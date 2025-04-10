package project

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// ProjectUnauthorized handles project unauthorized access
type ProjectUnauthorized struct {
	projectService *service.ProjectService
	authService    *service.AuthService
	userService    *service.UserService
	logger         *service.Logger
}

// NewProjectUnauthorized creates a new instance of ProjectUnauthorized
func NewProjectUnauthorized(
	projectService *service.ProjectService,
	authService *service.AuthService,
	userService *service.UserService,
	logger *service.Logger,
) *ProjectUnauthorized {
	return &ProjectUnauthorized{
		projectService: projectService,
		authService:    authService,
		userService:    userService,
		logger:         logger,
	}
}

// OnRoutableNotFound handles the case when a routable is not found
func (p *ProjectUnauthorized) OnRoutableNotFound(c *gin.Context, routable interface{}, fullPath string) error {
	// Check if routable is a project
	project, ok := routable.(*service.Project)
	if !ok {
		return nil
	}

	// Get external authorization classification label
	label, err := p.projectService.GetExternalAuthorizationClassificationLabel(project)
	if err != nil {
		return err
	}

	// Get current user
	user, err := p.userService.GetCurrentUser(c)
	if err != nil {
		return err
	}

	// Check if access is allowed
	accessAllowed, err := p.authService.IsExternalAuthorizationAccessAllowed(user, label)
	if err != nil {
		return err
	}

	// If access is allowed, return nil
	if accessAllowed {
		return nil
	}

	// Get rejection reason
	rejectionReason, err := p.authService.GetExternalAuthorizationRejectionReason(user, label)
	if err != nil {
		return err
	}

	// If no rejection reason, use default message
	if rejectionReason == "" {
		rejectionReason = "External authorization denied access to this project"
	}

	// Return 403 Forbidden
	c.JSON(http.StatusForbidden, gin.H{
		"message": rejectionReason,
	})

	return nil
}
