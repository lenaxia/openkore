package mail

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// TestRodexGetItem tests the HandleRodexGetItem method
func TestRodexGetItem(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create rodex manager
	rodexManager := NewRodexManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("rodex_get_item", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Initialize rodexList with a mail that has items
	rodexManager.rodexList = make(map[string]interface{})
	rodexManager.rodexList["mails"] = make(map[uint32]map[string]interface{})

	// Create test items
	items := []map[string]interface{}{
		{
			"ID":         uint32(1),
			"nameID":     uint16(501),
			"amount":     uint16(5),
			"identified": uint8(1),
			"broken":     uint8(0),
			"upgrade":    uint8(7),
		},
		{
			"ID":         uint32(2),
			"nameID":     uint16(502),
			"amount":     uint16(1),
			"identified": uint8(1),
			"broken":     uint8(0),
			"upgrade":    uint8(0),
		},
	}

	// Add mail with items
	rodexManager.rodexList["mails"].(map[uint32]map[string]interface{})[1] = map[string]interface{}{
		"mailID1": uint32(1),
		"sender":  "TestSender",
		"title":   "TestTitle",
		"items":   items,
	}

	// Test rodex get item
	args := map[string]interface{}{
		"mailID1":   uint32(1),
		"mailID2":   uint32(0),
		"itemCount": uint8(2),
	}

	err := rodexManager.HandleRodexGetItem(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify mail items were cleared
	if mail, ok := rodexManager.rodexList["mails"].(map[uint32]map[string]interface{})[1]; ok {
		// Check if items were cleared
		if items, ok := mail["items"].([]map[string]interface{}); !ok || len(items) != 0 {
			t.Errorf("Expected items to be empty, got %d items", len(items))
		}
	} else {
		t.Errorf("Mail with ID 1 not found in rodexList")
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected rodex_get_item hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"mailID1":   "invalid", // Invalid type
		"mailID2":   uint32(0),
		"itemCount": uint8(2),
	}

	err = rodexManager.HandleRodexGetItem(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}
