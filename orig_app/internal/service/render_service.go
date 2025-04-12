package service

import (
	"context"

	"github.com/jmadden/gitlab-redux/internal/model"
)

// RenderService defines the interface for rendering operations
type RenderService interface {
	// RenderEvents renders the given events
	RenderEvents(ctx context.Context, events []*model.Event, currentUser *model.User, isAtomRequest bool) ([]*model.RenderedEvent, error)

	// RenderProjects renders the given projects
	RenderProjects(ctx context.Context, projects []*model.Project, currentUser *model.User) ([]*model.RenderedProject, error)

	// RenderGroups renders the given groups
	RenderGroups(ctx context.Context, groups []*model.Group, currentUser *model.User) ([]*model.RenderedGroup, error)

	// RenderUsers renders the given users
	RenderUsers(ctx context.Context, users []*model.User, currentUser *model.User) ([]*model.RenderedUser, error)

	// GetNoteableMetaData gets the metadata for noteable items
	GetNoteableMetaData(ctx context.Context, items interface{}, itemType string) (map[string]interface{}, error)
}
