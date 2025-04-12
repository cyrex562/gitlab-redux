package markdown

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// PreviewMarkdown handles markdown preview functionality
type PreviewMarkdown struct {
	markdownService *service.MarkdownService
	viewService    *service.ViewService
	projectService *service.ProjectService
	wikiService    *service.WikiService
	groupService   *service.GroupService
	browserService *service.BrowserService
	logger         *service.Logger
}

// NewPreviewMarkdown creates a new instance of PreviewMarkdown
func NewPreviewMarkdown(
	markdownService *service.MarkdownService,
	viewService *service.ViewService,
	projectService *service.ProjectService,
	wikiService *service.WikiService,
	groupService *service.GroupService,
	browserService *service.BrowserService,
	logger *service.Logger,
) *PreviewMarkdown {
	return &PreviewMarkdown{
		markdownService: markdownService,
		viewService:    viewService,
		projectService: projectService,
		wikiService:    wikiService,
		groupService:   groupService,
		browserService: browserService,
		logger:         logger,
	}
}

// PreviewMarkdown handles the preview markdown action
func (p *PreviewMarkdown) PreviewMarkdown(c *gin.Context) error {
	// Get resource parent
	resourceParent, err := p.getResourceParent(c)
	if err != nil {
		return err
	}

	// Get current user
	user, err := p.getCurrentUser(c)
	if err != nil {
		return err
	}

	// Get markdown service params
	markdownServiceParams := p.getMarkdownServiceParams(c)

	// Execute markdown service
	result, err := p.markdownService.PreviewMarkdown(c, resourceParent, user, markdownServiceParams)
	if err != nil {
		return err
	}

	// Get markdown context params
	markdownContextParams, err := p.getMarkdownContextParams(c)
	if err != nil {
		return err
	}

	// Render markdown
	body, err := p.viewService.RenderMarkdown(c, result.Text, markdownContextParams)
	if err != nil {
		return err
	}

	// Render commands
	commands, err := p.viewService.RenderMarkdown(c, result.Commands, nil)
	if err != nil {
		return err
	}

	// Serialize suggestions
	suggestions, err := p.serializeSuggestions(c, result.Suggestions)
	if err != nil {
		return err
	}

	// Return JSON response
	c.JSON(http.StatusOK, gin.H{
		"body": body,
		"references": gin.H{
			"users":       result.Users,
			"suggestions": suggestions,
			"commands":    commands,
		},
	})

	return nil
}

// GetResourceParent gets the resource parent
func (p *PreviewMarkdown) getResourceParent(c *gin.Context) (interface{}, error) {
	// Get project from context
	project, err := p.projectService.GetProjectFromContext(c)
	if err != nil {
		return nil, err
	}

	return project, nil
}

// GetCurrentUser gets the current user
func (p *PreviewMarkdown) getCurrentUser(c *gin.Context) (*service.User, error) {
	// Get current user from context
	user, err := p.getCurrentUserFromContext(c)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetProjectsFilterParams gets the projects filter params
func (p *PreviewMarkdown) getProjectsFilterParams(c *gin.Context) map[string]interface{} {
	// Get preview suggestions from query parameters
	previewSuggestions := c.Query("preview_suggestions") != ""

	return map[string]interface{}{
		"issuable_reference_expansion_enabled": true,
		"suggestions_filter_enabled":           previewSuggestions,
	}
}

// GetTimelineEventsFilterParams gets the timeline events filter params
func (p *PreviewMarkdown) getTimelineEventsFilterParams(c *gin.Context) map[string]interface{} {
	return map[string]interface{}{
		"issuable_reference_expansion_enabled": true,
		"pipeline":                            "incident_management/timeline_event",
	}
}

// GetWikisFilterParams gets the wikis filter params
func (p *PreviewMarkdown) getWikisFilterParams(c *gin.Context) (map[string]interface{}, error) {
	// Get wiki from context
	wiki, err := p.wikiService.GetWikiFromContext(c)
	if err != nil {
		return nil, err
	}

	// Get page slug from query parameters
	pageSlug := c.Query("id")

	// Get repository
	repository, err := wiki.GetRepository()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"pipeline":                            "wiki",
		"wiki":                                wiki,
		"page_slug":                           pageSlug,
		"repository":                          repository,
		"issuable_reference_expansion_enabled": true,
	}, nil
}

// GetMarkdownServiceParams gets the markdown service params
func (p *PreviewMarkdown) getMarkdownServiceParams(c *gin.Context) map[string]interface{} {
	// Get all query parameters
	params := make(map[string]interface{})
	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}

	return params
}

// GetMarkdownContextParams gets the markdown context params
func (p *PreviewMarkdown) getMarkdownContextParams(c *gin.Context) (map[string]interface{}, error) {
	// Get controller name from context
	controllerName := c.GetString("controller")

	// Get base params
	baseParams := map[string]interface{}{
		"requested_path": c.Query("path"),
		"ref":            c.Query("ref"),
		"allow_comments": !p.browserService.IsIE(c),
	}

	// Get specific params based on controller name
	var specificParams map[string]interface{}
	var err error

	switch controllerName {
	case "wikis":
		// Get wiki page
		wiki, err := p.wikiService.GetWikiFromContext(c)
		if err != nil {
			return nil, err
		}

		pageSlug := c.Query("id")
		wikiPage, err := wiki.FindPage(pageSlug)
		if err != nil {
			return nil, err
		}

		// Get wikis filter params
		specificParams, err = p.getWikisFilterParams(c)
		if err != nil {
			return nil, err
		}

		// Add wiki page path to base params
		if wikiPage != nil {
			baseParams["requested_path"] = wikiPage.Path
		}

	case "snippets":
		specificParams = map[string]interface{}{
			"skip_project_check": true,
		}

	case "groups":
		// Get group from context
		group, err := p.groupService.GetGroupFromContext(c)
		if err != nil {
			return nil, err
		}

		specificParams = map[string]interface{}{
			"group":                                group,
			"issuable_reference_expansion_enabled": true,
		}

	case "projects":
		specificParams = p.getProjectsFilterParams(c)

	case "timeline_events":
		specificParams = p.getTimelineEventsFilterParams(c)

	case "organizations":
		specificParams = map[string]interface{}{
			"pipeline": "description",
		}

	default:
		specificParams = map[string]interface{}{}
	}

	// Merge base params and specific params
	for key, value := range specificParams {
		baseParams[key] = value
	}

	return baseParams, nil
}

// SerializeSuggestions serializes suggestions
func (p *PreviewMarkdown) serializeSuggestions(c *gin.Context, suggestions []interface{}) (interface{}, error) {
	// Create suggestion serializer
	suggestionSerializer := p.createSuggestionSerializer()

	// Serialize suggestions
	return suggestionSerializer.SerializeDiff(suggestions)
}

// CreateSuggestionSerializer creates a suggestion serializer
func (p *PreviewMarkdown) createSuggestionSerializer() *service.SuggestionSerializer {
	return service.NewSuggestionSerializer()
}

// GetCurrentUserFromContext gets the current user from context
func (p *PreviewMarkdown) getCurrentUserFromContext(c *gin.Context) (*service.User, error) {
	// This is a placeholder for actual implementation
	// In a real implementation, this would get the current user from the context
	return &service.User{}, nil
}
