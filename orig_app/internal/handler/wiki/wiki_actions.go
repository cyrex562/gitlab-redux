package wiki

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmadden/gitlab-redux/internal/model"
	"github.com/jmadden/gitlab-redux/internal/service"
)

// WikiActionsHandler handles wiki page operations
type WikiActionsHandler struct {
	wikiService      *service.WikiService
	wikiPageService  *service.WikiPageService
	authService      *service.AuthService
	featureService   *service.FeatureService
	analyticsService *service.AnalyticsService
}

// NewWikiActionsHandler creates a new wiki actions handler
func NewWikiActionsHandler(
	wikiService *service.WikiService,
	wikiPageService *service.WikiPageService,
	authService *service.AuthService,
	featureService *service.FeatureService,
	analyticsService *service.AnalyticsService,
) *WikiActionsHandler {
	return &WikiActionsHandler{
		wikiService:      wikiService,
		wikiPageService:  wikiPageService,
		authService:      authService,
		featureService:   featureService,
		analyticsService: analyticsService,
	}
}

// RegisterRoutes registers the wiki routes
func (h *WikiActionsHandler) RegisterRoutes(router *gin.RouterGroup) {
	wiki := router.Group("/wiki")
	{
		wiki.GET("", h.pages)
		wiki.GET("/templates", h.templates)
		wiki.GET("/new", h.new)
		wiki.GET("/:id", h.show)
		wiki.GET("/:id/raw", h.raw)
		wiki.GET("/:id/edit", h.edit)
		wiki.PUT("/:id", h.update)
		wiki.POST("", h.create)
		wiki.GET("/:id/history", h.history)
		wiki.GET("/:id/diff", h.diff)
		wiki.DELETE("/:id", h.destroy)
		wiki.GET("/git_access", h.gitAccess)
	}
}

