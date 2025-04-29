package item

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// TestUnequipItemSwitch tests the HandleUnequipItemSwitch method
func TestUnequipItemSwitch(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create equipment manager
	equipmentManager := NewEquipmentManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("item_unequipped_switch", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test successful unequip switch for arrow
	args := map[string]interface{}{
		"ID":      uint16(1001),
		"type":    uint16(10), // Arrow
		"success": uint8(1),
	}

	err := equipmentManager.HandleUnequipItemSwitch(args)

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
		t.Errorf("Expected item_unequipped_switch hook to be called")
	}

	// Reset hook called flag
	hookCalled = false

	// Test successful unequip switch for regular equipment
	args = map[string]interface{}{
		"ID":      uint16(1002),
		"type":    uint16(2), // Upper headgear
		"success": uint8(1),
	}

	err = equipmentManager.HandleUnequipItemSwitch(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 2 {
		t.Errorf("Expected 2 info messages, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected item_unequipped_switch hook to be called")
	}

	// Test failed unequip switch
	args = map[string]interface{}{
		"ID":      uint16(1003),
		"type":    uint16(2), // Upper headgear
		"success": uint8(0),
	}

	err = equipmentManager.HandleUnequipItemSwitch(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify error message was created
	if len(mockLogger.errorMessages) != 1 {
		t.Errorf("Expected 1 error message, got %d", len(mockLogger.errorMessages))
	}
}
