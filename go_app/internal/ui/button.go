package ui

import (
	"fmt"
	"strings"
)

// ButtonCategory represents the possible categories of a button
type ButtonCategory string

const (
	// ButtonCategoryPrimary represents a primary button
	ButtonCategoryPrimary ButtonCategory = "primary"
	// ButtonCategorySecondary represents a secondary button
	ButtonCategorySecondary ButtonCategory = "secondary"
	// ButtonCategoryTertiary represents a tertiary button
	ButtonCategoryTertiary ButtonCategory = "tertiary"
)

// ButtonVariant represents the possible variants of a button
type ButtonVariant string

const (
	// ButtonVariantDefault represents a default button
	ButtonVariantDefault ButtonVariant = "default"
	// ButtonVariantConfirm represents a confirm button
	ButtonVariantConfirm ButtonVariant = "confirm"
	// ButtonVariantDanger represents a danger button
	ButtonVariantDanger ButtonVariant = "danger"
	// ButtonVariantDashed represents a dashed button
	ButtonVariantDashed ButtonVariant = "dashed"
	// ButtonVariantLink represents a link button
	ButtonVariantLink ButtonVariant = "link"
	// ButtonVariantReset represents a reset button
	ButtonVariantReset ButtonVariant = "reset"
)

// ButtonSize represents the possible sizes of a button
type ButtonSize string

const (
	// ButtonSizeSmall represents a small button
	ButtonSizeSmall ButtonSize = "small"
	// ButtonSizeMedium represents a medium button
	ButtonSizeMedium ButtonSize = "medium"
)

// ButtonType represents the possible types of a button
type ButtonType string

const (
	// ButtonTypeButton represents a button type
	ButtonTypeButton ButtonType = "button"
	// ButtonTypeReset represents a reset type
	ButtonTypeReset ButtonType = "reset"
	// ButtonTypeSubmit represents a submit type
	ButtonTypeSubmit ButtonType = "submit"
)

// ButtonTarget represents the possible targets of a button
type ButtonTarget string

const (
	// ButtonTargetSelf represents a self target
	ButtonTargetSelf ButtonTarget = "_self"
	// ButtonTargetBlank represents a blank target
	ButtonTargetBlank ButtonTarget = "_blank"
	// ButtonTargetParent represents a parent target
	ButtonTargetParent ButtonTarget = "_parent"
	// ButtonTargetTop represents a top target
	ButtonTargetTop ButtonTarget = "_top"
)

// ButtonMethod represents the possible methods of a button
type ButtonMethod string

const (
	// ButtonMethodGet represents a get method
	ButtonMethodGet ButtonMethod = "get"
	// ButtonMethodPost represents a post method
	ButtonMethodPost ButtonMethod = "post"
	// ButtonMethodPut represents a put method
	ButtonMethodPut ButtonMethod = "put"
	// ButtonMethodDelete represents a delete method
	ButtonMethodDelete ButtonMethod = "delete"
	// ButtonMethodPatch represents a patch method
	ButtonMethodPatch ButtonMethod = "patch"
)

// Button represents a button component
type Button struct {
	Category          ButtonCategory
	Variant           ButtonVariant
	Size              ButtonSize
	Type              ButtonType
	Disabled          bool
	Loading           bool
	Block             bool
	Label             bool
	Selected          bool
	Icon              string
	Href              string
	Form              bool
	Target            ButtonTarget
	Method            ButtonMethod
	ButtonOptions     map[string]interface{}
	ButtonTextClasses string
	IconClasses       string
	IconContent       string
	Content           string
}

// NewButton creates a new button with the given parameters
func NewButton(category ButtonCategory, variant ButtonVariant, size ButtonSize, buttonType ButtonType, disabled bool, loading bool, block bool, label bool, selected bool, icon string, href string, form bool, target ButtonTarget, method ButtonMethod, buttonOptions map[string]interface{}, buttonTextClasses string, iconClasses string, iconContent string, content string) *Button {
	if buttonOptions == nil {
		buttonOptions = make(map[string]interface{})
	}

	return &Button{
		Category:          category,
		Variant:           variant,
		Size:              size,
		Type:              buttonType,
		Disabled:          disabled,
		Loading:           loading,
		Block:             block,
		Label:             label,
		Selected:          selected,
		Icon:              icon,
		Href:              href,
		Form:              form,
		Target:            target,
		Method:            method,
		ButtonOptions:     buttonOptions,
		ButtonTextClasses: buttonTextClasses,
		IconClasses:       iconClasses,
		IconContent:       iconContent,
		Content:           content,
	}
}

