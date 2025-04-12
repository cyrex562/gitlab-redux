package layouts

import (
	"fmt"
	"html/template"
)

// NewHorizontalSectionComponent creates a new HorizontalSectionComponent instance
func NewHorizontalSectionComponent(title, description, content string) *HorizontalSectionComponent {
	return &HorizontalSectionComponent{
		BaseComponent: BaseComponent{
			IsVisible: true,
		},
		Title:       title,
		Description: description,
		Content:     content,
	}
}

// Render renders the horizontal section component
func (c *HorizontalSectionComponent) Render() (template.HTML, error) {
	if !c.IsVisible {
		return "", nil
	}

	html := fmt.Sprintf(`
		<div class="horizontal-section">
			<div class="horizontal-section-header">
				<h3>%s</h3>
				<p class="description">%s</p>
			</div>
			<div class="horizontal-section-content">
				%s
			</div>
		</div>
	`, c.Title, c.Description, c.Content)

	return template.HTML(html), nil
}

// SetContent sets the component's content
func (c *HorizontalSectionComponent) SetContent(content string) {
	c.Content = content
}

// GetTitle returns the component's title
func (c *HorizontalSectionComponent) GetTitle() string {
	return c.Title
}

// GetDescription returns the component's description
func (c *HorizontalSectionComponent) GetDescription() string {
	return c.Description
}
