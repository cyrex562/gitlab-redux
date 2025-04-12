package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmadden/gitlab-redux/internal/model"
)

var (
	ErrDismissalNotFound = errors.New("dismissal not found")
	ErrInvalidMessageID  = errors.New("invalid message ID")
)

// BroadcastMessageDismissalService handles business logic for broadcast message dismissals
type BroadcastMessageDismissalService struct {
	db *sql.DB
}

// NewBroadcastMessageDismissalService creates a new broadcast message dismissal service
func NewBroadcastMessageDismissalService(db *sql.DB) *BroadcastMessageDismissalService {
	return &BroadcastMessageDismissalService{
		db: db,
	}
}

// GetDismissalsForUser retrieves all dismissals for a user
func (s *BroadcastMessageDismissalService) GetDismissalsForUser(ctx context.Context, userID int64) ([]*model.BroadcastMessageDismissal, error) {
	var dismissals []*model.BroadcastMessageDismissal
	err := s.db.QueryRowContext(ctx, "SELECT id, user_id, message_id, expires_at, created_at, updated_at FROM broadcast_message_dismissals WHERE user_id = $1", userID).Scan(
		&dismissals[0].ID,
		&dismissals[0].UserID,
		&dismissals[0].MessageID,
		&dismissals[0].ExpiresAt,
		&dismissals[0].CreatedAt,
		&dismissals[0].UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return dismissals, nil
}

// CreateDismissal creates a new dismissal for a message
func (s *BroadcastMessageDismissalService) CreateDismissal(ctx context.Context, userID, messageID int64) (*model.BroadcastMessageDismissal, error) {
	if messageID <= 0 {
		return nil, ErrInvalidMessageID
	}

	// Set expiration to 30 days from now
	expiresAt := time.Now().AddDate(0, 0, 30)

	query := `
		INSERT INTO broadcast_message_dismissals (user_id, message_id, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $4)
		RETURNING id, user_id, message_id, expires_at, created_at, updated_at
	`

	dismissal := &model.BroadcastMessageDismissal{}
	err := s.db.QueryRowContext(ctx, query, userID, messageID, expiresAt, time.Now()).
		Scan(&dismissal.ID, &dismissal.UserID, &dismissal.MessageID, &dismissal.ExpiresAt, &dismissal.CreatedAt, &dismissal.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return dismissal, nil
}

// GetDismissal retrieves a dismissal by user ID and message ID
func (s *BroadcastMessageDismissalService) GetDismissal(ctx context.Context, userID, messageID int64) (*model.BroadcastMessageDismissal, error) {
	query := `
		SELECT id, user_id, message_id, expires_at, created_at, updated_at
		FROM broadcast_message_dismissals
		WHERE user_id = $1 AND message_id = $2
	`

	dismissal := &model.BroadcastMessageDismissal{}
	err := s.db.QueryRowContext(ctx, query, userID, messageID).
		Scan(&dismissal.ID, &dismissal.UserID, &dismissal.MessageID, &dismissal.ExpiresAt, &dismissal.CreatedAt, &dismissal.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrDismissalNotFound
	}
	if err != nil {
		return nil, err
	}

	return dismissal, nil
}

// IsDismissed checks if a message has been dismissed by a user
func (s *BroadcastMessageDismissalService) IsDismissed(ctx context.Context, userID, messageID int64) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM broadcast_message_dismissals
			WHERE user_id = $1 AND message_id = $2 AND expires_at > $3
		)
	`

	var exists bool
	err := s.db.QueryRowContext(ctx, query, userID, messageID, time.Now()).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// SynchronizeDismissals synchronizes dismissals between the database and cookies
func (s *BroadcastMessageDismissalService) SynchronizeDismissals(ctx *gin.Context, userID int64) error {
	// Get dismissals from database
	dismissals, err := s.GetDismissalsForUser(ctx, userID)
	if err != nil {
		return err
	}

	// Create cookies for dismissals that don't have them
	for _, dismissal := range dismissals {
		cookieKey := dismissal.CookieKey()
		if _, err := ctx.Cookie(cookieKey); err != nil {
			// Cookie doesn't exist, create it
			ctx.SetCookie(
				cookieKey,
				"true",
				int(time.Until(dismissal.ExpiresAt).Seconds()),
				"/",
				"",
				false, // Not secure by default, can be configured based on environment
				true,  // HTTP only
			)
		}
	}

	return nil
}
