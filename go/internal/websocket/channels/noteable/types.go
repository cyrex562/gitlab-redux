package noteable

import (
	"context"
)

// NoteableType represents the type of noteable object
type NoteableType string

const (
	NoteableTypeIssue        NoteableType = "issue"
	NoteableTypeMergeRequest NoteableType = "merge_request"
	NoteableTypeCommit       NoteableType = "commit"
	NoteableTypeSnippet      NoteableType = "snippet"
)

// Noteable represents an object that can have notes
type Noteable interface {
	GetID() int64
	GetType() NoteableType
	GetProjectID() int64
	GetGroupID() int64
}

// NotesFinder interface defines methods for finding noteable objects
type NotesFinder interface {
	Find(ctx context.Context, params NotesFinderParams) (Noteable, error)
}

// NotesFinderParams contains parameters for finding a noteable object
type NotesFinderParams struct {
	ProjectID   int64
	GroupID     int64
	NoteableType NoteableType
	NoteableID  int64
}

// NotesChannel represents the WebSocket channel for notes
type NotesChannel struct {
	noteable Noteable
	params   NotesFinderParams
}
