package mail

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// TestMailWindow tests the HandleMailWindow method
func TestMailWindow(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create mail manager
	mailManager := NewMailManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("mail_window", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test mail window open
	args := map[string]interface{}{
		"flag": uint8(0), // Window opened
	}

	err := mailManager.HandleMailWindow(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify correct message
	if mockLogger.infoMessages[0] != "Mail window is now opened." {
		t.Errorf("Expected 'Mail window is now opened.', got '%s'", mockLogger.infoMessages[0])
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected mail_window hook to be called")
	}

	// Test mail window close
	hookCalled = false
	mockLogger = NewMockLogger()
	mailManager.logger = mockLogger

	args = map[string]interface{}{
		"flag": uint8(1), // Window closed
	}

	err = mailManager.HandleMailWindow(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify correct message
	if mockLogger.infoMessages[0] != "Mail window is now closed." {
		t.Errorf("Expected 'Mail window is now closed.', got '%s'", mockLogger.infoMessages[0])
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected mail_window hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"flag": "invalid", // Invalid type
	}

	err = mailManager.HandleMailWindow(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}
