package minimap

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

// TestMinimapIndicator tests the HandleMinimapIndicator method
func TestMinimapIndicator(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create minimap manager
	minimapManager := NewMinimapManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("minimap_indicator", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test minimap indicator (show)
	args := map[string]interface{}{
		"npcID":  uint32(12345),
		"x":      uint16(100),
		"y":      uint16(200),
		"type":   uint8(1),
		"effect": uint16(0),
		"red":    uint8(255),
		"green":  uint8(0),
		"blue":   uint8(0),
		"alpha":  uint8(128),
	}

	err := minimapManager.HandleMinimapIndicator(args)

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
		t.Errorf("Expected minimap_indicator hook to be called")
	}

	// Test minimap indicator (clear)
	hookCalled = false
	mockLogger = NewMockLogger()
	minimapManager.logger = mockLogger

	args = map[string]interface{}{
		"npcID":  uint32(12345),
		"x":      uint16(100),
		"y":      uint16(200),
		"type":   uint8(2), // Type 2 means clear
		"effect": uint16(0),
		"red":    uint8(255),
		"green":  uint8(0),
		"blue":   uint8(0),
		"alpha":  uint8(128),
	}

	err = minimapManager.HandleMinimapIndicator(args)

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
		t.Errorf("Expected minimap_indicator hook to be called")
	}

	// Test minimap indicator (quest)
	hookCalled = false
	mockLogger = NewMockLogger()
	minimapManager.logger = mockLogger

	args = map[string]interface{}{
		"npcID":  uint32(12345),
		"x":      uint16(100),
		"y":      uint16(200),
		"type":   uint8(1),
		"effect": uint16(1), // Effect 1 means quest
		"red":    uint8(255),
		"green":  uint8(0),
		"blue":   uint8(0),
		"alpha":  uint8(128),
	}

	err = minimapManager.HandleMinimapIndicator(args)

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
		t.Errorf("Expected minimap_indicator hook to be called")
	}

	// Test minimap indicator (special case)
	hookCalled = false
	mockLogger = NewMockLogger()
	minimapManager.logger = mockLogger

	args = map[string]interface{}{
		"npcID":  uint32(12345),
		"x":      uint16(100),
		"y":      uint16(200),
		"type":   uint8(1),
		"effect": uint16(9999), // Effect 9999 is a special case
		"red":    uint8(255),
		"green":  uint8(0),
		"blue":   uint8(0),
		"alpha":  uint8(128),
	}

	err = minimapManager.HandleMinimapIndicator(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify no info message was created (special case)
	if len(mockLogger.infoMessages) != 0 {
		t.Errorf("Expected 0 info messages, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was not called (special case)
	if hookCalled {
		t.Errorf("Expected minimap_indicator hook to not be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"npcID": "invalid", // Invalid type
		"x":     uint16(100),
		"y":     uint16(200),
	}

	err = minimapManager.HandleMinimapIndicator(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestColorString tests the String method of the Color struct
func TestColorString(t *testing.T) {
	// Create a color
	color := Color{
		Red:   255,
		Green: 0,
		Blue:  0,
		Alpha: 128,
	}

	// Test string representation
	expected := "[R:255, G:0, B:0, A:128]"
	if color.String() != expected {
		t.Errorf("Expected %s, got %s", expected, color.String())
	}
}
