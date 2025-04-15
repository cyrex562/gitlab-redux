package jiraconnect

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/gitlab-org/gitlab-redux/app/services"
)

const (
	homeURL = "https://gitlab.com"
	docURL  = "https://docs.gitlab.com/ee/integration/jira/"
)

// AppDescriptorHandler handles Jira Connect app descriptor
type AppDescriptorHandler struct {
	*BaseHandler
	jiraService *services.JiraService
}

// NewAppDescriptorHandler creates a new AppDescriptorHandler
func NewAppDescriptorHandler(baseHandler *BaseHandler, jiraService *services.JiraService) *AppDescriptorHandler {
	return &AppDescriptorHandler{
		BaseHandler: baseHandler,
		jiraService: jiraService,
	}
}

// RegisterRoutes registers the routes for the app descriptor handler
func (h *AppDescriptorHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/jira_connect/app_descriptor", h.Show).Methods("GET")
}

// Show returns the Jira Connect app descriptor
func (h *AppDescriptorHandler) Show(w http.ResponseWriter, r *http.Request) {
	// Get base URL for Jira Connect
	baseURL := h.jiraService.GetJiraConnectBaseURL(r)

	// Get logo URL
	logoURL := h.jiraService.GetLogoURL()

	// Create app descriptor
	descriptor := map[string]interface{}{
		"name":        h.jiraService.GetAppName(),
		"description": "Integrate commits, branches and merge requests from GitLab into Jira",
		"key":         h.jiraService.GetAppKey(),
		"baseUrl":     baseURL,
		"lifecycle": map[string]interface{}{
			"installed":   h.relativeToBasePath(baseURL, "/jira_connect/events/installed"),
			"uninstalled": h.relativeToBasePath(baseURL, "/jira_connect/events/uninstalled"),
		},
		"vendor": map[string]interface{}{
			"name": "GitLab",
			"url":  homeURL,
		},
		"links": map[string]interface{}{
			"documentation": h.jiraService.GetHelpPageURL("integration/jira/development_panel.md"),
		},
		"authentication": map[string]interface{}{
			"type": "jwt",
		},
		"modules": h.getModules(baseURL, logoURL),
		"scopes":  []string{"READ", "WRITE", "DELETE"},
		"apiVersion": 1,
		"apiMigrations": map[string]interface{}{
			"context-qsh":     true,
			"signed-install":  true,
			"gdpr":            true,
		},
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(descriptor)
}

// getModules returns the modules for the app descriptor
func (h *AppDescriptorHandler) getModules(baseURL, logoURL string) map[string]interface{} {
	modules := map[string]interface{}{
		"postInstallPage": map[string]interface{}{
			"key":   "gitlab-configuration",
			"name":  map[string]interface{}{"value": "GitLab Configuration"},
			"url":   h.relativeToBasePath(baseURL, "/jira_connect/subscriptions"),
			"conditions": []map[string]interface{}{
				{
					"condition": "user_is_admin",
					"invert":    false,
				},
			},
		},
	}

	// Add development tool module
	modules["jiraDevelopmentTool"] = h.getDevelopmentToolModule(baseURL, logoURL)

	// Add build information module
	modules["jiraBuildInfoProvider"] = h.getBuildInformationModule(baseURL, logoURL)

	// Add deployment information module
	modules["jiraDeploymentInfoProvider"] = h.getDeploymentInformationModule(baseURL, logoURL)

	// Add feature flag module
	modules["jiraFeatureFlagInfoProvider"] = h.getFeatureFlagModule(baseURL, logoURL)

	return modules
}

// getDevelopmentToolModule returns the development tool module
func (h *AppDescriptorHandler) getDevelopmentToolModule(baseURL, logoURL string) map[string]interface{} {
	return map[string]interface{}{
		"actions": map[string]interface{}{
			"createBranch": map[string]interface{}{
				"templateUrl": baseURL + "/jira_connect/branches" + h.createBranchParams(),
			},
			"searchConnectedWorkspaces": map[string]interface{}{
				"templateUrl": baseURL + "/jira_connect/workspaces/search",
			},
			"searchRepositories": map[string]interface{}{
				"templateUrl": baseURL + "/jira_connect/repositories/search",
			},
			"associateRepository": map[string]interface{}{
				"templateUrl": baseURL + "/jira_connect/repositories/associate",
			},
		},
		"key":         "gitlab-development-tool",
		"application": map[string]interface{}{"value": h.jiraService.GetDisplayName()},
		"name":        map[string]interface{}{"value": h.jiraService.GetDisplayName()},
		"url":         homeURL,
		"logoUrl":     logoURL,
		"capabilities": []string{"branch", "commit", "pull_request"},
	}
}

// getBuildInformationModule returns the build information module
func (h *AppDescriptorHandler) getBuildInformationModule(baseURL, logoURL string) map[string]interface{} {
	return h.getCommonModuleProperties(baseURL, logoURL, map[string]interface{}{
		"actions": map[string]interface{}{},
		"name":    map[string]interface{}{"value": "GitLab CI"},
		"key":     "gitlab-ci",
	})
}

// getDeploymentInformationModule returns the deployment information module
func (h *AppDescriptorHandler) getDeploymentInformationModule(baseURL, logoURL string) map[string]interface{} {
	return h.getCommonModuleProperties(baseURL, logoURL, map[string]interface{}{
		"actions": map[string]interface{}{}, // TODO: list deployments
		"name":    map[string]interface{}{"value": "GitLab Deployments"},
		"key":     "gitlab-deployments",
	})
}

// getFeatureFlagModule returns the feature flag module
func (h *AppDescriptorHandler) getFeatureFlagModule(baseURL, logoURL string) map[string]interface{} {
	return h.getCommonModuleProperties(baseURL, logoURL, map[string]interface{}{
		"actions": map[string]interface{}{}, // TODO: create, link and list feature flags
		"name":    map[string]interface{}{"value": "GitLab Feature Flags"},
		"key":     "gitlab-feature-flags",
	})
}

// getCommonModuleProperties returns common properties for modules
func (h *AppDescriptorHandler) getCommonModuleProperties(baseURL, logoURL string, additionalProps map[string]interface{}) map[string]interface{} {
	props := map[string]interface{}{
		"homeUrl":         homeURL,
		"logoUrl":         logoURL,
		"documentationUrl": docURL,
	}

	// Merge additional properties
	for k, v := range additionalProps {
		props[k] = v
	}

	return props
}

// relativeToBasePath returns a path relative to the base path
func (h *AppDescriptorHandler) relativeToBasePath(baseURL, fullPath string) string {
	// This is a simplified version - in a real implementation, you would need to handle
	// the path manipulation more carefully
	return fullPath
}

// createBranchParams returns the parameters for creating a branch
func (h *AppDescriptorHandler) createBranchParams() string {
	return "?issue_key={issue.key}&issue_summary={issue.summary}&jwt={jwt}&addonkey=" + h.jiraService.GetAppKey()
} 