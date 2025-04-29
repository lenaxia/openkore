package item

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// TestArrowNone tests the HandleArrowNone method
func TestArrowNone(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create equipment manager
	equipmentManager := NewEquipmentManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("arrow_none", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test arrow none (no arrows)
	args := map[string]interface{}{
		"type": uint8(0),
	}

	err := equipmentManager.HandleArrowNone(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify error message was created
	if len(mockLogger.errorMessages) != 1 {
		t.Errorf("Expected 1 error message, got %d", len(mockLogger.errorMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected arrow_none hook to be called")
	}

	// Reset hook called flag and messages
	hookCalled = false
	mockLogger.errorMessages = []string{}
	mockLogger.debugMessages = []string{}

	// Test arrow none (weight limit exceeded - can't attack or use skills)
	args = map[string]interface{}{
		"type": uint8(1),
	}

	err = equipmentManager.HandleArrowNone(args)

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
		t.Errorf("Expected arrow_none hook to be called")
	}

	// Reset hook called flag and messages
	hookCalled = false
	mockLogger.debugMessages = []string{}

	// Test arrow none (weight limit exceeded - can't use skills)
	args = map[string]interface{}{
		"type": uint8(2),
	}

	err = equipmentManager.HandleArrowNone(args)

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
		t.Errorf("Expected arrow_none hook to be called")
	}

	// Reset hook called flag and messages
	hookCalled = false
	mockLogger.debugMessages = []string{}

	// Test arrow none (arrow equipped)
	args = map[string]interface{}{
		"type": uint8(3),
	}

	err = equipmentManager.HandleArrowNone(args)

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
		t.Errorf("Expected arrow_none hook to be called")
	}
}
