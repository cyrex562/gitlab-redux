package layouts

import (
	"fmt"
	"html/template"
)

// NewSettingsBlockComponent creates a new SettingsBlockComponent instance
func NewSettingsBlockComponent(title, description, content, helpText string) *SettingsBlockComponent {
	return &SettingsBlockComponent{
		BaseComponent: BaseComponent{
			IsVisible: true,
		},
		Title:       title,
		Description: description,
		Content:     content,
		HelpText:    helpText,
	}
}

// Render renders the settings block component
func (c *SettingsBlockComponent) Render() (template.HTML, error) {
	if !c.IsVisible {
		return "", nil
	}

	// Build help text HTML if present
	helpTextHTML := ""
	if c.HelpText != "" {
		helpTextHTML = fmt.Sprintf(`
			<div class="help-text">
				%s
			</div>
		`, c.HelpText)
	}

	// Build the complete component HTML
	html := fmt.Sprintf(`
		<div class="settings-block">
			<div class="settings-block-header">
				<h3>%s</h3>
				<p class="description">%s</p>
				%s
			</div>
			<div class="settings-block-content">
				%s
			</div>
		</div>
	`, c.Title, c.Description, helpTextHTML, c.Content)

	return template.HTML(html), nil
}

// SetContent sets the component's content
func (c *SettingsBlockComponent) SetContent(content string) {
	c.Content = content
}

// SetHelpText sets the component's help text
func (c *SettingsBlockComponent) SetHelpText(helpText string) {
	c.HelpText = helpText
}

// GetTitle returns the component's title
func (c *SettingsBlockComponent) GetTitle() string {
	return c.Title
}

// GetDescription returns the component's description
func (c *SettingsBlockComponent) GetDescription() string {
	return c.Description
}

// GetHelpText returns the component's help text
func (c *SettingsBlockComponent) GetHelpText() string {
	return c.HelpText
}
