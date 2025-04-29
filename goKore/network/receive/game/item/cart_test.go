package item

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// TestCartOff tests the HandleCartOff method
func TestCartOff(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create cart manager
	cartManager := NewCartManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("cart_off", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test cart off
	args := map[string]interface{}{}

	err := cartManager.HandleCartOff(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify success message was created
	if len(mockLogger.successMessages) != 1 {
		t.Errorf("Expected 1 success message, got %d", len(mockLogger.successMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected cart_off hook to be called")
	}
}

// TestCartItemsStackable tests the HandleCartItemsStackable method
func TestCartItemsStackable(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create cart manager
	cartManager := NewCartManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCallCount := 0
	hookManager.AddHook("packet_cart", func(hookName string, data interface{}, userData interface{}) {
		hookCallCount++
	}, nil)

	// Test stackable items
	itemInfo := []map[string]interface{}{
		{
			"ID":         "12345",
			"nameID":     uint16(501),
			"amount":     uint16(10),
			"type":       uint8(3),
			"identified": uint8(1),
		},
		{
			"ID":         "67890",
			"nameID":     uint16(502),
			"amount":     uint16(5),
			"type":       uint8(3),
			"identified": uint8(1),
		},
	}

	args := map[string]interface{}{
		"itemInfo": itemInfo,
	}

	err := cartManager.HandleCartItemsStackable(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify hooks were called for each item
	if hookCallCount != 2 {
		t.Errorf("Expected 2 hook calls, got %d", hookCallCount)
	}

	// Verify debug message was created
	if len(mockLogger.debugMessages) != 1 {
		t.Errorf("Expected 1 debug message, got %d", len(mockLogger.debugMessages))
	}
}

// TestCartItemsNonstackable tests the HandleCartItemsNonstackable method
func TestCartItemsNonstackable(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create cart manager
	cartManager := NewCartManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCallCount := 0
	hookManager.AddHook("packet_cart", func(hookName string, data interface{}, userData interface{}) {
		hookCallCount++
	}, nil)

	// Test non-stackable items
	itemInfo := []map[string]interface{}{
		{
			"ID":         "12345",
			"nameID":     uint16(501),
			"amount":     uint16(1),
			"type":       uint8(4),
			"identified": uint8(1),
		},
		{
			"ID":         "67890",
			"nameID":     uint16(502),
			"amount":     uint16(1),
			"type":       uint8(4),
			"identified": uint8(1),
		},
	}

	args := map[string]interface{}{
		"itemInfo": itemInfo,
	}

	err := cartManager.HandleCartItemsNonstackable(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify hooks were called for each item
	if hookCallCount != 2 {
		t.Errorf("Expected 2 hook calls, got %d", hookCallCount)
	}

	// Verify debug message was created
	if len(mockLogger.debugMessages) != 1 {
		t.Errorf("Expected 1 debug message, got %d", len(mockLogger.debugMessages))
	}
}

// TestCartItemAdded tests the HandleCartItemAdded method
func TestCartItemAdded(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create cart manager
	cartManager := NewCartManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("packet_cart_add", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test item added to cart
	args := map[string]interface{}{
		"ID":     "12345",
		"amount": uint16(10),
		"nameID": uint16(501),
	}

	err := cartManager.HandleCartItemAdded(args)

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
		t.Errorf("Expected packet_cart_add hook to be called")
	}
}

// TestCartItemRemoved tests the HandleCartItemRemoved method
func TestCartItemRemoved(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create cart manager
	cartManager := NewCartManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("packet_cart_remove", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test item removed from cart
	args := map[string]interface{}{
		"ID":     "12345",
		"amount": uint16(10),
	}

	err := cartManager.HandleCartItemRemoved(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected packet_cart_remove hook to be called")
	}
}

// TestCartInfo tests the HandleCartInfo method
func TestCartInfo(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create cart manager
	cartManager := NewCartManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("cart_info", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test cart info
	args := map[string]interface{}{
		"items":      uint16(10),
		"items_max":  uint16(100),
		"weight":     uint32(1000),
		"weight_max": uint32(2000),
	}

	err := cartManager.HandleCartInfo(args)

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
		t.Errorf("Expected cart_info hook to be called")
	}
}

// TestCartAddFailed tests the HandleCartAddFailed method
func TestCartAddFailed(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create cart manager
	cartManager := NewCartManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("cart_add_failed", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test overweight
	args := map[string]interface{}{
		"fail": uint8(0),
	}

	err := cartManager.HandleCartAddFailed(args)

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
		t.Errorf("Expected cart_add_failed hook to be called")
	}

	// Reset
	hookCalled = false
	mockLogger.errorMessages = []string{}

	// Test too many items
	args = map[string]interface{}{
		"fail": uint8(1),
	}

	err = cartManager.HandleCartAddFailed(args)

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
		t.Errorf("Expected cart_add_failed hook to be called")
	}
}

// TestRegisterCartHandlers tests the RegisterHandlers method for cart
func TestRegisterCartHandlers(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create cart manager
	cartManager := NewCartManager(mockParser, hookManager, mockLogger)

	// Register handlers
	cartManager.RegisterHandlers()

	// Verify handlers were registered
	expectedHandlers := []string{
		"cart_off",
		"cart_items_stackable",
		"cart_items_nonstackable",
		"cart_item_added",
		"cart_item_removed",
		"cart_info",
		"cart_add_failed",
	}

	for _, handler := range expectedHandlers {
		if _, exists := mockParser.handlers[handler]; !exists {
			t.Errorf("Expected handler %s to be registered", handler)
		}
	}
}
