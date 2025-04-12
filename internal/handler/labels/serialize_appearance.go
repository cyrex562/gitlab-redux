package labels

import (
	"gitlab.com/gitlab-org/gitlab-redux/internal/model"
)

// SerializeAppearance converts a slice of labels to their appearance representation
func SerializeAppearance(labels []*model.Label) []*model.LabelAppearance {
	result := make([]*model.LabelAppearance, 0, len(labels))

	for _, label := range labels {
		appearance := &model.LabelAppearance{
			ID:          label.ID,
			Title:       label.Title,
			Color:       label.Color,
			Description: label.Description,
			Priority:    label.Priority,
			Type:        label.Type,
			TextColor:   calculateTextColor(label.Color),
		}
		result = append(result, appearance)
	}

	return result
}

// calculateTextColor determines the appropriate text color (black or white) based on the background color
func calculateTextColor(backgroundColor string) string {
	// This is a simplified implementation
	// In a real application, you would parse the hex color and calculate luminance
	// For now, we'll return a default value
	return "#FFFFFF" // White text by default
}
