package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger wraps logrus.Logger with additional methods
type Logger struct {
	*logrus.Logger
}

// New creates a new logger instance
func New(level string) *Logger {
	log := logrus.New()

	// Set log level
	switch level {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}

	// Set JSON formatter for production
	if os.Getenv("ENVIRONMENT") == "production" {
		log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	// Set output to stdout
	log.SetOutput(os.Stdout)

	return &Logger{log}
}

// WithField creates a new logger entry with a field
func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.Logger.WithField(key, value)
}

// WithFields creates a new logger entry with multiple fields
func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.Logger.WithFields(fields)
}
