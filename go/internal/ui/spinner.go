package ui

import (
	"fmt"
	"strings"
)

// SpinnerColor represents the available color options for the spinner
type SpinnerColor string

const (
	SpinnerColorLight SpinnerColor = "light"
	SpinnerColorDark  SpinnerColor = "dark"
)

// SpinnerSize represents the available size options for the spinner
type SpinnerSize string

const (
	SpinnerSizeSM SpinnerSize = "sm"
	SpinnerSizeMD SpinnerSize = "md"
	SpinnerSizeLG SpinnerSize = "lg"
	SpinnerSizeXL SpinnerSize = "xl"
)

// Spinner represents a spinner component
type Spinner struct {
	Color       SpinnerColor
	Inline      bool
	Label       string
	Size        SpinnerSize
	HTMLOptions map[string]interface{}
}

// NewSpinner creates a new spinner with the given parameters
func NewSpinner(
	color SpinnerColor,
	inline bool,
	label string,
	size SpinnerSize,
	htmlOptions map[string]interface{},
) *Spinner {
	if color == "" {
		color = SpinnerColorDark
	}
	if size == "" {
		size = SpinnerSizeSM
	}
	if label == "" {
		label = "Loading"
	}
	if htmlOptions == nil {
		htmlOptions = make(map[string]interface{})
	}

	return &Spinner{
		Color:       color,
		Inline:      inline,
		Label:       label,
		Size:        size,
		HTMLOptions: htmlOptions,
	}
}

// GetSpinnerClass returns the CSS classes for the spinner
func (s *Spinner) GetSpinnerClass() string {
	return fmt.Sprintf("gl-spinner gl-spinner-%s gl-spinner-%s !gl-align-text-bottom",
		s.Size, s.Color)
}

// GetHTMLOptions returns the formatted HTML options
func (s *Spinner) GetHTMLOptions() string {
	options := formatOptions(s.HTMLOptions, map[string]interface{}{
		"class": "gl-spinner-container",
		"role":  "status",
	})
	return formatAttributes(options)
}

// Render generates the HTML for the spinner
func (s *Spinner) Render() string {
	var parts []string

	// Determine the container tag based on inline property
	containerTag := "div"
	if s.Inline {
		containerTag = "span"
	}

	// Start container
	parts = append(parts, fmt.Sprintf("<%s %s>", containerTag, s.GetHTMLOptions()))

	// Spinner element
	parts = append(parts, fmt.Sprintf(`<span class="%s" aria-hidden="true"></span>`, s.GetSpinnerClass()))

	// Label
	parts = append(parts, fmt.Sprintf(`<span class="gl-sr-only !gl-absolute">%s</span>`, s.Label))

	// End container
	parts = append(parts, fmt.Sprintf("</%s>", containerTag))

	return strings.Join(parts, "\n")
}
