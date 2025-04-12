package explore

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmadden/gitlab-redux/internal/model"
	"github.com/jmadden/gitlab-redux/internal/service"
)

// SnippetsController handles snippet-related actions
type SnippetsController struct {
	*ApplicationController
	snippetService service.SnippetService
	noteService    service.NoteService
}

// NewSnippetsController creates a new SnippetsController
func NewSnippetsController(
	appController *ApplicationController,
	snippetService service.SnippetService,
	noteService service.NoteService,
) *SnippetsController {
	return &SnippetsController{
		ApplicationController: appController,
		snippetService:       snippetService,
		noteService:          noteService,
	}
}

// RegisterRoutes registers the routes for the SnippetsController
func (c *SnippetsController) RegisterRoutes(router *gin.Engine) {
	explore := router.Group("/explore")
	{
		explore.GET("/snippets", c.Index)
	}
}

// Index handles the index action
func (c *SnippetsController) Index(ctx *gin.Context) {
	// Get current user
	user := ctx.MustGet("current_user").(*model.User)

	// Get pagination parameters
	page := ctx.DefaultQuery("page", "1")
	perPage := ctx.DefaultQuery("per_page", "20")

	// Find snippets
	snippets, err := c.snippetService.FindSnippets(ctx, user, true, page, perPage)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get noteable metadata
	noteableMetaData, err := c.noteService.GetNoteableMetaData(ctx, snippets, "Snippet")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Render the index template
	ctx.HTML(http.StatusOK, "explore/snippets/index", gin.H{
		"layout":           c.GetLayout(),
		"snippets":         snippets,
		"noteableMetaData": noteableMetaData,
	})
}
