package admin

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// FeatureCategory represents the feature category for instance review
const FeatureCategory = "devops_reports"

// Urgency represents the urgency level for instance review
const Urgency = "low"

// InstanceReviewHandler handles instance review requests
type InstanceReviewHandler struct {
	usageService *service.UsageService
}

// NewInstanceReviewHandler creates a new InstanceReviewHandler instance
func NewInstanceReviewHandler(usageService *service.UsageService) *InstanceReviewHandler {
	return &InstanceReviewHandler{
		usageService: usageService,
	}
}

// Index handles the GET request to redirect to the subscription portal
func (h *InstanceReviewHandler) Index(c *gin.Context) {
	// Set feature category and urgency
	c.Set("feature_category", FeatureCategory)
	c.Set("urgency", Urgency)

	// Get current user from context
	currentUser, exists := c.Get("current_user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Build instance review parameters
	params := h.buildInstanceReviewParams(c, currentUser)

	// Get subscription portal URL from configuration
	portalURL := c.GetString("subscription_portal_url")
	if portalURL == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Subscription portal URL not configured"})
		return
	}

	// Build redirect URL
	redirectURL := portalURL + "/instance_review?" + params.Encode()

	// Redirect to subscription portal
	c.Redirect(http.StatusFound, redirectURL)
}

// buildInstanceReviewParams builds the parameters for instance review
func (h *InstanceReviewHandler) buildInstanceReviewParams(c *gin.Context, currentUser interface{}) url.Values {
	params := url.Values{}
	instanceReview := make(map[string]interface{})

	// Add basic user information
	instanceReview["email"] = currentUser.(*model.User).Email
	instanceReview["last_name"] = currentUser.(*model.User).Name
	instanceReview["version"] = c.GetString("gitlab_version")

	// Check if usage ping is enabled
	if c.GetBool("usage_ping_enabled") {
		// Get usage data
		usageData, err := h.usageService.GetServicePingData(c)
		if err == nil {
			// Add usage metrics
			instanceReview["users_count"] = strconv.Itoa(usageData.ActiveUserCount)
			instanceReview["projects_count"] = strconv.Itoa(usageData.Counts.Projects)
			instanceReview["groups_count"] = strconv.Itoa(usageData.Counts.Groups)
			instanceReview["issues_count"] = strconv.Itoa(usageData.Counts.Issues)
			instanceReview["merge_requests_count"] = strconv.Itoa(usageData.Counts.MergeRequests)
			instanceReview["internal_pipelines_count"] = strconv.Itoa(usageData.Counts.CIInternalPipelines)
			instanceReview["external_pipelines_count"] = strconv.Itoa(usageData.Counts.CIExternalPipelines)
			instanceReview["labels_count"] = strconv.Itoa(usageData.Counts.Labels)
			instanceReview["milestones_count"] = strconv.Itoa(usageData.Counts.Milestones)
			instanceReview["snippets_count"] = strconv.Itoa(usageData.Counts.Snippets)
			instanceReview["notes_count"] = strconv.Itoa(usageData.Counts.Notes)
		}
	}

	// Add instance review data to parameters
	params.Add("instance_review", instanceReview)

	return params
}
