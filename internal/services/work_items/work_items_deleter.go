package work_items

import (
	"context"

	"gitlab.com/gitlab-org/gitlab-redux/internal/models"
)

// WorkItemsDeleter handles deleting work items
type WorkItemsDeleter struct {
	// Add any dependencies here
}

// NewWorkItemsDeleter creates a new work items deleter
func NewWorkItemsDeleter() *WorkItemsDeleter {
	return &WorkItemsDeleter{}
}

// Execute deletes a work item
func (d *WorkItemsDeleter) Execute(ctx context.Context, workItem *models.WorkItem) error {
	// TODO: Implement the actual logic
	// This should:
	// 1. Delete the work item
	// 2. Return any errors
	return nil
} 