package service

import (
	"context"

	"github.com/jmadden/gitlab-redux/internal/model"
)

// FeatureFlagService defines the interface for feature flag operations
type FeatureFlagService interface {
	// IsFeatureEnabled checks if a feature is enabled for a user
	IsFeatureEnabled(ctx context.Context, featureName string, user *model.User) bool

	// GetFeatureFlags gets all feature flags for a user
	GetFeatureFlags(ctx context.Context, user *model.User) map[string]bool

	// SetFeatureFlag sets a feature flag for a user
	SetFeatureFlag(ctx context.Context, featureName string, enabled bool, user *model.User) error
}
