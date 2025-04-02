package models

import (
	"errors"
	"os"
	"path/filepath"
	"time"
)

// Topic represents a project topic
type Topic struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Avatar      string    `json:"avatar"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GetTopicByID retrieves a topic by its ID
func GetTopicByID(id uint) (*Topic, error) {
	var topic Topic
	// TODO: Implement database query to get topic by ID
	if topic.ID == 0 {
		return nil, errors.New("topic not found")
	}
	return &topic, nil
}

// Save saves the topic to the database
func (t *Topic) Save() error {
	// TODO: Implement database save
	return nil
}

// RemoveAvatar removes the avatar file and updates the topic
func (t *Topic) RemoveAvatar() error {
	if t.Avatar == "" {
		return nil
	}

	// Get the avatar file path
	avatarPath := filepath.Join("public", "uploads", "topic", "avatar", t.Avatar)

	// Remove the file if it exists
	if err := os.Remove(avatarPath); err != nil && !os.IsNotExist(err) {
		return err
	}

	// Clear the avatar field
	t.Avatar = ""
	return nil
}
