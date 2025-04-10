package issuable

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// IssuesCalendar handles rendering issues in a calendar format
type IssuesCalendar struct {
	configService *service.ConfigService
	logger        *service.Logger
}

// NewIssuesCalendar creates a new instance of IssuesCalendar
func NewIssuesCalendar(
	configService *service.ConfigService,
	logger *service.Logger,
) *IssuesCalendar {
	return &IssuesCalendar{
		configService: configService,
		logger:        logger,
	}
}

// RenderIssuesCalendar renders issues in a calendar format
func (i *IssuesCalendar) RenderIssuesCalendar(ctx *gin.Context, issuables []model.Issuable) {
	// Filter non-archived issues with due dates and limit to 100
	issues := i.filterIssuesForCalendar(issuables)

	// Set the issues in the context
	ctx.Set("issues", issues)

	// Check if the request is for ICS format
	format := ctx.DefaultQuery("format", "html")
	if format == "ics" {
		// Check if the request is from GitLab
		referer := ctx.Request.Referer()
		baseURL := i.configService.GetBaseURL()

		// If the referer starts with the GitLab base URL, set Content-Type to text/plain
		// to display the content inline instead of downloading it
		if referer != "" && strings.HasPrefix(referer, baseURL) {
			ctx.Header("Content-Type", "text/plain")
		} else {
			ctx.Header("Content-Type", "text/calendar")
		}

		// Render the ICS template
		ctx.HTML(http.StatusOK, "issues_calendar.ics", gin.H{
			"issues": issues,
		})
		return
	}

	// For other formats, render the default template
	ctx.HTML(http.StatusOK, "issues_calendar.html", gin.H{
		"issues": issues,
	})
}

// filterIssuesForCalendar filters issues for the calendar view
func (i *IssuesCalendar) filterIssuesForCalendar(issuables []model.Issuable) []model.Issuable {
	// Create a new slice for the filtered issues
	filteredIssues := make([]model.Issuable, 0, len(issuables))

	// Filter the issues
	for _, issuable := range issuables {
		// Skip archived issues
		if issuable.Archived {
			continue
		}

		// Skip issues without due dates
		if issuable.DueDate.IsZero() {
			continue
		}

		// Add the issue to the filtered list
		filteredIssues = append(filteredIssues, issuable)

		// Limit to 100 issues
		if len(filteredIssues) >= 100 {
			break
		}
	}

	return filteredIssues
}
