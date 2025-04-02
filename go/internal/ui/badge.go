package ui

import (
	"fmt"
	"strings"
)

// BadgeVariant represents the possible variants of a badge
type BadgeVariant string

const (
	// BadgeVariantMuted represents a muted badge
	BadgeVariantMuted BadgeVariant = "muted"
	// BadgeVariantNeutral represents a neutral badge
	BadgeVariantNeutral BadgeVariant = "neutral"
	// BadgeVariantInfo represents an info badge
	BadgeVariantInfo BadgeVariant = "info"
	// BadgeVariantSuccess represents a success badge
	BadgeVariantSuccess BadgeVariant = "success"
	// BadgeVariantWarning represents a warning badge
	BadgeVariantWarning BadgeVariant = "warning"
	// BadgeVariantDanger represents a danger badge
	BadgeVariantDanger BadgeVariant = "danger"
	// BadgeVariantTier represents a tier badge
	BadgeVariantTier BadgeVariant = "tier"
)

// Badge represents a badge component
type Badge struct {
	Text         string
	Icon         string
	IconClasses  []string
	IconOnly     bool
	Href         string
	Variant      BadgeVariant
	HTMLOptions  map[string]interface{}
	Content      string
}

// NewBadge creates a new badge with the given parameters
func NewBadge(text string, icon string, iconClasses []string, iconOnly bool, href string, variant BadgeVariant, htmlOptions map[string]interface{}, content string) *Badge {
	if iconClasses == nil {
		iconClasses = []string{}
	}
	if htmlOptions == nil {
		htmlOptions = make(map[string]interface{})
	}

	return &Badge{
		Text:         text,
		Icon:         icon,
		IconClasses:  iconClasses,
		IconOnly:     iconOnly,
		Href:         href,
		Variant:      variant,
		HTMLOptions:  htmlOptions,
		Content:      content,
	}
}

// BadgeClasses returns the CSS classes for the badge
func (b *Badge) BadgeClasses() string {
	var classes []string

	// Add base classes
	classes = append(classes, "gl-badge", "badge", "badge-pill", fmt.Sprintf("badge-%s", b.Variant))

	// Add icon-only class if needed
	if b.IconOnly {
		classes = append(classes, "!gl-px-2")
	}

	return strings.Join(classes, " ")
}

// IconClasses returns the CSS classes for the badge icon
func (b *Badge) IconClasses() string {
	var classes []string

	// Add base classes
	classes = append(classes, "gl-icon", "gl-badge-icon")
	classes = append(classes, b.IconClasses...)

	// Add margin class if needed
	if b.CircularIcon() {
		classes = append(classes, "-gl-ml-2")
	}

	return strings.Join(classes, " ")
}

// IconOnly returns whether the badge should only show the icon
func (b *Badge) IconOnly() bool {
	return b.IconOnly
}

// Link returns whether the badge should be rendered as a link
func (b *Badge) Link() bool {
	return b.Href != ""
}

// Text returns the text content of the badge
func (b *Badge) Text() string {
	if b.Content != "" {
		return b.Content
	}
	return b.Text
}

// HasIcon returns whether the badge has an icon
func (b *Badge) HasIcon() bool {
	return b.IconOnly() || b.Icon != ""
}

// CircularIcon returns whether the icon should be circular
func (b *Badge) CircularIcon() bool {
	return b.Icon == "issue-open-m" || b.Icon == "issue-close"
}

// HTMLOptions returns the HTML options for the badge
func (b *Badge) HTMLOptions() map[string]interface{} {
	options := make(map[string]interface{})
	for k, v := range b.HTMLOptions {
		options[k] = v
	}

	// Add CSS classes
	options["class"] = b.BadgeClasses()

	// Add aria attributes for icon-only badges
	if b.IconOnly() {
		if _, ok := options["aria"]; !ok {
			options["aria"] = make(map[string]interface{})
		}
		aria := options["aria"].(map[string]interface{})
		aria["label"] = b.Text()
		options["role"] = "img"
	}

	return options
}

// Render generates the HTML for the badge
func (b *Badge) Render() string {
	// Build HTML options
	htmlOptions := b.HTMLOptions()

	// Determine the tag to use
	tag := "span"
	if b.Link() {
		tag = "a"
		htmlOptions["href"] = b.Href
	}

	// Build the content
	var contentParts []string

	// Add icon if needed
	if b.HasIcon() {
		contentParts = append(contentParts, fmt.Sprintf(`<span class="%s">%s</span>`, b.IconClasses(), b.Icon))
	}

	// Add text if needed
	if text := b.Text(); text != "" {
		contentParts = append(contentParts, fmt.Sprintf(`<span class="gl-badge-content">%s</span>`, text))
	}

	// Render the badge
	return fmt.Sprintf(`<%s %s>%s</%s>`,
		tag,
		formatAttributes(htmlOptions),
		strings.Join(contentParts, ""),
		tag)
}
