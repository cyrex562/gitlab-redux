package ui

import (
	"fmt"
	"strings"
)

// AccordionState represents the possible states of an accordion item
type AccordionState string

const (
	// AccordionStateOpened represents an opened accordion item
	AccordionStateOpened AccordionState = "opened"
	// AccordionStateClosed represents a closed accordion item
	AccordionStateClosed AccordionState = "closed"
)

// AccordionItem represents an accordion item component
type AccordionItem struct {
	Title         string
	State         AccordionState
	ButtonOptions map[string]interface{}
	Content       string
}

// NewAccordionItem creates a new accordion item with the given parameters
func NewAccordionItem(title string, state AccordionState, buttonOptions map[string]interface{}, content string) *AccordionItem {
	if buttonOptions == nil {
		buttonOptions = make(map[string]interface{})
	}

	// Set default button options if not provided
	if _, ok := buttonOptions["aria-controls"]; !ok {
		buttonOptions["aria-controls"] = "accordion-item"
	}

	return &AccordionItem{
		Title:         title,
		State:         state,
		ButtonOptions: buttonOptions,
		Content:       content,
	}
}

// Icon returns the appropriate icon name based on the accordion state
func (a *AccordionItem) Icon() string {
	if a.State == AccordionStateOpened {
		return "chevron-down"
	}
	return "chevron-right"
}

// BodyClass returns the CSS classes for the accordion body
func (a *AccordionItem) BodyClass() map[string]interface{} {
	if a.State == AccordionStateOpened {
		return map[string]interface{}{
			"class": "show",
		}
	}
	return map[string]interface{}{}
}

// IsExpanded returns whether the accordion item is expanded
func (a *AccordionItem) IsExpanded() bool {
	return a.State == AccordionStateOpened
}

// Render generates the HTML for the accordion item
func (a *AccordionItem) Render() string {
	// Build button attributes
	buttonAttrs := make(map[string]interface{})
	for k, v := range a.ButtonOptions {
		buttonAttrs[k] = v
	}
	buttonAttrs["aria-expanded"] = a.IsExpanded()

	// Build body attributes
	bodyAttrs := a.BodyClass()
	bodyAttrs["class"] = fmt.Sprintf("accordion-item gl-mt-3 gl-text-base collapse %s",
		strings.Join(bodyAttrs["class"].([]string), " "))

	return fmt.Sprintf(`<div class="gl-accordion-item">
	<h3 class="gl-accordion-item-header">
		<button type="button" %s>
			<span class="js-chevron-icon">%s</span>
			%s
		</button>
	</h3>
	<div %s>
		%s
	</div>
</div>`,
		formatAttributes(buttonAttrs),
		a.Icon(),
		a.Title,
		formatAttributes(bodyAttrs),
		a.Content)
}
