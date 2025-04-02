package layouts

import (
	"fmt"
	"html/template"
)

// NewSettingsSectionComponent creates a new SettingsSectionComponent instance
func NewSettingsSectionComponent(title, description, content string) *SettingsSectionComponent {
	return &SettingsSectionComponent{
		BaseComponent: BaseComponent{
			IsVisible: true,
		},
		Title:       title,
		Description: description,
		Content:     content,
	}
}

// Render renders the settings section component
func (c *SettingsSectionComponent) Render() (template.HTML, error) {
	if !c.IsVisible {
		return "", nil
	}

	html := fmt.Sprintf(`
		<div class="settings-section">
			<div class="settings-section-header">
				<h2>%s</h2>
				<p class="description">%s</p>
			</div>
			<div class="settings-section-content">
				%s
			</div>
		</div>
	`, c.Title, c.Description, c.Content)

	return template.HTML(html), nil
}

// SetContent sets the component's content
func (c *SettingsSectionComponent) SetContent(content string) {
	c.Content = content
}

// GetTitle returns the component's title
func (c *SettingsSectionComponent) GetTitle() string {
	return c.Title
}

// GetDescription returns the component's description
func (c *SettingsSectionComponent) GetDescription() string {
	return c.Description
}
