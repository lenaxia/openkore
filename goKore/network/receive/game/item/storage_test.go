package item

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// TestStorageOpened tests the HandleStorageOpened method
func TestStorageOpened(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create storage manager
	storageManager := NewStorageManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("storage_opened", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test storage opened
	args := map[string]interface{}{}

	err := storageManager.HandleStorageOpened(args)

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
		t.Errorf("Expected storage_opened hook to be called")
	}
}

// TestStorageClosed tests the HandleStorageClosed method
func TestStorageClosed(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create storage manager
	storageManager := NewStorageManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("storage_closed", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test storage closed
	args := map[string]interface{}{}

	err := storageManager.HandleStorageClosed(args)

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
		t.Errorf("Expected storage_closed hook to be called")
	}
}

// TestStorageItemAdded tests the HandleStorageItemAdded method
func TestStorageItemAdded(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create storage manager
	storageManager := NewStorageManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("packet_storage_add", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test item added to storage
	args := map[string]interface{}{
		"ID":     "12345",
		"amount": uint16(10),
		"nameID": uint16(501),
	}

	err := storageManager.HandleStorageItemAdded(args)

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
		t.Errorf("Expected packet_storage_add hook to be called")
	}
}

// TestStorageItemRemoved tests the HandleStorageItemRemoved method
func TestStorageItemRemoved(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create storage manager
	storageManager := NewStorageManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("packet_storage_remove", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test item removed from storage
	args := map[string]interface{}{
		"ID":     "12345",
		"amount": uint16(10),
	}

	err := storageManager.HandleStorageItemRemoved(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected packet_storage_remove hook to be called")
	}
}

// TestStorageItemsStackable tests the HandleStorageItemsStackable method
func TestStorageItemsStackable(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create storage manager
	storageManager := NewStorageManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCallCount := 0
	hookManager.AddHook("packet_storage", func(hookName string, data interface{}, userData interface{}) {
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

	err := storageManager.HandleStorageItemsStackable(args)

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

// TestStorageItemsNonstackable tests the HandleStorageItemsNonstackable method
func TestStorageItemsNonstackable(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create storage manager
	storageManager := NewStorageManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCallCount := 0
	hookManager.AddHook("packet_storage", func(hookName string, data interface{}, userData interface{}) {
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

	err := storageManager.HandleStorageItemsNonstackable(args)

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

// TestStoragePasswordRequest tests the HandleStoragePasswordRequest method
func TestStoragePasswordRequest(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create storage manager
	storageManager := NewStorageManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("storage_password_request", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test new password request
	args := map[string]interface{}{
		"flag": uint8(0),
	}

	err := storageManager.HandleStoragePasswordRequest(args)

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
		t.Errorf("Expected storage_password_request hook to be called")
	}

	// Test password verification request
	args = map[string]interface{}{
		"flag": uint8(1),
	}

	err = storageManager.HandleStoragePasswordRequest(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 2 {
		t.Errorf("Expected 2 info messages, got %d", len(mockLogger.infoMessages))
	}

	// Test too many incorrect attempts
	args = map[string]interface{}{
		"flag": uint8(8),
	}

	err = storageManager.HandleStoragePasswordRequest(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify error message was created
	if len(mockLogger.errorMessages) != 1 {
		t.Errorf("Expected 1 error message, got %d", len(mockLogger.errorMessages))
	}
}

// TestStoragePasswordResult tests the HandleStoragePasswordResult method
func TestStoragePasswordResult(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create storage manager
	storageManager := NewStorageManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("storage_password_result", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test password change success
	args := map[string]interface{}{
		"type": uint8(4), // STORE_PASSWORD_CHANGE_OK
	}

	err := storageManager.HandleStoragePasswordResult(args)

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
		t.Errorf("Expected storage_password_result hook to be called")
	}

	// Test password check failure
	args = map[string]interface{}{
		"type": uint8(7), // STORE_PASSWORD_CHECK_NG
	}

	err = storageManager.HandleStoragePasswordResult(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify error message was created
	if len(mockLogger.errorMessages) != 1 {
		t.Errorf("Expected 1 error message, got %d", len(mockLogger.errorMessages))
	}
}

// TestRegisterStorageHandlers tests the RegisterHandlers method for storage
func TestRegisterStorageHandlers(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create storage manager
	storageManager := NewStorageManager(mockParser, hookManager, mockLogger)

	// Register handlers
	storageManager.RegisterHandlers()

	// Verify handlers were registered
	expectedHandlers := []string{
		"storage_opened",
		"storage_closed",
		"storage_items_stackable",
		"storage_items_nonstackable",
		"storage_item_added",
		"storage_item_removed",
		"storage_password_request",
		"storage_password_result",
	}

	for _, handler := range expectedHandlers {
		if _, exists := mockParser.handlers[handler]; !exists {
			t.Errorf("Expected handler %s to be registered", handler)
		}
	}
}
