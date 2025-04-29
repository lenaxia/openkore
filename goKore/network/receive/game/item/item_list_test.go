package item

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// TestItemListStart tests the HandleItemListStart method
func TestItemListStart(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create inventory manager
	inventoryManager := NewInventoryManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("item_list_start", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test cases for different inventory types
	testCases := []struct {
		name        string
		invType     uint8
		invTypeName string
		expectHook  bool
	}{
		{
			name:        "Inventory type",
			invType:     0, // INVTYPE_INVENTORY
			invTypeName: "inventory",
			expectHook:  true,
		},
		{
			name:        "Cart type",
			invType:     1, // INVTYPE_CART
			invTypeName: "cart",
			expectHook:  true,
		},
		{
			name:        "Storage type",
			invType:     2, // INVTYPE_STORAGE
			invTypeName: "storage",
			expectHook:  true,
		},
		{
			name:        "Guild storage type",
			invType:     3, // INVTYPE_GUILD_STORAGE
			invTypeName: "guild storage",
			expectHook:  true,
		},
		{
			name:        "Unknown type",
			invType:     99,
			invTypeName: "unknown",
			expectHook:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset hook called flag
			hookCalled = false

			// Create test args
			args := map[string]interface{}{
				"type": tc.invType,
				"name": "Test Container",
			}

			// Call handler
			err := inventoryManager.HandleItemListStart(args)

			// Verify no error occurred
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			// Verify debug message was created
			if len(mockLogger.debugMessages) == 0 {
				t.Errorf("Expected debug message, got none")
			}

			// Verify hook was called if expected
			if tc.expectHook && !hookCalled {
				t.Errorf("Expected item_list_start hook to be called")
			}
		})
	}

	// Test with missing type parameter
	args := map[string]interface{}{
		"name": "Test Container",
	}

	err := inventoryManager.HandleItemListStart(args)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for missing type parameter, got none")
	}
}

// TestItemListStackable tests the HandleItemListStackable method
func TestItemListStackable(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create inventory manager
	inventoryManager := NewInventoryManager(mockParser, hookManager, mockLogger)

	// Test cases for different inventory types
	testCases := []struct {
		name       string
		invType    uint8
		hookName   string
		itemsCount int
	}{
		{
			name:       "Inventory type",
			invType:    0, // INVTYPE_INVENTORY
			hookName:   "packet_inventory",
			itemsCount: 2,
		},
		{
			name:       "Cart type",
			invType:    1, // INVTYPE_CART
			hookName:   "packet_cart",
			itemsCount: 2,
		},
		{
			name:       "Storage type",
			invType:    2, // INVTYPE_STORAGE
			hookName:   "packet_storage",
			itemsCount: 2,
		},
		{
			name:       "Guild storage type",
			invType:    3, // INVTYPE_GUILD_STORAGE
			hookName:   "packet_storage",
			itemsCount: 2,
		},
		{
			name:       "Unknown type",
			invType:    99,
			hookName:   "",
			itemsCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset hook call count
			hookCallCount := 0

			// Only register hook if we expect it to be called
			if tc.hookName != "" {
				hookManager.AddHook(tc.hookName, func(hookName string, data interface{}, userData interface{}) {
					hookCallCount++
				}, nil)
			}

			// Create sample items
			items := []map[string]interface{}{
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
					"identified": uint8(0),
				},
			}

			// Create test args
			args := map[string]interface{}{
				"type":     tc.invType,
				"itemInfo": items,
			}

			// Call handler
			err := inventoryManager.HandleItemListStackable(args)

			// Verify no error occurred
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			// Verify debug message was created
			if len(mockLogger.debugMessages) == 0 {
				t.Errorf("Expected debug message, got none")
			}

			// Verify hooks were called the expected number of times
			if tc.hookName != "" && hookCallCount != tc.itemsCount {
				t.Errorf("Expected %d hook calls for %s, got %d", tc.itemsCount, tc.hookName, hookCallCount)
			}
		})
	}

	// Test with missing type parameter
	args := map[string]interface{}{
		"itemInfo": []map[string]interface{}{},
	}

	err := inventoryManager.HandleItemListStackable(args)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for missing type parameter, got none")
	}

	// Test with missing itemInfo parameter
	args = map[string]interface{}{
		"type": uint8(0),
	}

	err = inventoryManager.HandleItemListStackable(args)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for missing itemInfo parameter, got none")
	}
}

