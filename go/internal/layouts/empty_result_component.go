package layouts

import (
	"fmt"
	"html/template"
)

// NewEmptyResultComponent creates a new EmptyResultComponent instance
func NewEmptyResultComponent(title, description, icon string, action *Action) *EmptyResultComponent {
	return &EmptyResultComponent{
		BaseComponent: BaseComponent{
			IsVisible: true,
		},
		Title:       title,
		Description: description,
		Icon:        icon,
		Action:      action,
	}
}

// Render renders the empty result component
func (c *EmptyResultComponent) Render() (template.HTML, error) {
	if !c.IsVisible {
		return "", nil
	}

	// Build action HTML if present
	actionHTML := ""
	if c.Action != nil {
		dataAttrs := ""
		for key, value := range c.Action.Data {
			dataAttrs += fmt.Sprintf(` data-%s="%s"`, key, value)
		}

		actionHTML = fmt.Sprintf(`
			<div class="action">
				<a href="%s" class="%s" data-method="%s"%s>
					<i class="%s"></i> %s
				</a>
			</div>
		`, c.Action.URL, c.Action.Class, c.Action.Method, dataAttrs, c.Action.Icon, c.Action.Label)
	}

	// Build the complete component HTML
	html := fmt.Sprintf(`
		<div class="empty-result">
			<div class="empty-result-content">
				<i class="%s"></i>
				<h3>%s</h3>
				<p>%s</p>
				%s
			</div>
		</div>
	`, c.Icon, c.Title, c.Description, actionHTML)

	return template.HTML(html), nil
}

// SetAction sets the action for the component
func (c *EmptyResultComponent) SetAction(action *Action) {
	c.Action = action
}

// GetTitle returns the component's title
func (c *EmptyResultComponent) GetTitle() string {
	return c.Title
}

// GetDescription returns the component's description
func (c *EmptyResultComponent) GetDescription() string {
	return c.Description
}

// GetIcon returns the component's icon
func (c *EmptyResultComponent) GetIcon() string {
	return c.Icon
}
