package models

import (
	"time"
)

// BatchedMigration represents a background migration that processes data in batches
type BatchedMigration struct {
	ID          uint64    `json:"id" gorm:"primaryKey"`
	JobClass    string    `json:"job_class"`
	TableName   string    `json:"table_name"`
	ColumnName  string    `json:"column_name"`
	JobArguments string   `json:"job_arguments"`
	Status      string    `json:"status"`
	Database    string    `json:"database"`
	BatchSize   int64     `json:"batch_size"`
	SubBatchSize int64    `json:"sub_batch_size"`
	Interval    int64     `json:"interval"`
	MinValue    int64     `json:"min_value"`
	MaxValue    int64     `json:"max_value"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BatchedJob represents a single batch job in a batched migration
type BatchedJob struct {
	ID           uint64    `json:"id" gorm:"primaryKey"`
	MigrationID  uint64    `json:"migration_id"`
	BatchSize    int64     `json:"batch_size"`
	SubBatchSize int64     `json:"sub_batch_size"`
	MinValue     int64     `json:"min_value"`
	MaxValue     int64     `json:"max_value"`
	Status       string    `json:"status"`
	Attempts     int64     `json:"attempts"`
	Error        string    `json:"error"`
	StartedAt    time.Time `json:"started_at"`
	FinishedAt   time.Time `json:"finished_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName returns the table name for the BatchedMigration model
func (BatchedMigration) TableName() string {
	return "batched_background_migrations"
}

// TableName returns the table name for the BatchedJob model
func (BatchedJob) TableName() string {
	return "batched_background_migration_jobs"
}
