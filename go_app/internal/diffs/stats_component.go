package diffs

import (
	"fmt"
	"html/template"
)

// StatsComponent represents the component for displaying diff statistics
type StatsComponent struct {
	BaseComponent
}

// NewStatsComponent creates a new StatsComponent instance
func NewStatsComponent(opts ComponentOptions) *StatsComponent {
	stats := &DiffStats{
		Additions:    opts.Additions,
		Deletions:    opts.Deletions,
		TotalChanges: opts.Additions + opts.Deletions,
	}

	return &StatsComponent{
		BaseComponent: BaseComponent{
			Stats:     stats,
			IsVisible: opts.IsVisible,
		},
	}
}

// Render renders the stats component
func (c *StatsComponent) Render() (template.HTML, error) {
	if !c.IsVisible {
		return "", nil
	}

	stats := c.Stats
	html := fmt.Sprintf(`
		<div class="diff-stats">
			<span class="diff-stats-additions">+%d</span>
			<span class="diff-stats-deletions">-%d</span>
			<span class="diff-stats-total">%d</span>
		</div>
	`, stats.Additions, stats.Deletions, stats.TotalChanges)

	return template.HTML(html), nil
}

// GetStats returns the diff statistics
func (c *StatsComponent) GetStats() *DiffStats {
	return c.Stats
}
