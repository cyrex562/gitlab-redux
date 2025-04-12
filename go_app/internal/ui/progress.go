package ui

import (
	"fmt"
	"strings"
)

// ProgressVariant represents the possible variants of a progress bar
type ProgressVariant string

const (
	// ProgressVariantPrimary represents a primary progress bar
	ProgressVariantPrimary ProgressVariant = "primary"
	// ProgressVariantSuccess represents a success progress bar
	ProgressVariantSuccess ProgressVariant = "success"
)

// Progress represents a progress component
type Progress struct {
	Value   int
	Variant ProgressVariant
}

// NewProgress creates a new progress with the given parameters
func NewProgress(value int, variant ProgressVariant) *Progress {
	// Filter variant to ensure it's valid
	allowedVariants := []ProgressVariant{ProgressVariantPrimary, ProgressVariantSuccess}
	isValidVariant := false
	for _, v := range allowedVariants {
		if variant == v {
			isValidVariant = true
			break
		}
	}
	if !isValidVariant {
		variant = ProgressVariantPrimary
	}

	// Ensure value is between 0 and 100
	if value < 0 {
		value = 0
	} else if value > 100 {
		value = 100
	}

	return &Progress{
		Value:   value,
		Variant: variant,
	}
}

// Render generates the HTML for the progress
func (p *Progress) Render() string {
	var parts []string

	// Start progress bar div
	parts = append(parts, `<div class="gl-progress-bar progress">`)

	// Add progress div
	progressClass := fmt.Sprintf("gl-progress gl-progress-bar-%s", p.Variant)
	progressStyle := fmt.Sprintf("width: %d%%;", p.Value)
	parts = append(parts, fmt.Sprintf(`<div class="%s" style="%s"></div>`, progressClass, progressStyle))

	// Close progress bar div
	parts = append(parts, "</div>")

	return strings.Join(parts, "\n")
}
