package rapid_diffs

import (
	"fmt"
)

// LineBorder represents the border style for a line number
type LineBorder string

const (
	LineBorderRight LineBorder = "right"
	LineBorderBoth  LineBorder = "both"
)

// LineNumber represents a line number component
type LineNumber struct {
	Line       *Line
	Position   LinePosition
	FileHash   string
	FilePath   string
	Border     LineBorder
	OldPos     int
	NewPos     int
	ID         string
	LegacyID   string
}

// NewLineNumber creates a new line number component with the given parameters
func NewLineNumber(line *Line, position LinePosition, fileHash string, filePath string, border LineBorder, oldPos int, newPos int, id string, legacyID string) *LineNumber {
	return &LineNumber{
		Line:       line,
		Position:   position,
		FileHash:   fileHash,
		FilePath:   filePath,
		Border:     border,
		OldPos:     oldPos,
		NewPos:     newPos,
		ID:         id,
		LegacyID:   legacyID,
	}
}

// GetLineNumber returns the line number based on the position
func (l *LineNumber) GetLineNumber() int {
	if l.Position == LinePositionOld {
		return l.OldPos
	}
	return l.NewPos
}

// GetChangeType returns the change type of the line
func (l *LineNumber) GetChangeType() string {
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

// GetBorderClass returns the CSS class for the border
func (l *LineNumber) GetBorderClass() string {
	switch l.Border {
	case LineBorderRight:
		return "rd-line-number-border-right"
	case LineBorderBoth:
		return "rd-line-number-border-both"
	default:
		return ""
	}
}

// IsVisible returns whether the line number should be visible
func (l *LineNumber) IsVisible() bool {
	if l.Line == nil {
		return false
	}

	switch l.Position {
	case LinePositionOld:
		return !l.Line.Added
	case LinePositionNew:
		return !l.Line.Removed
	default:
		return false
	}
}

// Render generates the HTML for the line number component
func (l *LineNumber) Render() string {
	if !l.IsVisible() {
		return fmt.Sprintf(`<td class="rd-line-number %s" data-change="%s"></td>`, l.GetBorderClass(), l.GetChangeType())
	}

	return fmt.Sprintf(`<td class="rd-line-number %s" id="%s" data-legacy-id="%s" data-change="%s" data-position="%s">
		<a href="#%s" class="rd-line-link" data-line-number="%d" aria-label="Line number %d"></a>
	</td>`,
		l.GetBorderClass(), l.ID, l.LegacyID, l.GetChangeType(), l.Position, l.ID, l.GetLineNumber(), l.GetLineNumber())
}
