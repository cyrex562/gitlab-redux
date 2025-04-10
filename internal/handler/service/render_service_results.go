package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// RenderServiceResults handles rendering of service results
type RenderServiceResults struct {
	i18nService *service.I18nService
	logger      *service.Logger
}

// NewRenderServiceResults creates a new instance of RenderServiceResults
func NewRenderServiceResults(
	i18nService *service.I18nService,
	logger *service.Logger,
) *RenderServiceResults {
	return &RenderServiceResults{
		i18nService: i18nService,
		logger:      logger,
	}
}

// SuccessResponse renders a success response
func (r *RenderServiceResults) SuccessResponse(c *gin.Context, result map[string]interface{}) {
	httpStatus, ok := result["http_status"].(int)
	if !ok {
		httpStatus = http.StatusOK
	}

	body, ok := result["body"]
	if !ok {
		body = gin.H{}
	}

	c.JSON(httpStatus, body)
}

// ContinuePollingResponse renders a continue polling response
func (r *RenderServiceResults) ContinuePollingResponse(c *gin.Context) {
	c.JSON(http.StatusNoContent, gin.H{
		"status":  r.i18nService.Translate("processing"),
		"message": r.i18nService.Translate("Not ready yet. Try again later."),
	})
}

// ErrorResponse renders an error response
func (r *RenderServiceResults) ErrorResponse(c *gin.Context, result map[string]interface{}) {
	httpStatus, ok := result["http_status"].(int)
	if !ok {
		httpStatus = http.StatusBadRequest
	}

	status, _ := result["status"].(string)
	message, _ := result["message"].(string)

	c.JSON(httpStatus, gin.H{
		"status":  status,
		"message": message,
	})
}
