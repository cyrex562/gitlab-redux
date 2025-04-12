package groups

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

// ApplicationController is the base controller for all group controllers
type ApplicationController struct {
	// Add any dependencies here
	groupFinder        *GroupFinder
	groupProjectsFinder *GroupProjectsFinder
	authorizationService *AuthorizationService
	sortingService      *SortingService
}

// NewApplicationController creates a new ApplicationController
func NewApplicationController(
	groupFinder *GroupFinder,
	groupProjectsFinder *GroupProjectsFinder,
	authorizationService *AuthorizationService,
	sortingService *SortingService,
) *ApplicationController {
	return &ApplicationController{
		groupFinder:        groupFinder,
		groupProjectsFinder: groupProjectsFinder,
		authorizationService: authorizationService,
		sortingService:      sortingService,
	}
}

// RegisterMiddleware registers the middleware for the ApplicationController
func (c *ApplicationController) RegisterMiddleware(router *gin.RouterGroup) {
	// Skip authentication middleware
	// router.Use(c.AuthenticateUser())

	// Add group middleware
	router.Use(c.SetGroup())

	// Add sorting middleware
	router.Use(c.SetSorting())

	// Add cross project access middleware
	router.Use(c.RequireCrossProjectAccess())

	// Add math rendering limits middleware
	router.Use(c.PushNamespaceSetting())
}

// SetGroup middleware sets the group in the context
func (c *ApplicationController) SetGroup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the group ID from the URL parameters
		groupID := ctx.Param("group_id")
		if groupID == "" {
			groupID = ctx.Param("id")
		}

		// Find the group
		group, err := c.groupFinder.FindRoutable(groupID, ctx.Request.URL.Path)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
			ctx.Abort()
			return
		}

		// Set the group in the context
		ctx.Set("group", group)

		ctx.Next()
	}
}

// SetSorting middleware sets the sorting in the context
func (c *ApplicationController) SetSorting() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Check if the controller has a project list
		if c.HasProjectList() {
			// Set the sorting preference
			sortValue := c.sortingService.GetSortValue(ctx)
			ctx.Set("group_projects_sort", sortValue)
		}

		ctx.Next()
	}
}

// RequireCrossProjectAccess middleware requires cross project access
func (c *ApplicationController) RequireCrossProjectAccess() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// This is a placeholder for the actual implementation
		// In the actual implementation, you would check if the user has cross project access

		ctx.Next()
	}
}

// PushNamespaceSetting middleware pushes the namespace setting
func (c *ApplicationController) PushNamespaceSetting() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the group from the context
		group, _ := ctx.Get("group")

		// Push the namespace setting
		// This is a placeholder for the actual implementation
		// In the actual implementation, you would push the math rendering limits setting

		ctx.Next()
	}
}

// GetGroup gets the group from the context
func (c *ApplicationController) GetGroup(ctx *gin.Context) interface{} {
	group, _ := ctx.Get("group")
	return group
}

// GetGroupProjects gets the group projects
func (c *ApplicationController) GetGroupProjects(ctx *gin.Context) []interface{} {
	// Get the group from the context
	group := c.GetGroup(ctx)

	// Get the current user from the context
	currentUser, _ := ctx.Get("current_user")

	// Find the group projects
	projects := c.groupProjectsFinder.Execute(group, currentUser, false)

	return projects
}

// GetGroupProjectsWithSubgroups gets the group projects with subgroups
func (c *ApplicationController) GetGroupProjectsWithSubgroups(ctx *gin.Context) []interface{} {
	// Get the group from the context
	group := c.GetGroup(ctx)

	// Get the current user from the context
	currentUser, _ := ctx.Get("current_user")

	// Find the group projects with subgroups
	projects := c.groupProjectsFinder.Execute(group, currentUser, true)

	return projects
}

// AuthorizeAdminGroup checks if the user has permission to admin the group
func (c *ApplicationController) AuthorizeAdminGroup(ctx *gin.Context) bool {
	// Get the current user from the context
	currentUser, _ := ctx.Get("current_user")

	// Get the group from the context
	group := c.GetGroup(ctx)

	// Check if the user has permission to admin the group
	return c.authorizationService.Can(currentUser, "admin_group", group)
}

// AuthorizeCreateDeployToken checks if the user has permission to create deploy tokens
func (c *ApplicationController) AuthorizeCreateDeployToken(ctx *gin.Context) bool {
	// Get the current user from the context
	currentUser, _ := ctx.Get("current_user")

	// Get the group from the context
	group := c.GetGroup(ctx)

	// Check if the user has permission to create deploy tokens
	return c.authorizationService.Can(currentUser, "create_deploy_token", group)
}

