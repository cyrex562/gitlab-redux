package clusters

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

const (
	statusPollingInterval = 10 * time.Second
	clustersPerPage      = 20
)

// ClustersController handles cluster management operations
type ClustersController struct {
	*BaseController
	clusterService *service.ClusterService
	buildService   *service.BuildService
	updateService  *service.UpdateService
	destroyService *service.DestroyService
	migrationService *service.MigrationService
}

// NewClustersController creates a new instance of ClustersController
func NewClustersController(
	baseController *BaseController,
	clusterService *service.ClusterService,
	buildService *service.BuildService,
	updateService *service.UpdateService,
	destroyService *service.DestroyService,
	migrationService *service.MigrationService,
) *ClustersController {
	return &ClustersController{
		BaseController:    baseController,
		clusterService:    clusterService,
		buildService:      buildService,
		updateService:     updateService,
		destroyService:    destroyService,
		migrationService:  migrationService,
	}
}

// RegisterRoutes registers the routes for the ClustersController
func (c *ClustersController) RegisterRoutes(r *gin.RouterGroup) {
	clusters := r.Group("/clusters")
	{
		// Apply middleware
		clusters.Use(c.ensureFeatureEnabled)

		// Register routes
		clusters.GET("/", c.index)
		clusters.GET("/:id", c.show)
		clusters.GET("/:id/status", c.clusterStatus)
		clusters.PUT("/:id", c.update)
		clusters.DELETE("/:id", c.destroy)
		clusters.POST("/:id/clear_cache", c.clearCache)
		clusters.POST("/:id/migrate", c.migrate)
		clusters.PUT("/:id/migration", c.updateMigration)
		clusters.POST("/create_user", c.createUser)
		clusters.POST("/connect", c.connect)
	}
}

// index handles the GET /clusters endpoint
func (c *ClustersController) index(ctx *gin.Context) {
	// Get clusterable from context
	clusterable, exists := ctx.Get("clusterable")
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Clusterable not found"})
		return
	}

	// Get clusters list
	clusters, hasAncestorClusters, err := c.clusterService.ListClusters(ctx, clusterable)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list clusters"})
		return
	}

	// Handle different response formats
	if ctx.GetHeader("Accept") == "application/json" {
		// Set polling interval header
		ctx.Header("Poll-Interval", strconv.Itoa(int(statusPollingInterval.Milliseconds())))

		// Get pagination parameters
		page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
		perPage := clustersPerPage

		// Paginate clusters
		paginatedClusters := c.paginateClusters(clusters, page, perPage)

		ctx.JSON(http.StatusOK, gin.H{
			"clusters":            paginatedClusters,
			"has_ancestor_clusters": hasAncestorClusters,
		})
	} else {
		// Render HTML view
		ctx.HTML(http.StatusOK, "clusters/index.html", gin.H{
			"clusters":            clusters,
			"has_ancestor_clusters": hasAncestorClusters,
		})
	}
}

// clusterStatus handles the GET /clusters/:id/status endpoint
func (c *ClustersController) clusterStatus(ctx *gin.Context) {
	cluster, exists := ctx.Get("cluster")
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	// Set polling interval header
	ctx.Header("Poll-Interval", strconv.Itoa(int(statusPollingInterval.Milliseconds())))

	status, err := c.clusterService.GetClusterStatus(ctx, cluster.(*model.Cluster))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cluster status"})
		return
	}

	ctx.JSON(http.StatusOK, status)
}

// show handles the GET /clusters/:id endpoint
func (c *ClustersController) show(ctx *gin.Context) {
	cluster, exists := ctx.Get("cluster")
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	// Handle integrations tab
	if ctx.Query("tab") == "integrations" {
		prometheusIntegration, err := c.clusterService.GetPrometheusIntegration(ctx, cluster.(*model.Cluster))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get prometheus integration"})
			return
		}
		ctx.Set("prometheus_integration", prometheusIntegration)
	}

	// Render show view
	ctx.HTML(http.StatusOK, "clusters/show.html", gin.H{
		"cluster": cluster,
	})
}

