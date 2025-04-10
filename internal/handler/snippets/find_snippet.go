package snippets

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
	"gitlab.com/gitlab-org/gitlab-redux/internal/util"
)

// FindSnippet handles finding and retrieving snippets
type FindSnippet struct {
	snippetService *service.SnippetService
	logger         *util.Logger
}

// NewFindSnippet creates a new instance of FindSnippet
func NewFindSnippet(snippetService *service.SnippetService, logger *util.Logger) *FindSnippet {
	return &FindSnippet{
		snippetService: snippetService,
		logger:         logger,
	}
}

// GetSnippet gets a snippet by ID with all necessary relations for viewing
func (f *FindSnippet) GetSnippet(ctx *gin.Context) (*model.Snippet, error) {
	// Check if the snippet is already in the context
	if snippet, exists := ctx.Get("snippet"); exists {
		return snippet.(*model.Snippet), nil
	}

	// Get the snippet ID from the request parameters
	snippetID := f.GetSnippetID(ctx)

	// Get the snippet find parameters
	findParams := f.GetSnippetFindParams(ctx)

	// Find the snippet with all necessary relations for viewing
	snippet, err := f.snippetService.FindSnippetWithRelations(ctx, findParams)
	if err != nil {
		return nil, err
	}

	// Store the snippet in the context
	ctx.Set("snippet", snippet)

	return snippet, nil
}

// GetSnippetID gets the snippet ID from the request parameters
func (f *FindSnippet) GetSnippetID(ctx *gin.Context) string {
	return ctx.Param("id")
}

// GetSnippetFindParams gets the parameters for finding a snippet
func (f *FindSnippet) GetSnippetFindParams(ctx *gin.Context) map[string]interface{} {
	return map[string]interface{}{
		"id": f.GetSnippetID(ctx),
	}
}

// GetSnippetClass is an interface method that should be implemented by concrete types
// to specify which snippet class to use
func (f *FindSnippet) GetSnippetClass() (interface{}, error) {
	return nil, util.NewNotImplementedError("GetSnippetClass must be implemented by concrete types")
}
