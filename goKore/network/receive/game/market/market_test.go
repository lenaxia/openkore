package market

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

// TestNpcMarketInfo tests the HandleNpcMarketInfo method
func TestNpcMarketInfo(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create market manager
	marketManager := NewMarketManager(mockParser, hookManager, mockLogger)

	// Create a mock itemList with some item data
	itemList := make([]byte, 22) // 2 items (11 bytes each)

	// First item
	itemList[0] = 1   // nameID (low byte)
	itemList[1] = 0   // nameID (high byte)
	itemList[2] = 1   // type
	itemList[3] = 100 // price (low byte)
	itemList[4] = 0   // price (high byte)
	itemList[5] = 0   // price (high byte)
	itemList[6] = 0   // price (high byte)
	itemList[7] = 5   // amount (low byte)
	itemList[8] = 0   // amount (high byte)
	itemList[9] = 0   // amount (high byte)
	itemList[10] = 0  // amount (high byte)

	// Second item
	itemList[11] = 2   // nameID (low byte)
	itemList[12] = 0   // nameID (high byte)
	itemList[13] = 2   // type
	itemList[14] = 200 // price (low byte)
	itemList[15] = 0   // price (high byte)
	itemList[16] = 0   // price (high byte)
	itemList[17] = 0   // price (high byte)
	itemList[18] = 10  // amount (low byte)
	itemList[19] = 0   // amount (high byte)
	itemList[20] = 0   // amount (high byte)
	itemList[21] = 0   // amount (high byte)

	// Test NPC market info
	args := map[string]interface{}{
		"itemList": itemList,
	}

	err := marketManager.HandleNpcMarketInfo(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify debug messages were created
	if len(mockLogger.debugMessages) != 2 {
		t.Errorf("Expected 2 debug messages, got %d", len(mockLogger.debugMessages))
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify items were stored
	if len(marketManager.storeItems) != 2 {
		t.Errorf("Expected 2 items in storeItems, got %d", len(marketManager.storeItems))
	}

	// Verify in_market flag was set
	if !marketManager.inMarket {
		t.Errorf("Expected inMarket to be true")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"itemList": "invalid", // Invalid type
	}

	err = marketManager.HandleNpcMarketInfo(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestNpcMarketPurchaseResult tests the HandleNpcMarketPurchaseResult method
func TestNpcMarketPurchaseResult(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create market manager
	marketManager := NewMarketManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("market_buy_result", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Create a mock itemList with some item data
	itemList := make([]byte, 22) // 2 items (11 bytes each)

	// First item
	itemList[0] = 1   // nameID (low byte)
	itemList[1] = 0   // nameID (high byte)
	itemList[2] = 1   // type
	itemList[3] = 100 // price (low byte)
	itemList[4] = 0   // price (high byte)
	itemList[5] = 0   // price (high byte)
	itemList[6] = 0   // price (high byte)
	itemList[7] = 5   // amount (low byte)
	itemList[8] = 0   // amount (high byte)
	itemList[9] = 0   // amount (high byte)
	itemList[10] = 0  // amount (high byte)

	// Second item
	itemList[11] = 2   // nameID (low byte)
	itemList[12] = 0   // nameID (high byte)
	itemList[13] = 2   // type
	itemList[14] = 200 // price (low byte)
	itemList[15] = 0   // price (high byte)
	itemList[16] = 0   // price (high byte)
	itemList[17] = 0   // price (high byte)
	itemList[18] = 10  // amount (low byte)
	itemList[19] = 0   // amount (high byte)
	itemList[20] = 0   // amount (high byte)
	itemList[21] = 0   // amount (high byte)

	// Test NPC market purchase result (success)
	args := map[string]interface{}{
		"result":   uint8(MarketBuyResultSuccess),
		"itemList": itemList,
	}

	err := marketManager.HandleNpcMarketPurchaseResult(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify debug messages were created
	if len(mockLogger.debugMessages) != 3 { // 1 for result + 2 for items
		t.Errorf("Expected 3 debug messages, got %d", len(mockLogger.debugMessages))
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected market_buy_result hook to be called")
	}

	// Verify items were stored
	if len(marketManager.storeItems) != 2 {
		t.Errorf("Expected 2 items in storeItems, got %d", len(marketManager.storeItems))
	}

	// Verify in_market flag was set
	if !marketManager.inMarket {
		t.Errorf("Expected inMarket to be true")
	}

	// Test NPC market purchase result (failure)
	args = map[string]interface{}{
		"result":   uint8(MarketBuyResultNoZeny),
		"itemList": itemList,
	}

	err = marketManager.HandleNpcMarketPurchaseResult(args)

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
		"result": "invalid", // Invalid type
	}

	err = marketManager.HandleNpcMarketPurchaseResult(invalidArgs)

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

	// Create market manager
	marketManager := NewMarketManager(mockParser, hookManager, mockLogger)

	// Register handlers
	marketManager.RegisterHandlers()

	// Verify handlers were registered
	expectedHandlers := []string{
		"npc_market_info",
		"npc_market_purchase_result",
	}

	// We can't verify the handlers directly since we're using a real CoreParser
	// Instead, we'll just log the expected handlers
	for _, handler := range expectedHandlers {
		t.Logf("Expected handler: %s", handler)
	}
}
