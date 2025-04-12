package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmadden/gitlab-redux/internal/model"
)

var (
	ErrInvalidIssuableType = errors.New("invalid issuable type")
	ErrInvalidIssuableID   = errors.New("invalid issuable ID")
	ErrTodoNotFound        = errors.New("todo not found")
)

// TodoService handles business logic for todos
type TodoService struct {
	db *sql.DB
}

// NewTodoService creates a new todo service
func NewTodoService(db *sql.DB) *TodoService {
	return &TodoService{
		db: db,
	}
}

// MarkTodo creates a new todo for an issuable
func (s *TodoService) MarkTodo(ctx context.Context, issuableType string, issuableID, userID int64) (*model.Todo, error) {
	if issuableID <= 0 {
		return nil, ErrInvalidIssuableID
	}

	// Validate issuable type
	switch issuableType {
	case "issue", "merge_request":
		// Valid types
	default:
		return nil, ErrInvalidIssuableType
	}

	// Check if todo already exists
	exists, err := s.todoExists(ctx, issuableType, issuableID, userID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, nil // Todo already exists, return silently
	}

	// Create new todo
	query := `
		INSERT INTO todos (user_id, issuable_type, issuable_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $4)
		RETURNING id, user_id, issuable_type, issuable_id, created_at, updated_at
	`

	todo := &model.Todo{}
	err = s.db.QueryRowContext(ctx, query, userID, issuableType, issuableID, time.Now()).
		Scan(&todo.ID, &todo.UserID, &todo.IssuableType, &todo.IssuableID, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

// GetPendingCount returns the number of pending todos for a user
func (s *TodoService) GetPendingCount(ctx context.Context, userID int64) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM todos
		WHERE user_id = $1 AND state = 'pending'
	`

	var count int
	err := s.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// todoExists checks if a todo already exists for the given parameters
func (s *TodoService) todoExists(ctx context.Context, issuableType string, issuableID, userID int64) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM todos
			WHERE user_id = $1 AND issuable_type = $2 AND issuable_id = $3
		)
	`

	var exists bool
	err := s.db.QueryRowContext(ctx, query, userID, issuableType, issuableID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
