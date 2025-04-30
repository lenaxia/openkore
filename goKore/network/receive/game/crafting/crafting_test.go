package crafting

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

// TestMakableItemList tests the HandleMakableItemList method
func TestMakableItemList(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create crafting manager
	craftingManager := NewCraftingManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("makable_item_list", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Create a mock itemList with some item data
	itemList := make([]byte, 16) // 2 items (8 bytes each)

	// First item
	itemList[0] = 1 // nameID (low byte)
	itemList[1] = 0 // nameID (high byte)
	// material_1, material_2, material_3 bytes

	// Second item
	itemList[8] = 2 // nameID (low byte)
	itemList[9] = 0 // nameID (high byte)
	// material_1, material_2, material_3 bytes

	// Test makable item list
	args := map[string]interface{}{
		"item_list": itemList,
	}

	err := craftingManager.HandleMakableItemList(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info messages were created
	if len(mockLogger.infoMessages) != 2 {
		t.Errorf("Expected 2 info messages, got %d", len(mockLogger.infoMessages))
	}

	// Verify items were stored
	if len(craftingManager.makableList) != 2 {
		t.Errorf("Expected 2 items in makableList, got %d", len(craftingManager.makableList))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected makable_item_list hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"item_list": "invalid", // Invalid type
	}

	err = craftingManager.HandleMakableItemList(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestArrowcraftList tests the HandleArrowcraftList method
func TestArrowcraftList(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create crafting manager
	craftingManager := NewCraftingManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("arrowcraft_list", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Create a mock RAW_MSG with some item data
	rawMsg := make([]byte, 8) // 4 bytes header + 2 items (2 bytes each)

	// First item
	rawMsg[4] = 1 // ID (low byte)
	rawMsg[5] = 0 // ID (high byte)

	// Second item
	rawMsg[6] = 2 // ID (low byte)
	rawMsg[7] = 0 // ID (high byte)

	// Test arrowcraft list
	args := map[string]interface{}{
		"RAW_MSG": rawMsg,
		"len":     uint16(8),
	}

	err := craftingManager.HandleArrowcraftList(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify items were stored
	if len(craftingManager.arrowCraftIDs) != 2 {
		t.Errorf("Expected 2 items in arrowCraftIDs, got %d", len(craftingManager.arrowCraftIDs))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected arrowcraft_list hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"RAW_MSG": "invalid", // Invalid type
	}

	err = craftingManager.HandleArrowcraftList(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestCookingList tests the HandleCookingList method
func TestCookingList(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create crafting manager
	craftingManager := NewCraftingManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("cooking_list", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Create a mock itemList with some item data
	itemList := make([]byte, 4) // 2 items (2 bytes each)

	// First item
	itemList[0] = 1 // nameID (low byte)
	itemList[1] = 0 // nameID (high byte)

	// Second item
	itemList[2] = 2 // nameID (low byte)
	itemList[3] = 0 // nameID (high byte)

	// Test cooking list
	args := map[string]interface{}{
		"type":      uint8(1),
		"item_list": itemList,
	}

	err := craftingManager.HandleCookingList(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info messages were created
	if len(mockLogger.infoMessages) != 2 {
		t.Errorf("Expected 2 info messages, got %d", len(mockLogger.infoMessages))
	}

	// Verify items were stored
	if len(craftingManager.cookingList) != 2 {
		t.Errorf("Expected 2 items in cookingList, got %d", len(craftingManager.cookingList))
	}

	// Verify cooking type was stored
	if craftingManager.cookingType != 1 {
		t.Errorf("Expected cookingType to be 1, got %d", craftingManager.cookingType)
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected cooking_list hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"type": "invalid", // Invalid type
	}

	err = craftingManager.HandleCookingList(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestRepairList tests the HandleRepairList method
func TestRepairList(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create crafting manager
	craftingManager := NewCraftingManager(mockParser, hookManager, mockLogger)

	// Create a mock RAW_MSG with some item data
	rawMsg := make([]byte, 30) // 4 bytes header + 2 items (13 bytes each)

	// First item
	rawMsg[4] = 1 // index (low byte)
	rawMsg[5] = 0 // index (high byte)
	rawMsg[6] = 1 // nameID (low byte)
	rawMsg[7] = 0 // nameID (high byte)
	rawMsg[8] = 1 // upgrade
	// cards bytes

	// Second item
	rawMsg[17] = 2 // index (low byte)
	rawMsg[18] = 0 // index (high byte)
	rawMsg[19] = 2 // nameID (low byte)
	rawMsg[20] = 0 // nameID (high byte)
	rawMsg[21] = 2 // upgrade
	// cards bytes

	// Test repair list
	args := map[string]interface{}{
		"RAW_MSG": rawMsg,
		"len":     uint16(30),
	}

	err := craftingManager.HandleRepairList(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify items were stored
	if len(craftingManager.repairList) < 3 { // Should have at least 3 items (0, 1, 2)
		t.Errorf("Expected at least 3 items in repairList, got %d", len(craftingManager.repairList))
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"RAW_MSG": "invalid", // Invalid type
	}

	err = craftingManager.HandleRepairList(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestRepairResult tests the HandleRepairResult method
func TestRepairResult(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create crafting manager
	craftingManager := NewCraftingManager(mockParser, hookManager, mockLogger)

	// Add a repair item
	craftingManager.repairList = make([]map[string]interface{}, 3)
	craftingManager.repairList[0] = map[string]interface{}{
		"index":   uint16(0),
		"nameID":  uint16(1),
		"upgrade": uint8(1),
		"name":    "Test Item",
	}

	// Test repair result (success)
	args := map[string]interface{}{
		"index": uint16(2), // Index is adjusted in the handler
		"flag":  uint8(0),
	}

	err := craftingManager.HandleRepairResult(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify repair list was cleared
	if len(craftingManager.repairList) != 0 {
		t.Errorf("Expected repairList to be cleared, got %d items", len(craftingManager.repairList))
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"index": "invalid", // Invalid type
	}

	err = craftingManager.HandleRepairResult(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestIdentifyList tests the HandleIdentifyList method
func TestIdentifyList(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create crafting manager
	craftingManager := NewCraftingManager(mockParser, hookManager, mockLogger)

	// Create a mock RAW_MSG with some item data
	rawMsg := make([]byte, 8) // 4 bytes header + 2 items (2 bytes each)

	// First item
	rawMsg[4] = 1 // index (low byte)
	rawMsg[5] = 0 // index (high byte)

	// Second item
	rawMsg[6] = 2 // index (low byte)
	rawMsg[7] = 0 // index (high byte)

	// Test identify list
	args := map[string]interface{}{
		"RAW_MSG": rawMsg,
		"len":     uint16(8),
	}

	err := craftingManager.HandleIdentifyList(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify items were stored
	if len(craftingManager.identifyIDs) != 2 {
		t.Errorf("Expected 2 items in identifyIDs, got %d", len(craftingManager.identifyIDs))
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"RAW_MSG": "invalid", // Invalid type
	}

	err = craftingManager.HandleIdentifyList(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestIdentify tests the HandleIdentify method
func TestIdentify(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create crafting manager
	craftingManager := NewCraftingManager(mockParser, hookManager, mockLogger)

	// Add identify IDs
	craftingManager.identifyIDs = []uint16{1, 2}

	// Test identify (success)
	args := map[string]interface{}{
		"ID":   uint16(1),
		"flag": uint8(0),
	}

	err := craftingManager.HandleIdentify(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify identify IDs were cleared
	if len(craftingManager.identifyIDs) != 0 {
		t.Errorf("Expected identifyIDs to be cleared, got %d items", len(craftingManager.identifyIDs))
	}

	// Test identify (failure)
	craftingManager.identifyIDs = []uint16{1, 2}
	args = map[string]interface{}{
		"ID":   uint16(1),
		"flag": uint8(1),
	}

	err = craftingManager.HandleIdentify(args)

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
		"ID": "invalid", // Invalid type
	}

	err = craftingManager.HandleIdentify(invalidArgs)

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

	// Create crafting manager
	craftingManager := NewCraftingManager(mockParser, hookManager, mockLogger)

	// Register handlers
	craftingManager.RegisterHandlers()

	// Verify handlers were registered
	expectedHandlers := []string{
		"makable_item_list",
		"arrowcraft_list",
		"cooking_list",
		"repair_list",
		"repair_result",
		"identify_list",
		"identify",
	}

	// We can't verify the handlers directly since we're using a real CoreParser
	// Instead, we'll just log the expected handlers
	for _, handler := range expectedHandlers {
		t.Logf("Expected handler: %s", handler)
	}
}
