package rapid_diffs

import (
	"fmt"
	"strings"
)

// InlineView represents an inline view component
type InlineView struct {
	Hunks     []*InlineHunk
	FileHash  string
	FilePath  string
}

// NewInlineView creates a new inline view component with the given parameters
func NewInlineView(hunks []*InlineHunk, fileHash string, filePath string) *InlineView {
	return &InlineView{
		Hunks:     hunks,
		FileHash:  fileHash,
		FilePath:  filePath,
	}
}

// Render generates the HTML for the inline view component
func (v *InlineView) Render() string {
	var parts []string

	parts = append(parts, `<table class="rd-table rd-inline-view">`)
	parts = append(parts, "<tbody>")

	for i, hunk := range v.Hunks {
		if i > 0 {
			expandLines := NewExpandLines([]ExpandDirection{ExpandDirectionBoth})
			parts = append(parts, fmt.Sprintf(`<tr class="rd-expand-lines"><td colspan="3">%s</td></tr>`, expandLines.Render()))
		}

		parts = append(parts, hunk.Render())
	}

	parts = append(parts, "</tbody>")
	parts = append(parts, "</table>")

	return strings.Join(parts, "\n")
}
