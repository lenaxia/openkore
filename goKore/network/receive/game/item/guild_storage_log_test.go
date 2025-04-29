package item

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// TestGuildStorageLog tests the HandleGuildStorageLog method
func TestGuildStorageLog(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create storage manager
	storageManager := NewStorageManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("guild_storage_log", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test guild storage log with items
	args := map[string]interface{}{
		"result": uint8(0),                             // Get action
		"log":    []byte{0x01, 0x02, 0x03, 0x04, 0x05}, // Sample log data
	}

	err := storageManager.HandleGuildStorageLog(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected guild_storage_log hook to be called")
	}

	// Reset hook called flag and messages
	hookCalled = false
	mockLogger.infoMessages = []string{}

	// Test empty guild storage
	args = map[string]interface{}{
		"result": uint8(2), // Empty storage
	}

	err = storageManager.HandleGuildStorageLog(args)

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
		t.Errorf("Expected guild_storage_log hook to be called")
	}

	// Reset hook called flag and messages
	hookCalled = false
	mockLogger.infoMessages = []string{}

	// Test not using guild storage
	args = map[string]interface{}{
		"result": uint8(3), // Not using guild storage
	}

	err = storageManager.HandleGuildStorageLog(args)

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
		t.Errorf("Expected guild_storage_log hook to be called")
	}
}
