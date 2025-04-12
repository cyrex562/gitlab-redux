package service

import (
	"context"
	"time"

	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gorm.io/gorm"
)

// UsageCounts represents various usage metrics
type UsageCounts struct {
	Projects            int `json:"projects"`
	Groups             int `json:"groups"`
	Issues             int `json:"issues"`
	MergeRequests      int `json:"merge_requests"`
	CIInternalPipelines int `json:"ci_internal_pipelines"`
	CIExternalPipelines int `json:"ci_external_pipelines"`
	Labels             int `json:"labels"`
	Milestones         int `json:"milestones"`
	Snippets           int `json:"snippets"`
	Notes              int `json:"notes"`
}

// UsageData represents the complete usage data
type UsageData struct {
	ActiveUserCount int         `json:"active_user_count"`
	Counts         UsageCounts `json:"counts"`
	Timestamp      time.Time   `json:"timestamp"`
}

// UsageService handles business logic for usage data
type UsageService struct {
	db *gorm.DB
}

// NewUsageService creates a new UsageService instance
func NewUsageService(db *gorm.DB) *UsageService {
	return &UsageService{
		db: db,
	}
}

// GetServicePingData retrieves the service ping data
func (s *UsageService) GetServicePingData(ctx context.Context) (*UsageData, error) {
	// Get active user count
	var activeUserCount int
	err := s.db.WithContext(ctx).
		Model(&model.User{}).
		Where("state = ? AND last_activity_on > ?", "active", time.Now().AddDate(0, -1, 0)).
		Count(&activeUserCount).Error
	if err != nil {
		return nil, err
	}

	// Get project count
	var projectCount int
	err = s.db.WithContext(ctx).
		Model(&model.Project{}).
		Count(&projectCount).Error
	if err != nil {
		return nil, err
	}

	// Get group count
	var groupCount int
	err = s.db.WithContext(ctx).
		Model(&model.Group{}).
		Count(&groupCount).Error
	if err != nil {
		return nil, err
	}

	// Get issue count
	var issueCount int
	err = s.db.WithContext(ctx).
		Model(&model.Issue{}).
		Count(&issueCount).Error
	if err != nil {
		return nil, err
	}

	// Get merge request count
	var mergeRequestCount int
	err = s.db.WithContext(ctx).
		Model(&model.MergeRequest{}).
		Count(&mergeRequestCount).Error
	if err != nil {
		return nil, err
	}

	// Get pipeline counts
	var internalPipelineCount, externalPipelineCount int
	err = s.db.WithContext(ctx).
		Model(&model.Pipeline{}).
		Where("source = ?", "internal").
		Count(&internalPipelineCount).Error
	if err != nil {
		return nil, err
	}

	err = s.db.WithContext(ctx).
		Model(&model.Pipeline{}).
		Where("source = ?", "external").
		Count(&externalPipelineCount).Error
	if err != nil {
		return nil, err
	}

	// Get label count
	var labelCount int
	err = s.db.WithContext(ctx).
		Model(&model.Label{}).
		Count(&labelCount).Error
	if err != nil {
		return nil, err
	}

	// Get milestone count
	var milestoneCount int
	err = s.db.WithContext(ctx).
		Model(&model.Milestone{}).
		Count(&milestoneCount).Error
	if err != nil {
		return nil, err
	}

	// Get snippet count
	var snippetCount int
	err = s.db.WithContext(ctx).
		Model(&model.Snippet{}).
		Count(&snippetCount).Error
	if err != nil {
		return nil, err
	}

	// Get note count
	var noteCount int
	err = s.db.WithContext(ctx).
		Model(&model.Note{}).
		Count(&noteCount).Error
	if err != nil {
		return nil, err
	}

	// Build usage data
	usageData := &UsageData{
		ActiveUserCount: activeUserCount,
		Counts: UsageCounts{
			Projects:             projectCount,
			Groups:              groupCount,
			Issues:              issueCount,
			MergeRequests:       mergeRequestCount,
			CIInternalPipelines: internalPipelineCount,
			CIExternalPipelines: externalPipelineCount,
			Labels:              labelCount,
			Milestones:          milestoneCount,
			Snippets:            snippetCount,
			Notes:               noteCount,
		},
		Timestamp: time.Now(),
	}

	return usageData, nil
}
