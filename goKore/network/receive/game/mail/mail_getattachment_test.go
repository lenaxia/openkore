package mail

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// TestMailGetAttachment tests the HandleMailGetAttachment method
func TestMailGetAttachment(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create mail manager
	mailManager := NewMailManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("mail_getattachment", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test successful attachment retrieval
	args := map[string]interface{}{
		"fail": uint8(0), // Success
	}

	err := mailManager.HandleMailGetAttachment(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify correct message
	expectedMsg := "Successfully added attachment to inventory."
	if mockLogger.infoMessages[0] != expectedMsg {
		t.Errorf("Expected '%s', got '%s'", expectedMsg, mockLogger.infoMessages[0])
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected mail_getattachment hook to be called")
	}

	// Test failure due to weight
	hookCalled = false
	mockLogger = NewMockLogger()
	mailManager.logger = mockLogger

	args = map[string]interface{}{
		"fail": uint8(2), // Failure due to weight
	}

	err = mailManager.HandleMailGetAttachment(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify error message was created
	if len(mockLogger.errorMessages) != 1 {
		t.Errorf("Expected 1 error message, got %d", len(mockLogger.errorMessages))
	}

	// Verify correct message
	expectedMsg = "Failed to get the attachment to inventory due to your weight."
	if mockLogger.errorMessages[0] != expectedMsg {
		t.Errorf("Expected '%s', got '%s'", expectedMsg, mockLogger.errorMessages[0])
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected mail_getattachment hook to be called")
	}

	// Test general failure
	hookCalled = false
	mockLogger = NewMockLogger()
	mailManager.logger = mockLogger

	args = map[string]interface{}{
		"fail": uint8(1), // General failure
	}

	err = mailManager.HandleMailGetAttachment(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify error message was created
	if len(mockLogger.errorMessages) != 1 {
		t.Errorf("Expected 1 error message, got %d", len(mockLogger.errorMessages))
	}

	// Verify correct message
	expectedMsg = "Failed to get the attachment to inventory."
	if mockLogger.errorMessages[0] != expectedMsg {
		t.Errorf("Expected '%s', got '%s'", expectedMsg, mockLogger.errorMessages[0])
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected mail_getattachment hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"fail": "invalid", // Invalid type
	}

	err = mailManager.HandleMailGetAttachment(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}
