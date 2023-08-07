package model

type LogEngine uint8

// Fields is a map of strings to any type. It is used to pass to Logger.WithFields.
type Fields map[string]interface{}

// A Logger provides methods for logging messages.
type Logger interface {

	// Debug emits a "DEBUG" level log message.
	Debug(msg string, args ...any)

	// Info emits an "INFO" level log message.
	Info(msg string, args ...any)

	Infof(format string, args ...interface{})
	// Warn emits a "WARN" level log message.
	Warn(msg string, args ...any)

	// Error emits an "ERROR" level log message.
	Error(msg string, args ...any)

	// Fatal emits a "FATAL" level log message.
	Fatal(msg string, args ...any)

	// WithField adds a field to the logger and returns a new Logger.
	WithField(key string, value any) Logger

	WithFields(field Fields) Logger
}
