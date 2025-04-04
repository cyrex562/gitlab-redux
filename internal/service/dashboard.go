package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type DashboardService struct {
	db *sql.DB
}

func NewDashboardService(db *sql.DB) *DashboardService {
	return &DashboardService{
		db: db,
	}
}

type ApproximateCounts struct {
	Projects int64 `json:"projects"`
	Users    int64 `json:"users"`
	Groups   int64 `json:"groups"`
}

type Project struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Group struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SystemNotice struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type UserStatistics struct {
	TotalUsers     int64 `json:"total_users"`
	ActiveUsers    int64 `json:"active_users"`
	InactiveUsers  int64 `json:"inactive_users"`
	BlockedUsers   int64 `json:"blocked_users"`
	AdminUsers     int64 `json:"admin_users"`
	ExternalUsers  int64 `json:"external_users"`
	BotUsers       int64 `json:"bot_users"`
	BlockedUsers   int64 `json:"blocked_users"`
	CreatedToday   int64 `json:"created_today"`
	CreatedThisWeek int64 `json:"created_this_week"`
	CreatedThisMonth int64 `json:"created_this_month"`
}

// GetApproximateCounts gets approximate counts for projects, users, and groups
func (s *DashboardService) GetApproximateCounts(ctx context.Context) (*ApproximateCounts, error) {
	var counts ApproximateCounts

	// Get project count
	if err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM projects WHERE deleted_at IS NULL
	`).Scan(&counts.Projects); err != nil {
		return nil, fmt.Errorf("failed to get project count: %w", err)
	}

	// Get user count
	if err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM users WHERE state = 'active'
	`).Scan(&counts.Users); err != nil {
		return nil, fmt.Errorf("failed to get user count: %w", err)
	}

	// Get group count
	if err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM namespaces WHERE type = 'Group'
	`).Scan(&counts.Groups); err != nil {
		return nil, fmt.Errorf("failed to get group count: %w", err)
	}

	return &counts, nil
}

// GetRecentProjects gets the most recent projects
func (s *DashboardService) GetRecentProjects(ctx context.Context, limit int) ([]Project, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT p.id, p.name, p.path, p.created_at, p.updated_at
		FROM projects p
		WHERE p.deleted_at IS NULL
		ORDER BY p.id DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query projects: %w", err)
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Path, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}
		projects = append(projects, p)
	}

	return projects, nil
}

// GetRecentUsers gets the most recent users
func (s *DashboardService) GetRecentUsers(ctx context.Context, limit int) ([]User, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, username, email, created_at, updated_at
		FROM users
		WHERE state = 'active'
		ORDER BY id DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, u)
	}

	return users, nil
}

// GetRecentGroups gets the most recent groups
func (s *DashboardService) GetRecentGroups(ctx context.Context, limit int) ([]Group, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, path, created_at, updated_at
		FROM namespaces
		WHERE type = 'Group'
		ORDER BY id DESC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query groups: %w", err)
	}
	defer rows.Close()

	var groups []Group
	for rows.Next() {
		var g Group
		if err := rows.Scan(&g.ID, &g.Name, &g.Path, &g.CreatedAt, &g.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan group: %w", err)
		}
		groups = append(groups, g)
	}

	return groups, nil
}

// GetSystemNotices gets system notices and warnings
func (s *DashboardService) GetSystemNotices(ctx context.Context) ([]SystemNotice, error) {
	var notices []SystemNotice

	// Check database connection
	if err := s.db.PingContext(ctx); err != nil {
		notices = append(notices, SystemNotice{
			Type:    "warning",
			Message: "Database connection issues detected",
		})
	}

	// Check Redis connection
	// This would be implemented in the Redis service

	// Check external services
	// This would be implemented in a separate service

	return notices, nil
}

// GetUserStatistics gets detailed user statistics
func (s *DashboardService) GetUserStatistics(ctx context.Context) (*UserStatistics, error) {
	var stats UserStatistics

	// Get total users
	if err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM users
	`).Scan(&stats.TotalUsers); err != nil {
		return nil, fmt.Errorf("failed to get total users: %w", err)
	}

	// Get active users
	if err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM users
		WHERE state = 'active' AND last_activity_at >= NOW() - INTERVAL '30 days'
	`).Scan(&stats.ActiveUsers); err != nil {
		return nil, fmt.Errorf("failed to get active users: %w", err)
	}

	// Get inactive users
	if err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM users
		WHERE state = 'active' AND last_activity_at < NOW() - INTERVAL '30 days'
	`).Scan(&stats.InactiveUsers); err != nil {
		return nil, fmt.Errorf("failed to get inactive users: %w", err)
	}

	// Get admin users
	if err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM users
		WHERE admin = true
	`).Scan(&stats.AdminUsers); err != nil {
		return nil, fmt.Errorf("failed to get admin users: %w", err)
	}

	// Get external users
	if err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM users
		WHERE external = true
	`).Scan(&stats.ExternalUsers); err != nil {
		return nil, fmt.Errorf("failed to get external users: %w", err)
	}

	// Get bot users
	if err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM users
		WHERE bot = true
	`).Scan(&stats.BotUsers); err != nil {
		return nil, fmt.Errorf("failed to get bot users: %w", err)
	}

	// Get blocked users
	if err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM users
		WHERE state = 'blocked'
	`).Scan(&stats.BlockedUsers); err != nil {
		return nil, fmt.Errorf("failed to get blocked users: %w", err)
	}

	// Get users created today
	if err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM users
		WHERE created_at >= CURRENT_DATE
	`).Scan(&stats.CreatedToday); err != nil {
		return nil, fmt.Errorf("failed to get users created today: %w", err)
	}

	// Get users created this week
	if err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM users
		WHERE created_at >= CURRENT_DATE - INTERVAL '7 days'
	`).Scan(&stats.CreatedThisWeek); err != nil {
		return nil, fmt.Errorf("failed to get users created this week: %w", err)
	}

	// Get users created this month
	if err := s.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM users
		WHERE created_at >= CURRENT_DATE - INTERVAL '30 days'
	`).Scan(&stats.CreatedThisMonth); err != nil {
		return nil, fmt.Errorf("failed to get users created this month: %w", err)
	}

	return &stats, nil
}
