package clusters

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// BaseController provides common functionality for cluster-related controllers
type BaseController struct {
	clusterService *service.ClusterService
	authService   *service.AuthService
}

// NewBaseController creates a new instance of BaseController
func NewBaseController(clusterService *service.ClusterService, authService *service.AuthService) *BaseController {
	return &BaseController{
		clusterService: clusterService,
		authService:   authService,
	}
}

// RegisterRoutes registers the base routes for cluster controllers
func (c *BaseController) RegisterRoutes(r *gin.RouterGroup) {
	clusters := r.Group("/clusters")
	{
		// Apply middleware
		clusters.Use(c.clusterable)
		clusters.Use(c.authorizeAdminCluster)

		// Register common routes
		clusters.GET("/", c.index)
		clusters.GET("/:id", c.show)
		clusters.GET("/:id/environments", c.environments)
		clusters.GET("/:id/status", c.clusterStatus)
		clusters.DELETE("/:id", c.destroy)
		clusters.GET("/docs", c.newClusterDocs)
		clusters.POST("/:id/connect", c.connect)
		clusters.GET("/new", c.new)
		clusters.POST("/:id/create_user", c.createUser)
	}
}

// clusterable middleware ensures the clusterable object is available
func (c *BaseController) clusterable(ctx *gin.Context) {
	// This should be implemented by child controllers
	// It will set the clusterable object in the context
	ctx.Set("clusterable", nil)
}

// authorizeAdminCluster middleware checks admin cluster permissions
func (c *BaseController) authorizeAdminCluster(ctx *gin.Context) {
	// Skip for specific routes that don't require admin access
	if c.isPublicRoute(ctx.Request.URL.Path) {
		return
	}

	user, exists := ctx.Get("current_user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		ctx.Abort()
		return
	}

	clusterable, exists := ctx.Get("clusterable")
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Clusterable not found"})
		ctx.Abort()
		return
	}

	if !c.authService.CanAdminCluster(ctx, user.(*model.User), clusterable) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		ctx.Abort()
		return
	}
}

// isPublicRoute checks if the route is public (doesn't require admin access)
func (c *BaseController) isPublicRoute(path string) bool {
	publicRoutes := []string{
		"/clusters/:id",
		"/clusters",
		"/clusters/new",
		"/clusters/:id/authorize_aws_role",
		"/clusters/:id/update",
	}
	for _, route := range publicRoutes {
		if route == path {
			return true
		}
	}
	return false
}

// Helper methods for authorization checks
func (c *BaseController) authorizeUpdateCluster(ctx *gin.Context) error {
	user, exists := ctx.Get("current_user")
	if !exists {
		return fmt.Errorf("user not found")
	}

	clusterable, exists := ctx.Get("clusterable")
	if !exists {
		return fmt.Errorf("clusterable not found")
	}

	if !c.authService.CanUpdateCluster(ctx, user.(*model.User), clusterable) {
		return fmt.Errorf("access denied")
	}

	return nil
}

func (c *BaseController) authorizeReadCluster(ctx *gin.Context) error {
	user, exists := ctx.Get("current_user")
	if !exists {
		return fmt.Errorf("user not found")
	}

	clusterable, exists := ctx.Get("clusterable")
	if !exists {
		return fmt.Errorf("clusterable not found")
	}

	if !c.authService.CanReadCluster(ctx, user.(*model.User), clusterable) {
		return fmt.Errorf("access denied")
	}

	return nil
}

func (c *BaseController) authorizeCreateCluster(ctx *gin.Context) error {
	user, exists := ctx.Get("current_user")
	if !exists {
		return fmt.Errorf("user not found")
	}

	clusterable, exists := ctx.Get("clusterable")
	if !exists {
		return fmt.Errorf("clusterable not found")
	}

	if !c.authService.CanCreateCluster(ctx, user.(*model.User), clusterable) {
		return fmt.Errorf("access denied")
	}

	return nil
}

func (c *BaseController) authorizeReadPrometheus(ctx *gin.Context) error {
	user, exists := ctx.Get("current_user")
	if !exists {
		return fmt.Errorf("user not found")
	}

	clusterable, exists := ctx.Get("clusterable")
	if !exists {
		return fmt.Errorf("clusterable not found")
	}

	if !c.authService.CanReadPrometheus(ctx, user.(*model.User), clusterable) {
		return fmt.Errorf("access denied")
	}

	return nil
}

// Route handlers
func (c *BaseController) index(ctx *gin.Context) {
	// TODO: Implement index view
}

func (c *BaseController) show(ctx *gin.Context) {
	// TODO: Implement show view
}

func (c *BaseController) environments(ctx *gin.Context) {
	// TODO: Implement environments view
}

func (c *BaseController) clusterStatus(ctx *gin.Context) {
	// TODO: Implement cluster status view
}

func (c *BaseController) destroy(ctx *gin.Context) {
	// TODO: Implement cluster deletion
}

func (c *BaseController) newClusterDocs(ctx *gin.Context) {
	// TODO: Implement cluster docs view
}

func (c *BaseController) connect(ctx *gin.Context) {
	// TODO: Implement cluster connection
}

func (c *BaseController) new(ctx *gin.Context) {
	// TODO: Implement new cluster view
}

func (c *BaseController) createUser(ctx *gin.Context) {
	// TODO: Implement user creation
}
