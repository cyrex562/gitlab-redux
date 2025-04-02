package onboarding

import (
	"fmt"
	"strings"
)

// Component represents the base onboarding component
type Component struct {
	HTMLOptions map[string]interface{}
}

// NewComponent creates a new base component
func NewComponent() *Component {
	return &Component{
		HTMLOptions: make(map[string]interface{}),
	}
}

// FilterAttribute filters a given value against a list of allowed values
// If no value is given or value is not allowed return default one
func FilterAttribute(value interface{}, allowedValues []interface{}, defaultValue interface{}) interface{} {
	if value == nil {
		return defaultValue
	}

	for _, allowed := range allowedValues {
		if value == allowed {
			return value
		}
	}

	return defaultValue
}

// FormatOptions adds CSS classes and additional options to an existing options map
func FormatOptions(options map[string]interface{}, cssClasses []string, additionalOptions map[string]interface{}) map[string]interface{} {
	// Create a copy of the options map
	result := make(map[string]interface{})
	for k, v := range options {
		result[k] = v
	}

	// Handle CSS classes
	if existingClasses, ok := result["class"].(string); ok {
		classes := strings.Split(existingClasses, " ")
		classes = append(classes, cssClasses...)
		result["class"] = strings.Join(classes, " ")
	} else {
		result["class"] = strings.Join(cssClasses, " ")
	}

	// Add additional options
	for k, v := range additionalOptions {
		result[k] = v
	}

	return result
}

// GetHTMLOptions returns the HTML options map
func (c *Component) GetHTMLOptions() map[string]interface{} {
	return c.HTMLOptions
}

// SetHTMLOptions sets the HTML options map
func (c *Component) SetHTMLOptions(options map[string]interface{}) {
	c.HTMLOptions = options
}

// AddHTMLOption adds a single HTML option
func (c *Component) AddHTMLOption(key string, value interface{}) {
	c.HTMLOptions[key] = value
}

// GetHTMLAttributes returns a string of HTML attributes
func (c *Component) GetHTMLAttributes() string {
	var attrs []string
	for k, v := range c.HTMLOptions {
		attrs = append(attrs, fmt.Sprintf(`%s="%v"`, k, v))
	}
	return strings.Join(attrs, " ")
}
