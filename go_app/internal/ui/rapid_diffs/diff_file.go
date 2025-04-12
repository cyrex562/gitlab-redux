package rapid_diffs

import (
	"fmt"
	"strings"
)

// DiffFileComponent represents a diff file component
type DiffFileComponent struct {
	DiffFile     *DiffFile
	ParallelView bool
}

// NewDiffFileComponent creates a new diff file component with the given parameters
func NewDiffFileComponent(diffFile *DiffFile, parallelView bool) *DiffFileComponent {
	return &DiffFileComponent{
		DiffFile:     diffFile,
		ParallelView: parallelView,
	}
}

// GetID returns the file hash as the component ID
func (d *DiffFileComponent) GetID() string {
	return d.DiffFile.FileHash
}

// GetServerData returns the server data for the component
func (d *DiffFileComponent) GetServerData() map[string]interface{} {
	project := d.DiffFile.Repository.Project
	params := fmt.Sprintf("%s/%s", d.DiffFile.ContentSha, d.DiffFile.FilePath)

	return map[string]interface{}{
		"viewer":           d.GetViewerComponent().ViewerName(),
		"diff_lines_path": fmt.Sprintf("/projects/%d/blob/%s/diff_lines", project.ID, params),
	}
}

// GetViewerComponent returns the appropriate viewer component
func (d *DiffFileComponent) GetViewerComponent() ViewerComponentInterface {
	if d.DiffFile.Collapsed || !d.DiffFile.ModifiedFile {
		return NewNoPreviewComponent(d.DiffFile)
	}

	if d.DiffFile.DiffableText {
		if d.ParallelView {
			return NewParallelViewComponent(d.DiffFile)
		}
		return NewInlineViewComponent(d.DiffFile)
	}

	return NewNoPreviewComponent(d.DiffFile)
}

// Render generates the HTML for the diff file component
func (d *DiffFileComponent) Render() string {
	var parts []string

	// Start diff file component
	serverData := ""
	for key, value := range d.GetServerData() {
		serverData += fmt.Sprintf(` data-%s="%v"`, key, value)
	}

	parts = append(parts, fmt.Sprintf(`<diff-file class="rd-diff-file-component" id="%s"%s>`, d.GetID(), serverData))
	parts = append(parts, `<div class="rd-diff-file">`)

	// Header
	header := NewDiffFileHeaderComponent(d.DiffFile)
	parts = append(parts, header.Render())

	// Body
	parts = append(parts, `<div data-file-body="">`)
	parts = append(parts, `<div class="rd-diff-file-body">`)

	viewer := d.GetViewerComponent()
	parts = append(parts, viewer.Render())

	parts = append(parts, "</div>")
	parts = append(parts, "</div>")
	parts = append(parts, "</div>")
	parts = append(parts, "<diff-file-mounted></diff-file-mounted>")
	parts = append(parts, "</diff-file>")

	return strings.Join(parts, "\n")
}
