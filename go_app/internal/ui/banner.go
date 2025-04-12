package ui

import (
	"fmt"
	"strings"
)

// BannerVariant represents the possible variants of a banner
type BannerVariant string

const (
	// BannerVariantIntroduction represents an introduction banner
	BannerVariantIntroduction BannerVariant = "introduction"
	// BannerVariantPromotion represents a promotion banner
	BannerVariantPromotion BannerVariant = "promotion"
)

// Banner represents a banner component
type Banner struct {
	ButtonText     string
	ButtonLink     string
	Variant        BannerVariant
	SVGPath        string
	BannerOptions  map[string]interface{}
	ButtonOptions  map[string]interface{}
	CloseOptions   map[string]interface{}
	Title          string
	Illustration   string
	PrimaryAction  string
	Actions        []string
	Content        string
}

// NewBanner creates a new banner with the given parameters
func NewBanner(buttonText string, buttonLink string, variant BannerVariant, svgPath string, bannerOptions map[string]interface{}, buttonOptions map[string]interface{}, closeOptions map[string]interface{}, title string, illustration string, primaryAction string, actions []string, content string) *Banner {
	if bannerOptions == nil {
		bannerOptions = make(map[string]interface{})
	}
	if buttonOptions == nil {
		buttonOptions = make(map[string]interface{})
	}
	if closeOptions == nil {
		closeOptions = make(map[string]interface{})
	}
	if actions == nil {
		actions = []string{}
	}

	// Add default close button classes
	if _, ok := closeOptions["class"]; !ok {
		closeOptions["class"] = "js-close gl-banner-close"
	} else {
		closeOptions["class"] = fmt.Sprintf("js-close gl-banner-close %s", closeOptions["class"])
	}

	return &Banner{
		ButtonText:     buttonText,
		ButtonLink:     buttonLink,
		Variant:        variant,
		SVGPath:        svgPath,
		BannerOptions:  bannerOptions,
		ButtonOptions:  buttonOptions,
		CloseOptions:   closeOptions,
		Title:          title,
		Illustration:   illustration,
		PrimaryAction:  primaryAction,
		Actions:        actions,
		Content:        content,
	}
}

// BannerClass returns the CSS classes for the banner
func (b *Banner) BannerClass() string {
	var classes []string

	// Add variant class if needed
	if b.Introduction() {
		classes = append(classes, "gl-banner-introduction")
	}

	return strings.Join(classes, " ")
}

// CloseButtonVariant returns the variant for the close button
func (b *Banner) CloseButtonVariant() string {
	if b.Introduction() {
		return "confirm"
	}
	return "default"
}

// Introduction returns whether the banner is an introduction banner
func (b *Banner) Introduction() bool {
	return b.Variant == BannerVariantIntroduction
}

// HasIllustration returns whether the banner has an illustration
func (b *Banner) HasIllustration() bool {
	return b.Illustration != ""
}

// HasSVGPath returns whether the banner has an SVG path
func (b *Banner) HasSVGPath() bool {
	return b.SVGPath != ""
}

// HasPrimaryAction returns whether the banner has a primary action
func (b *Banner) HasPrimaryAction() bool {
	return b.PrimaryAction != ""
}

// Render generates the HTML for the banner
func (b *Banner) Render() string {
	// Build banner attributes
	bannerAttrs := make(map[string]interface{})
	for k, v := range b.BannerOptions {
		bannerAttrs[k] = v
	}

	// Add base classes
	bannerClass := "gl-banner gl-card gl-pl-6 gl-pr-8 gl-py-6"
	if bannerClass := b.BannerClass(); bannerClass != "" {
		bannerClass = fmt.Sprintf("%s %s", bannerClass, bannerClass)
	}
	bannerAttrs["class"] = bannerClass

	// Build the content
	var contentParts []string

	// Add illustration if needed
	if b.HasIllustration() {
		contentParts = append(contentParts, fmt.Sprintf(`<div class="gl-banner-illustration">%s</div>`, b.Illustration))
	} else if b.HasSVGPath() {
		contentParts = append(contentParts, fmt.Sprintf(`<div class="gl-banner-illustration"><img src="%s" alt=""></div>`, b.SVGPath))
	}

	// Add content
	contentParts = append(contentParts, `<div class="gl-banner-content">`)

	// Add title if needed
	if b.Title != "" {
		contentParts = append(contentParts, fmt.Sprintf(`<h2 class="gl-banner-title">%s</h2>`, b.Title))
	}

	// Add main content
	contentParts = append(contentParts, b.Content)

	// Add primary action or default button
	if b.HasPrimaryAction() {
		contentParts = append(contentParts, b.PrimaryAction)
	} else {
		buttonAttrs := make(map[string]interface{})
		for k, v := range b.ButtonOptions {
			buttonAttrs[k] = v
		}
		buttonAttrs["class"] = "js-close-callout"
		buttonAttrs["variant"] = "confirm"
		contentParts = append(contentParts, fmt.Sprintf(`<a href="%s" %s>%s</a>`, b.ButtonLink, formatAttributes(buttonAttrs), b.ButtonText))
	}

	// Add actions
	for _, action := range b.Actions {
		contentParts = append(contentParts, action)
	}

	contentParts = append(contentParts, `</div>`)

	// Add close button
	closeAttrs := make(map[string]interface{})
	for k, v := range b.CloseOptions {
		closeAttrs[k] = v
	}
	closeAttrs["category"] = "tertiary"
	closeAttrs["variant"] = b.CloseButtonVariant()
	closeAttrs["size"] = "small"
	closeAttrs["icon"] = "close"
	contentParts = append(contentParts, fmt.Sprintf(`<button type="button" %s><span class="js-close-icon">close</span></button>`, formatAttributes(closeAttrs)))

	// Render the banner
	return fmt.Sprintf(`<div %s><div class="gl-card-body gl-flex !gl-p-0">%s</div></div>`,
		formatAttributes(bannerAttrs),
		strings.Join(contentParts, "\n"))
}
