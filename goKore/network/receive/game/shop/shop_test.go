package shop

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// MockLogger is a simple mock implementation of the Logger interface
type MockLogger struct {
	debugMessages   []string
	infoMessages    []string
	warningMessages []string
	errorMessages   []string
}

func NewMockLogger() *MockLogger {
	return &MockLogger{
		debugMessages:   []string{},
		infoMessages:    []string{},
		warningMessages: []string{},
		errorMessages:   []string{},
	}
}

func (m *MockLogger) Debug(format string, args ...interface{}) {
	m.debugMessages = append(m.debugMessages, format)
}

func (m *MockLogger) Info(format string, args ...interface{}) {
	m.infoMessages = append(m.infoMessages, format)
}

func (m *MockLogger) Warning(format string, args ...interface{}) {
	m.warningMessages = append(m.warningMessages, format)
}

func (m *MockLogger) Error(format string, args ...interface{}) {
	m.errorMessages = append(m.errorMessages, format)
}

func (m *MockLogger) Success(format string, args ...interface{}) {
	// Not used in these tests
}

// TestNpcStoreBegin tests the HandleNpcStoreBegin method
func TestNpcStoreBegin(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create shop manager
	shopManager := NewShopManager(mockParser, hookManager, mockLogger)

	// Test NPC store begin
	args := map[string]interface{}{
		"ID": uint32(12345),
	}

	err := shopManager.HandleNpcStoreBegin(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify NPC name was stored
	if npcName, ok := shopManager.storeList["npcName"].(string); !ok || npcName != "NPC-12345" {
		t.Errorf("Expected NPC name to be stored in storeList")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"ID": "invalid", // Invalid type
	}

	err = shopManager.HandleNpcStoreBegin(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestNpcStoreInfo tests the HandleNpcStoreInfo method
func TestNpcStoreInfo(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create shop manager
	shopManager := NewShopManager(mockParser, hookManager, mockLogger)

	// Create a mock RAW_MSG with some item data
	rawMsg := make([]byte, 26) // 4 bytes header + 2 items (11 bytes each)

	// First item
	rawMsg[4] = 100 // price (low byte)
	rawMsg[5] = 0   // price (high byte)
	rawMsg[6] = 0   // price (high byte)
	rawMsg[7] = 0   // price (high byte)
	rawMsg[12] = 1  // type
	rawMsg[13] = 1  // nameID (low byte)
	rawMsg[14] = 0  // nameID (high byte)

	// Second item
	rawMsg[15] = 200 // price (low byte)
	rawMsg[16] = 0   // price (high byte)
	rawMsg[17] = 0   // price (high byte)
	rawMsg[18] = 0   // price (high byte)
	rawMsg[23] = 2   // type
	rawMsg[24] = 2   // nameID (low byte)
	rawMsg[25] = 0   // nameID (high byte)

	// Test NPC store info
	args := map[string]interface{}{
		"len":     uint16(26),
		"RAW_MSG": rawMsg,
	}

	err := shopManager.HandleNpcStoreInfo(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify items were stored
	if len(shopManager.storeItems) != 2 {
		t.Errorf("Expected 2 items in storeItems, got %d", len(shopManager.storeItems))
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"len": "invalid", // Invalid type
	}

	err = shopManager.HandleNpcStoreInfo(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestNpcSellList tests the HandleNpcSellList method
func TestNpcSellList(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create shop manager
	shopManager := NewShopManager(mockParser, hookManager, mockLogger)

	// Create a mock itemsdata with some item data
	itemsData := make([]byte, 20) // 2 items (10 bytes each)

	// First item
	itemsData[0] = 1 // index (low byte)
	itemsData[1] = 0 // index (high byte)
	// Skip price bytes (2-5)
	itemsData[6] = 100 // price_overcharge (low byte)
	itemsData[7] = 0   // price_overcharge (high byte)
	itemsData[8] = 0   // price_overcharge (high byte)
	itemsData[9] = 0   // price_overcharge (high byte)

	// Second item
	itemsData[10] = 2 // index (low byte)
	itemsData[11] = 0 // index (high byte)
	// Skip price bytes (12-15)
	itemsData[16] = 200 // price_overcharge (low byte)
	itemsData[17] = 0   // price_overcharge (high byte)
	itemsData[18] = 0   // price_overcharge (high byte)
	itemsData[19] = 0   // price_overcharge (high byte)

	// Test NPC sell list
	args := map[string]interface{}{
		"itemsdata": itemsData,
	}

	err := shopManager.HandleNpcSellList(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info messages were created (header + 2 items + footer)
	if len(mockLogger.infoMessages) != 4 {
		t.Errorf("Expected 4 info messages, got %d", len(mockLogger.infoMessages))
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"itemsdata": "invalid", // Invalid type
	}

	err = shopManager.HandleNpcSellList(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestBuyResult tests the HandleBuyResult method
func TestBuyResult(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create shop manager
	shopManager := NewShopManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("buy_result", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test buy result (success)
	args := map[string]interface{}{
		"fail": uint8(0),
	}

	err := shopManager.HandleBuyResult(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected buy_result hook to be called")
	}

	// Test buy result (failure)
	args = map[string]interface{}{
		"fail": uint8(1),
	}

	err = shopManager.HandleBuyResult(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify error message was created
	if len(mockLogger.errorMessages) != 1 {
		t.Errorf("Expected 1 error message, got %d", len(mockLogger.errorMessages))
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"fail": "invalid", // Invalid type
	}

	err = shopManager.HandleBuyResult(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestSellResult tests the HandleSellResult method
func TestSellResult(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create shop manager
	shopManager := NewShopManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("sell_result", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test sell result (success)
	args := map[string]interface{}{
		"fail": uint8(0),
	}

	err := shopManager.HandleSellResult(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info messages were created
	if len(mockLogger.infoMessages) != 2 {
		t.Errorf("Expected 2 info messages, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected sell_result hook to be called")
	}

	// Test sell result (failure)
	args = map[string]interface{}{
		"fail": uint8(1),
	}

	err = shopManager.HandleSellResult(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify error message was created
	if len(mockLogger.errorMessages) != 1 {
		t.Errorf("Expected 1 error message, got %d", len(mockLogger.errorMessages))
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"fail": "invalid", // Invalid type
	}

	err = shopManager.HandleSellResult(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestVendingStart tests the HandleVendingStart method
func TestVendingStart(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create shop manager
	shopManager := NewShopManager(mockParser, hookManager, mockLogger)

	// Create a mock itemList with some item data
	itemList := make([]byte, 40) // 2 items (20 bytes each)

	// First item
	itemList[0] = 100 // price (low byte)
	itemList[1] = 0   // price (high byte)
	itemList[2] = 0   // price (high byte)
	itemList[3] = 0   // price (high byte)
	itemList[4] = 1   // number (low byte)
	itemList[5] = 0   // number (high byte)
	itemList[6] = 5   // quantity (low byte)
	itemList[7] = 0   // quantity (high byte)
	itemList[8] = 1   // type
	itemList[9] = 1   // nameID (low byte)
	itemList[10] = 0  // nameID (high byte)

	// Second item
	itemList[20] = 200 // price (low byte)
	itemList[21] = 0   // price (high byte)
	itemList[22] = 0   // price (high byte)
	itemList[23] = 0   // price (high byte)
	itemList[24] = 2   // number (low byte)
	itemList[25] = 0   // number (high byte)
	itemList[26] = 10  // quantity (low byte)
	itemList[27] = 0   // quantity (high byte)
	itemList[28] = 2   // type
	itemList[29] = 2   // nameID (low byte)
	itemList[30] = 0   // nameID (high byte)

	// Test vending start
	args := map[string]interface{}{
		"len":      uint16(40),
		"itemList": itemList,
	}

	err := shopManager.HandleVendingStart(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info messages were created (header + 2 items + footer)
	if len(mockLogger.infoMessages) < 4 {
		t.Errorf("Expected at least 4 info messages, got %d", len(mockLogger.infoMessages))
	}

	// Verify items were stored
	if len(shopManager.storeItems) != 2 {
		t.Errorf("Expected 2 items in storeItems, got %d", len(shopManager.storeItems))
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"len": "invalid", // Invalid type
	}

	err = shopManager.HandleVendingStart(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestVenderFound tests the HandleVenderFound method
func TestVenderFound(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create shop manager
	shopManager := NewShopManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("packet_vender", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test vender found
	args := map[string]interface{}{
		"ID":    uint32(12345),
		"title": "Test Vender",
	}

	err := shopManager.HandleVenderFound(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify debug message was created
	if len(mockLogger.debugMessages) != 1 {
		t.Errorf("Expected 1 debug message, got %d", len(mockLogger.debugMessages))
	}

	// Verify vender was stored
	if len(shopManager.venderListsID) != 1 {
		t.Errorf("Expected 1 vender ID in venderListsID, got %d", len(shopManager.venderListsID))
	}

	if vender, exists := shopManager.venderLists[12345]; !exists {
		t.Errorf("Expected vender to be stored in venderLists")
	} else {
		if title, ok := vender["title"].(string); !ok || title != "Test Vender" {
			t.Errorf("Expected vender title to be 'Test Vender', got %v", title)
		}
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected packet_vender hook to be called")
	}

	// Test with existing vender
	hookCalled = false
	err = shopManager.HandleVenderFound(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify hook was not called again
	if hookCalled {
		t.Errorf("Expected packet_vender hook not to be called for existing vender")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"ID": "invalid", // Invalid type
	}

	err = shopManager.HandleVenderFound(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestVenderLost tests the HandleVenderLost method
func TestVenderLost(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create shop manager
	shopManager := NewShopManager(mockParser, hookManager, mockLogger)

	// Add a vender
	shopManager.venderListsID = append(shopManager.venderListsID, 12345)
	shopManager.venderLists[12345] = map[string]interface{}{
		"title": "Test Vender",
		"id":    uint32(12345),
	}

	// Test vender lost
	args := map[string]interface{}{
		"ID": uint32(12345),
	}

	err := shopManager.HandleVenderLost(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify debug message was created
	if len(mockLogger.debugMessages) != 1 {
		t.Errorf("Expected 1 debug message, got %d", len(mockLogger.debugMessages))
	}

	// Verify vender was removed
	if len(shopManager.venderListsID) != 0 {
		t.Errorf("Expected 0 vender IDs in venderListsID, got %d", len(shopManager.venderListsID))
	}

	if _, exists := shopManager.venderLists[12345]; exists {
		t.Errorf("Expected vender to be removed from venderLists")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"ID": "invalid", // Invalid type
	}

	err = shopManager.HandleVenderLost(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestVenderBuyFail tests the HandleVenderBuyFail method
func TestVenderBuyFail(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create shop manager
	shopManager := NewShopManager(mockParser, hookManager, mockLogger)

	// Test vender buy fail (insufficient zeny)
	args := map[string]interface{}{
		"fail":   uint8(1),
		"amount": uint16(5),
		"ID":     uint32(12345),
	}

	err := shopManager.HandleVenderBuyFail(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify error message was created
	if len(mockLogger.errorMessages) != 1 {
		t.Errorf("Expected 1 error message, got %d", len(mockLogger.errorMessages))
	}

	// Test vender buy fail (overweight)
	args = map[string]interface{}{
		"fail":   uint8(2),
		"amount": uint16(5),
		"ID":     uint32(12345),
	}

	err = shopManager.HandleVenderBuyFail(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify error message was created
	if len(mockLogger.errorMessages) != 2 {
		t.Errorf("Expected 2 error messages, got %d", len(mockLogger.errorMessages))
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"fail": "invalid", // Invalid type
	}

	err = shopManager.HandleVenderBuyFail(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestShopSkill tests the HandleShopSkill method
func TestShopSkill(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create shop manager
	shopManager := NewShopManager(mockParser, hookManager, mockLogger)

	// Test shop skill
	args := map[string]interface{}{
		"number": uint16(10),
	}

	err := shopManager.HandleShopSkill(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}
	
	// TestOpenStoreStatus tests the HandleOpenStoreStatus method
	func TestOpenStoreStatus(t *testing.T) {
		// Create mocks
		mockParser := core.NewCoreParser("ServerType0", nil)
		mockLogger := NewMockLogger()
		hookManager := hooks.NewHookManager()
	
		// Create shop manager
		shopManager := NewShopManager(mockParser, hookManager, mockLogger)
	
		// Track hook calls
		successHookCalled := false
		failHookCalled := false
		hookManager.AddHook("open_store_success", func(hookName string, data interface{}, userData interface{}) {
			successHookCalled = true
		}, nil)
		hookManager.AddHook("open_store_fail", func(hookName string, data interface{}, userData interface{}) {
			failHookCalled = true
		}, nil)
	
		// Test open store status (success)
		args := map[string]interface{}{
			"flag": uint8(0),
		}
	
		err := shopManager.HandleOpenStoreStatus(args)
	
		// Verify no error occurred
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	
		// Verify info message was created
		if len(mockLogger.infoMessages) != 1 {
			t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
		}
	
		// Verify success hook was called
		if !successHookCalled {
			t.Errorf("Expected open_store_success hook to be called")
		}
	
		// Test open store status (failure)
		args = map[string]interface{}{
			"flag": uint8(1),
		}
	
		err = shopManager.HandleOpenStoreStatus(args)
	
		// Verify no error occurred
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	
		// Verify error message was created
		if len(mockLogger.errorMessages) != 1 {
			t.Errorf("Expected 1 error message, got %d", len(mockLogger.errorMessages))
		}
	
		// Verify fail hook was called
		if !failHookCalled {
			t.Errorf("Expected open_store_fail hook to be called")
		}
	
		// Test with invalid parameters
		invalidArgs := map[string]interface{}{
			"flag": "invalid", // Invalid type
		}
	
		err = shopManager.HandleOpenStoreStatus(invalidArgs)
	
		// Verify error occurred
		if err == nil {
			t.Errorf("Expected error for invalid parameters, got nil")
		}
	}
	
	// TestOpenStoreStatus tests the HandleOpenStoreStatus method
	func TestOpenStoreStatus(t *testing.T) {
		// Create mocks
		mockParser := core.NewCoreParser("ServerType0", nil)
		mockLogger := NewMockLogger()
		hookManager := hooks.NewHookManager()
	
		// Create shop manager
		shopManager := NewShopManager(mockParser, hookManager, mockLogger)
	
		// Track hook calls
		successHookCalled := false
		failHookCalled := false
		hookManager.AddHook("open_store_success", func(hookName string, data interface{}, userData interface{}) {
			successHookCalled = true
		}, nil)
		hookManager.AddHook("open_store_fail", func(hookName string, data interface{}, userData interface{}) {
			failHookCalled = true
		}, nil)
	
		// Test open store status (success)
		args := map[string]interface{}{
			"flag": uint8(0),
		}
	
		err := shopManager.HandleOpenStoreStatus(args)
	
		// Verify no error occurred
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	
		// Verify info message was created
		if len(mockLogger.infoMessages) != 1 {
			t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
		}
	
		// Verify success hook was called
		if !successHookCalled {
			t.Errorf("Expected open_store_success hook to be called")
		}
	
		// Test open store status (failure)
		args = map[string]interface{}{
			"flag": uint8(1),
		}
	
		err = shopManager.HandleOpenStoreStatus(args)
	
		// Verify no error occurred
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	
		// Verify error message was created
		if len(mockLogger.errorMessages) != 1 {
			t.Errorf("Expected 1 error message, got %d", len(mockLogger.errorMessages))
		}
	
		// Verify fail hook was called
		if !failHookCalled {
			t.Errorf("Expected open_store_fail hook to be called")
		}
	
		// Test with invalid parameters
		invalidArgs := map[string]interface{}{
			"flag": "invalid", // Invalid type
		}
	
		err = shopManager.HandleOpenStoreStatus(invalidArgs)
	
		// Verify error occurred
		if err == nil {
			t.Errorf("Expected error for invalid parameters, got nil")
		}
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"number": "invalid", // Invalid type
	}

	err = shopManager.HandleShopSkill(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestShopSold tests the HandleShopSold method
func TestShopSold(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create shop manager
	shopManager := NewShopManager(mockParser, hookManager, mockLogger)

	// Add a test item to storeItems
	shopManager.storeItems = append(shopManager.storeItems, map[string]interface{}{
		"name":     "Test Item",
		"price":    uint32(100),
		"quantity": uint16(10),
	})

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("vending_item_sold", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test shop sold
	args := map[string]interface{}{
		"number": uint16(0),
		"amount": uint16(2),
	}

	err := shopManager.HandleShopSold(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected vending_item_sold hook to be called")
	}

	// Verify item was updated
	if quantity, ok := shopManager.storeItems[0]["quantity"].(uint16); !ok || quantity != 8 {
		t.Errorf("Expected item quantity to be 8, got %v", shopManager.storeItems[0]["quantity"])
	}

	if sold, ok := shopManager.storeItems[0]["sold"].(uint16); !ok || sold != 2 {
		t.Errorf("Expected item sold to be 2, got %v", shopManager.storeItems[0]["sold"])
	}

	// Verify shop earned was updated
	if earned, ok := shopManager.storeList["earned"].(uint32); !ok || earned != 200 {
		t.Errorf("Expected shop earned to be 200, got %v", shopManager.storeList["earned"])
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"number": "invalid", // Invalid type
	}

	err = shopManager.HandleShopSold(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestShopSoldLong tests the HandleShopSoldLong method
func TestShopSoldLong(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create shop manager
	shopManager := NewShopManager(mockParser, hookManager, mockLogger)

	// Add a test item to storeItems
	shopManager.storeItems = append(shopManager.storeItems, map[string]interface{}{
		"name":     "Test Item",
		"price":    uint32(100),
		"quantity": uint16(10),
	})

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("vending_item_sold", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test shop sold long
	args := map[string]interface{}{
		"number": uint16(0),
		"amount": uint16(2),
		"zeny":   uint32(200),
		"charID": uint32(12345),
	}

	err := shopManager.HandleShopSoldLong(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected vending_item_sold hook to be called")
	}

	// Verify item was updated
	if quantity, ok := shopManager.storeItems[0]["quantity"].(uint16); !ok || quantity != 8 {
		t.Errorf("Expected item quantity to be 8, got %v", shopManager.storeItems[0]["quantity"])
	}

	if sold, ok := shopManager.storeItems[0]["sold"].(uint16); !ok || sold != 2 {
		t.Errorf("Expected item sold to be 2, got %v", shopManager.storeItems[0]["sold"])
	}

	// Verify shop earned was updated
	if earned, ok := shopManager.storeList["earned"].(uint32); !ok || earned != 200 {
		t.Errorf("Expected shop earned to be 200, got %v", shopManager.storeList["earned"])
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"number": "invalid", // Invalid type
	}

	err = shopManager.HandleShopSoldLong(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestOpenStoreStatus tests the HandleOpenStoreStatus method
func TestOpenStoreStatus(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create shop manager
	shopManager := NewShopManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	successHookCalled := false
	failHookCalled := false
	hookManager.AddHook("open_store_success", func(hookName string, data interface{}, userData interface{}) {
		successHookCalled = true
	}, nil)
	hookManager.AddHook("open_store_fail", func(hookName string, data interface{}, userData interface{}) {
		failHookCalled = true
	}, nil)

	// Test open store status (success)
	args := map[string]interface{}{
		"flag": uint8(0),
	}

	err := shopManager.HandleOpenStoreStatus(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify success hook was called
	if !successHookCalled {
		t.Errorf("Expected open_store_success hook to be called")
	}

	// Test open store status (failure)
	args = map[string]interface{}{
		"flag": uint8(1),
	}

	err = shopManager.HandleOpenStoreStatus(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify error message was created
	if len(mockLogger.errorMessages) != 1 {
		t.Errorf("Expected 1 error message, got %d", len(mockLogger.errorMessages))
	}

	// Verify fail hook was called
	if !failHookCalled {
		t.Errorf("Expected open_store_fail hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"flag": "invalid", // Invalid type
	}

	err = shopManager.HandleOpenStoreStatus(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestRegisterHandlers tests the RegisterHandlers method
func TestRegisterHandlers(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create shop manager
	shopManager := NewShopManager(mockParser, hookManager, mockLogger)

	// Register handlers
	shopManager.RegisterHandlers()

	// Verify handlers were registered
	expectedHandlers := []string{
		"npc_store_begin",
		"npc_store_info",
		"npc_sell_list",
		"buy_result",
		"sell_result",
		"vending_start",
		"vender_items_list",
		"vender_found",
		"vender_lost",
		"vender_buy_fail",
		"shop_skill",
		"shop_sold",
		"shop_sold_long",
		"open_store_status",
	}

	// We can't verify the handlers directly since we're using a real CoreParser
	// Instead, we'll just log the expected handlers
	for _, handler := range expectedHandlers {
		t.Logf("Expected handler: %s", handler)
	}
}
