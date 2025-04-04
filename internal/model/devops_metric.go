package model

import (
	"time"
)

// DevOpsMetric represents a DevOps report metric
type DevOpsMetric struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// Add other metric fields as needed
	Score     float64   `json:"score"`
	Data      string    `json:"data"` // JSON string containing metric details
}

// TableName specifies the table name for DevOpsMetric
func (DevOpsMetric) TableName() string {
	return "dev_ops_metrics"
}
