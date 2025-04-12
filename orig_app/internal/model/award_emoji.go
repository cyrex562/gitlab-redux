package model

import (
	"time"
)

// AwardEmoji represents an emoji award on an awardable item
type AwardEmoji struct {
	ID           int64     `json:"id"`
	AwardableType string    `json:"awardable_type"`
	AwardableID   int64     `json:"awardable_id"`
	Name         string    `json:"name"`
	UserID       int64     `json:"user_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName returns the database table name for award emojis
func (a *AwardEmoji) TableName() string {
	return "award_emojis"
}

// Awardable is an interface that can receive award emojis
type Awardable interface {
	GetID() int64
	GetType() string
}
