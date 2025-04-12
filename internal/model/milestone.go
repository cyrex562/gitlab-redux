package model

import (
	"time"
)

// MilestoneState represents the state of a milestone
type MilestoneState string

const (
	// MilestoneStateActive represents an active milestone
	MilestoneStateActive MilestoneState = "active"
	// MilestoneStateClosed represents a closed milestone
	MilestoneStateClosed MilestoneState = "closed"
)

// Milestone represents a milestone in the system
type Milestone struct {
	ID          int64         `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	State       MilestoneState `json:"state"`
	DueDate     *time.Time    `json:"due_date"`
	StartDate   *time.Time    `json:"start_date"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	ProjectID   int64         `json:"project_id"`
	GroupID     int64         `json:"group_id"`
}

// MilestoneStateCount represents the count of milestones by state
type MilestoneStateCount struct {
	Active int `json:"active"`
	Closed int `json:"closed"`
}

// MilestoneJSON represents a milestone for JSON serialization
type MilestoneJSON struct {
	ID      int64      `json:"id"`
	Title   string     `json:"title"`
	DueDate *time.Time `json:"due_date"`
	Name    string     `json:"name"`
}
