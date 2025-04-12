package ui

import (
	"fmt"
	"strings"
)

// BreadcrumbItem represents a single item in a breadcrumb
type BreadcrumbItem struct {
	Href string
	Text string
}

// Breadcrumb represents a breadcrumb component
type Breadcrumb struct {
	HTMLOptions map[string]interface{}
	Items       []BreadcrumbItem
}

// NewBreadcrumb creates a new breadcrumb with the given parameters
func NewBreadcrumb(htmlOptions map[string]interface{}, items []BreadcrumbItem) *Breadcrumb {
	if htmlOptions == nil {
		htmlOptions = make(map[string]interface{})
	}
	if items == nil {
		items = []BreadcrumbItem{}
	}

	return &Breadcrumb{
		HTMLOptions: htmlOptions,
		Items:       items,
	}
}

// AddItem adds a new item to the breadcrumb
func (b *Breadcrumb) AddItem(href string, text string) {
	b.Items = append(b.Items, BreadcrumbItem{
		Href: href,
		Text: text,
	})
}

// Render generates the HTML for the breadcrumb
func (b *Breadcrumb) Render() string {
	// Build breadcrumb attributes
	breadcrumbAttrs := make(map[string]interface{})
	for k, v := range b.HTMLOptions {
		breadcrumbAttrs[k] = v
	}

	// Add default attributes
	breadcrumbAttrs["aria-label"] = "Breadcrumbs"
	breadcrumbAttrs["class"] = "gl-breadcrumbs"

	// Build the content
	var contentParts []string

	// Add items
	for _, item := range b.Items {
		contentParts = append(contentParts, fmt.Sprintf(`<li class="gl-breadcrumb-item"><a href="%s">%s</a></li>`, item.Href, item.Text))
	}

	// Render the breadcrumb
	return fmt.Sprintf(`<nav %s><ul class="gl-breadcrumb-list breadcrumb js-breadcrumbs-list">%s</ul></nav>`,
		formatAttributes(breadcrumbAttrs),
		strings.Join(contentParts, "\n"))
}
