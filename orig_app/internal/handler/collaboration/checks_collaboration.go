package collaboration

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
	"gitlab.com/gitlab-org/gitlab-redux/internal/util"
)

// ChecksCollaboration handles checking if a user can collaborate with a project
type ChecksCollaboration struct {
	authService *service.AuthService
	userService *service.UserService
	logger      *util.Logger
	userAccess  map[*model.Project]*model.UserAccess
}

// NewChecksCollaboration creates a new instance of ChecksCollaboration
func NewChecksCollaboration(
	authService *service.AuthService,
	userService *service.UserService,
	logger *util.Logger,
) *ChecksCollaboration {
	return &ChecksCollaboration{
		authService: authService,
		userService: userService,
		logger:      logger,
		userAccess:  make(map[*model.Project]*model.UserAccess),
	}
}

// CanCollaborateWithProject checks if a user can collaborate with a project
func (c *ChecksCollaboration) CanCollaborateWithProject(ctx *gin.Context, project *model.Project, ref string) bool {
	// Get the current user from the context
	currentUser, exists := ctx.Get("current_user")
	if !exists {
		return false
	}
	user := currentUser.(*model.User)

	// Check if the user can push code to the project
	if c.authService.Can(ctx, user, "push_code", project) {
		return true
	}

	// Check if the user can create merge requests and has already forked the project
	canCreateMergeRequest := c.authService.Can(ctx, user, "create_merge_request_in", project) &&
		c.userService.AlreadyForked(ctx, user, project)

	// Check if the user can push to the branch
	canPushToBranch := c.getUserAccess(ctx, project).CanPushToBranch(ref)

	return canCreateMergeRequest || canPushToBranch
}

// getUserAccess gets the user access for a project
func (c *ChecksCollaboration) getUserAccess(ctx *gin.Context, project *model.Project) *model.UserAccess {
	// Check if we already have the user access cached
	if access, exists := c.userAccess[project]; exists {
		return access
	}

	// Get the current user from the context
	currentUser, exists := ctx.Get("current_user")
	if !exists {
		return nil
	}
	user := currentUser.(*model.User)

	// Create a new user access
	access := model.NewUserAccess(user, project)

	// Cache the user access
	c.userAccess[project] = access

	return access
}
