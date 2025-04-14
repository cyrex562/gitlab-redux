package models

import (
	"time"
)

// ImportFailure represents an import failure
type ImportFailure struct {
	ID          string
	ProjectID   string
	Error       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
} 