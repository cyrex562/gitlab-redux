package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type BatchedJobService struct {
	db *sql.DB
}

func NewBatchedJobService(db *sql.DB) *BatchedJobService {
	return &BatchedJobService{
		db: db,
	}
}

type JobStatus string

const (
	JobStatusQueued     JobStatus = "queued"
	JobStatusRunning    JobStatus = "running"
	JobStatusSucceeded  JobStatus = "succeeded"
	JobStatusFailed     JobStatus = "failed"
	JobStatusCancelled  JobStatus = "cancelled"
)

type BatchedJob struct {
	ID           int64
	MigrationID  int64
	BatchSize    int
	SubBatchSize int
	Status       JobStatus
	Error        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	StartedAt    *time.Time
	FinishedAt   *time.Time
	MinValue     int64
	MaxValue     int64
}

type JobTransitionLog struct {
	ID        int64
	JobID     int64
	FromState JobStatus
	ToState   JobStatus
	CreatedAt time.Time
}

func (s *BatchedJobService) GetJob(ctx context.Context, id int64, database string) (*BatchedJob, []JobTransitionLog, error) {
	// Get job details
	var job BatchedJob
	err := s.db.QueryRowContext(ctx, `
		SELECT id, migration_id, batch_size, sub_batch_size, status, error,
		       created_at, updated_at, started_at, finished_at, min_value, max_value
		FROM batched_jobs
		WHERE id = $1
	`, id).Scan(
		&job.ID, &job.MigrationID, &job.BatchSize, &job.SubBatchSize,
		&job.Status, &job.Error, &job.CreatedAt, &job.UpdatedAt,
		&job.StartedAt, &job.FinishedAt, &job.MinValue, &job.MaxValue,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get job: %w", err)
	}

	// Get transition logs
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, job_id, from_state, to_state, created_at
		FROM batched_job_transition_logs
		WHERE job_id = $1
		ORDER BY created_at DESC
	`, id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query transition logs: %w", err)
	}
	defer rows.Close()

	var transitionLogs []JobTransitionLog
	for rows.Next() {
		var log JobTransitionLog
		err := rows.Scan(
			&log.ID, &log.JobID, &log.FromState, &log.ToState, &log.CreatedAt,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to scan transition log: %w", err)
		}
		transitionLogs = append(transitionLogs, log)
	}

	return &job, transitionLogs, nil
}
