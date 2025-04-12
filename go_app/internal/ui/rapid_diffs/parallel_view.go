package rapid_diffs

import (
	"fmt"
	"strings"
)

// ParallelView represents a parallel view component
type ParallelView struct {
	Hunks     []*ParallelHunk
	FileHash  string
	FilePath  string
}

// NewParallelView creates a new parallel view component with the given parameters
func NewParallelView(hunks []*ParallelHunk, fileHash string, filePath string) *ParallelView {
	return &ParallelView{
		Hunks:     hunks,
		FileHash:  fileHash,
		FilePath:  filePath,
	}
}

// Render generates the HTML for the parallel view component
func (v *ParallelView) Render() string {
	var parts []string

	parts = append(parts, `<table class="rd-table rd-parallel-view">`)
	parts = append(parts, "<tbody>")

	for i, hunk := range v.Hunks {
		if i > 0 {
			expandLines := NewExpandLines([]ExpandDirection{ExpandDirectionBoth})
			parts = append(parts, fmt.Sprintf(`<tr class="rd-expand-lines"><td colspan="4">%s</td></tr>`, expandLines.Render()))
		}

		parts = append(parts, hunk.Render())
	}

	parts = append(parts, "</tbody>")
	parts = append(parts, "</table>")

	return strings.Join(parts, "\n")
}
