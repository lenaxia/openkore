package mail

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// TestMailSend tests the HandleMailSend method
func TestMailSend(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create mail manager
	mailManager := NewMailManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("mail_send", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test successful mail send
	args := map[string]interface{}{
		"fail": uint8(0), // Success
	}

	err := mailManager.HandleMailSend(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify correct message
	expectedMsg := "Mail sent successfully."
	if mockLogger.infoMessages[0] != expectedMsg {
		t.Errorf("Expected '%s', got '%s'", expectedMsg, mockLogger.infoMessages[0])
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected mail_send hook to be called")
	}

	// Test failed mail send
	hookCalled = false
	mockLogger = NewMockLogger()
	mailManager.logger = mockLogger

	args = map[string]interface{}{
		"fail": uint8(1), // Failure
	}

	err = mailManager.HandleMailSend(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify error message was created
	if len(mockLogger.errorMessages) != 1 {
		t.Errorf("Expected 1 error message, got %d", len(mockLogger.errorMessages))
	}

	// Verify correct message
	expectedMsg = "Failed to send mail, the recipient does not exist."
	if mockLogger.errorMessages[0] != expectedMsg {
		t.Errorf("Expected '%s', got '%s'", expectedMsg, mockLogger.errorMessages[0])
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected mail_send hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"fail": "invalid", // Invalid type
	}

	err = mailManager.HandleMailSend(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}
