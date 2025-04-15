package models

import (
	"time"
)

// DeviceGrant represents an OAuth device authorization grant
type DeviceGrant struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserCode  string    `json:"user_code" gorm:"uniqueIndex"`
	DeviceCode string   `json:"device_code" gorm:"uniqueIndex"`
	ClientID  string    `json:"client_id"`
	Scopes    string    `json:"scopes"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
} 