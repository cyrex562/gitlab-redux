package ui

import (
	"fmt"
	"strings"
)

// Checkbox represents a checkbox component
type Checkbox struct {
	Form                 interface{} // This would be an ActionView::Helpers::FormBuilder in Ruby
	Method              string
	Label               string
	HelpText            string
	LabelOptions        map[string]interface{}
	CheckboxOptions     map[string]interface{}
	ContentWrapperOptions map[string]interface{}
	CheckedValue        string
	UncheckedValue      string
	Value               string
}

// NewCheckbox creates a new checkbox with the given parameters
func NewCheckbox(form interface{}, method string, label string, helpText string, labelOptions map[string]interface{}, checkboxOptions map[string]interface{}, contentWrapperOptions map[string]interface{}, checkedValue string, uncheckedValue string) *Checkbox {
	if labelOptions == nil {
		labelOptions = make(map[string]interface{})
	}
	if checkboxOptions == nil {
		checkboxOptions = make(map[string]interface{})
	}
	if contentWrapperOptions == nil {
		contentWrapperOptions = make(map[string]interface{})
	}

	// Set default values if not provided
	if checkedValue == "" {
		checkedValue = "1"
	}
	if uncheckedValue == "" {
		uncheckedValue = "0"
	}

	// Set value if multiple is true in checkbox options
	value := checkedValue
	if multiple, ok := checkboxOptions["multiple"].(bool); ok && multiple {
		value = checkedValue
	}

	return &Checkbox{
		Form:                 form,
		Method:              method,
		Label:               label,
		HelpText:            helpText,
		LabelOptions:        labelOptions,
		CheckboxOptions:     checkboxOptions,
		ContentWrapperOptions: contentWrapperOptions,
		CheckedValue:        checkedValue,
		UncheckedValue:      uncheckedValue,
		Value:               value,
	}
}

// GetLabelContent returns the label content
func (c *Checkbox) GetLabelContent() string {
	return c.Label
}

// GetHelpTextContent returns the help text content
func (c *Checkbox) GetHelpTextContent() string {
	return c.HelpText
}

// RenderLabelWithHelpText generates the HTML for the label and help text
func (c *Checkbox) RenderLabelWithHelpText() string {
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
	c.LabelOptions["for"] = c.Method

	// Add label
	parts = append(parts, fmt.Sprintf(`<label %s>%s</label>`, formatAttributes(c.LabelOptions), c.GetLabelContent()))

	// Add help text if present
	if c.GetHelpTextContent() != "" {
		parts = append(parts, fmt.Sprintf(`<div class="gl-form-text">%s</div>`, c.GetHelpTextContent()))
	}

	return strings.Join(parts, "\n")
}

// Render generates the HTML for the checkbox
func (c *Checkbox) Render() string {
	var parts []string

	// Add content wrapper class to content wrapper options
	if _, ok := c.ContentWrapperOptions["class"]; !ok {
		c.ContentWrapperOptions["class"] = "gl-form-checkbox custom-control custom-checkbox"
	} else {
		class := c.ContentWrapperOptions["class"].(string)
		if !strings.Contains(class, "gl-form-checkbox") {
			c.ContentWrapperOptions["class"] = fmt.Sprintf("gl-form-checkbox custom-control custom-checkbox %s", class)
		}
	}

	// Start content wrapper div
	parts = append(parts, fmt.Sprintf(`<div %s>`, formatAttributes(c.ContentWrapperOptions)))

	// Add checkbox input
	// In a real implementation, this would use the form builder to generate the checkbox
	// For now, we'll just generate a basic checkbox input
	checkboxAttrs := make(map[string]interface{})
	for k, v := range c.CheckboxOptions {
		checkboxAttrs[k] = v
	}
	checkboxAttrs["type"] = "checkbox"
	checkboxAttrs["name"] = c.Method
	checkboxAttrs["id"] = c.Method
	checkboxAttrs["value"] = c.Value
	parts = append(parts, fmt.Sprintf(`<input %s>`, formatAttributes(checkboxAttrs)))

	// Add label and help text
	parts = append(parts, c.RenderLabelWithHelpText())

	// Close content wrapper div
	parts = append(parts, "</div>")

	return strings.Join(parts, "\n")
}
