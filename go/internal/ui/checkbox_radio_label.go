package ui

// LabelOptions represents the configuration options for a label
type LabelOptions struct {
	Value       string
	CSSClasses  []string
	Additional  map[string]interface{}
}

// LabelWithHelpText is the interface that defines the behavior for labels with help text
type LabelWithHelpText interface {
	// RenderLabelWithHelpText renders a label with help text for a form field
	RenderLabelWithHelpText() string
	// RenderLabelTagWithHelpText renders a standalone label tag with help text
	RenderLabelTagWithHelpText() string
}

// BaseLabelWithHelpText provides the base implementation for labels with help text
type BaseLabelWithHelpText struct {
	Method          string
	Name            string
	LabelContent    string
	HelpTextContent string
	Value           string
	LabelOptions    LabelOptions
}

// FormatOptions formats the label options with default values
func FormatOptions(options LabelOptions, defaultCSSClasses []string) LabelOptions {
	if options.CSSClasses == nil {
		options.CSSClasses = defaultCSSClasses
	}
	if options.Additional == nil {
		options.Additional = make(map[string]interface{})
	}
	if options.Value == "" {
		options.Value = ""
	}
	return options
}
