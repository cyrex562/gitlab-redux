package ui

import (
	"fmt"
	"strings"
)

// LabelPosition represents the available label position options for the toggle
type LabelPosition string

const (
	LabelPositionTop    LabelPosition = "top"
	LabelPositionLeft   LabelPosition = "left"
	LabelPositionHidden LabelPosition = "hidden"
)

// Toggle represents a toggle component
type Toggle struct {
	Classes       string
	Label         string
	LabelPosition LabelPosition
	ID            string
	Name          string
	Help          string
	Data          map[string]interface{}
	IsDisabled    bool
	IsChecked     bool
	IsLoading     bool
	Content       string
}

// NewToggle creates a new toggle with the given parameters
func NewToggle(
	classes string,
	label string,
	labelPosition LabelPosition,
	id string,
	name string,
	help string,
	data map[string]interface{},
	isDisabled bool,
	isChecked bool,
	isLoading bool,
	content string,
) *Toggle {
	if data == nil {
		data = make(map[string]interface{})
	}

	return &Toggle{
		Classes:       classes,
		Label:         label,
		LabelPosition: labelPosition,
		ID:            id,
		Name:          name,
		Help:          help,
		Data:          data,
		IsDisabled:    isDisabled,
		IsChecked:     isChecked,
		IsLoading:     isLoading,
		Content:       content,
	}
}

// GetDataAttributes returns the data attributes for the toggle
func (t *Toggle) GetDataAttributes() string {
	data := map[string]interface{}{
		"name":           t.Name,
		"id":             t.ID,
		"is_checked":     fmt.Sprintf("%v", t.IsChecked),
		"disabled":       fmt.Sprintf("%v", t.IsDisabled),
		"is_loading":     fmt.Sprintf("%v", t.IsLoading),
		"label":          t.Label,
		"help":           t.Help,
		"label_position": t.LabelPosition,
	}

	// Merge with additional data attributes
	for k, v := range t.Data {
		data[k] = v
	}

	return formatAttributes(data)
}

// Render generates the HTML for the toggle
func (t *Toggle) Render() string {
	var parts []string

	// Start toggle span
	parts = append(parts, fmt.Sprintf(`<span class="%s" %s>`, t.Classes, t.GetDataAttributes()))

	// Add help label if content is present
	if t.Content != "" {
		parts = append(parts, fmt.Sprintf(`<div class="gl-help-label">%s</div>`, t.Content))
	}

	// End toggle span
	parts = append(parts, "</span>")

	return strings.Join(parts, "\n")
}
