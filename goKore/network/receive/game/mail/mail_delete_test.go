package mail

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// TestMailDelete tests the HandleMailDelete method
func TestMailDelete(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create mail manager
	mailManager := NewMailManager(mockParser, hookManager, mockLogger)

	// Initialize mailList with a test mail
	mailManager.mailList[1] = map[string]interface{}{
		"mailID": uint32(1),
		"title":  "Test Mail",
		"sender": "TestSender",
	}

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("mail_delete", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test mail delete success
	args := map[string]interface{}{
		"fail":   uint8(0), // Success
		"mailID": uint32(1),
	}

	err := mailManager.HandleMailDelete(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify mail was deleted from mailList
	if _, exists := mailManager.mailList[1]; exists {
		t.Errorf("Expected mail with ID 1 to be deleted from mailList")
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected mail_delete hook to be called")
	}

	// Test mail delete failure
	hookCalled = false

	// Add mail back to mailList
	mailManager.mailList[1] = map[string]interface{}{
		"mailID": uint32(1),
		"title":  "Test Mail",
		"sender": "TestSender",
	}

	args = map[string]interface{}{
		"fail":   uint8(1), // Failure
		"mailID": uint32(1),
	}

	err = mailManager.HandleMailDelete(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 2 {
		t.Errorf("Expected 2 info messages, got %d", len(mockLogger.infoMessages))
	}

	// Verify mail was not deleted from mailList
	if _, exists := mailManager.mailList[1]; !exists {
		t.Errorf("Expected mail with ID 1 to still exist in mailList")
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected mail_delete hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"fail":   "invalid", // Invalid type
		"mailID": uint32(1),
	}

	err = mailManager.HandleMailDelete(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}
