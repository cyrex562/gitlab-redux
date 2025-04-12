package agents

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// DashboardController handles the Kubernetes dashboard interface for cluster agents
type DashboardController struct {
	agentService *service.AgentService
	featureService *service.FeatureService
}

// NewDashboardController creates a new instance of DashboardController
func NewDashboardController(agentService *service.AgentService, featureService *service.FeatureService) *DashboardController {
	return &DashboardController{
		agentService: agentService,
		featureService: featureService,
	}
}

// RegisterRoutes registers the routes for the DashboardController
func (c *DashboardController) RegisterRoutes(r *gin.RouterGroup) {
	dashboard := r.Group("/clusters/agents/dashboard")
	{
		dashboard.GET("/", c.index)
		dashboard.GET("/:agent_id", c.show)
	}
}

// index handles the GET /clusters/agents/dashboard endpoint
func (c *DashboardController) index(ctx *gin.Context) {
	// TODO: Implement index view rendering
	ctx.HTML(http.StatusOK, "clusters/agents/dashboard/index.html", gin.H{
		"title": "Cluster Agents Dashboard",
	})
}

// show handles the GET /clusters/agents/dashboard/:agent_id endpoint
func (c *DashboardController) show(ctx *gin.Context) {
	// Check feature flag
	if err := c.checkFeatureFlag(ctx); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Feature not enabled"})
		return
	}

	// Find the agent
	agent, err := c.findAgent(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Agent not found"})
		return
	}

	// Check authorization
	if err := c.authorizeReadClusterAgent(ctx, agent); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Set KAS cookie
	if err := c.setKasCookie(ctx); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set KAS cookie"})
		return
	}

	// TODO: Implement show view rendering
	ctx.HTML(http.StatusOK, "clusters/agents/dashboard/show.html", gin.H{
		"title": "Cluster Agent Details",
		"agent": agent,
	})
}

// checkFeatureFlag verifies if the k8s_dashboard feature is enabled
func (c *DashboardController) checkFeatureFlag(ctx *gin.Context) error {
	user, exists := ctx.Get("current_user")
	if !exists {
		return fmt.Errorf("user not found")
	}

	if !c.featureService.IsEnabled("k8s_dashboard", user.(*model.User)) {
		return fmt.Errorf("feature not enabled")
	}

	return nil
}

// findAgent retrieves the agent by ID
func (c *DashboardController) findAgent(ctx *gin.Context) (*model.Agent, error) {
	agentID := ctx.Param("agent_id")
	return c.agentService.GetAgent(ctx, agentID)
}

// authorizeReadClusterAgent checks if the user has permission to read the cluster agent
func (c *DashboardController) authorizeReadClusterAgent(ctx *gin.Context, agent *model.Agent) error {
	user, exists := ctx.Get("current_user")
	if !exists {
		return fmt.Errorf("user not found")
	}

	if !c.agentService.CanReadClusterAgent(ctx, user.(*model.User), agent) {
		return fmt.Errorf("access denied")
	}

	return nil
}

// setKasCookie sets the KAS cookie for the current session
func (c *DashboardController) setKasCookie(ctx *gin.Context) error {
	// TODO: Implement KAS cookie setting logic
	return nil
}
