package ui

import (
	"fmt"
	"strings"
)

// AlertVariant represents the possible variants of an alert
type AlertVariant string

const (
	// AlertVariantInfo represents an info alert
	AlertVariantInfo AlertVariant = "info"
	// AlertVariantWarning represents a warning alert
	AlertVariantWarning AlertVariant = "warning"
	// AlertVariantSuccess represents a success alert
	AlertVariantSuccess AlertVariant = "success"
	// AlertVariantDanger represents a danger alert
	AlertVariantDanger AlertVariant = "danger"
	// AlertVariantTip represents a tip alert
	AlertVariantTip AlertVariant = "tip"
)

// Alert represents an alert component
type Alert struct {
	Title              string
	Variant            AlertVariant
	Dismissible        bool
	ShowIcon           bool
	AlertOptions       map[string]interface{}
	CloseButtonOptions map[string]interface{}
	Body               string
	Actions            string
}

// variantIcons maps alert variants to their corresponding icons
var variantIcons = map[AlertVariant]string{
	AlertVariantInfo:    "information-o",
	AlertVariantWarning: "warning",
	AlertVariantSuccess: "check-circle",
	AlertVariantDanger:  "error",
	AlertVariantTip:     "bulb",
}

// NewAlert creates a new alert with the given parameters
func NewAlert(title string, variant AlertVariant, dismissible bool, showIcon bool, alertOptions map[string]interface{}, closeButtonOptions map[string]interface{}, body string, actions string) *Alert {
	if alertOptions == nil {
		alertOptions = make(map[string]interface{})
	}
	if closeButtonOptions == nil {
		closeButtonOptions = make(map[string]interface{})
	}

	return &Alert{
		Title:              title,
		Variant:            variant,
		Dismissible:        dismissible,
		ShowIcon:           showIcon,
		AlertOptions:       alertOptions,
		CloseButtonOptions: closeButtonOptions,
		Body:               body,
		Actions:            actions,
	}
}

// BaseClass returns the base CSS classes for the alert
func (a *Alert) BaseClass() string {
	var classes []string

	// Add variant class
	classes = append(classes, fmt.Sprintf("gl-alert-%s", a.Variant))

	// Add dismissible class
	if !a.Dismissible {
		classes = append(classes, "gl-alert-not-dismissible")
	}

	// Add icon class
	if !a.ShowIcon {
		classes = append(classes, "gl-alert-no-icon")
	}

	// Add title class
	if a.Title != "" {
		classes = append(classes, "gl-alert-has-title")
	}

	return strings.Join(classes, " ")
}

// Icon returns the icon name for the alert variant
func (a *Alert) Icon() string {
	return variantIcons[a.Variant]
}

// IconClasses returns the CSS classes for the alert icon
func (a *Alert) IconClasses() string {
	classes := []string{"gl-alert-icon"}
	if a.Title == "" {
		classes = append(classes, "gl-alert-icon-no-title")
	}
	return strings.Join(classes, " ")
}

// DismissibleButtonOptions returns the options for the dismissible button
func (a *Alert) DismissibleButtonOptions() map[string]interface{} {
	options := make(map[string]interface{})
	for k, v := range a.CloseButtonOptions {
		options[k] = v
	}

	// Add default classes
	if class, ok := options["class"].(string); ok {
		options["class"] = fmt.Sprintf("js-close gl-dismiss-btn %s", class)
	} else {
		options["class"] = "js-close gl-dismiss-btn"
	}

	// Add aria label
	if _, ok := options["aria"]; !ok {
		options["aria"] = make(map[string]interface{})
	}
	aria := options["aria"].(map[string]interface{})
	aria["label"] = "Dismiss"

	return options
}

// Render generates the HTML for the alert
func (a *Alert) Render() string {
	// Build alert attributes
	alertAttrs := make(map[string]interface{})
	for k, v := range a.AlertOptions {
		alertAttrs[k] = v
	}
	alertAttrs["role"] = "alert"
	alertAttrs["class"] = a.BaseClass()

	var parts []string

	// Add icon if needed
	if a.ShowIcon {
		parts = append(parts, fmt.Sprintf(`<div class="gl-alert-icon-container">
			<span class="%s">%s</span>
		</div>`, a.IconClasses(), a.Icon()))
	}

	// Add dismissible button if needed
	if a.Dismissible {
		buttonAttrs := a.DismissibleButtonOptions()
		parts = append(parts, fmt.Sprintf(`<button type="button" %s>
			<span class="js-close-icon">close</span>
		</button>`, formatAttributes(buttonAttrs)))
	}

	// Add content
	contentParts := []string{}

	// Add title if present
	if a.Title != "" {
		contentParts = append(contentParts, fmt.Sprintf(`<h2 class="gl-alert-title">%s</h2>`, a.Title))
	}

	// Add body if present
	if a.Body != "" {
		contentParts = append(contentParts, fmt.Sprintf(`<div class="gl-alert-body">%s</div>`, a.Body))
	}

	// Add actions if present
	if a.Actions != "" {
		contentParts = append(contentParts, fmt.Sprintf(`<div class="gl-alert-actions">%s</div>`, a.Actions))
	}

	parts = append(parts, fmt.Sprintf(`<div class="gl-alert-content" role="alert">%s</div>`,
		strings.Join(contentParts, "\n")))

	return fmt.Sprintf(`<div %s>%s</div>`,
		formatAttributes(alertAttrs),
		strings.Join(parts, "\n"))
}
