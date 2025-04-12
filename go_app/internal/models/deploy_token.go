package models

// DeployToken represents a group's deploy token
type DeployToken struct {
	ID      int64
	GroupID int64
	Active  bool
	// Add other necessary fields like token value, expiry, etc.
} 