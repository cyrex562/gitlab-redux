package rapid_diffs

import (
	"fmt"
	"strings"
)

// ParallelHunk represents a parallel hunk component
type ParallelHunk struct {
	Lines     []*Line
	FileHash  string
	FilePath  string
	OldStart  int
	NewStart  int
	OldLines  int
	NewLines  int
}

// NewParallelHunk creates a new parallel hunk component with the given parameters
func NewParallelHunk(lines []*Line, fileHash string, filePath string, oldStart int, newStart int, oldLines int, newLines int) *ParallelHunk {
	return &ParallelHunk{
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
func (h *ParallelHunk) GetLineNumber(line *Line, position LinePosition, oldPos int, newPos int, border LineBorder) *LineNumber {
	id := fmt.Sprintf("%s_%d", h.FileHash, position)
	legacyID := fmt.Sprintf("%s_%s_%d", h.FileHash, h.FilePath, position)

	return NewLineNumber(
		line,
		position,
		h.FileHash,
		h.FilePath,
		border,
		oldPos,
		newPos,
		id,
		legacyID,
	)
}

// GetLineContent creates a line content component for the given line and position
func (h *ParallelHunk) GetLineContent(line *Line, position LinePosition) *LineContent {
	return NewLineContent(line, position)
}

// Render generates the HTML for the parallel hunk component
func (h *ParallelHunk) Render() string {
	var parts []string

	oldPos := h.OldStart
	newPos := h.NewStart

	for _, line := range h.Lines {
		parts = append(parts, "<tr>")

		// Old line number and content
		if !line.Added {
			lineNumber := h.GetLineNumber(line, LinePositionOld, oldPos, newPos, LineBorderRight)
			parts = append(parts, lineNumber.Render())

			lineContent := h.GetLineContent(line, LinePositionOld)
			parts = append(parts, lineContent.Render())

			oldPos++
		} else {
			parts = append(parts, `<td class="rd-line-number rd-line-number-border-right"></td>`)
			parts = append(parts, `<td class="rd-line-content"></td>`)
		}

		// New line number and content
		if !line.Removed {
			lineNumber := h.GetLineNumber(line, LinePositionNew, oldPos, newPos, LineBorderRight)
			parts = append(parts, lineNumber.Render())

			lineContent := h.GetLineContent(line, LinePositionNew)
			parts = append(parts, lineContent.Render())

			newPos++
		} else {
			parts = append(parts, `<td class="rd-line-number rd-line-number-border-right"></td>`)
			parts = append(parts, `<td class="rd-line-content"></td>`)
		}

		parts = append(parts, "</tr>")
	}

	return strings.Join(parts, "\n")
}
