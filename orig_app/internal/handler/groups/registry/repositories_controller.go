package registry

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RepositoriesController handles container registry repositories
type RepositoriesController struct {
	containerRepositoriesFinder *ContainerRepositoriesFinder
	containerRepositoriesSerializer *ContainerRepositoriesSerializer
	configService *ConfigService
	featureFlagService *FeatureFlagService
	packageEventTracker *PackageEventTracker
}

// NewRepositoriesController creates a new RepositoriesController
func NewRepositoriesController(
	containerRepositoriesFinder *ContainerRepositoriesFinder,
	containerRepositoriesSerializer *ContainerRepositoriesSerializer,
	configService *ConfigService,
	featureFlagService *FeatureFlagService,
	packageEventTracker *PackageEventTracker,
) *RepositoriesController {
	return &RepositoriesController{
		containerRepositoriesFinder: containerRepositoriesFinder,
		containerRepositoriesSerializer: containerRepositoriesSerializer,
		configService: configService,
		featureFlagService: featureFlagService,
		packageEventTracker: packageEventTracker,
	}
}

// RegisterRoutes registers the routes for the RepositoriesController
func (c *RepositoriesController) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/", c.Index)
	router.GET("/:id", c.Show)
}

// RegisterMiddleware registers the middleware for the RepositoriesController
func (c *RepositoriesController) RegisterMiddleware(router *gin.RouterGroup) {
	router.Use(c.VerifyContainerRegistryEnabled())
	router.Use(c.AuthorizeReadContainerImage())

	// Register middleware for index and show actions
	router.GET("/", c.PushFrontendFeatureFlag("show_container_registry_tag_signatures"))
	router.GET("/:id", c.PushFrontendFeatureFlag("show_container_registry_tag_signatures"))
}

// Index handles the index action
func (c *RepositoriesController) Index(ctx *gin.Context) {
	// Check the Accept header to determine the response format
	acceptHeader := ctx.GetHeader("Accept")

	if acceptHeader == "application/json" {
		// Get the current user from the context
		currentUser, _ := ctx.Get("current_user")

		// Get the group from the context
		group, _ := ctx.Get("group")

		// Get the name parameter from the query
		name := ctx.Query("name")

		// Find the container repositories
		images, err := c.containerRepositoriesFinder.Execute(currentUser, group, map[string]string{"name": name})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Add API entity associations
		images = images.WithAPIEntityAssociations()

		// Track the package event
		c.packageEventTracker.TrackPackageEvent("list_repositories", "container", currentUser, group)

		// Serialize the images
		serializer := c.containerRepositoriesSerializer.New(currentUser)

		// Render the JSON response
		ctx.JSON(http.StatusOK, serializer.WithPagination(ctx.Request, ctx.Writer).RepresentReadOnly(images))
	} else {
		// Render the HTML response
		ctx.HTML(http.StatusOK, "groups/registry/repositories/index", gin.H{})
	}
}

// Show handles the show action
func (c *RepositoriesController) Show(ctx *gin.Context) {
	// Render the index template to allow frontend routing to work on page refresh
	ctx.HTML(http.StatusOK, "groups/registry/repositories/index", gin.H{})
}

// VerifyContainerRegistryEnabled middleware verifies that the container registry is enabled
func (c *RepositoriesController) VerifyContainerRegistryEnabled() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !c.configService.IsRegistryEnabled() {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Container registry is not enabled"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

// AuthorizeReadContainerImage middleware checks if the user has permission to read container images
func (c *RepositoriesController) AuthorizeReadContainerImage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the current user from the context
		currentUser, exists := ctx.Get("current_user")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			ctx.Abort()
			return
		}

		// Get the group from the context
		group, exists := ctx.Get("group")
		if !exists {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
			ctx.Abort()
			return
		}

		// Check if the user has permission to read container images
		// This is a placeholder for the actual implementation
		// In the actual implementation, you would check if the user has the read_container_image permission
		canReadContainerImage := true // Replace with actual check

		if !canReadContainerImage {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "You don't have permission to read container images"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// PushFrontendFeatureFlag middleware pushes a frontend feature flag
func (c *RepositoriesController) PushFrontendFeatureFlag(flagName string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the group from the context
		group, exists := ctx.Get("group")
		if !exists {
			ctx.Next()
			return
		}

		// Push the frontend feature flag
		c.featureFlagService.PushFrontendFeatureFlag(flagName, group)

		ctx.Next()
	}
}

// ContainerRepositoriesFinder finds container repositories
type ContainerRepositoriesFinder struct {
	// Add fields as needed
}

// Execute executes the finder
func (f *ContainerRepositoriesFinder) Execute(user interface{}, subject interface{}, params map[string]string) (interface{}, error) {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would find the container repositories
	return nil, nil
}

// ContainerRepositoriesSerializer serializes container repositories
type ContainerRepositoriesSerializer struct {
	// Add fields as needed
}

// New creates a new serializer
func (s *ContainerRepositoriesSerializer) New(currentUser interface{}) interface{} {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would create a new serializer
	return nil
}

// WithPagination adds pagination to the serializer
func (s *ContainerRepositoriesSerializer) WithPagination(request interface{}, response interface{}) interface{} {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would add pagination to the serializer
	return s
}

// RepresentReadOnly represents the repositories as read-only
func (s *ContainerRepositoriesSerializer) RepresentReadOnly(repositories interface{}) interface{} {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would represent the repositories as read-only
	return nil
}

// ConfigService provides configuration-related functionality
type ConfigService struct {
	// Add fields as needed
}

// IsRegistryEnabled checks if the registry is enabled
func (s *ConfigService) IsRegistryEnabled() bool {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would check if the registry is enabled
	return true
}

// FeatureFlagService provides feature flag-related functionality
type FeatureFlagService struct {
	// Add fields as needed
}

// PushFrontendFeatureFlag pushes a frontend feature flag
func (s *FeatureFlagService) PushFrontendFeatureFlag(flagName string, group interface{}) {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would push a frontend feature flag
}

// PackageEventTracker tracks package events
type PackageEventTracker struct {
	// Add fields as needed
}

// TrackPackageEvent tracks a package event
func (t *PackageEventTracker) TrackPackageEvent(eventName string, packageType string, user interface{}, namespace interface{}) {
	// This is a placeholder for the actual implementation
	// In the actual implementation, you would track a package event
}
