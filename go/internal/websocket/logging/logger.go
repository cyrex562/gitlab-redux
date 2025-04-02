package logging

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

var (
	defaultLogger Logger
	once         sync.Once
)

// GetLogger returns the default logger instance
func GetLogger() Logger {
	once.Do(func() {
		defaultLogger = NewDefaultLogger()
	})
	return defaultLogger
}

// Logger defines the interface for logging operations
type Logger interface {
	Info(payload LogPayload, message string)
	Error(payload LogPayload, message string, err error)
	Debug(payload LogPayload, message string)
}

// DefaultLogger implements the Logger interface using the standard log package
type DefaultLogger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
}

// NewDefaultLogger creates a new DefaultLogger instance
func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{
		infoLogger:  log.New(os.Stdout, "INFO: ", log.LstdFlags),
		errorLogger: log.New(os.Stderr, "ERROR: ", log.LstdFlags),
		debugLogger: log.New(os.Stdout, "DEBUG: ", log.LstdFlags),
	}
}

// Info logs an informational message with the given payload
func (l *DefaultLogger) Info(payload LogPayload, message string) {
	l.log(l.infoLogger, payload, message)
}

// Error logs an error message with the given payload and error
func (l *DefaultLogger) Error(payload LogPayload, message string, err error) {
	if err != nil {
		payload.Params["error"] = err.Error()
	}
	l.log(l.errorLogger, payload, message)
}

// Debug logs a debug message with the given payload
func (l *DefaultLogger) Debug(payload LogPayload, message string) {
	l.log(l.debugLogger, payload, message)
}

// log is a helper method to format and write log messages
func (l *DefaultLogger) log(logger *log.Logger, payload LogPayload, message string) {
	// Convert payload to JSON
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		logger.Printf("Error marshaling payload: %v", err)
		return
	}

	// Log the message with the payload
	logger.Printf("%s %s", message, string(payloadJSON))
}
