package diff

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// RapidDiffsResource handles rapid diffs functionality
type RapidDiffsResource struct {
	featureService *service.FeatureService
	userService    *service.UserService
	diffService    *service.DiffService
	logger         *service.Logger
}

// NewRapidDiffsResource creates a new instance of RapidDiffsResource
func NewRapidDiffsResource(
	featureService *service.FeatureService,
	userService *service.UserService,
	diffService *service.DiffService,
	logger *service.Logger,
) *RapidDiffsResource {
	return &RapidDiffsResource{
		featureService: featureService,
		userService:    userService,
		diffService:    diffService,
		logger:         logger,
	}
}

// DiffsStreamURL gets the URL for streaming diffs
func (r *RapidDiffsResource) DiffsStreamURL(c *gin.Context, resource interface{}, offset *int, diffView *string) (string, error) {
	// Get diffs for streaming
	diffsForStreaming, err := r.diffService.GetDiffsForStreaming(resource)
	if err != nil {
		return "", err
	}

	// Get diff files count
	diffFilesCount, err := r.diffService.GetDiffFilesCount(diffsForStreaming)
	if err != nil {
		return "", err
	}

	// Check if offset is valid
	if offset != nil && *offset > diffFilesCount {
		return "", nil
	}

	// Get diffs stream resource URL
	return r.diffsStreamResourceURL(c, resource, offset, diffView)
}

// DiffFilesMetadata gets the metadata for diff files
func (r *RapidDiffsResource) DiffFilesMetadata(c *gin.Context) error {
	// Check if rapid diffs is enabled
	enabled, err := r.isRapidDiffsEnabled(c)
	if err != nil {
		return err
	}

	if !enabled {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Rapid diffs is not enabled",
		})
		return nil
	}

	// Get diffs resource
	diffsResource, err := r.getDiffsResource(c)
	if err != nil {
		return err
	}

	if diffsResource == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Diffs resource not found",
		})
		return nil
	}

	// Get raw diff files
	rawDiffFiles, err := r.diffService.GetRawDiffFiles(diffsResource, true)
	if err != nil {
		return err
	}

	// Return JSON response
	c.JSON(http.StatusOK, gin.H{
		"diff_files": rawDiffFiles,
	})

	return nil
}

// IsRapidDiffsEnabled checks if rapid diffs is enabled
func (r *RapidDiffsResource) isRapidDiffsEnabled(c *gin.Context) (bool, error) {
	// Get current user
	user, err := r.userService.GetCurrentUser(c)
	if err != nil {
		return false, err
	}

	// Check if feature is enabled
	return r.featureService.IsEnabled("rapid_diffs", user, "wip")
}

// GetDiffsResource gets the diffs resource
func (r *RapidDiffsResource) getDiffsResource(c *gin.Context) (interface{}, error) {
	// This is a placeholder for actual implementation
	// In a real implementation, this would get the diffs resource from the context
	return nil, nil
}

// DiffsStreamResourceURL gets the URL for streaming diffs resource
func (r *RapidDiffsResource) diffsStreamResourceURL(c *gin.Context, resource interface{}, offset *int, diffView *string) (string, error) {
	// This is a placeholder for actual implementation
	// In a real implementation, this would get the URL for streaming diffs resource
	return "", nil
}
