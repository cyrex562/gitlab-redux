package service

import (
	"context"
	"database/sql"
	"fmt"
	"html/template"
	"time"
)

type BroadcastMessageService struct {
	db *sql.DB
}

func NewBroadcastMessageService(db *sql.DB) *BroadcastMessageService {
	return &BroadcastMessageService{
		db: db,
	}
}

type BroadcastTheme string

const (
	ThemeInfo    BroadcastTheme = "info"
	ThemeSuccess BroadcastTheme = "success"
	ThemeWarning BroadcastTheme = "warning"
	ThemeDanger  BroadcastTheme = "danger"
)

type BroadcastType string

const (
	TypeBanner BroadcastType = "banner"
	TypeNotification BroadcastType = "notification"
)

type BroadcastMessage struct {
	ID                int64
	Message           string
	Theme             BroadcastTheme
	BroadcastType     BroadcastType
	StartsAt          time.Time
	EndsAt            time.Time
	TargetPath        string
	Dismissable       bool
	ShowInCLI         bool
	TargetAccessLevels []string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (s *BroadcastMessageService) ListMessages(ctx context.Context, page int) ([]BroadcastMessage, error) {
	offset := (page - 1) * 20 // Assuming 20 items per page

	rows, err := s.db.QueryContext(ctx, `
		SELECT id, message, theme, broadcast_type, starts_at, ends_at,
		       target_path, dismissable, show_in_cli, target_access_levels,
		       created_at, updated_at
		FROM broadcast_messages
		ORDER BY ends_at DESC
		LIMIT 20 OFFSET $1
	`, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	var messages []BroadcastMessage
	for rows.Next() {
		var msg BroadcastMessage
		err := rows.Scan(
			&msg.ID, &msg.Message, &msg.Theme, &msg.BroadcastType,
			&msg.StartsAt, &msg.EndsAt, &msg.TargetPath, &msg.Dismissable,
			&msg.ShowInCLI, &msg.TargetAccessLevels, &msg.CreatedAt, &msg.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func (s *BroadcastMessageService) GetMessage(ctx context.Context, id int64) (*BroadcastMessage, error) {
	var msg BroadcastMessage
	err := s.db.QueryRowContext(ctx, `
		SELECT id, message, theme, broadcast_type, starts_at, ends_at,
		       target_path, dismissable, show_in_cli, target_access_levels,
		       created_at, updated_at
		FROM broadcast_messages
		WHERE id = $1
	`, id).Scan(
		&msg.ID, &msg.Message, &msg.Theme, &msg.BroadcastType,
		&msg.StartsAt, &msg.EndsAt, &msg.TargetPath, &msg.Dismissable,
		&msg.ShowInCLI, &msg.TargetAccessLevels, &msg.CreatedAt, &msg.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	return &msg, nil
}

func (s *BroadcastMessageService) CreateMessage(ctx context.Context, msg *BroadcastMessage) error {
	query := `
		INSERT INTO broadcast_messages (
			message, theme, broadcast_type, starts_at, ends_at,
			target_path, dismissable, show_in_cli, target_access_levels,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := s.db.QueryRowContext(ctx, query,
		msg.Message, msg.Theme, msg.BroadcastType, msg.StartsAt, msg.EndsAt,
		msg.TargetPath, msg.Dismissable, msg.ShowInCLI, msg.TargetAccessLevels,
	).Scan(&msg.ID, &msg.CreatedAt, &msg.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	return nil
}

func (s *BroadcastMessageService) UpdateMessage(ctx context.Context, msg *BroadcastMessage) error {
	query := `
		UPDATE broadcast_messages
		SET message = $1, theme = $2, broadcast_type = $3, starts_at = $4,
			ends_at = $5, target_path = $6, dismissable = $7, show_in_cli = $8,
			target_access_levels = $9, updated_at = NOW()
		WHERE id = $10
		RETURNING updated_at
	`

	err := s.db.QueryRowContext(ctx, query,
		msg.Message, msg.Theme, msg.BroadcastType, msg.StartsAt, msg.EndsAt,
		msg.TargetPath, msg.Dismissable, msg.ShowInCLI, msg.TargetAccessLevels,
		msg.ID,
	).Scan(&msg.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}

	return nil
}

func (s *BroadcastMessageService) DeleteMessage(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `
		DELETE FROM broadcast_messages
		WHERE id = $1
	`, id)

	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	return nil
}

func (s *BroadcastMessageService) GeneratePreview(ctx context.Context, msg *BroadcastMessage) (string, error) {
	// This is a simplified version of the preview generation
	// In a real implementation, you would want to use proper HTML templates
	// and potentially include more styling and formatting
	tmpl := `
		<div class="broadcast-message {{.Theme}}">
			{{if .Dismissable}}
				<button class="close" data-dismiss="alert">&times;</button>
			{{end}}
			<div class="message">{{.Message}}</div>
			<div class="meta">
				Starts: {{.StartsAt.Format "2006-01-02 15:04:05"}}
				Ends: {{.EndsAt.Format "2006-01-02 15:04:05"}}
			</div>
		</div>
	`

	t, err := template.New("preview").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var result string
	if err := t.Execute(&result, msg); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return result, nil
}
