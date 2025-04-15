package noteable

import (
	"context"
	"fmt"
	"strconv"

	"github.com/cyrex562/gitlab-redux/internal/websocket/channel"
	"github.com/cyrex562/gitlab-redux/internal/websocket/logging"
)

// NewNotesChannel creates a new NotesChannel instance
func NewNotesChannel(channel *channel.Channel, finder NotesFinder) *NotesChannel {
	return &NotesChannel{
		params: NotesFinderParams{},
	}
}

// Subscribe handles the subscription request
func (c *NotesChannel) Subscribe(ctx context.Context, params map[string]interface{}) error {
	// Parse project_id if present
	if projectIDStr, ok := params["project_id"].(string); ok {
		projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid project_id: %v", err)
		}
		c.params.ProjectID = projectID
	}

	// Parse group_id if present
	if groupIDStr, ok := params["group_id"].(string); ok {
		groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid group_id: %v", err)
		}
		c.params.GroupID = groupID
	}

	// Parse noteable_type
	if noteableType, ok := params["noteable_type"].(string); ok {
		c.params.NoteableType = NoteableType(noteableType)
	} else {
		return fmt.Errorf("noteable_type is required")
	}

	// Parse noteable_id
	if noteableIDStr, ok := params["noteable_id"].(string); ok {
		noteableID, err := strconv.ParseInt(noteableIDStr, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid noteable_id: %v", err)
		}
		c.params.NoteableID = noteableID
	} else {
		return fmt.Errorf("noteable_id is required")
	}

	// Find the noteable object
	finder := NewNotesFinder()
	noteable, err := finder.Find(ctx, c.params)
	if err != nil {
		return fmt.Errorf("failed to find noteable: %v", err)
	}

	if noteable == nil {
		return fmt.Errorf("noteable not found")
	}

	c.noteable = noteable

	// Log successful subscription
	logging.GetLogger().Info(logging.LogPayload{
		Params: map[string]interface{}{
			"noteable_type": c.params.NoteableType,
			"noteable_id":   c.params.NoteableID,
			"project_id":    c.params.ProjectID,
			"group_id":      c.params.GroupID,
		},
	}, "Subscribed to notes channel")

	return nil
}

// Unsubscribe handles the unsubscription request
func (c *NotesChannel) Unsubscribe(ctx context.Context) error {
	if c.noteable != nil {
		logging.GetLogger().Info(logging.LogPayload{
			Params: map[string]interface{}{
				"noteable_type": c.noteable.GetType(),
				"noteable_id":   c.noteable.GetID(),
			},
		}, "Unsubscribed from notes channel")
	}
	return nil
}

// StreamFor returns the stream identifier for the noteable
func (c *NotesChannel) StreamFor() string {
	if c.noteable == nil {
		return ""
	}
	return fmt.Sprintf("%s:%d", c.noteable.GetType(), c.noteable.GetID())
}
