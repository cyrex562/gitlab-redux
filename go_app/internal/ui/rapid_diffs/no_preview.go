package rapid_diffs

import (
	"fmt"
	"path"
	"strings"
)

// NoPreviewComponent represents a no preview component
type NoPreviewComponent struct {
	*ViewerComponent
}

// NewNoPreviewComponent creates a new no preview component with the given diff file
func NewNoPreviewComponent(diffFile *DiffFile) *NoPreviewComponent {
	return &NoPreviewComponent{
		ViewerComponent: NewViewerComponent(diffFile),
	}
}

// ViewerName returns the name of the viewer
func (n *NoPreviewComponent) ViewerName() string {
	return "no_preview"
}

// GetChangeDescription returns the description of the change
func (n *NoPreviewComponent) GetChangeDescription() string {
	if n.DiffFile.NewFile {
		return "File added."
	} else if n.DiffFile.DeletedFile {
		return "File deleted."
	} else if n.DiffFile.ContentChanged && n.DiffFile.RenamedFile {
		return "File changed and moved."
	} else if n.DiffFile.ContentChanged {
		return "File changed."
	} else if n.DiffFile.RenamedFile {
		return "File moved."
	}
	return ""
}

// GetModeChangedDescription returns the description of the mode change
func (n *NoPreviewComponent) GetModeChangedDescription() string {
	if !n.DiffFile.ModeChanged || n.DiffFile.NewFile || n.DiffFile.DeletedFile {
		return ""
	}
	return fmt.Sprintf("File mode changed from %s to %s.", n.DiffFile.AMode, n.DiffFile.BMode)
}

// GetNoPreviewReason returns the reason why there is no preview
func (n *NoPreviewComponent) GetNoPreviewReason() string {
	if n.DiffFile.TooLarge {
		return "File size exceeds preview limit."
	} else if n.DiffFile.Collapsed {
		return "Preview size limit exceeded, changes collapsed."
	} else if !n.DiffFile.Diffable {
		return "Preview suppressed by a .gitattributes entry or the file's encoding is unsupported."
	} else if n.DiffFile.NewFile || n.DiffFile.ContentChanged {
		return "No diff preview for this file type."
	}
	return ""
}

// IsExpandable returns whether the file is expandable
func (n *NoPreviewComponent) IsExpandable() bool {
	return n.DiffFile.DiffableText
}

// IsImportant returns whether the no preview message is important
func (n *NoPreviewComponent) IsImportant() bool {
	return n.DiffFile.Collapsed || n.DiffFile.TooLarge
}

// GetOldBlobPath returns the path to the old blob
func (n *NoPreviewComponent) GetOldBlobPath() string {
	return path.Join(n.DiffFile.Repository.Project.Path, "blob", n.DiffFile.OldContentSha, n.DiffFile.OldPath)
}

// GetNewBlobPath returns the path to the new blob
func (n *NoPreviewComponent) GetNewBlobPath() string {
	return path.Join(n.DiffFile.Repository.Project.Path, "blob", n.DiffFile.ContentSha, n.DiffFile.FilePath)
}

// GetBlobPath returns the appropriate blob path based on file state
func (n *NoPreviewComponent) GetBlobPath() string {
	if n.DiffFile.DeletedFile {
		return n.GetOldBlobPath()
	}
	return n.GetNewBlobPath()
}

// RenderActionButton renders a button with the given options and content
func (n *NoPreviewComponent) RenderActionButton(href string, dataClick string, content string) string {
	var attrs []string
	attrs = append(attrs, `class="gl-button btn btn-default btn-md"`)

	if href != "" {
		attrs = append(attrs, fmt.Sprintf(`href="%s"`, href))
	}
	if dataClick != "" {
		attrs = append(attrs, fmt.Sprintf(`data-click="%s"`, dataClick))
	}

	return fmt.Sprintf(`<div class="rd-no-preview-action"><button %s>%s</button></div>`, strings.Join(attrs, " "), content)
}

// Render generates the HTML for the no preview component
func (n *NoPreviewComponent) Render() string {
	var parts []string

	importantClass := ""
	if n.IsImportant() {
		importantClass = " rd-no-preview-important"
	}

	parts = append(parts, fmt.Sprintf(`<div class="rd-no-preview%s">`, importantClass))
	parts = append(parts, `<div class="rd-no-preview-body">`)

	// Change description and mode changed description
	parts = append(parts, `<p class="rd-no-preview-paragraph">`)
	if desc := n.GetChangeDescription(); desc != "" {
		parts = append(parts, desc)
	}
	if modeDesc := n.GetModeChangedDescription(); modeDesc != "" {
		parts = append(parts, modeDesc)
	}
	parts = append(parts, "</p>")

	// No preview reason
	if reason := n.GetNoPreviewReason(); reason != "" {
		parts = append(parts, fmt.Sprintf(`<p class="rd-no-preview-paragraph">%s</p>`, reason))
	}
	parts = append(parts, "</div>")

	// Actions
	parts = append(parts, `<div class="rd-no-preview-actions">`)
	if n.DiffFile.Collapsed && n.IsExpandable() {
		parts = append(parts, n.RenderActionButton("", "showChanges", "Show changes"))
	} else if n.IsExpandable() {
		parts = append(parts, n.RenderActionButton("", "showFileContents", "Show file contents"))
	} else if n.DiffFile.ContentChanged {
		parts = append(parts, n.RenderActionButton(n.GetOldBlobPath(), "", "View original file"))
		parts = append(parts, n.RenderActionButton(n.GetNewBlobPath(), "", "View changed file"))
	} else if n.DiffFile.DeletedFile || n.DiffFile.NewFile || n.DiffFile.RenamedFile || n.DiffFile.ModeChanged {
		parts = append(parts, n.RenderActionButton(n.GetBlobPath(), "", "View file"))
	}
	parts = append(parts, "</div>")

	parts = append(parts, "</div>")

	return strings.Join(parts, "\n")
}
