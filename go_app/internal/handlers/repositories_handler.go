package handlers

import (
	"net/http"
	"strconv"

	"github.com/cyrex562/gitlab-redux/internal/services"
	"github.com/gin-gonic/gin"
)

type RepositoriesHandler struct {
	jiraService *services.JiraService
}

func NewRepositoriesHandler(jiraService *services.JiraService) *RepositoriesHandler {
	return &RepositoriesHandler{
		jiraService: jiraService,
	}
}

// SearchRepositories handles GET /api/v4/jira/connect/repositories/search
// Searches for repositories by name with pagination
func (h *RepositoriesHandler) SearchRepositories(c *gin.Context) {
	// Get query parameters
	searchQuery := c.Query("searchQuery")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	// Get the current Jira installation
	installation, err := h.jiraService.GetCurrentInstallation(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Search for repositories
	repositories, err := h.jiraService.SearchRepositories(installation.ID, searchQuery, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the repositories as JSON
	c.JSON(http.StatusOK, gin.H{
		"containers": repositories,
	})
}

// AssociateRepository handles GET /api/v4/jira/connect/repositories/associate
// Associates a repository with a Jira installation
func (h *RepositoriesHandler) AssociateRepository(c *gin.Context) {
	// Get the repository ID
	repoID := c.Query("id")
	if repoID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Repository ID is required"})
		return
	}

	// Get the current Jira installation
	installation, err := h.jiraService.GetCurrentInstallation(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Find the repository
	repository, err := h.jiraService.FindRepository(installation.ID, repoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}

	// Return the repository as JSON
	c.JSON(http.StatusOK, repository)
} 