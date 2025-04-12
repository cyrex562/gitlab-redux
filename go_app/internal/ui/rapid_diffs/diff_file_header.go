package rapid_diffs

import (
	"encoding/json"
	"fmt"
	"strings"
)

// DiffFileHeaderComponent represents a diff file header component
type DiffFileHeaderComponent struct {
	DiffFile *DiffFile
}

// NewDiffFileHeaderComponent creates a new diff file header component
func NewDiffFileHeaderComponent(diffFile *DiffFile) *DiffFileHeaderComponent {
	return &DiffFileHeaderComponent{
		DiffFile: diffFile,
	}
}

// GetOptionsMenuItems returns the options menu items
func (d *DiffFileHeaderComponent) GetOptionsMenuItems() string {
	viewTitle := fmt.Sprintf("View file @ %s", d.DiffFile.ContentSha[:8])
	viewHref := fmt.Sprintf("/projects/%d/blob/%s/%s",
		d.DiffFile.Repository.Project.ID,
		d.DiffFile.ContentSha,
		d.DiffFile.NewPath)

	items := []map[string]string{
		{
			"text": viewTitle,
			"href": viewHref,
		},
	}

	jsonBytes, _ := json.Marshal(items)
	return string(jsonBytes)
}

// RenderCopyPathButton renders the copy path button
func (d *DiffFileHeaderComponent) RenderCopyPathButton() string {
	return fmt.Sprintf(`
		<button class="gl-button btn btn-default btn-icon btn-sm"
			data-clipboard-text="%s"
			data-clipboard-gfm="\`%s\`"
			title="Copy file path"
			data-placement="top"
			data-boundary="viewport"
			data-testid="rd-diff-file-copy-clipboard">
			<span class="gl-icon">copy</span>
		</button>
	`, d.DiffFile.FilePath, d.DiffFile.FilePath)
}

// Render generates the HTML for the diff file header component
func (d *DiffFileHeaderComponent) Render() string {
	var parts []string

	parts = append(parts, `<div class="rd-diff-file-header" data-testid="rd-diff-file-header">`)

	// Toggle buttons
	parts = append(parts, `<div class="rd-diff-file-toggle">`)
	parts = append(parts, `
		<button class="gl-button btn btn-default btn-sm rd-diff-file-toggle-button" data-opened data-click="toggleFile" aria-label="Hide file contents">
			<span class="gl-icon">chevron-down</span>
		</button>
		<button class="gl-button btn btn-default btn-sm rd-diff-file-toggle-button" data-closed data-click="toggleFile" aria-label="Show file contents">
			<span class="gl-icon">chevron-right</span>
		</button>
	`)
	parts = append(parts, "</div>")

	// Title section
	parts = append(parts, `<div class="rd-diff-file-title">`)
	if d.DiffFile.Submodule {
		parts = append(parts, fmt.Sprintf(`
			<span data-testid="rd-diff-file-header-submodule">
				<span class="gl-icon gl-file-icon">folder-git</span>
				<strong>%s</strong>
			</span>
			%s
		`, d.DiffFile.FilePath, d.RenderCopyPathButton()))
	} else {
		if d.DiffFile.RenamedFile {
			parts = append(parts, fmt.Sprintf(`
				<strong>%s</strong>
				→
				<strong>%s</strong>
			`, d.DiffFile.OldPath, d.DiffFile.NewPath))
		} else {
			parts = append(parts, fmt.Sprintf(`<strong>%s</strong>`, d.DiffFile.FilePath))
			if d.DiffFile.DeletedFile {
				parts = append(parts, " deleted")
			}
			parts = append(parts, d.RenderCopyPathButton())
		}

		if d.DiffFile.ModeChanged {
			parts = append(parts, fmt.Sprintf(`<small>%s → %s</small>`, d.DiffFile.AMode, d.DiffFile.BMode))
		}

		if d.DiffFile.StoredExternally && d.DiffFile.ExternalStorage == "lfs" {
			parts = append(parts, `<span class="gl-badge badge badge-neutral">LFS</span>`)
		}
	}
	parts = append(parts, "</div>")

	// Info section
	parts = append(parts, `<div class="rd-diff-file-info">`)
	parts = append(parts, `<div class="rd-diff-file-stats">`)
	parts = append(parts, fmt.Sprintf(`
		<span class="rd-lines-added">
			<span>+</span>
			<span data-testid="js-file-addition-line">%d</span>
		</span>
		<span class="rd-lines-removed">
			<span>−</span>
			<span data-testid="js-file-deletion-line">%d</span>
		</span>
	`, d.DiffFile.AddedLines, d.DiffFile.RemovedLines))
	parts = append(parts, "</div>")

	// Options menu
	parts = append(parts, `<div class="rd-diff-file-options-menu">`)
	parts = append(parts, `<div class="js-options-menu">`)
	parts = append(parts, fmt.Sprintf(`<script type="application/json">%s</script>`, d.GetOptionsMenuItems()))
	parts = append(parts, `
		<button class="gl-button btn btn-default btn-sm js-options-button" data-click="toggleOptionsMenu" aria-label="Options">
			<span class="gl-icon">ellipsis_v</span>
		</button>
	`)
	parts = append(parts, "</div>")
	parts = append(parts, "</div>")
	parts = append(parts, "</div>")

	parts = append(parts, "</div>")

	return strings.Join(parts, "\n")
}
