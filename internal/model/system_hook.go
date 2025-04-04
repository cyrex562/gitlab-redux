package model

import (
	"time"
)

// SystemHook represents a system-wide web hook
type SystemHook struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	URL       string    `json:"url"`
	Token     string    `json:"-"` // Not exposed in JSON
	Enabled   bool      `json:"enabled"`
	Events    []string  `json:"events" gorm:"type:json"`
}

// TableName specifies the table name for SystemHook
func (SystemHook) TableName() string {
	return "system_hooks"
}
