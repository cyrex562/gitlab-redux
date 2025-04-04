package model

import "time"

// User represents a GitLab user
type User struct {
	ID                    int64     `json:"id"`
	Username             string    `json:"username"`
	Email                string    `json:"email"`
	Name                 string    `json:"name"`
	State                string    `json:"state"`
	Admin                bool      `json:"admin"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	LastSignInAt         time.Time `json:"last_sign_in_at"`
	LastActivityOn       time.Time `json:"last_activity_on"`
	TwoFactorEnabled     bool      `json:"two_factor_enabled"`
	TrustedWithSpamCheck bool      `json:"trusted_with_spam_check"`
	// Add other necessary fields
}

// UserParams represents the parameters for creating or updating a user
type UserParams struct {
	Username         string `json:"username" binding:"required"`
	Email            string `json:"email" binding:"required,email"`
	Name             string `json:"name" binding:"required"`
	Password         string `json:"password"`
	ResetPassword    bool   `json:"reset_password"`
	SkipConfirmation bool   `json:"skip_confirmation"`
	Admin            bool   `json:"admin"`
	// Add other necessary fields
}

// IsAdmin checks if the user has admin privileges
func (u *User) IsAdmin() bool {
	return u.Admin
}

// IsActive checks if the user is active
func (u *User) IsActive() bool {
	return u.State == "active"
}

// IsBlocked checks if the user is blocked
func (u *User) IsBlocked() bool {
	return u.State == "blocked"
}

// IsBanned checks if the user is banned
func (u *User) IsBanned() bool {
	return u.State == "banned"
}

// IsDeactivated checks if the user is deactivated
func (u *User) IsDeactivated() bool {
	return u.State == "deactivated"
}

// IsPendingApproval checks if the user is pending approval
func (u *User) IsPendingApproval() bool {
	return u.State == "pending_approval"
}

// IsLocked checks if the user is locked
func (u *User) IsLocked() bool {
	return u.State == "locked"
}

// IsTrustedWithSpamCheck checks if the user is trusted with spam checks
func (u *User) IsTrustedWithSpamCheck() bool {
	return u.TrustedWithSpamCheck
}
