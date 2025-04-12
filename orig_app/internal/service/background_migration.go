package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type BackgroundMigrationService struct {
	db *sql.DB
}

func NewBackgroundMigrationService(db *sql.DB) *BackgroundMigrationService {
	return &BackgroundMigrationService{
		db: db,
	}
}

type MigrationStatus string

const (
	StatusQueued    MigrationStatus = "queued"
	StatusFinalizing MigrationStatus = "finalizing"
	StatusFailed    MigrationStatus = "failed"
	StatusFinished  MigrationStatus = "finished"
)

type Migration struct {
	ID              int64
	JobClass        string
	TableName       string
	ColumnName      string
	JobArguments    []interface{}
	BatchSize       int
	SubBatchSize    int
	Interval        int
	MinValue        int64
	MaxValue        int64
	Status          MigrationStatus
	CreatedAt       time.Time
	UpdatedAt       time.Time
	StartedAt       *time.Time
	FinishedAt      *time.Time
	SuccessfulRows  int64
	FailedJobsCount int64
}

type FailedJob struct {
	ID         int64
	MigrationID int64
	BatchSize   int
	SubBatchSize int
	Status      string
	Error       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (s *BackgroundMigrationService) ListMigrations(ctx context.Context, tab string, page int, database string) ([]Migration, error) {
	offset := (page - 1) * 20 // Assuming 20 items per page

	query := `
		SELECT id, job_class, table_name, column_name, job_arguments,
		       batch_size, sub_batch_size, interval, min_value, max_value,
		       status, created_at, updated_at, started_at, finished_at,
		       (SELECT COUNT(*) FROM batched_jobs WHERE migration_id = batched_migrations.id AND status = 'succeeded') as successful_rows,
		       (SELECT COUNT(*) FROM batched_jobs WHERE migration_id = batched_migrations.id AND status = 'failed') as failed_jobs_count
		FROM batched_migrations
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT 20 OFFSET $2
	`

	rows, err := s.db.QueryContext(ctx, query, tab, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query migrations: %w", err)
	}
	defer rows.Close()

	var migrations []Migration
	for rows.Next() {
		var m Migration
		err := rows.Scan(
			&m.ID, &m.JobClass, &m.TableName, &m.ColumnName, &m.JobArguments,
			&m.BatchSize, &m.SubBatchSize, &m.Interval, &m.MinValue, &m.MaxValue,
			&m.Status, &m.CreatedAt, &m.UpdatedAt, &m.StartedAt, &m.FinishedAt,
			&m.SuccessfulRows, &m.FailedJobsCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan migration: %w", err)
		}
		migrations = append(migrations, m)
	}

	return migrations, nil
}

func (s *BackgroundMigrationService) GetMigration(ctx context.Context, id int64, page int) (*Migration, []FailedJob, error) {
	offset := (page - 1) * 20

	// Get migration details
	var m Migration
	err := s.db.QueryRowContext(ctx, `
		SELECT id, job_class, table_name, column_name, job_arguments,
		       batch_size, sub_batch_size, interval, min_value, max_value,
		       status, created_at, updated_at, started_at, finished_at,
		       (SELECT COUNT(*) FROM batched_jobs WHERE migration_id = batched_migrations.id AND status = 'succeeded') as successful_rows,
		       (SELECT COUNT(*) FROM batched_jobs WHERE migration_id = batched_migrations.id AND status = 'failed') as failed_jobs_count
		FROM batched_migrations
		WHERE id = $1
	`, id).Scan(
		&m.ID, &m.JobClass, &m.TableName, &m.ColumnName, &m.JobArguments,
		&m.BatchSize, &m.SubBatchSize, &m.Interval, &m.MinValue, &m.MaxValue,
		&m.Status, &m.CreatedAt, &m.UpdatedAt, &m.StartedAt, &m.FinishedAt,
		&m.SuccessfulRows, &m.FailedJobsCount,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get migration: %w", err)
	}

	// Get failed jobs
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, migration_id, batch_size, sub_batch_size, status, error, created_at, updated_at
		FROM batched_jobs
		WHERE migration_id = $1 AND status = 'failed'
		ORDER BY created_at DESC
		LIMIT 20 OFFSET $2
	`, id, offset)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query failed jobs: %w", err)
	}
	defer rows.Close()

	var failedJobs []FailedJob
	for rows.Next() {
		var j FailedJob
		err := rows.Scan(
			&j.ID, &j.MigrationID, &j.BatchSize, &j.SubBatchSize,
			&j.Status, &j.Error, &j.CreatedAt, &j.UpdatedAt,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan failed job: %w", err)
		}
		failedJobs = append(failedJobs, j)
	}

	return &m, failedJobs, nil
}

func (s *BackgroundMigrationService) PauseMigration(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE batched_migrations
		SET status = 'paused'
		WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("failed to pause migration: %w", err)
	}
	return nil
}

func (s *BackgroundMigrationService) ResumeMigration(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE batched_migrations
		SET status = 'queued'
		WHERE id = $1
	`, id)
	if err != nil {
		return fmt.Errorf("failed to resume migration: %w", err)
	}
	return nil
}

func (s *BackgroundMigrationService) RetryFailedJobs(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE batched_jobs
		SET status = 'queued', error = NULL
		WHERE migration_id = $1 AND status = 'failed'
	`, id)
	if err != nil {
		return fmt.Errorf("failed to retry failed jobs: %w", err)
	}
	return nil
}

func (s *BackgroundMigrationService) ListDatabases(ctx context.Context) ([]string, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT datname
		FROM pg_database
		WHERE datname LIKE 'gitlab_%'
		ORDER BY datname
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query databases: %w", err)
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var db string
		if err := rows.Scan(&db); err != nil {
			return nil, fmt.Errorf("failed to scan database name: %w", err)
		}
		databases = append(databases, db)
	}

	return databases, nil
}
