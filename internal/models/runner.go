package models

import "time"

// Runner represents a GitLab CI runner
type Runner struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	GroupID     int64     `json:"group_id"`
	AuthorID    int64     `json:"author_id"`
	// Add other fields as needed
}

// RunnerUpdateParams represents parameters for updating a runner
type RunnerUpdateParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	// Add other fields as needed
}

// RegistrationAvailable checks if the runner can be registered
func (r *Runner) RegistrationAvailable() bool {
	// TODO: Implement the actual logic
	return true
} 