package models

import (
	"errors"
	"time"
)

// User represents a user in the system
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"uniqueIndex"`
	Email     string    `json:"email" gorm:"uniqueIndex"`
	Password  string    `json:"-"` // Password hash, not exposed in JSON
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// IsAdmin checks if the user has administrator privileges
func (u *User) IsAdmin() bool {
	return u.IsAdmin
}

// Ban bans a user from the system
func (u *User) Ban(bannedBy *User) error {
	if !bannedBy.IsAdmin() {
		return errors.New("only administrators can ban users")
	}
	// TODO: Implement user ban logic
	return nil
}

// Block blocks a user
func (u *User) Block(blockedBy *User) error {
	if !blockedBy.IsAdmin() {
		return errors.New("only administrators can block users")
	}
	// TODO: Implement user block logic
	return nil
}

// Warn sends a warning to a user
func (u *User) Warn(warnedBy *User, reason string) error {
	if !warnedBy.IsAdmin() {
		return errors.New("only administrators can warn users")
	}
	// TODO: Implement user warning logic
	return nil
}

// GetUserByID retrieves a user by their ID
func GetUserByID(id uint) (*User, error) {
	var user User
	// TODO: Implement database query to get user by ID
	if user.ID == 0 {
		return nil, errors.New("user not found")
	}
	return &user, nil
}
