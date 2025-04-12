package rapid_diffs

import (
	"fmt"
)

// StreamingErrorComponent represents a streaming error component
type StreamingErrorComponent struct {
	Message string
}

// NewStreamingErrorComponent creates a new streaming error component
func NewStreamingErrorComponent(message string) *StreamingErrorComponent {
	return &StreamingErrorComponent{
		Message: message,
	}
}

// Render generates the HTML for the streaming error component
func (s *StreamingErrorComponent) Render() string {
	return fmt.Sprintf(`<streaming-error message="%s"></streaming-error>`, s.Message)
}
