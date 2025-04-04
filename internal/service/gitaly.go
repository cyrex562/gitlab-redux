package service

import (
	"context"

	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gorm.io/gorm"
)

// GitalyService handles business logic for Gitaly servers
type GitalyService struct {
	db *gorm.DB
}

// NewGitalyService creates a new GitalyService instance
func NewGitalyService(db *gorm.DB) *GitalyService {
	return &GitalyService{
		db: db,
	}
}

// GetAllServers retrieves all Gitaly servers
func (s *GitalyService) GetAllServers(ctx context.Context) ([]model.GitalyServer, error) {
	var servers []model.GitalyServer
	err := s.db.WithContext(ctx).Find(&servers).Error
	if err != nil {
		return nil, err
	}
	return servers, nil
}
