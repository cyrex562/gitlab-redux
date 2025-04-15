package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	Username        string    `json:"username" gorm:"uniqueIndex"`
	Email           string    `json:"email" gorm:"uniqueIndex"`
	Password        string    `json:"-" gorm:"not null"` // "-" means don't include in JSON
	Admin           bool      `json:"admin" gorm:"default:false"`
	Confirmed       bool      `json:"confirmed" gorm:"default:false"`
	Name            string    `json:"name"`
	State           string    `json:"state" gorm:"default:'active'"`
	LastSignInAt    time.Time `json:"last_sign_in_at"`
	CurrentSignInAt time.Time `json:"current_sign_in_at"`
	LastSignInIP    string    `json:"last_sign_in_ip"`
	CurrentSignInIP string    `json:"current_sign_in_ip"`
	SignInCount     int       `json:"sign_in_count" gorm:"default:0"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       time.Time `json:"deleted_at" gorm:"index"`
} 