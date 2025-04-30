package battle

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

// TestBattlegroundMessage tests the HandleBattlegroundMessage method
func TestBattlegroundMessage(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create battle manager
	battleManager := NewBattleManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("battleground_message", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test battleground message
	args := map[string]interface{}{
		"message": "Test battleground message",
	}

	err := battleManager.HandleBattlegroundMessage(args)

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
		t.Errorf("Expected battleground_message hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"message": 123, // Invalid type
	}

	err = battleManager.HandleBattlegroundMessage(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestBattlegroundEmblem tests the HandleBattlegroundEmblem method
func TestBattlegroundEmblem(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create battle manager
	battleManager := NewBattleManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("battleground_emblem", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test battleground emblem
	args := map[string]interface{}{
		"emblemID": uint32(123),
	}

	err := battleManager.HandleBattlegroundEmblem(args)

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
		t.Errorf("Expected battleground_emblem hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"emblemID": "invalid", // Invalid type
	}

	err = battleManager.HandleBattlegroundEmblem(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestInstanceWindowStart tests the HandleInstanceWindowStart method
func TestInstanceWindowStart(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create battle manager
	battleManager := NewBattleManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("instance_window_start", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test instance window start
	args := map[string]interface{}{
		"instanceID": uint32(123),
		"name":       "Test Instance",
		"time":       uint32(3600),
		"progress":   uint32(50),
	}

	err := battleManager.HandleInstanceWindowStart(args)

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
		t.Errorf("Expected instance_window_start hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"instanceID": "invalid", // Invalid type
		"name":       "Test Instance",
		"time":       uint32(3600),
		"progress":   uint32(50),
	}

	err = battleManager.HandleInstanceWindowStart(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestInstanceWindowQueue tests the HandleInstanceWindowQueue method
func TestInstanceWindowQueue(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create battle manager
	battleManager := NewBattleManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("instance_window_queue", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test instance window queue
	args := map[string]interface{}{
		"instanceID": uint32(123),
	}

	err := battleManager.HandleInstanceWindowQueue(args)

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
		t.Errorf("Expected instance_window_queue hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"instanceID": "invalid", // Invalid type
	}

	err = battleManager.HandleInstanceWindowQueue(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestInstanceWindowJoin tests the HandleInstanceWindowJoin method
func TestInstanceWindowJoin(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create battle manager
	battleManager := NewBattleManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("instance_ready", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test instance window join
	args := map[string]interface{}{
		"instanceID": uint32(123),
		"result":     uint32(1),
	}

	err := battleManager.HandleInstanceWindowJoin(args)

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
		t.Errorf("Expected instance_ready hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"instanceID": "invalid", // Invalid type
		"result":     uint32(1),
	}

	err = battleManager.HandleInstanceWindowJoin(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestInstanceWindowLeave tests the HandleInstanceWindowLeave method
func TestInstanceWindowLeave(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create battle manager
	battleManager := NewBattleManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("instance_window_leave", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test instance window leave (notify)
	args := map[string]interface{}{
		"instanceID": uint32(123),
		"flag":       uint8(0),
	}

	err := battleManager.HandleInstanceWindowLeave(args)

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
		t.Errorf("Expected instance_window_leave hook to be called")
	}

	// Test instance window leave (expired)
	hookCalled = false
	mockLogger = NewMockLogger()
	battleManager.logger = mockLogger

	args = map[string]interface{}{
		"instanceID": uint32(123),
		"flag":       uint8(1),
	}

	err = battleManager.HandleInstanceWindowLeave(args)

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
		t.Errorf("Expected instance_window_leave hook to be called")
	}

	// Test instance window leave (unknown flag)
	hookCalled = false
	mockLogger = NewMockLogger()
	battleManager.logger = mockLogger

	args = map[string]interface{}{
		"instanceID": uint32(123),
		"flag":       uint8(10), // Unknown flag
	}

	err = battleManager.HandleInstanceWindowLeave(args)

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
		t.Errorf("Expected instance_window_leave hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"instanceID": "invalid", // Invalid type
		"flag":       uint8(0),
	}

	err = battleManager.HandleInstanceWindowLeave(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}
