package item

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// TestSpecialItemObtain tests the HandleSpecialItemObtain method
func TestSpecialItemObtain(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create inventory manager
	inventoryManager := NewInventoryManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("packet_special_item_obtain", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test box item obtain
	args := map[string]interface{}{
		"type":   uint8(0), // TYPE_BOXITEM
		"nameID": uint16(501),
		"holder": "PlayerName",
		"etc":    []byte{2, 0, 10}, // c/v format with box_nameID = 10
	}

	err := inventoryManager.HandleSpecialItemObtain(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected packet_special_item_obtain hook to be called")
	}

	// Reset hook called flag and messages
	hookCalled = false
	mockLogger.infoMessages = []string{}

	// Test monster item obtain
	args = map[string]interface{}{
		"type":   uint8(1), // TYPE_MONSTER_ITEM
		"nameID": uint16(501),
		"holder": "PlayerName",
		"etc":    []byte{10, 77, 111, 110, 115, 116, 101, 114, 78, 97, 109, 101}, // len=10, monster_name="MonsterName"
	}

	err = inventoryManager.HandleSpecialItemObtain(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected packet_special_item_obtain hook to be called")
	}

	// Reset hook called flag and messages
	hookCalled = false
	mockLogger.infoMessages = []string{}
	mockLogger.warningMessages = []string{}

	// Test unknown type
	args = map[string]interface{}{
		"type":   uint8(2), // Unknown type
		"nameID": uint16(501),
		"holder": "PlayerName",
		"etc":    []byte{},
	}

	err = inventoryManager.HandleSpecialItemObtain(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify warning message was created
	if len(mockLogger.warningMessages) != 1 {
		t.Errorf("Expected 1 warning message, got %d", len(mockLogger.warningMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected packet_special_item_obtain hook to be called")
	}
}
