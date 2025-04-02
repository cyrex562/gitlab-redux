package rapid_diffs

import (
	"fmt"
	"strings"
)

// InlineHunk represents an inline hunk component
type InlineHunk struct {
	Lines     []*Line
	FileHash  string
	FilePath  string
	OldStart  int
	NewStart  int
	OldLines  int
	NewLines  int
}

// NewInlineHunk creates a new inline hunk component with the given parameters
func NewInlineHunk(lines []*Line, fileHash string, filePath string, oldStart int, newStart int, oldLines int, newLines int) *InlineHunk {
	return &InlineHunk{
		Lines:     lines,
		FileHash:  fileHash,
		FilePath:  filePath,
		OldStart:  oldStart,
		NewStart:  newStart,
		OldLines:  oldLines,
		NewLines:  newLines,
	}
}

// GetLineNumber creates a line number component for the given line and position
func (h *InlineHunk) GetLineNumber(line *Line, position LinePosition, oldPos int, newPos int) *LineNumber {
	id := fmt.Sprintf("%s_%d", h.FileHash, position)
	legacyID := fmt.Sprintf("%s_%s_%d", h.FileHash, h.FilePath, position)

	return NewLineNumber(
		line,
		position,
		h.FileHash,
		h.FilePath,
		LineBorderBoth,
		oldPos,
		newPos,
		id,
		legacyID,
	)
}

// GetLineContent creates a line content component for the given line and position
func (h *InlineHunk) GetLineContent(line *Line, position LinePosition) *LineContent {
	return NewLineContent(line, position)
}

// Render generates the HTML for the inline hunk component
func (h *InlineHunk) Render() string {
	var parts []string

	oldPos := h.OldStart
	newPos := h.NewStart

	for _, line := range h.Lines {
		parts = append(parts, "<tr>")

		// Old line number
		lineNumber := h.GetLineNumber(line, LinePositionOld, oldPos, newPos)
		parts = append(parts, lineNumber.Render())

		// New line number
		lineNumber = h.GetLineNumber(line, LinePositionNew, oldPos, newPos)
		parts = append(parts, lineNumber.Render())

		// Line content
		lineContent := h.GetLineContent(line, LinePositionNew)
		parts = append(parts, lineContent.Render())

		parts = append(parts, "</tr>")

		if !line.Added {
			oldPos++
		}
		if !line.Removed {
			newPos++
		}
	}

	return strings.Join(parts, "\n")
}
