package groups

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ChildrenController handles listing child groups and their descendants
type ChildrenController struct {
	// Embed the ApplicationController to inherit its functionality
	*ApplicationController

	// Add any additional dependencies here
	groupFinder *GroupFinder
	groupDescendantsFinder *GroupDescendantsFinder
	groupChildSerializer *GroupChildSerializer
}

// NewChildrenController creates a new ChildrenController
func NewChildrenController(
	applicationController *ApplicationController,
	groupFinder *GroupFinder,
	groupDescendantsFinder *GroupDescendantsFinder,
	groupChildSerializer *GroupChildSerializer,
) *ChildrenController {
	return &ChildrenController{
		ApplicationController: applicationController,
		groupFinder:           groupFinder,
		groupDescendantsFinder: groupDescendantsFinder,
		groupChildSerializer:   groupChildSerializer,
	}
}

// RegisterRoutes registers the routes for the ChildrenController
func (c *ChildrenController) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/", c.Index)
}

// RegisterMiddleware registers the middleware for the ChildrenController
func (c *ChildrenController) RegisterMiddleware(router *gin.RouterGroup) {
	// Register the middleware from the ApplicationController
	c.ApplicationController.RegisterMiddleware(router)

	// Add group middleware
	router.Use(c.SetGroup())

	// Add validate per page middleware
	router.Use(c.ValidatePerPage())

	// Skip cross project access check for index action
	router.Use(c.SkipCrossProjectAccessCheck("index"))
}

// Index handles the index action
func (c *ChildrenController) Index(ctx *gin.Context) {
	// Get the sort parameter or use the default from group projects sort
	sort := ctx.DefaultQuery("sort", c.GetGroupProjectsSort(ctx))

	// Get the parent group
	var parent interface{}
	if parentID := ctx.Query("parent_id"); parentID != "" {
		// Find the parent group by ID
		parent = c.groupFinder.Execute(c.GetCurrentUser(ctx), map[string]interface{}{
			"id": parentID,
		})
	} else {
		// Use the current group
		parent = c.GetGroup(ctx)
	}

	// Check if the parent group exists
	if parent == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	}

	// Set up the children
	children := c.setupChildren(parent, ctx)

	// Serialize the children
	serializer := c.groupChildSerializer.New(c.GetCurrentUser(ctx))
	serializer.WithPagination(ctx.Request, ctx.Writer)

	// Expand the hierarchy if the filter parameter is present
	if ctx.Query("filter") != "" {
		serializer.ExpandHierarchy(parent)
	}

	// Render the JSON response
	ctx.JSON(http.StatusOK, serializer.Represent(children))
}

// setupChildren sets up the children for the parent group
func (c *ChildrenController) setupChildren(parent interface{}, ctx *gin.Context) interface{} {
	return c.groupDescendantsFinder.New(
		c.GetCurrentUser(ctx),
		parent,
		c.getGroupDescendantsParams(ctx),
	).Execute()
}

// getGroupDescendantsParams gets the parameters for the group descendants finder
func (c *ChildrenController) getGroupDescendantsParams(ctx *gin.Context) map[string]interface{} {
	// Get all parameters from the request
	params := make(map[string]interface{})
	for k, v := range ctx.Request.URL.Query() {
		if len(v) > 0 {
			params[k] = v[0]
		}
	}
	return params
}

// ValidatePerPage middleware validates the per_page parameter
func (c *ChildrenController) ValidatePerPage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the per_page parameter
		perPageStr := ctx.Query("per_page")
		if perPageStr == "" {
			ctx.Next()
			return
		}

		// Parse the per_page parameter
		perPage, err := strconv.Atoi(perPageStr)
		if err != nil || perPage < 1 {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "per_page does not have a valid value"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// SkipCrossProjectAccessCheck middleware skips the cross project access check for the specified action
func (c *ChildrenController) SkipCrossProjectAccessCheck(action string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Skip the cross project access check for the specified action
		ctx.Set("skip_cross_project_access_check", action)
		ctx.Next()
	}
}

// HasProjectList returns true if the controller has a project list
func (c *ChildrenController) HasProjectList() bool {
	return true
}

// GroupFinder finds groups
type GroupFinder struct {
	// Add any dependencies here
}

// Execute executes the finder
func (f *GroupFinder) Execute(currentUser interface{}, params map[string]interface{}) interface{} {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would find a group
	return nil
}

// GroupDescendantsFinder finds group descendants
type GroupDescendantsFinder struct {
	// Add any dependencies here
}

// New creates a new group descendants finder
func (f *GroupDescendantsFinder) New(currentUser interface{}, parentGroup interface{}, params map[string]interface{}) *GroupDescendantsFinder {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would create a new group descendants finder
	return &GroupDescendantsFinder{}
}

// Execute executes the finder
func (f *GroupDescendantsFinder) Execute() interface{} {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would find group descendants
	return nil
}

// GroupChildSerializer serializes group children
type GroupChildSerializer struct {
	// Add any dependencies here
	currentUser interface{}
	pagination  interface{}
	expanded    bool
}

// New creates a new group child serializer
func (s *GroupChildSerializer) New(currentUser interface{}) *GroupChildSerializer {
	return &GroupChildSerializer{
		currentUser: currentUser,
	}
}

// WithPagination adds pagination to the serializer
func (s *GroupChildSerializer) WithPagination(request interface{}, response interface{}) *GroupChildSerializer {
	s.pagination = map[string]interface{}{
		"request":  request,
		"response": response,
	}
	return s
}

// ExpandHierarchy expands the hierarchy
func (s *GroupChildSerializer) ExpandHierarchy(parent interface{}) *GroupChildSerializer {
	s.expanded = true
	return s
}

// Represent represents the children
func (s *GroupChildSerializer) Represent(children interface{}) interface{} {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would represent the children
	return children
}
