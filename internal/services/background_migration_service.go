package services

import (
	"errors"
	"fmt"

	"github.com/gitlab-org/gitlab-redux/internal/models"
	"gorm.io/gorm"
)

// BackgroundMigrationService handles background migration operations
type BackgroundMigrationService struct {
	db *gorm.DB
}

// NewBackgroundMigrationService creates a new instance of BackgroundMigrationService
func NewBackgroundMigrationService(db *gorm.DB) *BackgroundMigrationService {
	return &BackgroundMigrationService{
		db: db,
	}
}

// GetMigrationsByStatus retrieves migrations filtered by status and database
func (s *BackgroundMigrationService) GetMigrationsByStatus(status, database string, page int) ([]*models.BatchedMigration, error) {
	var migrations []*models.BatchedMigration

	query := s.db.Model(&models.BatchedMigration{})

	// Filter by status
	switch status {
	case "queued":
		query = query.Where("status = ?", "queued")
	case "finalizing":
		query = query.Where("status = ?", "finalizing")
	case "failed":
		query = query.Where("status = ?", "failed")
	case "finished":
		query = query.Where("status = ?", "finished").Order("created_at DESC")
	default:
		return nil, errors.New("invalid status")
	}

	// Filter by database
	if database != "main" {
		query = query.Where("database = ?", database)
	}

	// Apply pagination
	limit := 20
	offset := (page - 1) * limit
	query = query.Limit(limit).Offset(offset)

	if err := query.Find(&migrations).Error; err != nil {
		return nil, err
	}

	return migrations, nil
}

// GetSuccessfulRowsCounts retrieves the count of successful rows for each migration
func (s *BackgroundMigrationService) GetSuccessfulRowsCounts(migrations []*models.BatchedMigration) (map[uint64]int64, error) {
	if len(migrations) == 0 {
		return make(map[uint64]int64), nil
	}

	var ids []uint64
	for _, m := range migrations {
		ids = append(ids, m.ID)
	}

	var results []struct {
		MigrationID uint64
		Count       int64
	}

	if err := s.db.Model(&models.BatchedJob{}).
		Select("migration_id, COUNT(*) as count").
		Where("migration_id IN ? AND status = ?", ids, "succeeded").
		Group("migration_id").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	counts := make(map[uint64]int64)
	for _, r := range results {
		counts[r.MigrationID] = r.Count
	}

	return counts, nil
}

// GetAvailableDatabases retrieves the list of available databases
func (s *BackgroundMigrationService) GetAvailableDatabases() ([]string, error) {
	var databases []string
	if err := s.db.Model(&models.BatchedMigration{}).
		Distinct().
		Pluck("database", &databases).Error; err != nil {
		return nil, err
	}
	return databases, nil
}

// GetMigration retrieves a specific migration by ID
func (s *BackgroundMigrationService) GetMigration(id uint64) (*models.BatchedMigration, error) {
	var migration models.BatchedMigration
	if err := s.db.First(&migration, id).Error; err != nil {
		return nil, err
	}
	return &migration, nil
}

// GetFailedJobs retrieves failed jobs for a migration
func (s *BackgroundMigrationService) GetFailedJobs(migrationID uint64, page int) ([]*models.BatchedJob, error) {
	var jobs []*models.BatchedJob

	limit := 20
	offset := (page - 1) * limit

	if err := s.db.Where("migration_id = ? AND status = ?", migrationID, "failed").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&jobs).Error; err != nil {
		return nil, err
	}

	return jobs, nil
}

// PauseMigration pauses a running migration
func (s *BackgroundMigrationService) PauseMigration(id uint64) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var migration models.BatchedMigration
		if err := tx.First(&migration, id).Error; err != nil {
			return err
		}

		if migration.Status != "running" {
			return fmt.Errorf("cannot pause migration with status %s", migration.Status)
		}

		return tx.Model(&migration).Update("status", "paused").Error
	})
}

// ResumeMigration resumes a paused migration
func (s *BackgroundMigrationService) ResumeMigration(id uint64) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var migration models.BatchedMigration
		if err := tx.First(&migration, id).Error; err != nil {
			return err
		}

		if migration.Status != "paused" {
			return fmt.Errorf("cannot resume migration with status %s", migration.Status)
		}

		return tx.Model(&migration).Update("status", "running").Error
	})
}

// RetryFailedJobs retries failed jobs in a migration
func (s *BackgroundMigrationService) RetryFailedJobs(id uint64) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var migration models.BatchedMigration
		if err := tx.First(&migration, id).Error; err != nil {
			return err
		}

		if migration.Status != "failed" {
			return fmt.Errorf("cannot retry jobs in migration with status %s", migration.Status)
		}

		// Update failed jobs to queued
		if err := tx.Model(&models.BatchedJob{}).
			Where("migration_id = ? AND status = ?", id, "failed").
			Update("status", "queued").Error; err != nil {
			return err
		}

		// Update migration status to queued
		return tx.Model(&migration).Update("status", "queued").Error
	})
}
