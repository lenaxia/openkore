package core

// Logger defines the interface for logging functionality
type Logger interface {
	// Debug logs a debug message
	Debug(format string, args ...interface{})

	// Info logs an informational message
	Info(format string, args ...interface{})

	// Warning logs a warning message
	Warning(format string, args ...interface{})

	// Error logs an error message
	Error(format string, args ...interface{})

	// Success logs a success message
	Success(format string, args ...interface{})
}