// ButtonClass returns the CSS classes for the button
func (b *Button) ButtonClass() string {
	var classes []string

	// Add base classes
	classes = append(classes, "gl-button", "btn")

	// Add state classes
	if b.Disabled || b.Loading {
		classes = append(classes, "disabled")
	}
	if b.Selected {
		classes = append(classes, "selected")
	}
	if b.Block {
		classes = append(classes, "btn-block")
	}
	if b.Label {
		classes = append(classes, "btn-label")
	}
	if b.Icon != "" && b.Content == "" {
		classes = append(classes, "btn-icon")
	}

	// Add size class
	switch b.Size {
	case ButtonSizeSmall:
		classes = append(classes, "btn-sm")
	case ButtonSizeMedium:
		classes = append(classes, "btn-md")
	}

	// Add variant class
	switch b.Variant {
	case ButtonVariantDefault:
		classes = append(classes, "btn-default")
	case ButtonVariantConfirm:
		classes = append(classes, "btn-confirm")
	case ButtonVariantDanger:
		classes = append(classes, "btn-danger")
	case ButtonVariantDashed:
		classes = append(classes, "btn-dashed")
	case ButtonVariantLink:
		classes = append(classes, "btn-link")
	case ButtonVariantReset:
		classes = append(classes, "btn-gl-reset")
	}

	// Add category class for non-special variants
	nonCategoryVariants := []ButtonVariant{ButtonVariantDashed, ButtonVariantLink, ButtonVariantReset}
	isNonCategoryVariant := false
	for _, variant := range nonCategoryVariants {
		if b.Variant == variant {
			isNonCategoryVariant = true
			break
		}
	}

	if !isNonCategoryVariant && b.Category != ButtonCategoryPrimary {
		switch b.Category {
		case ButtonCategorySecondary:
			classes = append(classes, "btn-default-secondary")
		case ButtonCategoryTertiary:
			classes = append(classes, "btn-default-tertiary")
		}
	}

	// Add custom classes from button options
	if class, ok := b.ButtonOptions["class"].(string); ok && class != "" {
		classes = append(classes, class)
	}

	return strings.Join(classes, " ")
}

// BaseAttributes returns the base attributes for the button
func (b *Button) BaseAttributes() map[string]interface{} {
	attributes := make(map[string]interface{})

	// Add disabled attribute
	if b.Disabled || b.Loading {
		attributes["disabled"] = "disabled"
		attributes["aria-disabled"] = true
	}

	// Add type attribute for non-link buttons
	if b.Href == "" {
		attributes["type"] = string(b.Type)
	}

	// Add rel attribute for target="_blank" links
	if b.Link() && b.Target == ButtonTargetBlank {
		rel := "noopener noreferrer"
		if relValue, ok := b.ButtonOptions["rel"].(string); ok && relValue != "" {
			rel = fmt.Sprintf("%s %s", relValue, rel)
		}
		attributes["rel"] = rel
	}

	return attributes
}

// Link returns whether the button should be rendered as a link
func (b *Button) Link() bool {
	return b.Href != ""
}

// Form returns whether the button should be rendered as a form
func (b *Button) Form() bool {
	return b.Href != "" && b.Form
}

// ShowIcon returns whether the icon should be shown
func (b *Button) ShowIcon() bool {
	return !b.Loading || b.Content != ""
}

// Render generates the HTML for the button
func (b *Button) Render() string {
	// Build button content
	var contentParts []string

	// Add loading indicator if needed
	if b.Loading {
		contentParts = append(contentParts, `<span class="gl-button-icon gl-button-loading-indicator">Loading...</span>`)
	}

	// Add icon if needed
	if b.Icon != "" && b.ShowIcon() {
		iconClasses := "gl-icon gl-button-icon"
		if b.IconClasses != "" {
			iconClasses = fmt.Sprintf("%s %s", iconClasses, b.IconClasses)
		}
		contentParts = append(contentParts, fmt.Sprintf(`<span class="%s">%s</span>`, iconClasses, b.Icon))
	} else if b.IconContent != "" && b.ShowIcon() {
		contentParts = append(contentParts, b.IconContent)
	}

	// Add text content if needed
	if b.Content != "" {
		textClasses := "gl-button-text"
		if b.ButtonTextClasses != "" {
			textClasses = fmt.Sprintf("%s %s", textClasses, b.ButtonTextClasses)
		}
		contentParts = append(contentParts, fmt.Sprintf(`<span class="%s">%s</span>`, textClasses, b.Content))
	}

	// Determine the tag to use
	tag := "button"
	if b.Label {
		tag = "span"
	}

	// Build button attributes
	buttonAttrs := make(map[string]interface{})
	for k, v := range b.ButtonOptions {
		buttonAttrs[k] = v
	}

	// Add base attributes
	for k, v := range b.BaseAttributes() {
		buttonAttrs[k] = v
	}

	// Add class attribute
	buttonAttrs["class"] = b.ButtonClass()

	// Render the button based on its type
	if b.Form() {
		// Form button
		buttonAttrs["target"] = string(b.Target)
		buttonAttrs["method"] = string(b.Method)
		return fmt.Sprintf(`<form action="%s" method="%s"><button type="submit" %s>%s</button></form>`,
			b.Href,
			strings.ToUpper(string(b.Method)),
			formatAttributes(buttonAttrs),
			strings.Join(contentParts, ""))
	} else if b.Link() {
		// Link button
		buttonAttrs["href"] = b.Href
		buttonAttrs["target"] = string(b.Target)
		if b.Method != "" {
			buttonAttrs["data-method"] = string(b.Method)
		}
		return fmt.Sprintf(`<a %s>%s</a>`,
			formatAttributes(buttonAttrs),
			strings.Join(contentParts, ""))
	} else {
		// Regular button
		return fmt.Sprintf(`<%s %s>%s</%s>`,
			tag,
			formatAttributes(buttonAttrs),
			strings.Join(contentParts, ""),
			tag)
	}
}
