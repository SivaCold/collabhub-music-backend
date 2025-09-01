package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger is the global logger instance
var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
	
	// Set output to stdout
	Logger.SetOutput(os.Stdout)
	
	// Set default log level
	Logger.SetLevel(logrus.InfoLevel)
	
	// Use JSON formatter for structured logging
	Logger.SetFormatter(&logrus.JSONFormatter{})
}

// NewLogger creates a new logger instance with custom configuration
func NewLogger() *logrus.Logger {
	logger := logrus.New()
	
	// Configure based on environment
	env := os.Getenv("SERVER_ENV")
	if env == "production" {
		logger.SetLevel(logrus.WarnLevel)
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logger.SetLevel(logrus.DebugLevel)
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
		})
	}
	
	return logger
}

// Info logs an info message
func Info(args ...interface{}) {
	Logger.Info(args...)
}

// Debug logs a debug message
func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

// Warn logs a warning message
func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

// Error logs an error message
func Error(args ...interface{}) {
	Logger.Error(args...)
}

// Fatal logs a fatal message and exits
func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

// WithField creates a new entry with a single field
func WithField(key string, value interface{}) *logrus.Entry {
	return Logger.WithField(key, value)
}

// WithFields creates a new entry with multiple fields
func WithFields(fields logrus.Fields) *logrus.Entry {
	return Logger.WithFields(fields)
}
