package services

import (
	"gorm.io/gorm"
)

// FeatureService handles feature flag operations
type FeatureService struct {
	db *gorm.DB
}

// NewFeatureService creates a new FeatureService
func NewFeatureService(db *gorm.DB) *FeatureService {
	return &FeatureService{
		db: db,
	}
}

// IsEnabled checks if a feature is enabled for a user
func (s *FeatureService) IsEnabled(featureName string, userID uint) bool {
	// TODO: Implement proper feature flag check
	// For now, return true for all features
	return true
} 