// update handles the PUT /clusters/:id endpoint
func (c *ClustersController) update(ctx *gin.Context) {
	cluster, exists := ctx.Get("cluster")
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	// Parse update parameters
	var params model.ClusterParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters"})
		return
	}

	// Update cluster
	updatedCluster, err := c.updateService.UpdateCluster(ctx, cluster.(*model.Cluster), &params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update cluster"})
		return
	}

	// Handle different response formats
	if ctx.GetHeader("Accept") == "application/json" {
		ctx.Status(http.StatusNoContent)
	} else {
		ctx.Redirect(http.StatusFound, updatedCluster.ShowPath())
	}
}

// destroy handles the DELETE /clusters/:id endpoint
func (c *ClustersController) destroy(ctx *gin.Context) {
	cluster, exists := ctx.Get("cluster")
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	// Parse destroy parameters
	var params model.DestroyParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters"})
		return
	}

	// Destroy cluster
	response, err := c.destroyService.DestroyCluster(ctx, cluster.(*model.Cluster), &params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to destroy cluster"})
		return
	}

	// Redirect to index
	ctx.Redirect(http.StatusFound, response.RedirectPath)
}

// clearCache handles the POST /clusters/:id/clear_cache endpoint
func (c *ClustersController) clearCache(ctx *gin.Context) {
	cluster, exists := ctx.Get("cluster")
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	if err := c.clusterService.ClearCache(ctx, cluster.(*model.Cluster)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cache"})
		return
	}

	ctx.Redirect(http.StatusFound, cluster.(*model.Cluster).ShowPath())
}

// migrate handles the POST /clusters/:id/migrate endpoint
func (c *ClustersController) migrate(ctx *gin.Context) {
	cluster, exists := ctx.Get("cluster")
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	// Parse migration parameters
	var params model.MigrationParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters"})
		return
	}

	// Start migration
	response, err := c.migrationService.CreateMigration(ctx, cluster.(*model.Cluster), &params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start migration"})
		return
	}

	// Redirect to show page with migrate tab
	ctx.Redirect(http.StatusFound, cluster.(*model.Cluster).ShowPath()+"?tab=migrate")
}

// updateMigration handles the PUT /clusters/:id/migration endpoint
func (c *ClustersController) updateMigration(ctx *gin.Context) {
	cluster, exists := ctx.Get("cluster")
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Cluster not found"})
		return
	}

	// Parse migration update parameters
	var params model.MigrationUpdateParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters"})
		return
	}

	// Update migration
	response, err := c.migrationService.UpdateMigration(ctx, cluster.(*model.Cluster), &params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update migration"})
		return
	}

	// Redirect to show page with migrate tab
	ctx.Redirect(http.StatusFound, cluster.(*model.Cluster).ShowPath()+"?tab=migrate")
}

// createUser handles the POST /clusters/create_user endpoint
func (c *ClustersController) createUser(ctx *gin.Context) {
	// Parse create user parameters
	var params model.CreateUserClusterParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parameters"})
		return
	}

	// Get token from session
	token := ctx.GetString("google_api_token")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No token found"})
		return
	}

	// Create user cluster
	userCluster, err := c.clusterService.CreateUserCluster(ctx, &params, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user cluster"})
		return
	}

	if userCluster.IsPersisted() {
		ctx.Redirect(http.StatusFound, userCluster.ShowPath())
	} else {
		ctx.HTML(http.StatusOK, "clusters/connect.html", gin.H{
			"cluster": userCluster,
		})
	}
}

// connect handles the POST /clusters/connect endpoint
func (c *ClustersController) connect(ctx *gin.Context) {
	// Build user cluster
	userCluster, err := c.buildService.BuildUserCluster(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to build user cluster"})
		return
	}

	ctx.HTML(http.StatusOK, "clusters/connect.html", gin.H{
		"cluster": userCluster,
	})
}

// ensureFeatureEnabled middleware checks if certificate-based clusters are enabled
func (c *ClustersController) ensureFeatureEnabled(ctx *gin.Context) {
	clusterable, exists := ctx.Get("clusterable")
	if !exists {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Clusterable not found"})
		ctx.Abort()
		return
	}

	if !c.clusterService.AreCertificateBasedClustersEnabled(ctx, clusterable) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Feature not enabled"})
		ctx.Abort()
		return
	}
}

// Helper methods
func (c *ClustersController) paginateClusters(clusters []*model.Cluster, page, perPage int) []*model.Cluster {
	start := (page - 1) * perPage
	end := start + perPage
	if end > len(clusters) {
		end = len(clusters)
	}
	return clusters[start:end]
}
