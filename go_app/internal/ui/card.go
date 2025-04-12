package ui

import (
	"fmt"
	"strings"
)

// Card represents a card component
type Card struct {
	CardOptions   map[string]interface{}
	HeaderOptions map[string]interface{}
	BodyOptions   map[string]interface{}
	FooterOptions map[string]interface{}
	Card          map[string]interface{}
	Header        string
	Body          string
	Footer        string
}

// NewCard creates a new card with the given parameters
func NewCard(card map[string]interface{}, cardOptions map[string]interface{}, headerOptions map[string]interface{}, bodyOptions map[string]interface{}, footerOptions map[string]interface{}, header string, body string, footer string) *Card {
	if card == nil {
		card = make(map[string]interface{})
	}
	if cardOptions == nil {
		cardOptions = make(map[string]interface{})
	}
	if headerOptions == nil {
		headerOptions = make(map[string]interface{})
	}
	if bodyOptions == nil {
		bodyOptions = make(map[string]interface{})
	}
	if footerOptions == nil {
		footerOptions = make(map[string]interface{})
	}

	return &Card{
		CardOptions:   cardOptions,
		HeaderOptions: headerOptions,
		BodyOptions:   bodyOptions,
		FooterOptions: footerOptions,
		Card:          card,
		Header:        header,
		Body:          body,
		Footer:        footer,
	}
}

// HasHeader returns whether the card has a header
func (c *Card) HasHeader() bool {
	return c.Header != "" || (c.Card["header"] != nil && c.Card["header"] != "")
}

// GetHeader returns the header content
func (c *Card) GetHeader() string {
	if c.Header != "" {
		return c.Header
	}
	if header, ok := c.Card["header"].(string); ok {
		return header
	}
	return ""
}

// HasFooter returns whether the card has a footer
func (c *Card) HasFooter() bool {
	return c.Footer != "" || (c.Card["footer"] != nil && c.Card["footer"] != "")
}

// GetFooter returns the footer content
func (c *Card) GetFooter() string {
	if c.Footer != "" {
		return c.Footer
	}
	if footer, ok := c.Card["footer"].(string); ok {
		return footer
	}
	return ""
}

// GetBody returns the body content
func (c *Card) GetBody() string {
	if c.Body != "" {
		return c.Body
	}
	if body, ok := c.Card["body"].(string); ok {
		return body
	}
	return ""
}

// Render generates the HTML for the card
func (c *Card) Render() string {
	var parts []string

	// Add card class to card options
	if _, ok := c.CardOptions["class"]; !ok {
		c.CardOptions["class"] = "gl-card"
	} else {
		class := c.CardOptions["class"].(string)
		if !strings.Contains(class, "gl-card") {
			c.CardOptions["class"] = fmt.Sprintf("gl-card %s", class)
		}
	}

	// Start card div
	parts = append(parts, fmt.Sprintf(`<div %s>`, formatAttributes(c.CardOptions)))

	// Add header if present
	if c.HasHeader() {
		// Add header class to header options
		if _, ok := c.HeaderOptions["class"]; !ok {
			c.HeaderOptions["class"] = "gl-card-header"
		} else {
			class := c.HeaderOptions["class"].(string)
			if !strings.Contains(class, "gl-card-header") {
				c.HeaderOptions["class"] = fmt.Sprintf("gl-card-header %s", class)
			}
		}
		parts = append(parts, fmt.Sprintf(`<div %s>%s</div>`, formatAttributes(c.HeaderOptions), c.GetHeader()))
	}

	// Add body
	// Add body class to body options
	if _, ok := c.BodyOptions["class"]; !ok {
		c.BodyOptions["class"] = "gl-card-body"
	} else {
		class := c.BodyOptions["class"].(string)
		if !strings.Contains(class, "gl-card-body") {
			c.BodyOptions["class"] = fmt.Sprintf("gl-card-body %s", class)
		}
	}
	parts = append(parts, fmt.Sprintf(`<div %s>%s</div>`, formatAttributes(c.BodyOptions), c.GetBody()))

	// Add footer if present
	if c.HasFooter() {
		// Add footer class to footer options
		if _, ok := c.FooterOptions["class"]; !ok {
			c.FooterOptions["class"] = "gl-card-footer"
		} else {
			class := c.FooterOptions["class"].(string)
			if !strings.Contains(class, "gl-card-footer") {
				c.FooterOptions["class"] = fmt.Sprintf("gl-card-footer %s", class)
			}
		}
		parts = append(parts, fmt.Sprintf(`<div %s>%s</div>`, formatAttributes(c.FooterOptions), c.GetFooter()))
	}

	// Close card div
	parts = append(parts, "</div>")

	return strings.Join(parts, "\n")
}
