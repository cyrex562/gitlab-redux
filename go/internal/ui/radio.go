package ui

import (
	"fmt"
	"strings"
)

// Radio represents a radio component
type Radio struct {
	Form         interface{} // This would be an ActionView::Helpers::FormBuilder in Ruby
	Method       string
	Label        string
	HelpText     string
	LabelOptions map[string]interface{}
	RadioOptions map[string]interface{}
	Value        string
}

// NewRadio creates a new radio with the given parameters
func NewRadio(form interface{}, method string, label string, helpText string, labelOptions map[string]interface{}, radioOptions map[string]interface{}, value string) *Radio {
	if labelOptions == nil {
		labelOptions = make(map[string]interface{})
	}
	if radioOptions == nil {
		radioOptions = make(map[string]interface{})
	}

	return &Radio{
		Form:         form,
		Method:       method,
		Label:        label,
		HelpText:     helpText,
		LabelOptions: labelOptions,
		RadioOptions: radioOptions,
		Value:        value,
	}
}

// GetLabelContent returns the label content
func (r *Radio) GetLabelContent() string {
	return r.Label
}

// GetHelpTextContent returns the help text content
func (r *Radio) GetHelpTextContent() string {
	return r.HelpText
}

// RenderLabelWithHelpText generates the HTML for the label and help text
func (r *Radio) RenderLabelWithHelpText() string {
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

	// Add for attribute to label options
	r.LabelOptions["for"] = fmt.Sprintf("%s_%s", r.Method, r.Value)

	// Add label
	parts = append(parts, fmt.Sprintf(`<label %s>%s</label>`, formatAttributes(r.LabelOptions), r.GetLabelContent()))

	// Add help text if present
	if r.GetHelpTextContent() != "" {
		parts = append(parts, fmt.Sprintf(`<div class="gl-form-text">%s</div>`, r.GetHelpTextContent()))
	}

	return strings.Join(parts, "\n")
}

// Render generates the HTML for the radio
func (r *Radio) Render() string {
	var parts []string

	// Start content wrapper div
	parts = append(parts, `<div class="gl-form-radio custom-control custom-radio">`)

	// Add radio input
	// In a real implementation, this would use the form builder to generate the radio button
	// For now, we'll just generate a basic radio input
	radioAttrs := make(map[string]interface{})
	for k, v := range r.RadioOptions {
		radioAttrs[k] = v
	}
	radioAttrs["type"] = "radio"
	radioAttrs["name"] = r.Method
	radioAttrs["id"] = fmt.Sprintf("%s_%s", r.Method, r.Value)
	radioAttrs["value"] = r.Value
	parts = append(parts, fmt.Sprintf(`<input %s>`, formatAttributes(radioAttrs)))

	// Add label and help text
	parts = append(parts, r.RenderLabelWithHelpText())

	// Close content wrapper div
	parts = append(parts, "</div>")

	return strings.Join(parts, "\n")
}
