package item

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// MockParser implements the Parser interface for testing
type MockParser struct {
	handlers map[string]core.PacketHandler
}

func NewMockParser() *MockParser {
	return &MockParser{
		handlers: make(map[string]core.PacketHandler),
	}
}

func (m *MockParser) RegisterHandler(packetID, name, format string, paramNames []string, handler core.PacketHandler) {
	m.handlers[name] = handler
}

// MockLogger implements a simple logger for testing
type MockLogger struct {
	debugMessages   []string
	infoMessages    []string
	warningMessages []string
	errorMessages   []string
	successMessages []string
}

func NewMockLogger() *MockLogger {
	return &MockLogger{
		debugMessages:   make([]string, 0),
		infoMessages:    make([]string, 0),
		warningMessages: make([]string, 0),
		errorMessages:   make([]string, 0),
		successMessages: make([]string, 0),
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
	m.successMessages = append(m.successMessages, format)
}

// TestInventoryItemAdded tests the HandleInventoryItemAdded method
func TestInventoryItemAdded(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create inventory manager
	inventoryManager := NewInventoryManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("item_gathered", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test successful item addition
	args := map[string]interface{}{
		"ID":         "12345",
		"amount":     uint16(10),
		"fail":       uint8(0),
		"nameID":     uint16(501),
		"type":       uint8(3),
		"type_equip": uint16(0),
		"identified": uint8(1),
		"broken":     uint8(0),
		"upgrade":    uint8(0),
		"cards":      "0,0,0,0",
		"options":    "",
	}

	err := inventoryManager.HandleInventoryItemAdded(args)

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
		t.Errorf("Expected item_gathered hook to be called")
	}

	// Test failed item addition (inventory full)
	failArgs := map[string]interface{}{
		"ID":     "12345",
		"amount": uint16(10),
		"fail":   uint8(2),
	}

	err = inventoryManager.HandleInventoryItemAdded(failArgs)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 2 {
		t.Errorf("Expected 2 info messages, got %d", len(mockLogger.infoMessages))
	}
}

// TestInventoryItemRemoved tests the HandleInventoryItemRemoved method
func TestInventoryItemRemoved(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create inventory manager
	inventoryManager := NewInventoryManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("packet_item_removed", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test item removal
	args := map[string]interface{}{
		"ID":     "12345",
		"amount": uint16(10),
		"reason": uint8(1), // Used to cast a skill
	}

	err := inventoryManager.HandleInventoryItemRemoved(args)

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
		t.Errorf("Expected packet_item_removed hook to be called")
	}
}

// TestInventoryItemsStackable tests the HandleInventoryItemsStackable method
func TestInventoryItemsStackable(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create inventory manager
	inventoryManager := NewInventoryManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCallCount := 0
	hookManager.AddHook("packet_inventory", func(hookName string, data interface{}, userData interface{}) {
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

	err := inventoryManager.HandleInventoryItemsStackable(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify hooks were called for each item
	if hookCallCount != 2 {
		t.Errorf("Expected 2 hook calls, got %d", hookCallCount)
	}
}

// TestInventoryItemFavorite tests the HandleInventoryItemFavorite method
func TestInventoryItemFavorite(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create inventory manager
	inventoryManager := NewInventoryManager(mockParser, hookManager, mockLogger)

	// Test mark as favorite
	args := map[string]interface{}{
		"ID":   "12345",
		"flag": uint8(1),
	}

	err := inventoryManager.HandleInventoryItemFavorite(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Test unmark as favorite
	args = map[string]interface{}{
		"ID":   "12345",
		"flag": uint8(0),
	}

	err = inventoryManager.HandleInventoryItemFavorite(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 2 {
		t.Errorf("Expected 2 info messages, got %d", len(mockLogger.infoMessages))
	}
}

// TestInventoryExpansionResult tests the HandleInventoryExpansionResult method
func TestInventoryExpansionResult(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create inventory manager
	inventoryManager := NewInventoryManager(mockParser, hookManager, mockLogger)

	// Test successful expansion
	args := map[string]interface{}{
		"result": uint8(0), // EXPAND_INVENTORY_RESULT_SUCCESS
	}

	err := inventoryManager.HandleInventoryExpansionResult(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Test failed expansion
	args = map[string]interface{}{
		"result": uint8(1), // EXPAND_INVENTORY_RESULT_FAILED
	}

	err = inventoryManager.HandleInventoryExpansionResult(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 2 {
		t.Errorf("Expected 2 info messages, got %d", len(mockLogger.infoMessages))
	}
}

// TestRegisterInventoryHandlers tests the RegisterHandlers method for inventory
func TestRegisterInventoryHandlers(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create inventory manager
	inventoryManager := NewInventoryManager(mockParser, hookManager, mockLogger)

	// Register handlers
	inventoryManager.RegisterHandlers()

	// Verify handlers were registered
	expectedHandlers := []string{
		"inventory_item_added",
		"inventory_item_removed",
		"inventory_items_stackable",
		"inventory_items_nonstackable",
		"inventory_item_favorite",
		"inventory_expansion_result",
	}

	for _, handler := range expectedHandlers {
		if _, exists := mockParser.handlers[handler]; !exists {
			t.Errorf("Expected handler %s to be registered", handler)
		}
	}
}
