package ui

// InputOptions represents the configuration options for checkbox and radio inputs
type InputOptions struct {
	CSSClasses  []string
	Additional  map[string]interface{}
}

// CheckboxRadioOptions provides functionality for formatting checkbox and radio input options
type CheckboxRadioOptions struct {
	InputOptions InputOptions
}

// FormattedInputOptions returns the formatted input options with default values
func (c *CheckboxRadioOptions) FormattedInputOptions() InputOptions {
	return FormatInputOptions(c.InputOptions, []string{"custom-control-input"})
}

// FormatInputOptions formats the input options with default values
func FormatInputOptions(options InputOptions, defaultCSSClasses []string) InputOptions {
	if options.CSSClasses == nil {
		options.CSSClasses = defaultCSSClasses
	}
	if options.Additional == nil {
		options.Additional = make(map[string]interface{})
	}
	return options
}
