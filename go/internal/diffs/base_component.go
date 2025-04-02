package diffs

import (
	"html/template"
)

// NewBaseComponent creates a new BaseComponent instance
func NewBaseComponent(opts ComponentOptions) *BaseComponent {
	return &BaseComponent{
		Content:   opts.Content,
		IsVisible: opts.IsVisible,
	}
}

// Render renders the base component
func (c *BaseComponent) Render() (template.HTML, error) {
	if !c.IsVisible {
		return "", nil
	}

	return template.HTML(c.Content), nil
}

// GetContent returns the component's content
func (c *BaseComponent) GetContent() string {
	return c.Content
}

// SetContent sets the component's content
func (c *BaseComponent) SetContent(content string) {
	c.Content = content
}

// IsComponentVisible returns whether the component is visible
func (c *BaseComponent) IsComponentVisible() bool {
	return c.IsVisible
}

// SetVisibility sets the component's visibility
func (c *BaseComponent) SetVisibility(visible bool) {
	c.IsVisible = visible
}
