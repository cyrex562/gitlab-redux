package onboarding

import (
	"fmt"
	"html/template"
	"strings"
)

// Variant represents the visual variant of an action card
type Variant string

const (
	// VariantDefault represents the default variant
	VariantDefault Variant = "default"
	// VariantSuccess represents the success variant
	VariantSuccess Variant = "success"
	// VariantPromo represents the promo variant
	VariantPromo Variant = "promo"
)

// ActionCardComponent represents an action card in the onboarding flow
type ActionCardComponent struct {
	*Component
	Title       string
	Description string
	Icon        string
	Href        string
	Variant     Variant
	LinkOptions map[string]interface{}
	Content     template.HTML
}

// NewActionCardComponent creates a new action card component
func NewActionCardComponent(
	title string,
	description string,
	icon string,
	href string,
	variant Variant,
	linkOptions map[string]interface{},
	htmlOptions map[string]interface{},
) *ActionCardComponent {
	// Filter variant against allowed values
	filteredVariant := FilterAttribute(variant, []interface{}{
		VariantDefault,
		VariantSuccess,
		VariantPromo,
	}, VariantDefault).(Variant)

	return &ActionCardComponent{
		Component:   NewComponent(),
		Title:       title,
		Description: description,
		Icon:        icon,
		Href:        href,
		Variant:     filteredVariant,
		LinkOptions: linkOptions,
	}
}

// Render renders the action card component
func (c *ActionCardComponent) Render() (template.HTML, error) {
	// Build card classes
	cardClasses := []string{
		"action-card",
		fmt.Sprintf("action-card-%s", c.Variant),
	}

	// Set HTML options with card classes
	c.SetHTMLOptions(FormatOptions(c.GetHTMLOptions(), cardClasses, nil))

	// Build the HTML
	var html string
	if c.HasLink() {
		// Build link attributes
		linkAttrs := FormatOptions(c.LinkOptions, []string{"gl-link", "action-card-title"}, nil)
		linkHTML := fmt.Sprintf(`<a href="%s" %s>`, c.Href, c.getHTMLAttributes(linkAttrs))

		html = fmt.Sprintf(`
			<div %s>
				%s
					%s
					%s
					%s
				</a>
				<p class="action-card-text">%s</p>
				%s
			</div>
		`, c.GetHTMLAttributes(), linkHTML, c.getIconHTML(c.getCardIcon()), c.Title, c.getIconHTML("arrow-right", "action-card-arrow"), c.Description, c.Content)
	} else {
		html = fmt.Sprintf(`
			<div %s>
				<div class="action-card-title">
					%s
					%s
				</div>
				<p class="action-card-text">%s</p>
				%s
			</div>
		`, c.GetHTMLAttributes(), c.getIconHTML(c.getCardIcon()), c.Title, c.Description, c.Content)
	}

	return template.HTML(html), nil
}

// HasLink returns whether the card has a link
func (c *ActionCardComponent) HasLink() bool {
	return c.Href != ""
}

// getCardIcon returns the appropriate icon based on the variant
func (c *ActionCardComponent) getCardIcon() string {
	if c.Variant == VariantSuccess {
		return "check"
	}
	return c.Icon
}

// getIconHTML generates the HTML for an icon
func (c *ActionCardComponent) getIconHTML(icon string, classes ...string) string {
	iconClasses := append([]string{"sprite-icon"}, classes...)
	return fmt.Sprintf(`<i class="%s" data-icon="%s"></i>`, strings.Join(iconClasses, " "), icon)
}

// getHTMLAttributes generates HTML attributes from a map
func (c *ActionCardComponent) getHTMLAttributes(options map[string]interface{}) string {
	var attrs []string
	for k, v := range options {
		attrs = append(attrs, fmt.Sprintf(`%s="%v"`, k, v))
	}
	return strings.Join(attrs, " ")
}

// SetContent sets the content of the action card
func (c *ActionCardComponent) SetContent(content template.HTML) {
	c.Content = content
}

// GetTitle returns the title of the action card
func (c *ActionCardComponent) GetTitle() string {
	return c.Title
}

// GetDescription returns the description of the action card
func (c *ActionCardComponent) GetDescription() string {
	return c.Description
}

// GetIcon returns the icon of the action card
func (c *ActionCardComponent) GetIcon() string {
	return c.Icon
}

// GetVariant returns the variant of the action card
func (c *ActionCardComponent) GetVariant() Variant {
	return c.Variant
}