// TestItemListEnd tests the HandleItemListEnd method
func TestItemListEnd(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create inventory manager
	inventoryManager := NewInventoryManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("item_list_end", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test cases for different inventory types
	testCases := []struct {
		name        string
		invType     uint8
		invTypeName string
		expectHook  bool
	}{
		{
			name:        "Inventory type",
			invType:     0, // INVTYPE_INVENTORY
			invTypeName: "inventory",
			expectHook:  true,
		},
		{
			name:        "Cart type",
			invType:     1, // INVTYPE_CART
			invTypeName: "cart",
			expectHook:  true,
		},
		{
			name:        "Storage type",
			invType:     2, // INVTYPE_STORAGE
			invTypeName: "storage",
			expectHook:  true,
		},
		{
			name:        "Guild storage type",
			invType:     3, // INVTYPE_GUILD_STORAGE
			invTypeName: "guild storage",
			expectHook:  true,
		},
		{
			name:        "Unknown type",
			invType:     99,
			invTypeName: "unknown",
			expectHook:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset hook called flag
			hookCalled = false

			// Create test args
			args := map[string]interface{}{
				"type": tc.invType,
			}

			// Call handler
			err := inventoryManager.HandleItemListEnd(args)

			// Verify no error occurred
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			// Verify debug message was created
			if len(mockLogger.infoMessages) == 0 {
				t.Errorf("Expected info message, got none")
			}

			// Verify hook was called if expected
			if tc.expectHook && !hookCalled {
				t.Errorf("Expected item_list_end hook to be called")
			}
		})
	}

	// Test with missing type parameter
	args := map[string]interface{}{}

	err := inventoryManager.HandleItemListEnd(args)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for missing type parameter, got none")
	}
}

// TestItemListNonstackable tests the HandleItemListNonstackable method
func TestItemListNonstackable(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create inventory manager
	inventoryManager := NewInventoryManager(mockParser, hookManager, mockLogger)

	// Test cases for different inventory types
	testCases := []struct {
		name       string
		invType    uint8
		hookName   string
		itemsCount int
	}{
		{
			name:       "Inventory type",
			invType:    0, // INVTYPE_INVENTORY
			hookName:   "packet_inventory",
			itemsCount: 2,
		},
		{
			name:       "Cart type",
			invType:    1, // INVTYPE_CART
			hookName:   "packet_cart",
			itemsCount: 2,
		},
		{
			name:       "Storage type",
			invType:    2, // INVTYPE_STORAGE
			hookName:   "packet_storage",
			itemsCount: 2,
		},
		{
			name:       "Guild storage type",
			invType:    3, // INVTYPE_GUILD_STORAGE
			hookName:   "packet_storage",
			itemsCount: 2,
		},
		{
			name:       "Unknown type",
			invType:    99,
			hookName:   "",
			itemsCount: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset hook call count
			hookCallCount := 0

			// Only register hook if we expect it to be called
			if tc.hookName != "" {
				hookManager.AddHook(tc.hookName, func(hookName string, data interface{}, userData interface{}) {
					hookCallCount++
				}, nil)
			}

			// Create sample items with equipment info
			items := []map[string]interface{}{
				{
					"ID":         "12345",
					"nameID":     uint16(501),
					"amount":     uint16(1),
					"type":       uint8(4),
					"identified": uint8(1),
					"broken":     uint8(0),
					"upgrade":    uint8(7),
					"cards":      "0,0,0,0",
					"expire":     uint32(0),
					"equipped":   uint32(32), // Equipped as body armor
				},
				{
					"ID":         "67890",
					"nameID":     uint16(502),
					"amount":     uint16(1),
					"type":       uint8(4),
					"identified": uint8(0),
					"broken":     uint8(1),
					"upgrade":    uint8(0),
					"cards":      "0,0,0,0",
					"expire":     uint32(0),
					"equipped":   uint32(0), // Not equipped
				},
			}

			// Create test args
			args := map[string]interface{}{
				"type":     tc.invType,
				"itemInfo": items,
			}

			// Call handler
			err := inventoryManager.HandleItemListNonstackable(args)

			// Verify no error occurred
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			// Verify debug message was created
			if len(mockLogger.debugMessages) == 0 {
				t.Errorf("Expected debug message, got none")
			}

			// Verify hooks were called the expected number of times
			if tc.hookName != "" && hookCallCount != tc.itemsCount {
				t.Errorf("Expected %d hook calls for %s, got %d", tc.itemsCount, tc.hookName, hookCallCount)
			}
		})
	}

	// Test with missing type parameter
	args := map[string]interface{}{
		"itemInfo": []map[string]interface{}{},
	}

	err := inventoryManager.HandleItemListNonstackable(args)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for missing type parameter, got none")
	}

	// Test with missing itemInfo parameter
	args = map[string]interface{}{
		"type": uint8(0),
	}

	err = inventoryManager.HandleItemListNonstackable(args)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for missing itemInfo parameter, got none")
	}
}
