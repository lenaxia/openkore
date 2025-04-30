package mail

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// TestMailReturn tests the HandleMailReturn method
func TestMailReturn(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create mail manager
	mailManager := NewMailManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("mail_return", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test mail return success
	args := map[string]interface{}{
		"fail":   uint8(0), // Success
		"mailID": uint32(1),
	}

	err := mailManager.HandleMailReturn(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify correct message format
	expectedMsg := "The mail with ID: %d is returned to the sender."
	if mockLogger.infoMessages[0] != expectedMsg {
		t.Errorf("Expected '%s', got '%s'", expectedMsg, mockLogger.infoMessages[0])
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected mail_return hook to be called")
	}

	// Test mail return failure
	hookCalled = false
	mockLogger = NewMockLogger()
	mailManager.logger = mockLogger

	args = map[string]interface{}{
		"fail":   uint8(1), // Failure
		"mailID": uint32(1),
	}

	err = mailManager.HandleMailReturn(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify error message was created
	if len(mockLogger.errorMessages) != 1 {
		t.Errorf("Expected 1 error message, got %d", len(mockLogger.errorMessages))
	}

	// Verify correct message format
	expectedMsg = "The mail with ID: %d does not exist."
	if mockLogger.errorMessages[0] != expectedMsg {
		t.Errorf("Expected '%s', got '%s'", expectedMsg, mockLogger.errorMessages[0])
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected mail_return hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"fail":   "invalid", // Invalid type
		"mailID": uint32(1),
	}

	err = mailManager.HandleMailReturn(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}
