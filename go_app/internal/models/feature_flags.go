package models

// FeatureFlags represents feature flags configuration
type FeatureFlags struct {
	// Add any necessary fields
}

// IsEnabled checks if a feature flag is enabled for a user
func (f *FeatureFlags) IsEnabled(flag string, user *User) bool {
	// TODO: Implement feature flag checking logic
	return false
} 