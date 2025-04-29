package item

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// TestItemAppeared tests the HandleItemAppeared method
func TestItemAppeared(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create inventory manager
	inventoryManager := NewInventoryManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("item_appeared", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Create test args
	args := map[string]interface{}{
		"ID":         "12345",
		"nameID":     uint16(501),
		"amount":     uint16(10),
		"identified": uint8(1),
		"x":          uint16(100),
		"y":          uint16(200),
		"subX":       uint8(5),
		"subY":       uint8(5),
		"type":       uint16(3),
	}

	// Call handler
	err := inventoryManager.HandleItemAppeared(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) == 0 {
		t.Errorf("Expected info message, got none")
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected item_appeared hook to be called")
	}

	// Test with missing ID parameter
	args = map[string]interface{}{
		"nameID":     uint16(501),
		"amount":     uint16(10),
		"identified": uint8(1),
		"x":          uint16(100),
		"y":          uint16(200),
	}

	err = inventoryManager.HandleItemAppeared(args)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for missing ID parameter, got none")
	}
}

// TestItemExists tests the HandleItemExists method
func TestItemExists(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create inventory manager
	inventoryManager := NewInventoryManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("item_exists", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Create test args
	args := map[string]interface{}{
		"ID":          "12345",
		"nameID":      uint16(501),
		"amount":      uint16(10),
		"identified":  uint8(1),
		"x":           uint16(100),
		"y":           uint16(200),
		"type":        uint16(3),
		"show_effect": uint8(1),
		"effect_type": uint16(2),
	}

	// Call handler
	err := inventoryManager.HandleItemExists(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) == 0 {
		t.Errorf("Expected info message, got none")
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected item_exists hook to be called")
	}

	// Test with missing ID parameter
	args = map[string]interface{}{
		"nameID":     uint16(501),
		"amount":     uint16(10),
		"identified": uint8(1),
		"x":          uint16(100),
		"y":          uint16(200),
	}

	err = inventoryManager.HandleItemExists(args)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for missing ID parameter, got none")
	}
}

// TestItemDisappeared tests the HandleItemDisappeared method
func TestItemDisappeared(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create inventory manager
	inventoryManager := NewInventoryManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("item_disappeared", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Create test args
	args := map[string]interface{}{
		"ID": "12345",
	}

	// Call handler
	err := inventoryManager.HandleItemDisappeared(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify debug message was created
	if len(mockLogger.debugMessages) == 0 {
		t.Errorf("Expected debug message, got none")
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected item_disappeared hook to be called")
	}

	// Test with missing ID parameter
	args = map[string]interface{}{}

	err = inventoryManager.HandleItemDisappeared(args)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for missing ID parameter, got none")
	}
}
