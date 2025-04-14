package models

import "time"

// User represents a GitLab user
type User struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	// Add other fields as needed
}

// CanReadUsageQuotas checks if the user can read usage quotas for the group
func (u *User) CanReadUsageQuotas(group *Group) bool {
	// TODO: Implement the actual authorization logic
	return true
}

// CanAdminGroup checks if the user can admin the group
func (u *User) CanAdminGroup(group *Group) bool {
	// TODO: Implement the actual authorization logic
	return true
}

// CanAdminCicdVariables checks if the user can admin CI/CD variables for the group
func (u *User) CanAdminCicdVariables(group *Group) bool {
	// TODO: Implement the actual authorization logic
	return true
} 