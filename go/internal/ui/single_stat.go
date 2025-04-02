package ui

import (
	"fmt"
	"strings"
)

// SingleStat represents a single stat component
type SingleStat struct {
	Title           string
	TitleTag        string
	StatValue       string
	StatValueTestID string
	Unit            string
	TitleIcon       string
	MetaText        string
	MetaIcon        string
	TextColor       string
	Variant         string
}

// NewSingleStat creates a new single stat with the given parameters
func NewSingleStat(
	title string,
	titleTag string,
	statValue string,
	statValueTestID string,
	unit string,
	titleIcon string,
	metaText string,
	metaIcon string,
	textColor string,
	variant string,
) *SingleStat {
	if titleTag == "" {
		titleTag = "span"
	}
	if statValueTestID == "" {
		statValueTestID = "non-animated-value"
	}
	if variant == "" {
		variant = "muted"
	}

	return &SingleStat{
		Title:           title,
		TitleTag:        titleTag,
		StatValue:       statValue,
		StatValueTestID: statValueTestID,
		Unit:            unit,
		TitleIcon:       titleIcon,
		MetaText:        metaText,
		MetaIcon:        metaIcon,
		TextColor:       textColor,
		Variant:         variant,
	}
}

// HasUnit returns whether the component has a unit
func (s *SingleStat) HasUnit() bool {
	return s.Unit != ""
}

// HasTitleIcon returns whether the component has a title icon
func (s *SingleStat) HasTitleIcon() bool {
	return s.TitleIcon != ""
}

// HasMetaIcon returns whether the component has a meta icon
func (s *SingleStat) HasMetaIcon() bool {
	return s.MetaIcon != ""
}

// HasMetaText returns whether the component has meta text
func (s *SingleStat) HasMetaText() bool {
	return s.MetaText != ""
}

// GetUnitClass returns the unit class based on whether there is a unit
func (s *SingleStat) GetUnitClass() string {
	if !s.HasUnit() {
		return "gl-mr-2"
	}
	return ""
}

// Render generates the HTML for the single stat
func (s *SingleStat) Render() string {
	var parts []string

	// Start single stat div
	parts = append(parts, `<div class="gl-single-stat gl-flex gl-flex-col gl-py-2">`)

	// Title section
	parts = append(parts, `<div class="gl-flex gl-items-center gl-text-subtle gl-mb-4">`)
	if s.HasTitleIcon() {
		parts = append(parts, fmt.Sprintf(`<span class="gl-icon gl-mr-2 gl-fill-icon-subtle">%s</span>`, s.TitleIcon))
	}
	parts = append(parts, fmt.Sprintf(`<%s class="gl-text-base gl-font-normal gl-text-subtle gl-leading-reset gl-m-0" data-testid="title-text">%s</%s>`,
		s.TitleTag, s.Title, s.TitleTag))
	parts = append(parts, "</div>")

	// Content section
	parts = append(parts, `<div class="gl-single-stat-content gl-flex gl-items-baseline gl-font-bold gl-text-default gl-mb-4">`)
	parts = append(parts, fmt.Sprintf(`<span class="gl-single-stat-number gl-leading-1 %s" data-testid="displayValue">`, s.GetUnitClass()))
	parts = append(parts, fmt.Sprintf(`<span data-testid="%s">%s</span>`, s.StatValueTestID, s.StatValue))
	parts = append(parts, "</span>")

	if s.HasUnit() {
		parts = append(parts, fmt.Sprintf(`<span class="gl-text-sm gl-mx-2 gl-transition-all gl-opacity-10" data-testid="unit">%s</span>`, s.Unit))
	}

	if s.HasMetaIcon() && !s.HasMetaText() {
		parts = append(parts, fmt.Sprintf(`<span class="gl-icon %s">%s</span>`, s.TextColor, s.MetaIcon))
	} else if s.HasMetaText() {
		badge := NewBadge(s.MetaText, s.Variant, s.MetaIcon, map[string]interface{}{
			"data-testid": "meta-badge",
		})
		parts = append(parts, badge.Render())
	}

	parts = append(parts, "</div>")
	parts = append(parts, "</div>")

	return strings.Join(parts, "\n")
}
