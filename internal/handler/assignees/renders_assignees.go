package assignees

import (
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
)

// RendersAssignees handles rendering of assignees for merge requests
type RendersAssignees struct {
	teamService TeamService
}

// TeamService defines the interface for team-related operations
type TeamService interface {
	MaxMemberAccessForUserIDs(projectID int64, userIDs []int64) map[int64]int
}

// NewRendersAssignees creates a new instance of RendersAssignees
func NewRendersAssignees(teamService TeamService) *RendersAssignees {
	return &RendersAssignees{
		teamService: teamService,
	}
}

// PreloadAssigneesForRender returns the maximum member access level for the assignees of a merge request
func (r *RendersAssignees) PreloadAssigneesForRender(mergeRequest *model.MergeRequest) map[int64]int {
	// Extract assignee IDs from the merge request
	assigneeIDs := make([]int64, 0, len(mergeRequest.Assignees))
	for _, assignee := range mergeRequest.Assignees {
		assigneeIDs = append(assigneeIDs, assignee.ID)
	}

	// Get the maximum member access for the assignee IDs
	return r.teamService.MaxMemberAccessForUserIDs(mergeRequest.ProjectID, assigneeIDs)
}
