package groups

import (
	"github.com/gin-gonic/gin"
)

// BoardsController handles group boards
type BoardsController struct {
	// Embed the ApplicationController to inherit its functionality
	*ApplicationController

	// Add any additional dependencies here
	boardsActionsService *BoardsActionsService
	userActivityService  *RecordUserLastActivityService
	boardsFinder        *BoardsFinder
	boardsCreateService *BoardsCreateService
	authorizationService *AuthorizationService
	featureFlagService  *FeatureFlagService
}

// NewBoardsController creates a new BoardsController
func NewBoardsController(
	applicationController *ApplicationController,
	boardsActionsService *BoardsActionsService,
	userActivityService *RecordUserLastActivityService,
	boardsFinder *BoardsFinder,
	boardsCreateService *BoardsCreateService,
	authorizationService *AuthorizationService,
	featureFlagService *FeatureFlagService,
) *BoardsController {
	return &BoardsController{
		ApplicationController: applicationController,
		boardsActionsService:  boardsActionsService,
		userActivityService:   userActivityService,
		boardsFinder:         boardsFinder,
		boardsCreateService:  boardsCreateService,
		authorizationService: authorizationService,
		featureFlagService:   featureFlagService,
	}
}

// RegisterMiddleware registers the middleware for the BoardsController
func (c *BoardsController) RegisterMiddleware(router *gin.RouterGroup) {
	// Register the middleware from the ApplicationController
	c.ApplicationController.RegisterMiddleware(router)

	// Add feature flag middleware
	router.Use(c.PushFeatureFlags())

	// Add user activity tracking middleware
	router.Use(c.userActivityService.RecordUserLastActivity())
}

// PushFeatureFlags middleware pushes feature flags to the frontend
func (c *BoardsController) PushFeatureFlags() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the group from the context
		group := c.GetGroup(ctx)

		// Push the board multi select feature flag
		c.featureFlagService.PushFrontendFeatureFlag(ctx, "board_multi_select", group)

		// Push the issues list drawer feature flag
		c.featureFlagService.PushFrontendFeatureFlag(ctx, "issues_list_drawer", group)

		// Push the work items beta feature flag
		workItemsBetaEnabled := false
		if group != nil {
			workItemsBetaEnabled = c.featureFlagService.IsWorkItemsBetaEnabled(group)
		}
		c.featureFlagService.PushForceFrontendFeatureFlag(ctx, "work_items_beta", workItemsBetaEnabled)

		ctx.Next()
	}
}

// GetBoardFinder gets the board finder
func (c *BoardsController) GetBoardFinder(ctx *gin.Context) *BoardsFinder {
	// Get the group from the context
	group := c.GetGroup(ctx)

	// Get the current user from the context
	currentUser, _ := ctx.Get("current_user")

	// Get the board ID from the URL parameters
	boardID := ctx.Param("id")

	// Create a new board finder
	return c.boardsFinder.New(group, currentUser, boardID)
}

// GetBoardCreateService gets the board create service
func (c *BoardsController) GetBoardCreateService(ctx *gin.Context) *BoardsCreateService {
	// Get the group from the context
	group := c.GetGroup(ctx)

	// Get the current user from the context
	currentUser, _ := ctx.Get("current_user")

	// Create a new board create service
	return c.boardsCreateService.New(group, currentUser)
}

// AuthorizeReadBoard checks if the user has permission to read the board
func (c *BoardsController) AuthorizeReadBoard(ctx *gin.Context) bool {
	// Get the current user from the context
	currentUser, _ := ctx.Get("current_user")

	// Get the group from the context
	group := c.GetGroup(ctx)

	// Check if the user has permission to read the board
	return c.authorizationService.Can(currentUser, "read_issue_board", group)
}

// BoardsActionsService handles board actions
type BoardsActionsService struct {
	// Add any dependencies here
}

// BoardsFinder finds boards
type BoardsFinder struct {
	// Add any dependencies here
}

// New creates a new board finder
func (f *BoardsFinder) New(parent interface{}, currentUser interface{}, boardID string) *BoardsFinder {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would create a new board finder
	return &BoardsFinder{}
}

// BoardsCreateService creates boards
type BoardsCreateService struct {
	// Add any dependencies here
}

// New creates a new board create service
func (s *BoardsCreateService) New(parent interface{}, currentUser interface{}) *BoardsCreateService {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would create a new board create service
	return &BoardsCreateService{}
}

// RecordUserLastActivityService records user last activity
type RecordUserLastActivityService struct {
	// Add any dependencies here
}

// RecordUserLastActivity records the user's last activity
func (s *RecordUserLastActivityService) RecordUserLastActivity() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// This is a placeholder for the actual implementation
		// In the actual implementation, you would record the user's last activity

		ctx.Next()
	}
}

// FeatureFlagService handles feature flags
type FeatureFlagService struct {
	// Add any dependencies here
}

// PushFrontendFeatureFlag pushes a frontend feature flag
func (s *FeatureFlagService) PushFrontendFeatureFlag(ctx *gin.Context, flag string, subject interface{}) {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would push a frontend feature flag
}

// PushForceFrontendFeatureFlag pushes a forced frontend feature flag
func (s *FeatureFlagService) PushForceFrontendFeatureFlag(ctx *gin.Context, flag string, enabled bool) {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would push a forced frontend feature flag
}

// IsWorkItemsBetaEnabled checks if work items beta is enabled
func (s *FeatureFlagService) IsWorkItemsBetaEnabled(group interface{}) bool {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would check if work items beta is enabled
	return false
}
