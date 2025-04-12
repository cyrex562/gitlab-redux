package graphql

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// GraphqlExplorerController handles the GraphQL explorer interface
type GraphqlExplorerController struct {
	gonService *service.GonService
}

// NewGraphqlExplorerController creates a new instance of GraphqlExplorerController
func NewGraphqlExplorerController(gonService *service.GonService) *GraphqlExplorerController {
	return &GraphqlExplorerController{
		gonService: gonService,
	}
}

// RegisterRoutes registers the routes for the GraphqlExplorerController
func (c *GraphqlExplorerController) RegisterRoutes(r *gin.RouterGroup) {
	graphql := r.Group("/api/graphql")
	{
		graphql.GET("/explorer", c.show)
	}
}

// show handles the GET /api/graphql/explorer endpoint
func (c *GraphqlExplorerController) show(ctx *gin.Context) {
	// Add gon variables needed by Apollo client
	if err := c.gonService.AddVariables(ctx); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to setup GraphQL explorer"})
		return
	}

	// TODO: Render the GraphQL explorer template
	// This should render the HTML template that contains the Apollo client setup
	// and the GraphQL explorer interface
	ctx.HTML(http.StatusOK, "graphql/explorer.html", gin.H{
		"title": "GraphQL Explorer",
	})
}
