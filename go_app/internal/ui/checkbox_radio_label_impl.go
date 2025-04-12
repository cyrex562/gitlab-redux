package ui

import (
	"fmt"
	"strings"
)

// RenderLabelWithHelpText implements the LabelWithHelpText interface for form labels
func (b *BaseLabelWithHelpText) RenderLabelWithHelpText() string {
	options := b.formattedLabelOptions()
	return fmt.Sprintf(`<label for="%s" %s>%s</label>`,
		b.Method,
		formatAttributes(options),
		b.labelEntry())
}

// RenderLabelTagWithHelpText implements the LabelWithHelpText interface for standalone labels
func (b *BaseLabelWithHelpText) RenderLabelTagWithHelpText() string {
	options := b.formattedLabelOptions()
	return fmt.Sprintf(`<label for="%s" %s>%s</label>`,
		b.Name,
		formatAttributes(options),
		b.labelEntry())
}

// labelEntry returns the HTML content for the label, including help text if present
func (b *BaseLabelWithHelpText) labelEntry() string {
	if b.HelpTextContent != "" {
		return fmt.Sprintf(`<span>%s</span><p class="help-text" data-testid="pajamas-component-help-text">%s</p>`,
			b.LabelContent,
			b.HelpTextContent)
	}
	return fmt.Sprintf(`<span>%s</span>`, b.LabelContent)
}

// formattedLabelOptions returns the formatted label options with default values
func (b *BaseLabelWithHelpText) formattedLabelOptions() LabelOptions {
	return FormatOptions(b.LabelOptions, []string{"custom-control-label"})
}

// formatAttributes converts the label options into HTML attributes
func formatAttributes(options LabelOptions) string {
	var attrs []string

	// Add CSS classes
	if len(options.CSSClasses) > 0 {
		attrs = append(attrs, fmt.Sprintf(`class="%s"`, strings.Join(options.CSSClasses, " ")))
	}

	// Add value if present
	if options.Value != "" {
		attrs = append(attrs, fmt.Sprintf(`value="%s"`, options.Value))
	}

	// Add additional attributes
	for key, value := range options.Additional {
		attrs = append(attrs, fmt.Sprintf(`%s="%v"`, key, value))
	}

	return strings.Join(attrs, " ")
}
