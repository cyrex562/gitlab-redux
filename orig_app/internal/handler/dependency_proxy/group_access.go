package dependency_proxy

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// GroupAccess provides authorization checks for dependency proxy access
type GroupAccess struct {
	authService  *service.AuthService
	groupService *service.GroupService
}

// NewGroupAccess creates a new instance of GroupAccess
func NewGroupAccess(
	authService *service.AuthService,
	groupService *service.GroupService,
) *GroupAccess {
	return &GroupAccess{
		authService:  authService,
		groupService: groupService,
	}
}

// RegisterRoutes registers the routes for dependency proxy group access
func (g *GroupAccess) RegisterRoutes(r *gin.RouterGroup) {
	// Add middleware for dependency proxy access
	r.Use(g.verifyDependencyProxyAvailable)
	r.Use(g.authorizeReadDependencyProxy)
}

// verifyDependencyProxyAvailable middleware checks if dependency proxy is available for the group
func (g *GroupAccess) verifyDependencyProxyAvailable(ctx *gin.Context) {
	// Get the group from context (should be set by previous middleware)
	group, exists := ctx.Get("group")
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		ctx.Abort()
		return
	}

	groupObj, ok := group.(*model.Group)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid group object"})
		ctx.Abort()
		return
	}

	// Check if dependency proxy is available
	available, err := g.groupService.IsDependencyProxyAvailable(ctx, groupObj)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check dependency proxy availability"})
		ctx.Abort()
		return
	}

	if !available {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Dependency proxy is not available for this group"})
		ctx.Abort()
		return
	}
}

// authorizeReadDependencyProxy middleware checks if the user has permission to read dependency proxy
func (g *GroupAccess) authorizeReadDependencyProxy(ctx *gin.Context) {
	// Get the group from context
	group, exists := ctx.Get("group")
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		ctx.Abort()
		return
	}

	groupObj, ok := group.(*model.Group)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid group object"})
		ctx.Abort()
		return
	}

	// Get the authenticated user or token
	authUser, err := g.getAuthUserOrToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		ctx.Abort()
		return
	}

	// Check permissions based on auth type
	switch authUser.(type) {
	case *model.User:
		if err := g.authorizeReadDependencyProxyForUser(ctx, authUser.(*model.User), groupObj); err != nil {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			ctx.Abort()
			return
		}
	case *model.PersonalAccessToken:
		if err := g.authorizeReadDependencyProxyForToken(ctx, authUser.(*model.PersonalAccessToken), groupObj); err != nil {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			ctx.Abort()
			return
		}
	default:
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authentication type"})
		ctx.Abort()
		return
	}
}

// getAuthUserOrToken returns the authenticated user or token
func (g *GroupAccess) getAuthUserOrToken(ctx *gin.Context) (interface{}, error) {
	// Get the personal access token from context
	token, exists := ctx.Get("personal_access_token")
	if !exists {
		return nil, nil
	}

	pat, ok := token.(*model.PersonalAccessToken)
	if !ok {
		return nil, nil
	}

	// Get the user from context
	user, exists := ctx.Get("user")
	if !exists {
		return nil, nil
	}

	authUser, ok := user.(*model.User)
	if !ok {
		return nil, nil
	}

	// Check if user is a project bot, human, or service account
	if (authUser.IsProjectBot() && authUser.ResourceBotResource.IsGroup()) ||
		authUser.IsHuman() ||
		authUser.IsServiceAccount() {
		return pat, nil
	}

	return authUser, nil
}

// authorizeReadDependencyProxyForUser checks if a user has permission to read dependency proxy
func (g *GroupAccess) authorizeReadDependencyProxyForUser(ctx context.Context, user *model.User, group *model.Group) error {
	// TODO: Implement user permission check
	// This should:
	// 1. Check if user has read_dependency_proxy permission for the group
	// 2. Return error if permission is denied
	return nil
}

// authorizeReadDependencyProxyForToken checks if a token has permission to read dependency proxy
func (g *GroupAccess) authorizeReadDependencyProxyForToken(ctx context.Context, token *model.PersonalAccessToken, group *model.Group) error {
	// TODO: Implement token permission check
	// This should:
	// 1. Check if token has read_dependency_proxy permission for the group's dependency proxy policy subject
	// 2. Return error if permission is denied
	return nil
}
