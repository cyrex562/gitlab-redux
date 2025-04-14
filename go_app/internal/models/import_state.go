package models

// ImportState represents the state of a group import
type ImportState struct {
	ID        int64
	GroupID   int64
	Status    string
	LastError string
}

// IsFinished checks if the import is finished
func (s *ImportState) IsFinished() bool {
	return s.Status == "finished"
}

// IsFailed checks if the import failed
func (s *ImportState) IsFailed() bool {
	return s.Status == "failed"
}

// Group methods for import state

// ImportState returns the group's import state
func (g *Group) ImportState() *ImportState {
	// TODO: Implement import state retrieval logic
	return nil
} 