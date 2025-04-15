package models

// SourceUser represents a source user for import operations
type SourceUser struct {
	ID               string `json:"id"`
	Username         string `json:"username"`
	Email            string `json:"email"`
	ReassignToUserID string `json:"reassign_to_user_id"`
	ReassignedByUser string `json:"reassigned_by_user"`
	ReassignmentToken string `json:"reassignment_token"`
	AwaitingApproval bool   `json:"awaiting_approval"`
} 