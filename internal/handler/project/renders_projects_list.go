package project

import (
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
)

// RendersProjectsList handles rendering of projects in GitLab
type RendersProjectsList struct {
	memberAccess *MemberAccess
	projectService ProjectService
}

// ProjectService defines the interface for project-related operations
type ProjectService interface {
	PreloadForksCount(projects []*model.Project)
	PreloadOpenIssuesCount(projects []*model.Project)
	PreloadOpenMergeRequestsCount(projects []*model.Project)
}

// NewRendersProjectsList creates a new instance of RendersProjectsList
func NewRendersProjectsList(memberAccess *MemberAccess, projectService ProjectService) *RendersProjectsList {
	return &RendersProjectsList{
		memberAccess: memberAccess,
		projectService: projectService,
	}
}

// PrepareProjectsForRendering prepares projects for rendering by preloading member access and roles
func (r *RendersProjectsList) PrepareProjectsForRendering(projects []*model.Project, currentUser *model.User) []*model.Project {
	// Preload max member access for the collection
	r.memberAccess.PreloadMaxMemberAccessForCollection(model.ProjectType, projects)

	// Preload member roles if there's a current user
	if currentUser != nil {
		r.PreloadMemberRoles(projects)
	}

	// Preload counts for each project
	r.projectService.PreloadForksCount(projects)
	r.projectService.PreloadOpenIssuesCount(projects)
	r.projectService.PreloadOpenMergeRequestsCount(projects)

	return projects
}

// PreloadMemberRoles preloads member roles for projects
// This method can be overridden in EE (Enterprise Edition)
func (r *RendersProjectsList) PreloadMemberRoles(projects []*model.Project) {
	// Default implementation does nothing
	// This can be overridden in EE
}
