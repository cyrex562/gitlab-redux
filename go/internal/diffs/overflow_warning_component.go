package diffs

import (
	"fmt"
	"html/template"
)

// OverflowWarningComponent represents the component for displaying overflow warnings
type OverflowWarningComponent struct {
	BaseComponent
}

// NewOverflowWarningComponent creates a new OverflowWarningComponent instance
func NewOverflowWarningComponent(opts ComponentOptions) *OverflowWarningComponent {
	warning := &OverflowWarning{
		Message:     fmt.Sprintf("This diff is too large to show. Please download it to view the changes."),
		IsVisible:   opts.CurrentLines > opts.MaxLines,
		MaxLines:    opts.MaxLines,
		CurrentLines: opts.CurrentLines,
	}

	return &OverflowWarningComponent{
		BaseComponent: BaseComponent{
			Warning:   warning,
			IsVisible: opts.IsVisible,
		},
	}
}

// Render renders the overflow warning component
func (c *OverflowWarningComponent) Render() (template.HTML, error) {
	if !c.IsVisible || !c.Warning.IsVisible {
		return "", nil
	}

	warning := c.Warning
	html := fmt.Sprintf(`
		<div class="diff-overflow-warning">
			<div class="diff-overflow-warning-content">
				<p>%s</p>
				<p class="diff-overflow-warning-details">
					Current size: %d lines<br>
					Maximum size: %d lines
				</p>
			</div>
		</div>
	`, warning.Message, warning.CurrentLines, warning.MaxLines)

	return template.HTML(html), nil
}

// GetWarning returns the overflow warning
func (c *OverflowWarningComponent) GetWarning() *OverflowWarning {
	return c.Warning
}

// IsOverflowed returns whether the diff has overflowed
func (c *OverflowWarningComponent) IsOverflowed() bool {
	return c.Warning.IsVisible
}
