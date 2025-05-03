package core

// MockLogger is a mock implementation of the Logger interface for testing
type MockLogger struct {
	DebugMessages   []string
	InfoMessages    []string
	WarnMessages    []string
	ErrorMessages   []string
	SuccessMessages []string
}

// NewMockLogger creates a new MockLogger
func NewMockLogger() *MockLogger {
	return &MockLogger{
		DebugMessages:   []string{},
		InfoMessages:    []string{},
		WarnMessages:    []string{},
		ErrorMessages:   []string{},
		SuccessMessages: []string{},
	}
}

// Debug logs a debug message
func (l *MockLogger) Debug(message string, args ...interface{}) {
	l.DebugMessages = append(l.DebugMessages, message)
}

// Info logs an info message
func (l *MockLogger) Info(message string, args ...interface{}) {
	l.InfoMessages = append(l.InfoMessages, message)
}

// Warn logs a warning message
func (l *MockLogger) Warn(message string, args ...interface{}) {
	l.WarnMessages = append(l.WarnMessages, message)
}

// Error logs an error message
func (l *MockLogger) Error(message string, args ...interface{}) {
	l.ErrorMessages = append(l.ErrorMessages, message)
}

// Success logs a success message
func (l *MockLogger) Success(message string, args ...interface{}) {
	l.SuccessMessages = append(l.SuccessMessages, message)
}

// Warning logs a warning message
func (l *MockLogger) Warning(message string, args ...interface{}) {
	l.WarnMessages = append(l.WarnMessages, message)
}
