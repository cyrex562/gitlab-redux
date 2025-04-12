package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/handler/clusters"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

type AdminClustersController struct {
	*clusters.ClustersController
	featureService *service.FeatureService
}

func NewAdminClustersController(
	clusterService *service.ClusterService,
	featureService *service.FeatureService,
	user *model.User,
) *AdminClustersController {
	return &AdminClustersController{
		ClustersController: clusters.NewClustersController(clusterService, user),
		featureService:     featureService,
	}
}

// EnsureFeatureEnabled middleware ensures the clusters feature is enabled
func (c *AdminClustersController) EnsureFeatureEnabled() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		enabled, err := c.featureService.IsEnabled(ctx, "clusters")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			ctx.Abort()
			return
		}

		if !enabled {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Clusters feature is not enabled"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// getClusterable returns the instance clusterable for admin operations
func (c *AdminClustersController) getClusterable(ctx *gin.Context) clusters.Clusterable {
	instance := &model.Instance{}
	return &InstanceClusterable{
		instance: instance,
		user:     c.user,
	}
}

// InstanceClusterable represents an instance that can have clusters
type InstanceClusterable struct {
	instance *model.Instance
	user     *model.User
}

func (i *InstanceClusterable) GetID() int64 {
	return i.instance.ID
}

func (i *InstanceClusterable) GetName() string {
	return "Instance"
}

func (i *InstanceClusterable) GetType() string {
	return "instance"
}
