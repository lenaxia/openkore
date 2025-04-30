package ui

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

// TestMiscConfig tests the HandleMiscConfig method
func TestMiscConfig(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create UI manager
	uiManager := NewUIManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("misc_config", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test misc config
	args := map[string]interface{}{
		"show_eq_flag": uint8(1),
		"call_flag":    uint8(0),
	}

	err := uiManager.HandleMiscConfig(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info messages were created
	if len(mockLogger.infoMessages) != 2 {
		t.Errorf("Expected 2 info messages, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected misc_config hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"show_eq_flag": "invalid", // Invalid type
		"call_flag":    uint8(0),
	}

	err = uiManager.HandleMiscConfig(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestOpenUI tests the HandleOpenUI method
func TestOpenUI(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create UI manager
	uiManager := NewUIManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("open_ui", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test open UI (bank)
	args := map[string]interface{}{
		"type": uint8(BANK_UI),
	}

	err := uiManager.HandleOpenUI(args)

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
		t.Errorf("Expected open_ui hook to be called")
	}

	// Test open UI (unknown)
	hookCalled = false
	mockLogger = NewMockLogger()
	uiManager.logger = mockLogger

	args = map[string]interface{}{
		"type": uint8(100), // Unknown UI type
	}

	err = uiManager.HandleOpenUI(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify error message was created
	if len(mockLogger.errorMessages) != 1 {
		t.Errorf("Expected 1 error message, got %d", len(mockLogger.errorMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected open_ui hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"type": "invalid", // Invalid type
	}

	err = uiManager.HandleOpenUI(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestProgressBar tests the HandleProgressBar method
func TestProgressBar(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create UI manager
	uiManager := NewUIManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("progress_bar", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test progress bar
	args := map[string]interface{}{
		"time": uint32(10),
	}

	err := uiManager.HandleProgressBar(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify progress bar active
	if !uiManager.progressBarActive {
		t.Errorf("Expected progress bar to be active")
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected progress_bar hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"time": "invalid", // Invalid type
	}

	err = uiManager.HandleProgressBar(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestProgressBarStop tests the HandleProgressBarStop method
func TestProgressBarStop(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create UI manager
	uiManager := NewUIManager(mockParser, hookManager, mockLogger)
	uiManager.progressBarActive = true

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("progress_bar_stop", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test progress bar stop
	args := map[string]interface{}{}

	err := uiManager.HandleProgressBarStop(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify progress bar inactive
	if uiManager.progressBarActive {
		t.Errorf("Expected progress bar to be inactive")
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected progress_bar_stop hook to be called")
	}
}

// TestShowScript tests the HandleShowScript method
func TestShowScript(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create UI manager
	uiManager := NewUIManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("show_script", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test show script
	args := map[string]interface{}{
		"ID":      uint32(12345),
		"message": "Hello, world!",
	}

	err := uiManager.HandleShowScript(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify debug message was created
	if len(mockLogger.debugMessages) != 1 {
		t.Errorf("Expected 1 debug message, got %d", len(mockLogger.debugMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected show_script hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"ID":      "invalid", // Invalid type
		"message": "Hello, world!",
	}

	err = uiManager.HandleShowScript(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestInventoryExpansionResult tests the HandleInventoryExpansionResult method
func TestInventoryExpansionResult(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create UI manager
	uiManager := NewUIManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("inventory_expansion_result", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test inventory expansion result (success)
	args := map[string]interface{}{
		"result": uint16(EXPAND_INVENTORY_RESULT_SUCCESS),
	}

	err := uiManager.HandleInventoryExpansionResult(args)

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
		t.Errorf("Expected inventory_expansion_result hook to be called")
	}

	// Test inventory expansion result (unknown)
	hookCalled = false
	mockLogger = NewMockLogger()
	uiManager.logger = mockLogger

	args = map[string]interface{}{
		"result": uint16(100), // Unknown result
	}

	err = uiManager.HandleInventoryExpansionResult(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify warning message was created
	if len(mockLogger.warningMessages) != 1 {
		t.Errorf("Expected 1 warning message, got %d", len(mockLogger.warningMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected inventory_expansion_result hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"result": "invalid", // Invalid type
	}

	err = uiManager.HandleInventoryExpansionResult(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}
