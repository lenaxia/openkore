package marriage

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

// TestMarried tests the HandleMarried method
func TestMarried(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create marriage manager
	marriageManager := NewMarriageManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("married", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test married
	args := map[string]interface{}{
		"ID": uint32(12345),
	}

	err := marriageManager.HandleMarried(args)

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
		t.Errorf("Expected married hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"ID": "invalid", // Invalid type
	}

	err = marriageManager.HandleMarried(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestDivorced tests the HandleDivorced method
func TestDivorced(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create marriage manager
	marriageManager := NewMarriageManager(mockParser, hookManager, mockLogger)
	marriageManager.SetCharName("TestChar")

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("divorced", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test divorced
	args := map[string]interface{}{
		"name": "TestPartner",
	}

	err := marriageManager.HandleDivorced(args)

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
		t.Errorf("Expected divorced hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"name": 123, // Invalid type
	}

	err = marriageManager.HandleDivorced(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestMarriagePartnerName tests the HandleMarriagePartnerName method
func TestMarriagePartnerName(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create marriage manager
	marriageManager := NewMarriageManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("marriage_partner_name", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test marriage partner name
	args := map[string]interface{}{
		"name": "TestPartner",
	}

	err := marriageManager.HandleMarriagePartnerName(args)

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
		t.Errorf("Expected marriage_partner_name hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"name": 123, // Invalid type
	}

	err = marriageManager.HandleMarriagePartnerName(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestAdoptRequest tests the HandleAdoptRequest method
func TestAdoptRequest(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create marriage manager
	marriageManager := NewMarriageManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("adopt_request", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test adopt request
	args := map[string]interface{}{
		"name": "TestAdopter",
	}

	err := marriageManager.HandleAdoptRequest(args)

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
		t.Errorf("Expected adopt_request hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"name": 123, // Invalid type
	}

	err = marriageManager.HandleAdoptRequest(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestAdoptReply tests the HandleAdoptReply method
func TestAdoptReply(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create marriage manager
	marriageManager := NewMarriageManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("adopt_reply", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test adopt reply (cannot adopt more than 1 child)
	args := map[string]interface{}{
		"type": uint8(0),
	}

	err := marriageManager.HandleAdoptReply(args)

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
		t.Errorf("Expected adopt_reply hook to be called")
	}

	// Test adopt reply (must be at least level 70)
	hookCalled = false
	mockLogger = NewMockLogger()
	marriageManager.logger = mockLogger

	args = map[string]interface{}{
		"type": uint8(1),
	}

	err = marriageManager.HandleAdoptReply(args)

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
		t.Errorf("Expected adopt_reply hook to be called")
	}

	// Test adopt reply (unknown type)
	hookCalled = false
	mockLogger = NewMockLogger()
	marriageManager.logger = mockLogger

	args = map[string]interface{}{
		"type": uint8(10), // Unknown type
	}

	err = marriageManager.HandleAdoptReply(args)

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
		t.Errorf("Expected adopt_reply hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"type": "invalid", // Invalid type
	}

	err = marriageManager.HandleAdoptReply(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}
