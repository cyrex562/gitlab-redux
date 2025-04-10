package labels

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/gitlab-org/gitlab-redux/internal/service"
)

// LabelHash represents a label with its title, color, and set status
type LabelHash struct {
	Title string `json:"title"`
	Color string `json:"color"`
	Set   bool   `json:"set,omitempty"`
}

// LabelsAsHash handles converting labels to hash format
type LabelsAsHash struct {
	labelsFinder *service.LabelsFinder
	logger       *service.Logger
}

// NewLabelsAsHash creates a new instance of LabelsAsHash
func NewLabelsAsHash(
	labelsFinder *service.LabelsFinder,
	logger *service.Logger,
) *LabelsAsHash {
	return &LabelsAsHash{
		labelsFinder: labelsFinder,
		logger:       logger,
	}
}

// GetLabelsAsHash converts labels to hash format
func (l *LabelsAsHash) GetLabelsAsHash(ctx *gin.Context, target interface{}, params map[string]interface{}) ([]LabelHash, error) {
	// Get current user from context
	currentUser, err := ctx.Get("current_user")
	if err != nil {
		return nil, err
	}

	// Find available labels
	availableLabels, err := l.labelsFinder.Execute(currentUser, params)
	if err != nil {
		return nil, err
	}

	// Convert labels to hash format
	labelHashes := make([]LabelHash, 0, len(availableLabels))
	for _, label := range availableLabels {
		labelHash := LabelHash{
			Title: label.Title,
			Color: label.Color,
		}
		labelHashes = append(labelHashes, labelHash)
	}

	// Check if target has labels and update set status
	if target != nil {
		if targetWithLabels, ok := target.(interface{ GetLabels() []*service.Label }); ok {
			alreadySetLabels := l.findAlreadySetLabels(availableLabels, targetWithLabels.GetLabels())
			if len(alreadySetLabels) > 0 {
				// Update set status for matching labels
				alreadySetTitles := make(map[string]bool)
				for _, label := range alreadySetLabels {
					alreadySetTitles[label.Title] = true
				}

				for i := range labelHashes {
					if alreadySetTitles[labelHashes[i].Title] {
						labelHashes[i].Set = true
					}
				}
			}
		}
	}

	return labelHashes, nil
}

// findAlreadySetLabels finds labels that are already set on the target
func (l *LabelsAsHash) findAlreadySetLabels(availableLabels []*service.Label, targetLabels []*service.Label) []*service.Label {
	alreadySet := make([]*service.Label, 0)
	targetLabelMap := make(map[string]*service.Label)

	// Create a map of target labels for faster lookup
	for _, label := range targetLabels {
		targetLabelMap[label.Title] = label
	}

	// Find matching labels
	for _, label := range availableLabels {
		if _, exists := targetLabelMap[label.Title]; exists {
			alreadySet = append(alreadySet, label)
		}
	}

	return alreadySet
}
