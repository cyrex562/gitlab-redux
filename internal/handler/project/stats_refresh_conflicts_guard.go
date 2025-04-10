package project

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// ProjectStatsRefreshConflictsGuard handles project stats refresh conflicts
type ProjectStatsRefreshConflictsGuard struct {
	projectService *service.ProjectService
	logger         *service.Logger
}

// NewProjectStatsRefreshConflictsGuard creates a new instance of ProjectStatsRefreshConflictsGuard
func NewProjectStatsRefreshConflictsGuard(
	projectService *service.ProjectService,
	logger *service.Logger,
) *ProjectStatsRefreshConflictsGuard {
	return &ProjectStatsRefreshConflictsGuard{
		projectService: projectService,
		logger:         logger,
	}
}

// RejectIfBuildArtifactsSizeRefreshing checks if the project is refreshing build artifacts size and rejects the request if it is
func (p *ProjectStatsRefreshConflictsGuard) RejectIfBuildArtifactsSizeRefreshing(c *gin.Context) error {
	// Get project from context
	project, err := p.projectService.GetProjectFromContext(c)
	if err != nil {
		return err
	}

	// Check if project is refreshing build artifacts size
	isRefreshing, err := p.projectService.IsRefreshingBuildArtifactsSize(project)
	if err != nil {
		return err
	}

	// If not refreshing, return nil
	if !isRefreshing {
		return nil
	}

	// Log warning
	p.logger.Warn("Request rejected during stats refresh", map[string]interface{}{
		"project_id": project.ID,
	})

	// Return 409 Conflict
	c.JSON(http.StatusConflict, gin.H{
		"message": "Action temporarily disabled. The project this pipeline belongs to is undergoing stats refresh.",
	})

	return nil
}
