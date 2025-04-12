package feature_flags

import (
	"github.com/gin-gonic/gin"
	"github.com/jmadden/gitlab-redux/internal/model"
	"github.com/jmadden/gitlab-redux/internal/service"
)

// PushFrontendFeatureFlag pushes a frontend feature flag to the context
func PushFrontendFeatureFlag(ctx *gin.Context, flagName string, user *model.User) bool {
	// Get feature flag service from context
	featureFlagService, ok := ctx.MustGet("feature_flag_service").(service.FeatureFlagService)
	if !ok {
		return false
	}

	// Check if feature is enabled
	enabled := featureFlagService.IsFeatureEnabled(ctx, flagName, user)

	// Set feature flag in context
	ctx.Set(flagName, enabled)

	return enabled
}

// GetFeatureFlag gets a feature flag from the context
func GetFeatureFlag(ctx *gin.Context, flagName string) bool {
	flag, exists := ctx.Get(flagName)
	if !exists {
		return false
	}

	enabled, ok := flag.(bool)
	if !ok {
		return false
	}

	return enabled
}

// SetFeatureFlag sets a feature flag in the context
func SetFeatureFlag(ctx *gin.Context, flagName string, enabled bool) {
	ctx.Set(flagName, enabled)
}
