package params

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// ParamsBackwardCompatibility handles backward compatibility for parameters
type ParamsBackwardCompatibility struct {
	logger *service.Logger
}

// NewParamsBackwardCompatibility creates a new instance of ParamsBackwardCompatibility
func NewParamsBackwardCompatibility(
	logger *service.Logger,
) *ParamsBackwardCompatibility {
	return &ParamsBackwardCompatibility{
		logger: logger,
	}
}

// SetNonArchivedParam sets the non_archived parameter based on the archived parameter
func (p *ParamsBackwardCompatibility) SetNonArchivedParam(c *gin.Context) {
	// Get archived parameter
	archived := c.Query("archived")

	// Set non_archived parameter based on archived parameter
	if archived == "" {
		c.Set("non_archived", "true")
	} else {
		c.Set("non_archived", "false")
	}
}
