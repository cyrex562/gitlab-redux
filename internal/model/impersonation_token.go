package model

import (
	"time"
)

// ImpersonationToken represents a token used for user impersonation
type ImpersonationToken struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	UserID         uint      `json:"user_id"`
	OrganizationID uint      `json:"organization_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Token          string    `json:"token"`
	Scopes         []string  `json:"scopes" gorm:"type:json"`
	ExpiresAt      time.Time `json:"expires_at"`
	LastUsedAt     time.Time `json:"last_used_at"`
	Revoked        bool      `json:"revoked"`
	RevokedAt      time.Time `json:"revoked_at,omitempty"`
}

// TableName specifies the table name for ImpersonationToken
func (ImpersonationToken) TableName() string {
	return "impersonation_tokens"
}
