package mvp

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// MockLogger is a simple mock implementation of the Logger interface
type MockLogger struct {
	debugMessages   []string
	infoMessages    []string
	warningMessages []string
	errorMessages   []string
	successMessages []string
}

func NewMockLogger() *MockLogger {
	return &MockLogger{
		debugMessages:   []string{},
		infoMessages:    []string{},
		warningMessages: []string{},
		errorMessages:   []string{},
		successMessages: []string{},
	}
}

func (m *MockLogger) Debug(format string, args ...interface{}) {
	m.debugMessages = append(m.debugMessages, format)
}

func (m *MockLogger) Info(format string, args ...interface{}) {
	m.infoMessages = append(m.infoMessages, format)
}

func (m *MockLogger) Warning(format string, args ...interface{}) {
	m.warningMessages = append(m.warningMessages, format)
}

func (m *MockLogger) Error(format string, args ...interface{}) {
	m.errorMessages = append(m.errorMessages, format)
}

func (m *MockLogger) Success(format string, args ...interface{}) {
	m.successMessages = append(m.successMessages, format)
}

// TestMVPItem tests the HandleMVPItem method
func TestMVPItem(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create MVP manager
	mvpManager := NewMVPManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("mvp_item", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test MVP item
	args := map[string]interface{}{
		"itemID": uint16(501),
	}

	err := mvpManager.HandleMVPItem(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected mvp_item hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"itemID": "invalid", // Invalid type
	}

	err = mvpManager.HandleMVPItem(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestMVPOther tests the HandleMVPOther method
func TestMVPOther(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create MVP manager
	mvpManager := NewMVPManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("mvp_other", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test MVP other
	args := map[string]interface{}{
		"ID": uint32(12345),
	}

	err := mvpManager.HandleMVPOther(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected mvp_other hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"ID": "invalid", // Invalid type
	}

	err = mvpManager.HandleMVPOther(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestMVPYou tests the HandleMVPYou method
func TestMVPYou(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create MVP manager
	mvpManager := NewMVPManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("mvp_you", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test MVP you
	args := map[string]interface{}{
		"expAmount": uint32(1000),
	}

	err := mvpManager.HandleMVPYou(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected mvp_you hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"expAmount": "invalid", // Invalid type
	}

	err = mvpManager.HandleMVPYou(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}
