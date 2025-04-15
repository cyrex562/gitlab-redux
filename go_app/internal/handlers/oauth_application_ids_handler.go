package handlers

import (
	"net/http"

	"github.com/cyrex562/gitlab-redux/internal/services"
	"github.com/gin-gonic/gin"
)

type OauthApplicationIdsHandler struct {
	jiraService *services.JiraService
}

func NewOauthApplicationIdsHandler(jiraService *services.JiraService) *OauthApplicationIdsHandler {
	return &OauthApplicationIdsHandler{
		jiraService: jiraService,
	}
}

// ShowApplicationId handles GET /api/v4/jira/connect/oauth_application_ids
func (h *OauthApplicationIdsHandler) ShowApplicationId(c *gin.Context) {
	// Skip JWT verification as per the Ruby controller
	// This is handled by the router middleware

	applicationId, err := h.jiraService.GetJiraConnectApplicationKey()
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	if applicationId == "" {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"application_id": applicationId,
	})
} 