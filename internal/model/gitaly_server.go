package model

import (
	"time"
)

// GitalyServer represents a Gitaly server in the system
type GitalyServer struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Address   string    `json:"address"` // The server's address
	Token     string    `json:"-"`      // The server's authentication token (not exposed in JSON)
	Enabled   bool      `json:"enabled"` // Whether the server is enabled
}

// TableName specifies the table name for GitalyServer
func (GitalyServer) TableName() string {
	return "gitaly_servers"
}
