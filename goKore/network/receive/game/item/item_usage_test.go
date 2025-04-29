package item

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// TestItemUsed tests the HandleItemUsed method
func TestItemUsed(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create item usage manager
	itemUsageManager := NewItemUsageManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("packet_useitem", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test successful item usage by player
	args := map[string]interface{}{
		"ID":        uint16(1001),
		"itemID":    uint16(501),
		"actorID":   "12345", // Player's account ID
		"remaining": uint16(5),
		"success":   uint8(1),
	}

	// Set the account ID for testing
	itemUsageManager.accountID = "12345"

	err := itemUsageManager.HandleItemUsed(args)

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
		t.Errorf("Expected packet_useitem hook to be called")
	}

	// Reset hook called flag and messages
	hookCalled = false
	mockLogger.infoMessages = []string{}

	// Test failed item usage by player
	args = map[string]interface{}{
		"ID":        uint16(1001),
		"itemID":    uint16(501),
		"actorID":   "12345", // Player's account ID
		"remaining": uint16(10),
		"success":   uint8(0),
	}

	err = itemUsageManager.HandleItemUsed(args)

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
		t.Errorf("Expected packet_useitem hook to be called")
	}

	// Reset hook called flag and messages
	hookCalled = false
	mockLogger.infoMessages = []string{}

	// Test item usage by another actor
	args = map[string]interface{}{
		"ID":        uint16(1001),
		"itemID":    uint16(501),
		"actorID":   "67890", // Another actor's ID
		"remaining": uint16(5),
		"success":   uint8(1),
	}

	err = itemUsageManager.HandleItemUsed(args)

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
		t.Errorf("Expected packet_useitem hook to be called")
	}
}

// TestUseItem tests the HandleUseItem method
func TestUseItem(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create item usage manager
	itemUsageManager := NewItemUsageManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("item_used", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test successful item usage
	args := map[string]interface{}{
		"ID":     uint16(1001),
		"amount": uint16(1),
	}

	err := itemUsageManager.HandleUseItem(args)

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
		t.Errorf("Expected item_used hook to be called")
	}
}

// TestRegisterItemUsageHandlers tests the RegisterHandlers method for item usage
func TestRegisterItemUsageHandlers(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create item usage manager
	itemUsageManager := NewItemUsageManager(mockParser, hookManager, mockLogger)

	// Register handlers
	itemUsageManager.RegisterHandlers()

	// Verify handlers were registered
	expectedHandlers := []string{
		"item_used",
		"use_item",
	}

	for _, handler := range expectedHandlers {
		if _, exists := mockParser.handlers[handler]; !exists {
			t.Errorf("Expected handler %s to be registered", handler)
		}
	}
}
