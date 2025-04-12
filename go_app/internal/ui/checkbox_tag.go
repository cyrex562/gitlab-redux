package ui

import (
	"fmt"
	"strings"
)

// CheckboxTag represents a checkbox tag component
type CheckboxTag struct {
	Name            string
	LabelOptions    map[string]interface{}
	CheckboxOptions map[string]interface{}
	Value           string
	Checked         bool
	Label           string
	HelpText        string
}

// NewCheckboxTag creates a new checkbox tag with the given parameters
func NewCheckboxTag(name string, labelOptions map[string]interface{}, checkboxOptions map[string]interface{}, value string, checked bool, label string, helpText string) *CheckboxTag {
	if labelOptions == nil {
		labelOptions = make(map[string]interface{})
	}
	if checkboxOptions == nil {
		checkboxOptions = make(map[string]interface{})
	}

	// Set default value if not provided
	if value == "" {
		value = "1"
	}

	return &CheckboxTag{
		Name:            name,
		LabelOptions:    labelOptions,
		CheckboxOptions: checkboxOptions,
		Value:           value,
		Checked:         checked,
		Label:           label,
		HelpText:        helpText,
	}
}

// GetLabelContent returns the label content
func (c *CheckboxTag) GetLabelContent() string {
	return c.Label
}

// GetHelpTextContent returns the help text content
func (c *CheckboxTag) GetHelpTextContent() string {
	return c.HelpText
}

// RenderLabelTagWithHelpText generates the HTML for the label and help text
func (c *CheckboxTag) RenderLabelTagWithHelpText() string {
	var parts []string

	// Add label class to label options
	if _, ok := c.LabelOptions["class"]; !ok {
		c.LabelOptions["class"] = "custom-control-label"
	} else {
		class := c.LabelOptions["class"].(string)
		if !strings.Contains(class, "custom-control-label") {
			c.LabelOptions["class"] = fmt.Sprintf("custom-control-label %s", class)
		}
	}

	// Add for attribute to label options
	c.LabelOptions["for"] = c.Name

	// Add label
	parts = append(parts, fmt.Sprintf(`<label %s>%s</label>`, formatAttributes(c.LabelOptions), c.GetLabelContent()))

	// Add help text if present
	if c.GetHelpTextContent() != "" {
		parts = append(parts, fmt.Sprintf(`<div class="gl-form-text">%s</div>`, c.GetHelpTextContent()))
	}

	return strings.Join(parts, "\n")
}

// Render generates the HTML for the checkbox tag
func (c *CheckboxTag) Render() string {
	var parts []string

	// Start content wrapper div
	parts = append(parts, `<div class="gl-form-checkbox custom-control custom-checkbox">`)

	// Add checkbox input
	checkboxAttrs := make(map[string]interface{})
	for k, v := range c.CheckboxOptions {
		checkboxAttrs[k] = v
	}
	checkboxAttrs["type"] = "checkbox"
	checkboxAttrs["name"] = c.Name
	checkboxAttrs["id"] = c.Name
	checkboxAttrs["value"] = c.Value
	if c.Checked {
		checkboxAttrs["checked"] = "checked"
	}
	parts = append(parts, fmt.Sprintf(`<input %s>`, formatAttributes(checkboxAttrs)))

	// Add label and help text
	parts = append(parts, c.RenderLabelTagWithHelpText())

	// Close content wrapper div
	parts = append(parts, "</div>")

	return strings.Join(parts, "\n")
}
