package mail

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// TestMailSetAttachment tests the HandleMailSetAttachment method
func TestMailSetAttachment(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create mail manager
	mailManager := NewMailManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("mail_setattachment", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test successful item attachment
	args := map[string]interface{}{
		"fail":   uint8(0),    // Success
		"ID":     uint32(501), // Item ID
		"amount": uint32(5),
	}

	err := mailManager.HandleMailSetAttachment(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify correct message format
	expectedMsg := "Succeeded to attach item: %d."
	if mockLogger.infoMessages[0] != expectedMsg {
		t.Errorf("Expected '%s', got '%s'", expectedMsg, mockLogger.infoMessages[0])
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected mail_setattachment hook to be called")
	}

	// Test successful zeny attachment
	hookCalled = false
	mockLogger = NewMockLogger()
	mailManager.logger = mockLogger

	args = map[string]interface{}{
		"fail":   uint8(0),  // Success
		"ID":     uint32(0), // 0 for zeny
		"amount": uint32(1000),
	}

	err = mailManager.HandleMailSetAttachment(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify correct message format
	expectedMsg = "Succeeded to attach zeny: %d."
	if mockLogger.infoMessages[0] != expectedMsg {
		t.Errorf("Expected '%s', got '%s'", expectedMsg, mockLogger.infoMessages[0])
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected mail_setattachment hook to be called")
	}

	// Test failed item attachment
	hookCalled = false
	mockLogger = NewMockLogger()
	mailManager.logger = mockLogger

	args = map[string]interface{}{
		"fail":   uint8(1),    // Failure
		"ID":     uint32(501), // Item ID
		"amount": uint32(5),
	}

	err = mailManager.HandleMailSetAttachment(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify correct message format
	expectedMsg = "Failed to attach item: %d."
	if mockLogger.infoMessages[0] != expectedMsg {
		t.Errorf("Expected '%s', got '%s'", expectedMsg, mockLogger.infoMessages[0])
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected mail_setattachment hook to be called")
	}

	// Test failed zeny attachment
	hookCalled = false
	mockLogger = NewMockLogger()
	mailManager.logger = mockLogger

	args = map[string]interface{}{
		"fail":   uint8(1),  // Failure
		"ID":     uint32(0), // 0 for zeny
		"amount": uint32(1000),
	}

	err = mailManager.HandleMailSetAttachment(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify correct message format
	expectedMsg = "Failed to attach zeny."
	if mockLogger.infoMessages[0] != expectedMsg {
		t.Errorf("Expected '%s', got '%s'", expectedMsg, mockLogger.infoMessages[0])
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected mail_setattachment hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"fail": "invalid", // Invalid type
	}

	err = mailManager.HandleMailSetAttachment(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}
