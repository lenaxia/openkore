package item

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// TestNewHandler tests the NewHandler function
func TestNewHandler(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create handler
	handler := NewHandler(mockParser, hookManager, mockLogger)

	// Verify handler was created
	if handler == nil {
		t.Errorf("Expected handler to be created")
	}

	// Verify sub-managers were created
	if handler.inventoryManager == nil {
		t.Errorf("Expected inventory manager to be created")
	}

	if handler.storageManager == nil {
		t.Errorf("Expected storage manager to be created")
	}

	if handler.equipmentManager == nil {
		t.Errorf("Expected equipment manager to be created")
	}
}

// TestRegisterAllHandlers tests the RegisterHandlers method for the main handler
func TestRegisterAllHandlers(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create handler
	handler := NewHandler(mockParser, hookManager, mockLogger)

	// Register handlers
	handler.RegisterHandlers()

	// Verify handlers were registered
	// We can't directly check the handlers, but we can verify that the
	// sub-managers' RegisterHandlers methods were called by checking
	// if the expected handlers were registered

	// Expected inventory handlers
	inventoryHandlers := []string{
		"inventory_item_added",
		"inventory_item_removed",
		"inventory_items_stackable",
		"inventory_items_nonstackable",
		"inventory_item_favorite",
		"inventory_expansion_result",
	}

	// Expected storage handlers
	storageHandlers := []string{
		"storage_opened",
		"storage_closed",
		"storage_items_stackable",
		"storage_items_nonstackable",
		"storage_item_added",
		"storage_item_removed",
		"storage_password_request",
		"storage_password_result",
	}

	// Expected equipment handlers
	equipmentHandlers := []string{
		"equip_item",
		"unequip_item",
		"equip_item_switch",
		"equip_switch_run_res",
		"equip_switch_log",
		"arrow_equipped",
	}

	// Check if all expected handlers were registered
	for _, handlerName := range inventoryHandlers {
		if _, exists := mockParser.handlers[handlerName]; !exists {
			t.Errorf("Expected inventory handler %s to be registered", handlerName)
		}
	}

	for _, handlerName := range storageHandlers {
		if _, exists := mockParser.handlers[handlerName]; !exists {
			t.Errorf("Expected storage handler %s to be registered", handlerName)
		}
	}

	for _, handlerName := range equipmentHandlers {
		if _, exists := mockParser.handlers[handlerName]; !exists {
			t.Errorf("Expected equipment handler %s to be registered", handlerName)
		}
	}
}

// TestGetInventoryManager tests the GetInventoryManager method
func TestGetInventoryManager(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create handler
	handler := NewHandler(mockParser, hookManager, mockLogger)

	// Get inventory manager
	inventoryManager := handler.GetInventoryManager()

	// Verify inventory manager was returned
	if inventoryManager == nil {
		t.Errorf("Expected inventory manager to be returned")
	}

	// Verify it's the same instance
	if inventoryManager != handler.inventoryManager {
		t.Errorf("Expected returned inventory manager to be the same instance")
	}
}

// TestGetStorageManager tests the GetStorageManager method
func TestGetStorageManager(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create handler
	handler := NewHandler(mockParser, hookManager, mockLogger)

	// Get storage manager
	storageManager := handler.GetStorageManager()

	// Verify storage manager was returned
	if storageManager == nil {
		t.Errorf("Expected storage manager to be returned")
	}

	// Verify it's the same instance
	if storageManager != handler.storageManager {
		t.Errorf("Expected returned storage manager to be the same instance")
	}
}

// TestGetEquipmentManager tests the GetEquipmentManager method
func TestGetEquipmentManager(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create handler
	handler := NewHandler(mockParser, hookManager, mockLogger)

	// Get equipment manager
	equipmentManager := handler.GetEquipmentManager()

	// Verify equipment manager was returned
	if equipmentManager == nil {
		t.Errorf("Expected equipment manager to be returned")
	}

	// Verify it's the same instance
	if equipmentManager != handler.equipmentManager {
		t.Errorf("Expected returned equipment manager to be the same instance")
	}
}
