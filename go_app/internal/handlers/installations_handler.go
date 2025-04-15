package handlers

import (
	"net/http"

	"github.com/cyrex562/gitlab-redux/internal/models"
	"github.com/cyrex562/gitlab-redux/internal/services"
	"github.com/gin-gonic/gin"
)

type InstallationsHandler struct {
	jiraService *services.JiraService
}

func NewInstallationsHandler(jiraService *services.JiraService) *InstallationsHandler {
	return &InstallationsHandler{
		jiraService: jiraService,
	}
}

// GetInstallation handles GET /api/v4/jira/connect/installations
func (h *InstallationsHandler) GetInstallation(c *gin.Context) {
	installation, err := h.jiraService.GetCurrentInstallation(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, h.installationJSON(installation))
}

// UpdateInstallation handles PUT /api/v4/jira/connect/installations
func (h *InstallationsHandler) UpdateInstallation(c *gin.Context) {
	var params struct {
		Installation struct {
			InstanceURL string `json:"instance_url"`
		} `json:"installation"`
	}

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	installation, err := h.jiraService.GetCurrentInstallation(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update the installation
	installation.InstanceURL = params.Installation.InstanceURL
	if err := h.jiraService.UpdateInstallation(c, installation); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, h.installationJSON(installation))
}

// installationJSON returns the JSON representation of an installation
func (h *InstallationsHandler) installationJSON(installation *models.JiraConnectInstallation) gin.H {
	return gin.H{
		"gitlab_com": installation.InstanceURL == "",
		"instance_url": installation.InstanceURL,
	}
} 