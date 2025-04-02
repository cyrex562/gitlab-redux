package ui

import (
	"reflect"
	"strings"
)

// Component represents the base component class
type Component struct{}

// FilterAttribute filters a given value against a list of allowed values
// If no value is given or value is not allowed, it returns the default value
//
// @param value - The value to filter
// @param allowedValues - The list of allowed values
// @param defaultValue - The default value to return if the value is not allowed
func (c *Component) FilterAttribute(value interface{}, allowedValues []interface{}, defaultValue interface{}) interface{} {
	if value == nil {
		return defaultValue
	}

	for _, allowedValue := range allowedValues {
		if reflect.DeepEqual(value, allowedValue) {
			return value
		}
	}

	return defaultValue
}

// FormatOptions adds CSS classes and additional options to an existing options hash
//
// @param options - The existing options hash
// @param cssClasses - The CSS classes to add
// @param additionalOptions - Additional options to add
func (c *Component) FormatOptions(options map[string]interface{}, cssClasses []string, additionalOptions map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// Copy existing options
	for k, v := range options {
		result[k] = v
	}

	// Add CSS classes
	existingClasses := ""
	if classes, ok := options["class"].(string); ok {
		existingClasses = classes
	}

	allClasses := append(cssClasses, existingClasses)
	result["class"] = strings.Join(allClasses, " ")

	// Add additional options
	for k, v := range additionalOptions {
		result[k] = v
	}

	return result
}

// FormatAttributes formats a map of attributes into an HTML attribute string
//
// @param attributes - The map of attributes to format
func (c *Component) FormatAttributes(attributes map[string]interface{}) string {
	var parts []string

	for k, v := range range attributes {
		if v == nil {
			continue
		}

		switch val := v.(type) {
		case bool:
			if val {
				parts = append(parts, k)
			}
		case string:
			if val != "" {
				parts = append(parts, fmt.Sprintf(`%s="%s"`, k, val))
			}
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			parts = append(parts, fmt.Sprintf(`%s="%v"`, k, val))
		default:
			parts = append(parts, fmt.Sprintf(`%s="%v"`, k, val))
		}
	}

	return strings.Join(parts, " ")
}
