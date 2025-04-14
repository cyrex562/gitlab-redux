package models

// WorkItem represents a work item in the system
type WorkItem struct {
	ID          int64  `json:"id"`
	IID         string `json:"iid"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
	// Add other fields as needed
} 