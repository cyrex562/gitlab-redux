package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmadden/gitlab-redux/internal/model"
)

var (
	ErrInvalidEmojiName = errors.New("invalid emoji name")
	ErrAwardableNotFound = errors.New("awardable not found")
)

// AwardEmojiService handles business logic for award emojis
type AwardEmojiService struct {
	db *sql.DB
}

// NewAwardEmojiService creates a new award emoji service
func NewAwardEmojiService(db *sql.DB) *AwardEmojiService {
	return &AwardEmojiService{
		db: db,
	}
}

// Toggle toggles an award emoji on an awardable
func (s *AwardEmojiService) Toggle(ctx context.Context, awardable interface{}, name string, userID int64) (bool, error) {
	if name == "" {
		return false, ErrInvalidEmojiName
	}

	// Get awardable ID and type
	awardableID, awardableType, err := s.getAwardableInfo(awardable)
	if err != nil {
		return false, err
	}

	// Check if award already exists
	exists, err := s.awardExists(ctx, awardableType, awardableID, name, userID)
	if err != nil {
		return false, err
	}

	if exists {
		// Remove the award
		err = s.removeAward(ctx, awardableType, awardableID, name, userID)
		if err != nil {
			return false, err
		}
		return false, nil
	}

	// Add the award
	err = s.addAward(ctx, awardableType, awardableID, name, userID)
	if err != nil {
		return false, err
	}

	return true, nil
}

// getAwardableInfo extracts ID and type from an awardable
func (s *AwardEmojiService) getAwardableInfo(awardable interface{}) (int64, string, error) {
	switch a := awardable.(type) {
	case *model.Issue:
		return a.ID, "Issue", nil
	case *model.MergeRequest:
		return a.ID, "MergeRequest", nil
	case *model.Note:
		return a.ID, "Note", nil
	default:
		return 0, "", ErrAwardableNotFound
	}
}

// awardExists checks if an award emoji exists
func (s *AwardEmojiService) awardExists(ctx context.Context, awardableType string, awardableID int64, name string, userID int64) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM award_emojis
			WHERE awardable_type = $1 AND awardable_id = $2 AND name = $3 AND user_id = $4
		)
	`

	var exists bool
	err := s.db.QueryRowContext(ctx, query, awardableType, awardableID, name, userID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// addAward adds a new award emoji
func (s *AwardEmojiService) addAward(ctx context.Context, awardableType string, awardableID int64, name string, userID int64) error {
	query := `
		INSERT INTO award_emojis (awardable_type, awardable_id, name, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $5)
	`

	_, err := s.db.ExecContext(ctx, query, awardableType, awardableID, name, userID, time.Now())
	return err
}

// removeAward removes an award emoji
func (s *AwardEmojiService) removeAward(ctx context.Context, awardableType string, awardableID int64, name string, userID int64) error {
	query := `
		DELETE FROM award_emojis
		WHERE awardable_type = $1 AND awardable_id = $2 AND name = $3 AND user_id = $4
	`

	_, err := s.db.ExecContext(ctx, query, awardableType, awardableID, name, userID)
	return err
}
