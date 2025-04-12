package rapid_diffs

import (
	"fmt"
)

// LinePosition represents the position of a line
type LinePosition string

const (
	LinePositionOld LinePosition = "old"
	LinePositionNew LinePosition = "new"
)

// Line represents a line in a diff
type Line struct {
	TextContent string
	Added       bool
	Removed     bool
}

// LineContent represents a line content component
type LineContent struct {
	Line     *Line
	Position LinePosition
}

// NewLineContent creates a new line content component with the given line and position
func NewLineContent(line *Line, position LinePosition) *LineContent {
	return &LineContent{
		Line:     line,
		Position: position,
	}
}

// GetChangeType returns the change type of the line
func (l *LineContent) GetChangeType() string {
	if l.Line == nil {
		return ""
	}

	if l.Line.Added {
		return "added"
	}

	if l.Line.Removed {
		return "removed"
	}

	return ""
}

// Render generates the HTML for the line content component
func (l *LineContent) Render() string {
	if l.Line == nil {
		return ""
	}

	return fmt.Sprintf(`<td class="rd-line-content" data-change="%s" data-position="%s" tabindex="-1">%s</td>`,
		l.GetChangeType(), l.Position, l.Line.TextContent)
}
