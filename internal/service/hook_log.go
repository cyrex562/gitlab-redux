package service

import (
	"context"
	"time"

	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
	"gorm.io/gorm"
)

// HookLogService handles business logic for web hook logs
type HookLogService struct {
	db *gorm.DB
}

// NewHookLogService creates a new HookLogService instance
func NewHookLogService(db *gorm.DB) *HookLogService {
	return &HookLogService{
		db: db,
	}
}

// GetSystemHook retrieves a system hook by ID
func (s *HookLogService) GetSystemHook(ctx context.Context, hookID string) (*model.SystemHook, error) {
	var hook model.SystemHook
	err := s.db.WithContext(ctx).First(&hook, hookID).Error
	if err != nil {
		return nil, err
	}
	return &hook, nil
}

// GetHookLog retrieves a specific web hook log
func (s *HookLogService) GetHookLog(ctx context.Context, hookID, logID string) (*model.HookLog, error) {
	var log model.HookLog
	err := s.db.WithContext(ctx).
		Where("hook_id = ? AND id = ?", hookID, logID).
		First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// RetryHook retries a failed web hook
func (s *HookLogService) RetryHook(ctx context.Context, hookID, logID string) error {
	// Start a transaction
	tx := s.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get the original hook log
	var log model.HookLog
	err := tx.Where("hook_id = ? AND id = ?", hookID, logID).First(&log).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Create a new hook log entry for the retry
	newLog := model.HookLog{
		HookID:      log.HookID,
		TriggeredAt: time.Now(),
		Status:      "pending",
		RequestData: log.RequestData,
		ResponseData: log.ResponseData,
	}

	err = tx.Create(&newLog).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// TODO: Trigger the web hook asynchronously
	// This would typically involve:
	// 1. Creating a background job
	// 2. Sending the web hook request
	// 3. Updating the log status based on the result

	// Commit the transaction
	return tx.Commit().Error
}

// HookLog represents a web hook execution log
type HookLog struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	HookID        uint      `json:"hook_id"`
	TriggeredAt   time.Time `json:"triggered_at"`
	Status        string    `json:"status"`
	RequestData   string    `json:"request_data"`
	ResponseData  string    `json:"response_data"`
	ResponseCode  int       `json:"response_code"`
	ExecutionTime int64     `json:"execution_time"`
}

// TableName specifies the table name for HookLog
func (HookLog) TableName() string {
	return "hook_logs"
}
