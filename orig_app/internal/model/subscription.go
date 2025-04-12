package model

import (
	"time"
)

// Subscription represents a user's subscription to a subscribable resource
type Subscription struct {
	ID              int64     `json:"id"`
	SubscribableType string    `json:"subscribable_type"`
	SubscribableID   int64     `json:"subscribable_id"`
	ProjectID        int64     `json:"project_id"`
	UserID           int64     `json:"user_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// TableName returns the database table name for subscriptions
func (s *Subscription) TableName() string {
	return "subscriptions"
}

// Subscribable is an interface that can be subscribed to
type Subscribable interface {
	GetID() int64
	GetType() string
}
