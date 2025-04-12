package rapid_diffs

import (
	"fmt"
	"strings"
)

// ExpandDirection represents the available expand directions
type ExpandDirection string

const (
	ExpandDirectionUp   ExpandDirection = "up"
	ExpandDirectionDown ExpandDirection = "down"
	ExpandDirectionBoth ExpandDirection = "both"
)

// IconNames maps expand directions to their corresponding icon names
var IconNames = map[ExpandDirection]string{
	ExpandDirectionUp:   "expand-up",
	ExpandDirectionDown: "expand-down",
	ExpandDirectionBoth: "expand",
}

// ExpandLines represents an expand lines component
type ExpandLines struct {
	Directions []ExpandDirection
}

// NewExpandLines creates a new expand lines component with the given directions
func NewExpandLines(directions []ExpandDirection) *ExpandLines {
	return &ExpandLines{
		Directions: directions,
	}
}

// GetIconName returns the icon name for the given direction
func GetIconName(direction ExpandDirection) string {
	return IconNames[direction]
}

// Render generates the HTML for the expand lines component
func (e *ExpandLines) Render() string {
	var parts []string

	for _, direction := range e.Directions {
		parts = append(parts, fmt.Sprintf(`<button class="rd-expand-lines-button" type="button" data-click="expandLines" data-expand-direction="%s">`, direction))
		parts = append(parts, fmt.Sprintf(`<span data-visible-when="idle"><span class="gl-icon">%s</span></span>`, GetIconName(direction)))
		parts = append(parts, `<span data-visible-when="loading"><span class="gl-icon gl-spinner gl-spinner-sm gl-spinner-dark !gl-align-text-bottom"></span></span>')
		parts = append(parts, "</button>")
	}

	return strings.Join(parts, "\n")
}
