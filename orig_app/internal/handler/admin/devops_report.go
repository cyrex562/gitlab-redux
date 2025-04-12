package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// DevOpsReportHandler handles DevOps report related requests
type DevOpsReportHandler struct {
	devOpsReportService *service.DevOpsReportService
}

// NewDevOpsReportHandler creates a new DevOpsReportHandler instance
func NewDevOpsReportHandler(devOpsReportService *service.DevOpsReportService) *DevOpsReportHandler {
	return &DevOpsReportHandler{
		devOpsReportService: devOpsReportService,
	}
}

// Show handles the GET request to display the DevOps report
func (h *DevOpsReportHandler) Show(c *gin.Context) {
	metric, err := h.devOpsReportService.GetLatestMetric(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch DevOps metric"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"metric": metric,
	})
}

// ShouldTrackDevOpsScore determines if we should track the DevOps score
func (h *DevOpsReportHandler) ShouldTrackDevOpsScore() bool {
	return true
}
