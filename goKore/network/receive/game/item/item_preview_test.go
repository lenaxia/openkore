package item

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// TestItemPreview tests the HandleItemPreview method
func TestItemPreview(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create inventory manager
	inventoryManager := NewInventoryManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("item_preview", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test item preview
	args := map[string]interface{}{
		"index":   uint16(1001),
		"broken":  uint8(0),
		"upgrade": uint8(7),
		"cards":   "1234,5678,0,0",
		"options": "1:2,3:4",
	}

	err := inventoryManager.HandleItemPreview(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify debug message was created
	if len(mockLogger.debugMessages) != 1 {
		t.Errorf("Expected 1 debug message, got %d", len(mockLogger.debugMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected item_preview hook to be called")
	}

	// Test with missing parameters
	args = map[string]interface{}{
		"index": uint16(1001),
	}

	err = inventoryManager.HandleItemPreview(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify debug message was created
	if len(mockLogger.debugMessages) != 2 {
		t.Errorf("Expected 2 debug messages, got %d", len(mockLogger.debugMessages))
	}
}
