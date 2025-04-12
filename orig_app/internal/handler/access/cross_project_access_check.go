package access

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
	"gitlab.com/gitlab-org/gitlab-redux/internal/util"
)

// CrossProjectAccessCheck handles checking if a user has permission to access information across multiple projects
type CrossProjectAccessCheck struct {
	authService *service.AuthService
	crossProjectService *service.CrossProjectService
	logger      *util.Logger
}

// NewCrossProjectAccessCheck creates a new instance of CrossProjectAccessCheck
func NewCrossProjectAccessCheck(
	authService *service.AuthService,
	crossProjectService *service.CrossProjectService,
	logger *util.Logger,
) *CrossProjectAccessCheck {
	return &CrossProjectAccessCheck{
		authService: authService,
		crossProjectService: crossProjectService,
		logger:      logger,
	}
}

// CrossProjectCheck checks if the user has permission to access information across multiple projects
func (c *CrossProjectAccessCheck) CrossProjectCheck(ctx *gin.Context) error {
	// Get the current user from the context
	currentUser, exists := ctx.Get("current_user")
	if !exists {
		return util.NewUnauthorizedError("user not authenticated")
	}
	user := currentUser.(*model.User)

	// Check if we should run the cross project check
	check := c.crossProjectService.FindCheck(ctx)
	if check == nil || !check.ShouldRun(ctx) {
		return nil
	}

	// Authorize the cross project page
	return c.authorizeCrossProjectPage(ctx, user)
}

// authorizeCrossProjectPage authorizes access to the cross project page
func (c *CrossProjectAccessCheck) authorizeCrossProjectPage(ctx *gin.Context, user *model.User) error {
	// Check if the user can read across projects
	if c.authService.Can(ctx, user, "read_cross_project") {
		return nil
	}

	// Return an access denied error
	rejectionMessage := "This page is unavailable because you are not allowed to read information across multiple projects."
	return util.NewForbiddenError(rejectionMessage)
}
