package dashboard

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmadden/gitlab-redux/internal/handler/pagination"
	"github.com/jmadden/gitlab-redux/internal/handler/snippets"
	"github.com/jmadden/gitlab-redux/internal/service"
)

// SnippetsController handles snippet-related actions in the dashboard
type SnippetsController struct {
	ApplicationController
	snippetService service.SnippetService
	countService   service.SnippetCountService
	renderService  service.RenderService
}

// NewSnippetsController creates a new SnippetsController
func NewSnippetsController(
	appController ApplicationController,
	snippetService service.SnippetService,
	countService service.SnippetCountService,
	renderService service.RenderService,
) *SnippetsController {
	return &SnippetsController{
		ApplicationController: appController,
		snippetService:       snippetService,
		countService:         countService,
		renderService:        renderService,
	}
}

// RegisterRoutes registers the routes for the SnippetsController
func (c *SnippetsController) RegisterRoutes(router *gin.Engine) {
	dashboard := router.Group("/dashboard")
	{
		dashboard.GET("/snippets", c.Index)
	}
}

// Index handles the index action for snippets
func (c *SnippetsController) Index(ctx *gin.Context) {
	// Get current user
	currentUser, err := c.GetCurrentUser(ctx)
	if err != nil {
		c.HandleError(ctx, err)
		return
	}

	// Get snippet counts
	snippetCounts, err := c.countService.GetSnippetCounts(ctx, currentUser, currentUser)
	if err != nil {
		c.HandleError(ctx, err)
		return
	}
	ctx.Set("snippet_counts", snippetCounts)

	// Get pagination parameters
	page, perPage := pagination.GetPaginationParams(ctx)
	ctx.Set("page", page)
	ctx.Set("per_page", perPage)

	// Get sort parameter
	sortParam := snippets.GetSortParam(ctx)
	ctx.Set("sort", sortParam)

	// Get scope parameter
	scope := ctx.DefaultQuery("scope", "")
	ctx.Set("scope", scope)

	// Find snippets
	snippets, err := c.snippetService.FindSnippets(ctx, currentUser, currentUser, scope, sortParam, page, perPage)
	if err != nil {
		c.HandleError(ctx, err)
		return
	}

	// Check if we need to redirect out of range
	if pagination.ShouldRedirectOutOfRange(ctx, snippets) {
		return
	}

	// Get noteable metadata
	noteableMetaData, err := c.renderService.GetNoteableMetaData(ctx, snippets, "Snippet")
	if err != nil {
		c.HandleError(ctx, err)
		return
	}
	ctx.Set("noteable_meta_data", noteableMetaData)

	// Render the view
	ctx.HTML(http.StatusOK, "dashboard/snippets/index", gin.H{
		"snippets":         snippets,
		"snippet_counts":   snippetCounts,
		"noteable_meta_data": noteableMetaData,
		"page":             page,
		"per_page":         perPage,
		"sort":             sortParam,
		"scope":            scope,
	})
}
