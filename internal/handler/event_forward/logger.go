package event_forward

import (
	"github.com/jmadden/gitlab-redux/internal/service"
)

// Logger is a specialized logger for event collection
type Logger struct {
	service.Logger
}

// NewLogger creates a new EventForward Logger
func NewLogger(baseLogger service.Logger) *Logger {
	return &Logger{
		Logger: baseLogger,
	}
}

// FileNameNoExt returns the base filename without extension for the logger
func (l *Logger) FileNameNoExt() string {
	return "event_collection"
}
