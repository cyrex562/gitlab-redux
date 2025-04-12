package layouts

import (
	"fmt"
	"html/template"
)

// NewCRUDComponent creates a new CRUDComponent instance
func NewCRUDComponent(title, description string, actions []Action, form interface{}) *CRUDComponent {
	return &CRUDComponent{
		BaseComponent: BaseComponent{
			IsVisible: true,
		},
		Title:       title,
		Description: description,
		Actions:     actions,
		Form:        form,
	}
}

// Render renders the CRUD component
func (c *CRUDComponent) Render() (template.HTML, error) {
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

	// Build form HTML if present
	formHTML := ""
	if c.Form != nil {
		formHTML = fmt.Sprintf("%v", c.Form)
	}

	// Build the complete component HTML
	html := fmt.Sprintf(`
		<div class="crud-component">
			<div class="crud-header">
				<h2>%s</h2>
				<p class="description">%s</p>
				<div class="actions">
					%s
				</div>
			</div>
			<div class="crud-content">
				%s
			</div>
		</div>
	`, c.Title, c.Description, actionsHTML, formHTML)

	return template.HTML(html), nil
}

// AddAction adds a new action to the component
func (c *CRUDComponent) AddAction(action Action) {
	c.Actions = append(c.Actions, action)
}

// SetForm sets the form for the component
func (c *CRUDComponent) SetForm(form interface{}) {
	c.Form = form
}

// GetTitle returns the component's title
func (c *CRUDComponent) GetTitle() string {
	return c.Title
}

// GetDescription returns the component's description
func (c *CRUDComponent) GetDescription() string {
	return c.Description
}

// GetActions returns the component's actions
func (c *CRUDComponent) GetActions() []Action {
	return c.Actions
}
