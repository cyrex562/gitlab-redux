package milestone

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// MilestoneActions handles milestone-related actions
type MilestoneActions struct {
	milestoneService *service.MilestoneService
	viewService      *service.ViewService
	logger           *service.Logger
}

// NewMilestoneActions creates a new instance of MilestoneActions
func NewMilestoneActions(
	milestoneService *service.MilestoneService,
	viewService *service.ViewService,
	logger *service.Logger,
) *MilestoneActions {
	return &MilestoneActions{
		milestoneService: milestoneService,
		viewService:      viewService,
		logger:           logger,
	}
}

// Issues handles the issues action
func (m *MilestoneActions) Issues(ctx *gin.Context) error {
	// Get milestone from context
	milestone, err := ctx.Get("milestone")
	if err != nil {
		return err
	}

	// Get current user from context
	currentUser, err := ctx.Get("current_user")
	if err != nil {
		return err
	}

	// Handle response based on format
	switch ctx.GetHeader("Accept") {
	case "text/html":
		// Redirect to milestone page
		ctx.Redirect(http.StatusSeeOther, m.milestoneRedirectPath(ctx))
	case "application/json":
		// Get show project name param
		showProjectName := false
		if showProjectNameStr := ctx.DefaultQuery("show_project_name", "false"); showProjectNameStr == "true" {
			showProjectName = true
		}

		// Get sorted issues
		issues, err := milestone.(interface{ SortedIssues(user interface{}) ([]*service.Issue, error) }).SortedIssues(currentUser)
		if err != nil {
			return err
		}

		// Render JSON response
		ctx.JSON(http.StatusOK, m.tabsJSON(ctx, "shared/milestones/_issues_tab", map[string]interface{}{
			"issues":            issues,
			"show_project_name": showProjectName,
		}))
	default:
		ctx.JSON(http.StatusNotAcceptable, gin.H{
			"error": "Unsupported format",
		})
	}

	return nil
}

// MergeRequests handles the merge requests action
func (m *MilestoneActions) MergeRequests(ctx *gin.Context) error {
	// Get milestone from context
	milestone, err := ctx.Get("milestone")
	if err != nil {
		return err
	}

	// Get current user from context
	currentUser, err := ctx.Get("current_user")
	if err != nil {
		return err
	}

	// Handle response based on format
	switch ctx.GetHeader("Accept") {
	case "text/html":
		// Redirect to milestone page
		ctx.Redirect(http.StatusSeeOther, m.milestoneRedirectPath(ctx))
	case "application/json":
		// Get show project name param
		showProjectName := false
		if showProjectNameStr := ctx.DefaultQuery("show_project_name", "false"); showProjectNameStr == "true" {
			showProjectName = true
		}

		// Get sorted merge requests
		mergeRequests, err := milestone.(interface{ SortedMergeRequests(user interface{}) ([]*service.MergeRequest, error) }).SortedMergeRequests(currentUser)
		if err != nil {
			return err
		}

		// Preload milestoneish associations
		mergeRequests, err = m.milestoneService.PreloadMilestoneishAssociations(mergeRequests)
		if err != nil {
			return err
		}

		// Render JSON response
		ctx.JSON(http.StatusOK, m.tabsJSON(ctx, "shared/milestones/_merge_requests_tab", map[string]interface{}{
			"merge_requests":    mergeRequests,
			"show_project_name": showProjectName,
		}))
	default:
		ctx.JSON(http.StatusNotAcceptable, gin.H{
			"error": "Unsupported format",
		})
	}

	return nil
}

// Participants handles the participants action
func (m *MilestoneActions) Participants(ctx *gin.Context) error {
	// Get milestone from context
	milestone, err := ctx.Get("milestone")
	if err != nil {
		return err
	}

	// Get current user from context
	currentUser, err := ctx.Get("current_user")
	if err != nil {
		return err
	}

	// Handle response based on format
	switch ctx.GetHeader("Accept") {
	case "text/html":
		// Redirect to milestone page
		ctx.Redirect(http.StatusSeeOther, m.milestoneRedirectPath(ctx))
	case "application/json":
		// Get issue participants
		users, err := milestone.(interface{ IssueParticipantsVisibleByUser(user interface{}) ([]*service.User, error) }).IssueParticipantsVisibleByUser(currentUser)
		if err != nil {
			return err
		}

		// Render JSON response
		ctx.JSON(http.StatusOK, m.tabsJSON(ctx, "shared/milestones/_participants_tab", map[string]interface{}{
			"users": users,
		}))
	default:
		ctx.JSON(http.StatusNotAcceptable, gin.H{
			"error": "Unsupported format",
		})
	}

	return nil
}

// Labels handles the labels action
func (m *MilestoneActions) Labels(ctx *gin.Context) error {
	// Get milestone from context
	milestone, err := ctx.Get("milestone")
	if err != nil {
		return err
	}

	// Get current user from context
	currentUser, err := ctx.Get("current_user")
	if err != nil {
		return err
	}

	// Handle response based on format
	switch ctx.GetHeader("Accept") {
	case "text/html":
		// Redirect to milestone page
		ctx.Redirect(http.StatusSeeOther, m.milestoneRedirectPath(ctx))
	case "application/json":
		// Get issue labels
		milestoneLabels, err := milestone.(interface{ IssueLabelsVisibleByUser(user interface{}) ([]*service.Label, error) }).IssueLabelsVisibleByUser(currentUser)
		if err != nil {
			return err
		}

		// Get resource parent
		resourceParent := milestone.(interface{ ResourceParent() interface{} }).ResourceParent()

		// Present labels
		presentedLabels := make([]*service.LabelPresenter, len(milestoneLabels))
		for i, label := range milestoneLabels {
			presentedLabels[i] = label.Present(map[string]interface{}{
				"issuable_subject": resourceParent,
			})
		}

		// Render JSON response
		ctx.JSON(http.StatusOK, m.tabsJSON(ctx, "shared/milestones/_labels_tab", map[string]interface{}{
			"labels": presentedLabels,
		}))
	default:
		ctx.JSON(http.StatusNotAcceptable, gin.H{
			"error": "Unsupported format",
		})
	}

	return nil
}

// Private methods

// TabsJSON returns the tabs JSON
func (m *MilestoneActions) tabsJSON(ctx *gin.Context, partial string, data map[string]interface{}) map[string]interface{} {
	// Render partial to HTML
	html, err := m.viewService.RenderPartial(ctx, partial, data)
	if err != nil {
		m.logger.Error("Failed to render partial", "partial", partial, "error", err)
		return map[string]interface{}{
			"html": "<div class='error'>Failed to render content</div>",
		}
	}

	return map[string]interface{}{
		"html": html,
	}
}

// MilestoneRedirectPath returns the milestone redirect path
func (m *MilestoneActions) milestoneRedirectPath(ctx *gin.Context) string {
	// Get milestone from context
	milestone, err := ctx.Get("milestone")
	if err != nil {
		return "/"
	}

	// Get milestone ID
	milestoneID := milestone.(interface{ ID() string }).ID()

	// Get project ID from context
	projectID, err := ctx.Get("project_id")
	if err != nil {
		return "/"
	}

	// Get group ID from context
	groupID, err := ctx.Get("group_id")
	if err != nil {
		return "/"
	}

	// Build redirect path
	if projectID != "" {
		return "/projects/" + projectID.(string) + "/milestones/" + milestoneID
	} else if groupID != "" {
		return "/groups/" + groupID.(string) + "/milestones/" + milestoneID
	}

	return "/"
}
