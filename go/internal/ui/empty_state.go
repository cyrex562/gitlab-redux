package ui

import (
	"fmt"
	"strings"
)

// EmptyState represents an empty state component
type EmptyState struct {
	Compact              bool
	Title                string
	SVGPath              string
	PrimaryButtonText    string
	PrimaryButtonLink    string
	PrimaryButtonOptions map[string]interface{}
	SecondaryButtonText  string
	SecondaryButtonLink  string
	SecondaryButtonOptions map[string]interface{}
	EmptyStateOptions    map[string]interface{}
	Description          string
}

// NewEmptyState creates a new empty state with the given parameters
func NewEmptyState(compact bool, title string, svgPath string, primaryButtonText string, primaryButtonLink string, primaryButtonOptions map[string]interface{}, secondaryButtonText string, secondaryButtonLink string, secondaryButtonOptions map[string]interface{}, emptyStateOptions map[string]interface{}, description string) *EmptyState {
	if primaryButtonOptions == nil {
		primaryButtonOptions = make(map[string]interface{})
	}
	if secondaryButtonOptions == nil {
		secondaryButtonOptions = make(map[string]interface{})
	}
	if emptyStateOptions == nil {
		emptyStateOptions = make(map[string]interface{})
	}

	return &EmptyState{
		Compact:              compact,
		Title:                title,
		SVGPath:              svgPath,
		PrimaryButtonText:    primaryButtonText,
		PrimaryButtonLink:    primaryButtonLink,
		PrimaryButtonOptions: primaryButtonOptions,
		SecondaryButtonText:  secondaryButtonText,
		SecondaryButtonLink:  secondaryButtonLink,
		SecondaryButtonOptions: secondaryButtonOptions,
		EmptyStateOptions:    emptyStateOptions,
		Description:          description,
	}
}

// GetEmptyStateClass returns the CSS class for the empty state
func (e *EmptyState) GetEmptyStateClass() string {
	if e.Compact {
		return "gl-flex-row"
	}
	return "gl-text-center gl-flex-col"
}

// GetImageClass returns the CSS class for the image
func (e *EmptyState) GetImageClass() string {
	if e.Compact {
		return "gl-hidden sm:gl-block gl-px-4"
	}
	return "gl-max-w-full"
}

// GetContentWrapperClass returns the CSS class for the content wrapper
func (e *EmptyState) GetContentWrapperClass() string {
	if e.Compact {
		return "gl-grow gl-basis-0 gl-px-4"
	}
	return "gl-m-auto gl-p-5"
}

// GetTitleClass returns the CSS class for the title
func (e *EmptyState) GetTitleClass() string {
	if e.Compact {
		return "h5"
	}
	return "h4"
}

// GetButtonWrapperClass returns the CSS class for the button wrapper
func (e *EmptyState) GetButtonWrapperClass() string {
	if e.Compact {
		return ""
	}
	return "gl-justify-center"
}

// Render generates the HTML for the empty state
func (e *EmptyState) Render() string {
	var parts []string

	// Add empty state class to empty state options
	emptyStateClass := e.GetEmptyStateClass()
	if _, ok := e.EmptyStateOptions["class"]; !ok {
		e.EmptyStateOptions["class"] = fmt.Sprintf("gl-flex gl-empty-state %s", emptyStateClass)
	} else {
		class := e.EmptyStateOptions["class"].(string)
		if !strings.Contains(class, "gl-flex") {
			e.EmptyStateOptions["class"] = fmt.Sprintf("gl-flex gl-empty-state %s %s", emptyStateClass, class)
		}
	}

	// Start empty state section
	parts = append(parts, fmt.Sprintf(`<section %s>`, formatAttributes(e.EmptyStateOptions)))

	// Add SVG image if present
	if e.SVGPath != "" {
		imageClass := e.GetImageClass()
		parts = append(parts, fmt.Sprintf(`<div class="%s"><img src="%s" alt="" class="gl-dark-invert-keep-hue"></div>`, imageClass, e.SVGPath))
	}

	// Add content wrapper
	contentWrapperClass := e.GetContentWrapperClass()
	parts = append(parts, fmt.Sprintf(`<div class="gl-empty-state-content gl-mx-auto gl-my-0 %s">`, contentWrapperClass))

	// Add title
	titleClass := e.GetTitleClass()
	parts = append(parts, fmt.Sprintf(`<h1 class="gl-text-size-h-display gl-leading-36 gl-mt-0 gl-mb-0 %s">%s</h1>`, titleClass, e.Title))

	// Add description if present
	if e.Description != "" {
		parts = append(parts, fmt.Sprintf(`<p class="gl-mt-4 gl-mb-0" data-testid="empty-state-description">%s</p>`, e.Description))
	}

	// Add buttons if present
	if e.PrimaryButtonText != "" || e.SecondaryButtonText != "" {
		buttonWrapperClass := e.GetButtonWrapperClass()
		parts = append(parts, fmt.Sprintf(`<div class="gl-flex gl-flex-wrap gl-mt-5 gl-gap-3 %s">`, buttonWrapperClass))

		// Add primary button if present
		if e.PrimaryButtonText != "" {
			primaryButtonOptions := make(map[string]interface{})
			for k, v := range e.PrimaryButtonOptions {
				primaryButtonOptions[k] = v
			}
			primaryButtonOptions["class"] = "!gl-ml-0"

			primaryButton := NewButton(
				ButtonCategoryPrimary,
				ButtonVariantConfirm,
				ButtonSizeMedium,
				ButtonTypeButton,
				false, // disabled
				false, // loading
				false, // block
				false, // label
				false, // selected
				"", // icon
				e.PrimaryButtonLink, // href
				false, // form
				"", // target
				"", // method
				primaryButtonOptions, // button options
				"", // button text classes
				"", // icon classes
				"", // icon content
				e.PrimaryButtonText, // content
			)
			parts = append(parts, primaryButton.Render())
		}

		// Add secondary button if present
		if e.SecondaryButtonText != "" {
			secondaryButtonOptions := make(map[string]interface{})
			for k, v := range e.SecondaryButtonOptions {
				secondaryButtonOptions[k] = v
			}
			if e.PrimaryButtonText == "" {
				secondaryButtonOptions["class"] = "!gl-ml-0"
			}

			secondaryButton := NewButton(
				ButtonCategoryPrimary,
				ButtonVariantDefault,
				ButtonSizeMedium,
				ButtonTypeButton,
				false, // disabled
				false, // loading
				false, // block
				false, // label
				false, // selected
				"", // icon
				e.SecondaryButtonLink, // href
				false, // form
				"", // target
				"", // method
				secondaryButtonOptions, // button options
				"", // button text classes
				"", // icon classes
				"", // icon content
				e.SecondaryButtonText, // content
			)
			parts = append(parts, secondaryButton.Render())
		}

		// Close button wrapper
		parts = append(parts, "</div>")
	}

	// Close content wrapper
	parts = append(parts, "</div>")

	// Close empty state section
	parts = append(parts, "</section>")

	return strings.Join(parts, "\n")
}
