package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type CohortsService struct {
	db    *sql.DB
	cache *redis.Client
}

func NewCohortsService(db *sql.DB, cache *redis.Client) *CohortsService {
	return &CohortsService{
		db:    db,
		cache: cache,
	}
}

type Cohort struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	UserCount       int       `json:"user_count"`
	RetentionRate   float64   `json:"retention_rate"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	LastActivityAt  time.Time `json:"last_activity_at"`
	ActiveUsers     int       `json:"active_users"`
	InactiveUsers   int       `json:"inactive_users"`
	ChurnRate       float64   `json:"churn_rate"`
}

// Execute retrieves and processes cohort data
func (s *CohortsService) Execute(ctx context.Context) ([]Cohort, error) {
	// Query to get cohort data
	rows, err := s.db.QueryContext(ctx, `
		SELECT
			c.id,
			c.name,
			c.description,
			COUNT(DISTINCT u.id) as user_count,
			COALESCE(
				COUNT(DISTINCT CASE WHEN u.last_activity_at >= NOW() - INTERVAL '30 days' THEN u.id END)::float /
				NULLIF(COUNT(DISTINCT u.id), 0) * 100,
				0
			) as retention_rate,
			c.created_at,
			c.updated_at,
			MAX(u.last_activity_at) as last_activity_at,
			COUNT(DISTINCT CASE WHEN u.last_activity_at >= NOW() - INTERVAL '7 days' THEN u.id END) as active_users,
			COUNT(DISTINCT CASE WHEN u.last_activity_at < NOW() - INTERVAL '30 days' THEN u.id END) as inactive_users,
			COALESCE(
				COUNT(DISTINCT CASE WHEN u.last_activity_at < NOW() - INTERVAL '30 days' THEN u.id END)::float /
				NULLIF(COUNT(DISTINCT u.id), 0) * 100,
				0
			) as churn_rate
		FROM cohorts c
		LEFT JOIN users u ON u.cohort_id = c.id
		GROUP BY c.id, c.name, c.description, c.created_at, c.updated_at
		ORDER BY c.created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query cohorts: %w", err)
	}
	defer rows.Close()

	var cohorts []Cohort
	for rows.Next() {
		var c Cohort
		err := rows.Scan(
			&c.ID, &c.Name, &c.Description, &c.UserCount,
			&c.RetentionRate, &c.CreatedAt, &c.UpdatedAt,
			&c.LastActivityAt, &c.ActiveUsers, &c.InactiveUsers,
			&c.ChurnRate,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan cohort: %w", err)
		}
		cohorts = append(cohorts, c)
	}

	return cohorts, nil
}

// GetFromCache retrieves cohort data from Redis cache
func (s *CohortsService) GetFromCache(ctx context.Context, key string) ([]Cohort, error) {
	data, err := s.cache.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var cohorts []Cohort
	if err := json.Unmarshal([]byte(data), &cohorts); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached cohorts: %w", err)
	}

	return cohorts, nil
}

// SetCache stores cohort data in Redis cache
func (s *CohortsService) SetCache(ctx context.Context, key string, cohorts []Cohort, expiration time.Duration) error {
	data, err := json.Marshal(cohorts)
	if err != nil {
		return fmt.Errorf("failed to marshal cohorts for cache: %w", err)
	}

	if err := s.cache.Set(ctx, key, data, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}
