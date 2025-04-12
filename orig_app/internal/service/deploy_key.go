package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type DeployKeyService struct {
	db *sql.DB
}

func NewDeployKeyService(db *sql.DB) *DeployKeyService {
	return &DeployKeyService{
		db: db,
	}
}

type DeployKey struct {
	ID        int64      `json:"id"`
	Title     string     `json:"title"`
	Key       string     `json:"key"`
	Public    bool       `json:"public"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type CreateDeployKeyParams struct {
	Title     string     `json:"title" binding:"required"`
	Key       string     `json:"key" binding:"required"`
	Public    bool       `json:"public"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type UpdateDeployKeyParams struct {
	Title string `json:"title" binding:"required"`
}

// GetPublicKeys gets all public deploy keys
func (s *DeployKeyService) GetPublicKeys(ctx context.Context) ([]DeployKey, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, title, key, public, expires_at, created_at, updated_at
		FROM deploy_keys
		WHERE public = true
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query deploy keys: %w", err)
	}
	defer rows.Close()

	var deployKeys []DeployKey
	for rows.Next() {
		var dk DeployKey
		if err := rows.Scan(&dk.ID, &dk.Title, &dk.Key, &dk.Public, &dk.ExpiresAt, &dk.CreatedAt, &dk.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan deploy key: %w", err)
		}
		deployKeys = append(deployKeys, dk)
	}

	return deployKeys, nil
}

// GetByID gets a deploy key by ID
func (s *DeployKeyService) GetByID(ctx context.Context, id string) (*DeployKey, error) {
	var dk DeployKey
	err := s.db.QueryRowContext(ctx, `
		SELECT id, title, key, public, expires_at, created_at, updated_at
		FROM deploy_keys
		WHERE id = $1 AND public = true
	`, id).Scan(&dk.ID, &dk.Title, &dk.Key, &dk.Public, &dk.ExpiresAt, &dk.CreatedAt, &dk.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("deploy key not found")
		}
		return nil, fmt.Errorf("failed to query deploy key: %w", err)
	}

	return &dk, nil
}

// Create creates a new deploy key
func (s *DeployKeyService) Create(ctx context.Context, params CreateDeployKeyParams) (*DeployKey, error) {
	var dk DeployKey
	err := s.db.QueryRowContext(ctx, `
		INSERT INTO deploy_keys (title, key, public, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, title, key, public, expires_at, created_at, updated_at
	`, params.Title, params.Key, params.Public, params.ExpiresAt).Scan(
		&dk.ID, &dk.Title, &dk.Key, &dk.Public, &dk.ExpiresAt, &dk.CreatedAt, &dk.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create deploy key: %w", err)
	}

	return &dk, nil
}

// Update updates a deploy key
func (s *DeployKeyService) Update(ctx context.Context, id string, params UpdateDeployKeyParams) (*DeployKey, error) {
	var dk DeployKey
	err := s.db.QueryRowContext(ctx, `
		UPDATE deploy_keys
		SET title = $1, updated_at = NOW()
		WHERE id = $2 AND public = true
		RETURNING id, title, key, public, expires_at, created_at, updated_at
	`, params.Title, id).Scan(
		&dk.ID, &dk.Title, &dk.Key, &dk.Public, &dk.ExpiresAt, &dk.CreatedAt, &dk.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("deploy key not found")
		}
		return nil, fmt.Errorf("failed to update deploy key: %w", err)
	}

	return &dk, nil
}

// Delete deletes a deploy key
func (s *DeployKeyService) Delete(ctx context.Context, id string) error {
	result, err := s.db.ExecContext(ctx, `
		DELETE FROM deploy_keys
		WHERE id = $1 AND public = true
	`, id)
	if err != nil {
		return fmt.Errorf("failed to delete deploy key: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("deploy key not found")
	}

	return nil
}
