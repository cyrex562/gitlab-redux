package layouts

import (
	"fmt"
	"html/template"
)

// NewPageHeadingComponent creates a new PageHeadingComponent instance
func NewPageHeadingComponent(title, description string, actions []Action) *PageHeadingComponent {
	return &PageHeadingComponent{
		BaseComponent: BaseComponent{
			IsVisible: true,
		},
		Title:       title,
		Description: description,
		Actions:     actions,
	}
}

// Render renders the page heading component
func (c *PageHeadingComponent) Render() (template.HTML, error) {
	if !c.IsVisible {
		return "", nil
	}

	// Build actions HTML
	actionsHTML := ""
	for _, action := range c.Actions {
		dataAttrs := ""
		for key, value := range action.Data {
			dataAttrs += fmt.Sprintf(` data-%s="%s"`, key, value)
		}

		actionsHTML += fmt.Sprintf(`
			<a href="%s" class="%s" data-method="%s"%s>
				<i class="%s"></i> %s
			</a>
		`, action.URL, action.Class, action.Method, dataAttrs, action.Icon, action.Label)
	}

	// Build the complete component HTML
	html := fmt.Sprintf(`
		<div class="page-heading">
			<div class="page-heading-content">
				<h1>%s</h1>
				<p class="description">%s</p>
				<div class="actions">
					%s
				</div>
			</div>
		</div>
	`, c.Title, c.Description, actionsHTML)

	return template.HTML(html), nil
}

// AddAction adds a new action to the component
func (c *PageHeadingComponent) AddAction(action Action) {
	c.Actions = append(c.Actions, action)
}

// GetTitle returns the component's title
func (c *PageHeadingComponent) GetTitle() string {
	return c.Title
}

// GetDescription returns the component's description
func (c *PageHeadingComponent) GetDescription() string {
	return c.Description
}

// GetActions returns the component's actions
func (c *PageHeadingComponent) GetActions() []Action {
	return c.Actions
}
