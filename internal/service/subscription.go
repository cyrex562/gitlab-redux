package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmadden/gitlab-redux/internal/model"
)

var (
	ErrSubscribableNotFound = errors.New("subscribable not found")
	ErrProjectNotFound      = errors.New("project not found")
)

// SubscriptionService handles business logic for subscriptions
type SubscriptionService struct {
	db *sql.DB
}

// NewSubscriptionService creates a new subscription service
func NewSubscriptionService(db *sql.DB) *SubscriptionService {
	return &SubscriptionService{
		db: db,
	}
}

// ToggleSubscription toggles a subscription for a subscribable resource
func (s *SubscriptionService) ToggleSubscription(ctx context.Context, subscribable interface{}, project interface{}, userID int64) error {
	// Get subscribable info
	subscribableID, subscribableType, err := s.getSubscribableInfo(subscribable)
	if err != nil {
		return err
	}

	// Get project info
	projectID, err := s.getProjectInfo(project)
	if err != nil {
		return err
	}

	// Check if subscription exists
	exists, err := s.subscriptionExists(ctx, subscribableType, subscribableID, projectID, userID)
	if err != nil {
		return err
	}

	if exists {
		// Remove the subscription
		return s.removeSubscription(ctx, subscribableType, subscribableID, projectID, userID)
	}

	// Add the subscription
	return s.addSubscription(ctx, subscribableType, subscribableID, projectID, userID)
}

// getSubscribableInfo extracts ID and type from a subscribable
func (s *SubscriptionService) getSubscribableInfo(subscribable interface{}) (int64, string, error) {
	switch s := subscribable.(type) {
	case *model.Issue:
		return s.ID, "Issue", nil
	case *model.MergeRequest:
		return s.ID, "MergeRequest", nil
	case *model.Note:
		return s.ID, "Note", nil
	default:
		return 0, "", ErrSubscribableNotFound
	}
}

// getProjectInfo extracts ID from a project
func (s *SubscriptionService) getProjectInfo(project interface{}) (int64, error) {
	switch p := project.(type) {
	case *model.Project:
		return p.ID, nil
	default:
		return 0, ErrProjectNotFound
	}
}

// subscriptionExists checks if a subscription exists
func (s *SubscriptionService) subscriptionExists(ctx context.Context, subscribableType string, subscribableID, projectID, userID int64) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM subscriptions
			WHERE subscribable_type = $1 AND subscribable_id = $2 AND project_id = $3 AND user_id = $4
		)
	`

	var exists bool
	err := s.db.QueryRowContext(ctx, query, subscribableType, subscribableID, projectID, userID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// addSubscription adds a new subscription
func (s *SubscriptionService) addSubscription(ctx context.Context, subscribableType string, subscribableID, projectID, userID int64) error {
	query := `
		INSERT INTO subscriptions (subscribable_type, subscribable_id, project_id, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $5)
	`

	_, err := s.db.ExecContext(ctx, query, subscribableType, subscribableID, projectID, userID, time.Now())
	return err
}

// removeSubscription removes a subscription
func (s *SubscriptionService) removeSubscription(ctx context.Context, subscribableType string, subscribableID, projectID, userID int64) error {
	query := `
		DELETE FROM subscriptions
		WHERE subscribable_type = $1 AND subscribable_id = $2 AND project_id = $3 AND user_id = $4
	`

	_, err := s.db.ExecContext(ctx, query, subscribableType, subscribableID, projectID, userID)
	return err
}
