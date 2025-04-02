package ui

import (
	"fmt"
	"strings"
)

// RadioTag represents a radio tag component
type RadioTag struct {
	Name         string
	Value        string
	Checked      bool
	Label        string
	HelpText     string
	LabelOptions map[string]interface{}
	RadioOptions map[string]interface{}
}

// NewRadioTag creates a new radio tag with the given parameters
func NewRadioTag(name string, value string, checked bool, label string, helpText string, labelOptions map[string]interface{}, radioOptions map[string]interface{}) *RadioTag {
	if labelOptions == nil {
		labelOptions = make(map[string]interface{})
	}
	if radioOptions == nil {
		radioOptions = make(map[string]interface{})
	}

	// Set for attribute in label options if not already set
	if _, ok := labelOptions["for"]; !ok {
		labelOptions["for"] = labelFor(name, value)
	}

	return &RadioTag{
		Name:         name,
		Value:        value,
		Checked:      checked,
		Label:        label,
		HelpText:     helpText,
		LabelOptions: labelOptions,
		RadioOptions: radioOptions,
	}
}

// GetLabelContent returns the label content
func (r *RadioTag) GetLabelContent() string {
	return r.Label
}

// GetHelpTextContent returns the help text content
func (r *RadioTag) GetHelpTextContent() string {
	return r.HelpText
}

// RenderLabelTagWithHelpText generates the HTML for the label and help text
func (r *RadioTag) RenderLabelTagWithHelpText() string {
	var parts []string

	// Add label class to label options
	if _, ok := r.LabelOptions["class"]; !ok {
		r.LabelOptions["class"] = "custom-control-label"
	} else {
		class := r.LabelOptions["class"].(string)
		if !strings.Contains(class, "custom-control-label") {
			r.LabelOptions["class"] = fmt.Sprintf("custom-control-label %s", class)
		}
	}

	// Add label
	parts = append(parts, fmt.Sprintf(`<label %s>%s</label>`, formatAttributes(r.LabelOptions), r.GetLabelContent()))

	// Add help text if present
	if r.GetHelpTextContent() != "" {
		parts = append(parts, fmt.Sprintf(`<div class="gl-form-text">%s</div>`, r.GetHelpTextContent()))
	}

	return strings.Join(parts, "\n")
}

// Render generates the HTML for the radio tag
func (r *RadioTag) Render() string {
	var parts []string

	// Start content wrapper div
	parts = append(parts, `<div class="gl-form-radio custom-control custom-radio">`)

	// Add radio input
	radioAttrs := make(map[string]interface{})
	for k, v := range r.RadioOptions {
		radioAttrs[k] = v
	}
	radioAttrs["type"] = "radio"
	radioAttrs["name"] = r.Name
	radioAttrs["id"] = labelFor(r.Name, r.Value)
	radioAttrs["value"] = r.Value
	if r.Checked {
		radioAttrs["checked"] = "checked"
	}
	parts = append(parts, fmt.Sprintf(`<input %s>`, formatAttributes(radioAttrs)))

	// Add label and help text
	parts = append(parts, r.RenderLabelTagWithHelpText())

	// Close content wrapper div
	parts = append(parts, "</div>")

	return strings.Join(parts, "\n")
}

// labelFor generates the ID for a label
func labelFor(name string, value string) string {
	return fmt.Sprintf("%s_%s", sanitizeToID(name), value)
}

// sanitizeToID sanitizes a string to be used as an ID
func sanitizeToID(name string) string {
	// Replace any non-alphanumeric characters with underscores
	return strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return '_'
	}, name)
}