// pages handles the wiki pages list
func (h *WikiActionsHandler) pages(c *gin.Context) {
	// Check authorization
	if err := h.authService.AuthorizeReadWiki(c); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// Get wiki
	wiki, err := h.wikiService.GetWiki(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get pages list
	pages, err := h.wikiService.ListPages(c, wiki, c.Query("direction"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Filter out template pages
	var filteredPages []*model.WikiPage
	for _, page := range pages {
		if !strings.HasPrefix(page.Slug, "templates/") {
			filteredPages = append(filteredPages, page)
		}
	}

	// Group pages by directory
	groupedPages := h.wikiService.GroupPagesByDirectory(filteredPages)

	c.JSON(http.StatusOK, gin.H{
		"wiki_entries": groupedPages,
	})
}

// templates handles the wiki templates list
func (h *WikiActionsHandler) templates(c *gin.Context) {
	// Check authorization
	if err := h.authService.AuthorizeReadWiki(c); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// Get wiki
	wiki, err := h.wikiService.GetWiki(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get pages list
	pages, err := h.wikiService.ListPages(c, wiki, c.Query("direction"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Filter template pages
	var templatePages []*model.WikiPage
	for _, page := range pages {
		if strings.HasPrefix(page.Slug, "templates/") {
			templatePages = append(templatePages, page)
		}
	}

	// Group pages by directory
	groupedPages := h.wikiService.GroupPagesByDirectory(templatePages)

	c.JSON(http.StatusOK, gin.H{
		"wiki_entries":       groupedPages,
		"wiki_entries_count": len(templatePages),
	})
}

// new handles creating a new wiki page
func (h *WikiActionsHandler) new(c *gin.Context) {
	// Generate a random UUID for the new page
	pageID := uuid.New().String()

	// Redirect to the edit page with the random ID
	c.Redirect(http.StatusFound, fmt.Sprintf("/wiki/%s?random_title=true", pageID))
}

// show handles displaying a wiki page
func (h *WikiActionsHandler) show(c *gin.Context) {
	// Check authorization
	if err := h.authService.AuthorizeReadWiki(c); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// Get wiki
	wiki, err := h.wikiService.GetWiki(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get page
	page, err := h.wikiService.FindPage(c, wiki, c.Param("id"))
	if err != nil {
		// Check if it's a file blob
		fileBlob, err := h.wikiService.GetFileBlob(c, wiki, c.Param("id"))
		if err == nil && fileBlob != nil {
			h.sendWikiFileBlob(c, wiki, fileBlob)
			return
		}

		// Handle redirection
		redir, err := h.findRedirection(c, c.Param("id"))
		if err == nil && redir != nil {
			if redir.Error {
				message := fmt.Sprintf("The page at %s redirected too many times. You are now editing the page at %s.",
					c.Param("id"), c.Param("id"))
				c.Redirect(http.StatusFound, fmt.Sprintf("/wiki/%s?redirect_limit_reached=true", c.Param("id")))
				c.Set("flash", message)
				return
			}

			redirectedFrom := c.Query("redirected_from")
			if redirectedFrom == "" {
				redirectedFrom = c.Param("id")
			}

			message := fmt.Sprintf("The page at %s has been moved to %s.", redirectedFrom, redir.Path)
			c.Redirect(http.StatusFound, fmt.Sprintf("/wiki/%s?redirected_from=%s", redir.Path, redirectedFrom))
			c.Set("flash", message)
			return
		}

		// Show create form if authorized
		if h.canCreateWiki(c) && !wiki.Exists() {
			h.handleCreateForm(c, wiki, c.Param("id"))
			return
		}

		// Show 404 or empty wiki
		if wiki.Exists() {
			c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		} else {
			c.JSON(http.StatusOK, gin.H{"empty": true})
		}
		return
	}

	// Check encoding
	if !h.isValidEncoding(page) {
		c.Set("flash", "The content of this page is not encoded in UTF-8. Edits can only be made via the Git repository.")
	}

	// Get templates
	templates, err := h.wikiService.ListTemplates(c, wiki)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get sidebar
	sidebar, err := h.wikiService.FindSidebar(c, wiki, c.Query("version_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get pages count
	pagesCount, err := h.wikiService.GetPagesCount(c, wiki)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page":           page,
		"templates":      templates,
		"sidebar":        sidebar,
		"pages_count":    pagesCount,
		"ref":            c.Query("version_id"),
		"path":           page.Path,
		"valid_encoding": h.isValidEncoding(page),
	})
}

// handleCreateForm handles the create form display
func (h *WikiActionsHandler) handleCreateForm(c *gin.Context, wiki *model.Wiki, title string) {
	// Check if redirected
	redirectedFrom := c.Query("redirected_from")
	if redirectedFrom != "" {
		message := fmt.Sprintf("The page at %s tried to redirect to %s, but it does not exist. You are now editing the page at %s.",
			redirectedFrom, title, title)
		c.Set("flash", message)
	}

	// Build page
	page := h.wikiService.BuildPage(wiki, title)

	// Get templates
	templates, err := h.wikiService.ListTemplates(c, wiki)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page":      page,
		"templates": templates,
	})
}

// raw handles displaying raw wiki page content
func (h *WikiActionsHandler) raw(c *gin.Context) {
	// Check authorization
	if err := h.authService.AuthorizeReadWiki(c); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// Get wiki
	wiki, err := h.wikiService.GetWiki(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get page
	page, err := h.wikiService.FindPage(c, wiki, c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	// Set content type
	c.Header("Content-Type", "text/plain")
	c.String(http.StatusOK, page.RawContent)
}

// edit handles displaying the edit form
func (h *WikiActionsHandler) edit(c *gin.Context) {
	// Check authorization
	if err := h.authService.AuthorizeCreateWiki(c); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// Get wiki
	wiki, err := h.wikiService.GetWiki(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get page
	page, err := h.wikiService.FindPage(c, wiki, c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	// Check encoding
	if !h.isValidEncoding(page) {
		c.Redirect(http.StatusFound, fmt.Sprintf("/wiki/%s", page.Slug))
		return
	}

	// Get templates
	templates, err := h.wikiService.ListTemplates(c, wiki)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page":      page,
		"templates": templates,
	})
}

// update handles updating a wiki page
func (h *WikiActionsHandler) update(c *gin.Context) {
	// Check authorization
	if err := h.authService.AuthorizeCreateWiki(c); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// Get wiki
	wiki, err := h.wikiService.GetWiki(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get page
	page, err := h.wikiService.FindPage(c, wiki, c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	// Parse request body
	var params struct {
		Title         string `json:"title"`
		Content       string `json:"content"`
		Format        string `json:"format"`
		Message       string `json:"message"`
		LastCommitSHA string `json:"last_commit_sha"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update page
	updatedPage, err := h.wikiPageService.UpdatePage(c, wiki, page, params.Title, params.Content, params.Format, params.Message, params.LastCommitSHA)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Handle success
	h.handleActionSuccess(c, "updated", updatedPage)
}

// create handles creating a new wiki page
func (h *WikiActionsHandler) create(c *gin.Context) {
	// Check authorization
	if err := h.authService.AuthorizeCreateWiki(c); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// Get wiki
	wiki, err := h.wikiService.GetWiki(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Parse request body
	var params struct {
		Title         string `json:"title"`
		Content       string `json:"content"`
		Format        string `json:"format"`
		Message       string `json:"message"`
		LastCommitSHA string `json:"last_commit_sha"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create page
	page, err := h.wikiPageService.CreatePage(c, wiki, params.Title, params.Content, params.Format, params.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Handle success
	h.handleActionSuccess(c, "created", page)
}

// history handles displaying page history
func (h *WikiActionsHandler) history(c *gin.Context) {
	// Check authorization
	if err := h.authService.AuthorizeReadWiki(c); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// Get wiki
	wiki, err := h.wikiService.GetWiki(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get page
	page, err := h.wikiService.FindPage(c, wiki, c.Param("id"))
	if err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("/wiki?notice=%s", "Page not found"))
		return
	}

	// Get versions
	versions, err := h.wikiService.GetVersions(c, wiki, page, c.Query("page"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get versions count
	versionsCount, err := h.wikiService.GetVersionsCount(c, wiki, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page":           page,
		"versions":       versions,
		"versions_count": versionsCount,
	})
}

// diff handles displaying page diff
func (h *WikiActionsHandler) diff(c *gin.Context) {
	// Check authorization
	if err := h.authService.AuthorizeReadWiki(c); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// Get wiki
	wiki, err := h.wikiService.GetWiki(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get page
	page, err := h.wikiService.FindPage(c, wiki, c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	// Get diff options
	diffOptions := h.getDiffOptions(c)

	// Get diffs
	diffs, err := h.wikiService.GetDiffs(c, wiki, page, diffOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page":  page,
		"diffs": diffs,
	})
}

// destroy handles deleting a wiki page
func (h *WikiActionsHandler) destroy(c *gin.Context) {
	// Check authorization
	if err := h.authService.AuthorizeCreateWiki(c); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// Get wiki
	wiki, err := h.wikiService.GetWiki(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get page
	page, err := h.wikiService.FindPage(c, wiki, c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	// Delete page
	err = h.wikiPageService.DestroyPage(c, wiki, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set flash message
	c.Set("flash", "Wiki page was successfully deleted.")

	// Redirect to wiki index
	c.Redirect(http.StatusFound, "/wiki")
}

// gitAccess handles displaying git access information
func (h *WikiActionsHandler) gitAccess(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// handleActionSuccess handles successful actions
func (h *WikiActionsHandler) handleActionSuccess(c *gin.Context, action string, page *model.WikiPage) {
	var message string
	if page.Title == "sidebar" {
		if action == "created" {
			message = "Sidebar was successfully created."
		} else {
			message = "Sidebar was successfully updated."
		}
		c.Redirect(http.StatusFound, "/wiki")
	} else {
		if action == "created" {
			message = "Wiki page was successfully created."
		} else {
			message = "Wiki page was successfully updated."
		}
		c.Redirect(http.StatusFound, fmt.Sprintf("/wiki/%s", page.Slug))
	}
	c.Set("flash", message)
}

// canCreateWiki checks if the user can create wiki pages
func (h *WikiActionsHandler) canCreateWiki(c *gin.Context) bool {
	return h.authService.CanCreateWiki(c)
}

// isValidEncoding checks if the page content is valid UTF-8
func (h *WikiActionsHandler) isValidEncoding(page *model.WikiPage) bool {
	return page != nil && page.IsValidEncoding()
}

// sendWikiFileBlob sends a wiki file blob
func (h *WikiActionsHandler) sendWikiFileBlob(c *gin.Context, wiki *model.Wiki, blob *model.Blob) {
	// TODO: Implement sending file blob
	c.JSON(http.StatusOK, gin.H{"blob": blob})
}

// getDiffOptions gets diff options from request
func (h *WikiActionsHandler) getDiffOptions(c *gin.Context) *model.DiffOptions {
	return &model.DiffOptions{
		View: c.Query("view"),
	}
}

// findRedirection finds a redirection for a path
func (h *WikiActionsHandler) findRedirection(c *gin.Context, path string) (*model.Redirection, error) {
	// Get wiki
	wiki, err := h.wikiService.GetWiki(c)
	if err != nil {
		return nil, err
	}

	// Get redirections
	redirections, err := h.wikiService.GetRedirections(c, wiki)
	if err != nil {
		return nil, err
	}

	// Find redirection
	seen := make(map[string]bool)
	currentPath := path
	redirectLimit := 50

	for i := 0; i < redirectLimit; i++ {
		seen[currentPath] = true
		nextPath := h.findSingleRedirection(currentPath, redirections)

		if nextPath == "" {
			if currentPath == path {
				return nil, nil
			}
			return &model.Redirection{Path: currentPath}, nil
		}

		if seen[nextPath] {
			return &model.Redirection{Error: true, Reason: "loop"}, nil
		}

		currentPath = nextPath
	}

	return &model.Redirection{Error: true, Reason: "limit"}, nil
}

// findSingleRedirection finds a single redirection for a path
func (h *WikiActionsHandler) findSingleRedirection(path string, redirections map[string]string) string {
	current := path
	rest := []string{}

	for current != "." {
		redirect, ok := redirections[current]
		if ok {
			return filepath.Join(append([]string{redirect}, rest...)...)
		}

		dir, file := filepath.Split(current)
		rest = append([]string{file}, rest...)
		current = dir
	}

	return ""
}