// AuthorizeDestroyDeployToken checks if the user has permission to destroy deploy tokens
func (c *ApplicationController) AuthorizeDestroyDeployToken(ctx *gin.Context) bool {
	// Get the current user from the context
	currentUser, _ := ctx.Get("current_user")

	// Get the group from the context
	group := c.GetGroup(ctx)

	// Check if the user has permission to destroy deploy tokens
	return c.authorizationService.Can(currentUser, "destroy_deploy_token", group)
}

// AuthorizeAdminGroupMember checks if the user has permission to admin group members
func (c *ApplicationController) AuthorizeAdminGroupMember(ctx *gin.Context) bool {
	// Get the current user from the context
	currentUser, _ := ctx.Get("current_user")

	// Get the group from the context
	group := c.GetGroup(ctx)

	// Check if the user has permission to admin group members
	return c.authorizationService.Can(currentUser, "admin_group_member", group)
}

// AuthorizeOwnerAccess checks if the user has permission to access owner features
func (c *ApplicationController) AuthorizeOwnerAccess(ctx *gin.Context) bool {
	// Get the current user from the context
	currentUser, _ := ctx.Get("current_user")

	// Get the group from the context
	group := c.GetGroup(ctx)

	// Check if the user has permission to access owner features
	return c.authorizationService.Can(currentUser, "owner_access", group)
}

// AuthorizeBillingsPage checks if the user has permission to view the billings page
func (c *ApplicationController) AuthorizeBillingsPage(ctx *gin.Context) bool {
	// Get the current user from the context
	currentUser, _ := ctx.Get("current_user")

	// Get the group from the context
	group := c.GetGroup(ctx)

	// Check if the user has permission to view the billings page
	return c.authorizationService.Can(currentUser, "read_billing", group)
}

// AuthorizeReadGroupMember checks if the user has permission to read group members
func (c *ApplicationController) AuthorizeReadGroupMember(ctx *gin.Context) bool {
	// Get the current user from the context
	currentUser, _ := ctx.Get("current_user")

	// Get the group from the context
	group := c.GetGroup(ctx)

	// Check if the user has permission to read group members
	return c.authorizationService.Can(currentUser, "read_group_member", group)
}

// BuildCanonicalPath builds the canonical path for a group
func (c *ApplicationController) BuildCanonicalPath(ctx *gin.Context, group interface{}) string {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would build the canonical path for a group
	return ""
}

// HasProjectList checks if the controller has a project list
func (c *ApplicationController) HasProjectList() bool {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would check if the controller has a project list
	return false
}

// ValidateCrmGroup checks if the group is a CRM group
func (c *ApplicationController) ValidateCrmGroup(ctx *gin.Context) bool {
	// Get the group from the context
	group := c.GetGroup(ctx)

	// Check if the group is a CRM group
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would check if the group is a CRM group
	return true
}

// AuthorizeAction checks if the user has permission to perform an action
func (c *ApplicationController) AuthorizeAction(ctx *gin.Context, action string) bool {
	// Get the current user from the context
	currentUser, _ := ctx.Get("current_user")

	// Get the group from the context
	group := c.GetGroup(ctx)

	// Check if the user has permission to perform the action
	return c.authorizationService.Can(currentUser, action, group)
}

// HandleAuthorization dynamically handles authorization methods
func (c *ApplicationController) HandleAuthorization(ctx *gin.Context, method string) bool {
	// Check if the method is an authorization method
	authRegex := regexp.MustCompile(`^authorize_(.*)!$`)
	matches := authRegex.FindStringSubmatch(method)

	if len(matches) > 1 {
		// Get the action from the method name
		action := matches[1]

		// Check if the user has permission to perform the action
		return c.AuthorizeAction(ctx, action)
	}

	// Return false for unknown methods
	return false
}

// GroupFinder finds groups
type GroupFinder struct {
	// Add any dependencies here
}

// FindRoutable finds a routable group
func (f *GroupFinder) FindRoutable(groupID string, path string) (interface{}, error) {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would find a routable group
	return nil, nil
}

// GroupProjectsFinder finds group projects
type GroupProjectsFinder struct {
	// Add any dependencies here
}

// Execute executes the finder
func (f *GroupProjectsFinder) Execute(group interface{}, currentUser interface{}, includeSubgroups bool) []interface{} {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would find group projects
	return []interface{}{}
}

// AuthorizationService authorizes actions
type AuthorizationService struct {
	// Add any dependencies here
}

// Can checks if a user can perform an action
func (s *AuthorizationService) Can(user interface{}, action string, subject interface{}) bool {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would check if a user can perform an action
	return true
}

// SortingService handles sorting
type SortingService struct {
	// Add any dependencies here
}

// GetSortValue gets the sort value
func (s *SortingService) GetSortValue(ctx *gin.Context) string {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would get the sort value
	return ""
}
