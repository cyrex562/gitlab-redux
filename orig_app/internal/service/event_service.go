package service

import (
	"context"

	"github.com/jmadden/gitlab-redux/internal/model"
)

// EventService defines the interface for event-related operations
type EventService interface {
	// GetEventsForProjects gets events for the given projects
	GetEventsForProjects(ctx context.Context, projects []*model.Project, offset int, filter string) ([]*model.Event, error)

	// GetEventsForUser gets events for the given user
	GetEventsForUser(ctx context.Context, user *model.User, offset int, filter string) ([]*model.Event, error)

	// GetEventsForGroup gets events for the given group
	GetEventsForGroup(ctx context.Context, group *model.Group, offset int, filter string) ([]*model.Event, error)
}
