package mail

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// TestRodexDelete tests the HandleRodexDelete method
func TestRodexDelete(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create rodex manager
	rodexManager := NewRodexManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("rodex_delete", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Initialize rodexList with a mail
	rodexManager.rodexList = make(map[string]interface{})
	rodexManager.rodexList["mails"] = make(map[uint32]map[string]interface{})
	rodexManager.rodexList["mails"].(map[uint32]map[string]interface{})[1] = map[string]interface{}{
		"mailID1": uint32(1),
		"sender":  "TestSender",
		"title":   "TestTitle",
	}

	// Test rodex delete
	args := map[string]interface{}{
		"mailID1": uint32(1),
		"mailID2": uint32(0),
	}

	err := rodexManager.HandleRodexDelete(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify mail was deleted
	mails := rodexManager.rodexList["mails"].(map[uint32]map[string]interface{})
	if _, exists := mails[1]; exists {
		t.Errorf("Expected mail with ID 1 to be deleted")
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected rodex_delete hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"mailID1": "invalid", // Invalid type
		"mailID2": uint32(0),
	}

	err = rodexManager.HandleRodexDelete(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}
