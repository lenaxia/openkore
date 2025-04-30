package shop

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// TestCashShopList tests the HandleCashShopList method
func TestCashShopList(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create cash shop manager
	cashShopManager := NewCashShopManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("cash_shop_list", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Create test item info
	// Each item is 6 bytes: 2 bytes for ID, 4 bytes for price
	itemInfo := []byte{
		// Item 1: ID = 501, Price = 100
		245, 1, 100, 0, 0, 0,
		// Item 2: ID = 502, Price = 200
		246, 1, 200, 0, 0, 0,
	}

	// Test cash shop list
	args := map[string]interface{}{
		"tabcode":  uint8(1), // Popular tab
		"itemInfo": itemInfo,
	}

	err := cashShopManager.HandleCashShopList(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info messages were created
	if len(mockLogger.infoMessages) < 3 {
		t.Errorf("Expected at least 3 info messages, got %d", len(mockLogger.infoMessages))
	}

	// Verify items were added to cash shop list
	tabItems, ok := cashShopManager.cashShopList["1"].([]map[string]interface{})
	if !ok {
		t.Errorf("Expected tab items to be a slice of maps")
	} else if len(tabItems) != 2 {
		t.Errorf("Expected 2 items in tab, got %d", len(tabItems))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected cash_shop_list hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"tabcode":  "invalid", // Invalid type
		"itemInfo": itemInfo,
	}

	err = cashShopManager.HandleCashShopList(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestCashShopOpenResult tests the HandleCashShopOpenResult method
func TestCashShopOpenResult(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create cash shop manager
	cashShopManager := NewCashShopManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("cash_shop_open_result", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test cash shop open result
	args := map[string]interface{}{
		"cash_points":  uint32(1000),
		"kafra_points": uint32(500),
	}

	err := cashShopManager.HandleCashShopOpenResult(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify points were stored
	if cashShopManager.cashPoints["cash"] != 1000 {
		t.Errorf("Expected cash points to be 1000, got %d", cashShopManager.cashPoints["cash"])
	}
	if cashShopManager.cashPoints["kafra"] != 500 {
		t.Errorf("Expected kafra points to be 500, got %d", cashShopManager.cashPoints["kafra"])
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected cash_shop_open_result hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"cash_points":  "invalid", // Invalid type
		"kafra_points": uint32(500),
	}

	err = cashShopManager.HandleCashShopOpenResult(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestCashShopBuyResult tests the HandleCashShopBuyResult method
func TestCashShopBuyResult(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create cash shop manager
	cashShopManager := NewCashShopManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("cash_shop_buy_result", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test successful purchase
	args := map[string]interface{}{
		"result":         uint8(0), // Success
		"item_id":        uint32(501),
		"updated_points": uint32(900),
	}

	err := cashShopManager.HandleCashShopBuyResult(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify cash points were updated
	if cashShopManager.cashPoints["cash"] != 900 {
		t.Errorf("Expected cash points to be 900, got %d", cashShopManager.cashPoints["cash"])
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected cash_shop_buy_result hook to be called")
	}

	// Test failed purchase
	hookCalled = false
	mockLogger = NewMockLogger()
	cashShopManager.logger = mockLogger

	args = map[string]interface{}{
		"result":         uint8(2), // Shortage cash
		"item_id":        uint32(502),
		"updated_points": uint32(900),
	}

	err = cashShopManager.HandleCashShopBuyResult(args)

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
		t.Errorf("Expected cash_shop_buy_result hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"result":         "invalid", // Invalid type
		"item_id":        uint32(501),
		"updated_points": uint32(900),
	}

	err = cashShopManager.HandleCashShopBuyResult(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestMergeItemOpen tests the HandleMergeItemOpen method
func TestMergeItemOpen(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create cash shop manager
	cashShopManager := NewCashShopManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("merge_item_open", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Create test item list
	// Each item ID is 2 bytes
	itemList := []byte{
		// Item 1: ID = 501
		245, 1,
		// Item 2: ID = 502
		246, 1,
	}

	// Test merge item open
	args := map[string]interface{}{
		"itemList": itemList,
	}

	err := cashShopManager.HandleMergeItemOpen(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify debug message was created
	if len(mockLogger.debugMessages) != 1 {
		t.Errorf("Expected 1 debug message, got %d", len(mockLogger.debugMessages))
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected merge_item_open hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"itemList": "invalid", // Invalid type
	}

	err = cashShopManager.HandleMergeItemOpen(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestMergeItemResult tests the HandleMergeItemResult method
func TestMergeItemResult(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create cash shop manager
	cashShopManager := NewCashShopManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("merge_item_result", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test successful merge
	args := map[string]interface{}{
		"itemIndex": uint16(1),
		"total":     uint16(10),
		"result":    uint8(0), // Success
	}

	err := cashShopManager.HandleMergeItemResult(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info messages were created
	if len(mockLogger.infoMessages) != 2 {
		t.Errorf("Expected 2 info messages, got %d", len(mockLogger.infoMessages))
	}

	// Verify debug message was created
	if len(mockLogger.debugMessages) != 1 {
		t.Errorf("Expected 1 debug message, got %d", len(mockLogger.debugMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected merge_item_result hook to be called")
	}

	// Test failed merge
	hookCalled = false
	mockLogger = NewMockLogger()
	cashShopManager.logger = mockLogger

	args = map[string]interface{}{
		"itemIndex": uint16(1),
		"total":     uint16(10),
		"result":    uint8(1), // Cannot merge
	}

	err = cashShopManager.HandleMergeItemResult(args)

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
		t.Errorf("Expected merge_item_result hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"itemIndex": "invalid", // Invalid type
		"total":     uint16(10),
		"result":    uint8(0),
	}

	err = cashShopManager.HandleMergeItemResult(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}
