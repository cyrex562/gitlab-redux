package service

import (
	"context"

	"github.com/jmadden/gitlab-redux/internal/model"
)

// AnalyticsService defines the interface for analytics-related operations
type AnalyticsService interface {
	// TrackInternalEvent tracks an internal event
	TrackInternalEvent(ctx context.Context, eventName string, user *model.User) error

	// TrackUserEvent tracks a user event
	TrackUserEvent(ctx context.Context, eventName string, user *model.User, properties map[string]interface{}) error

	// TrackGroupEvent tracks a group event
	TrackGroupEvent(ctx context.Context, eventName string, group *model.Group, user *model.User, properties map[string]interface{}) error

	// TrackProjectEvent tracks a project event
	TrackProjectEvent(ctx context.Context, eventName string, project *model.Project, user *model.User, properties map[string]interface{}) error
}
