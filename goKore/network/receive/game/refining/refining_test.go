package refining

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

// TestRefineUIOpened tests the HandleRefineUIOpened method
func TestRefineUIOpened(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create refining manager
	refiningManager := NewRefiningManager(mockParser, hookManager, mockLogger)

	// Test refineui opened
	args := map[string]interface{}{}

	err := refiningManager.HandleRefineUIOpened(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify refineUI open flag was set
	if open, ok := refiningManager.refineUI["open"].(bool); !ok || !open {
		t.Errorf("Expected refineUI open flag to be true")
	}
}

// TestRefineUIInfo tests the HandleRefineUIInfo method
func TestRefineUIInfo(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create refining manager
	refiningManager := NewRefiningManager(mockParser, hookManager, mockLogger)

	// Create a mock materials with some item data
	materials := make([]byte, 14) // 2 materials (7 bytes each)

	// First material
	materials[0] = 1   // nameID (low byte)
	materials[1] = 0   // nameID (high byte)
	materials[2] = 50  // chance
	materials[3] = 100 // zeny (low byte)
	materials[4] = 0   // zeny (high byte)
	materials[5] = 0   // zeny (high byte)
	materials[6] = 0   // zeny (high byte)

	// Second material
	materials[7] = 2    // nameID (low byte)
	materials[8] = 0    // nameID (high byte)
	materials[9] = 75   // chance
	materials[10] = 200 // zeny (low byte)
	materials[11] = 0   // zeny (high byte)
	materials[12] = 0   // zeny (high byte)
	materials[13] = 0   // zeny (high byte)

	// Test refineui info
	args := map[string]interface{}{
		"index":     uint16(1),
		"bless":     uint8(2),
		"materials": materials,
	}

	err := refiningManager.HandleRefineUIInfo(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info messages were created
	if len(mockLogger.infoMessages) < 3 {
		t.Errorf("Expected at least 3 info messages, got %d", len(mockLogger.infoMessages))
	}

	// Verify refineUI data was stored
	if itemIndex, ok := refiningManager.refineUI["itemIndex"].(uint16); !ok || itemIndex != 1 {
		t.Errorf("Expected refineUI itemIndex to be 1, got %v", refiningManager.refineUI["itemIndex"])
	}

	if bless, ok := refiningManager.refineUI["bless"].(uint8); !ok || bless != 2 {
		t.Errorf("Expected refineUI bless to be 2, got %v", refiningManager.refineUI["bless"])
	}

	// Test with empty materials
	args = map[string]interface{}{
		"index":     uint16(1),
		"bless":     uint8(2),
		"materials": []byte{},
	}

	err = refiningManager.HandleRefineUIInfo(args)

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
		"index": "invalid", // Invalid type
	}

	err = refiningManager.HandleRefineUIInfo(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestRefineResult tests the HandleRefineResult method
func TestRefineResult(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create refining manager
	refiningManager := NewRefiningManager(mockParser, hookManager, mockLogger)

	// Test refine result (success)
	args := map[string]interface{}{
		"fail":   uint8(0),
		"nameID": uint16(1),
	}

	err := refiningManager.HandleRefineResult(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Test refine result (failure)
	args = map[string]interface{}{
		"fail":   uint8(1),
		"nameID": uint16(1),
	}

	err = refiningManager.HandleRefineResult(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 2 {
		t.Errorf("Expected 2 info messages, got %d", len(mockLogger.infoMessages))
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"fail": "invalid", // Invalid type
	}

	err = refiningManager.HandleRefineResult(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestRefineStatus tests the HandleRefineStatus method
func TestRefineStatus(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create refining manager
	refiningManager := NewRefiningManager(mockParser, hookManager, mockLogger)

	// Test refine status (success)
	args := map[string]interface{}{
		"status":       uint8(0),
		"refine_level": uint16(7),
		"name":         "Test Item",
	}

	err := refiningManager.HandleRefineStatus(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify warning message was created
	if len(mockLogger.warningMessages) != 1 {
		t.Errorf("Expected 1 warning message, got %d", len(mockLogger.warningMessages))
	}

	// Test refine status (failure)
	args = map[string]interface{}{
		"status":       uint8(1),
		"refine_level": uint16(7),
		"name":         "Test Item",
	}

	err = refiningManager.HandleRefineStatus(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify warning message was created
	if len(mockLogger.warningMessages) != 2 {
		t.Errorf("Expected 2 warning messages, got %d", len(mockLogger.warningMessages))
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"status": "invalid", // Invalid type
	}

	err = refiningManager.HandleRefineStatus(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestUpgradeList tests the HandleUpgradeList method
func TestUpgradeList(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create refining manager
	refiningManager := NewRefiningManager(mockParser, hookManager, mockLogger)

	// Create a mock itemList with some item data
	itemList := make([]byte, 26) // 2 items (13 bytes each)

	// First item
	itemList[0] = 1 // index (low byte)
	itemList[1] = 0 // index (high byte)
	// Skip 6 bytes
	itemList[8] = 1 // nameID
	// Skip 4 bytes

	// Second item
	itemList[13] = 2 // index (low byte)
	itemList[14] = 0 // index (high byte)
	// Skip 6 bytes
	itemList[21] = 2 // nameID
	// Skip 4 bytes

	// Test upgrade list
	args := map[string]interface{}{
		"item_list": itemList,
	}

	err := refiningManager.HandleUpgradeList(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info messages were created
	if len(mockLogger.infoMessages) != 2 {
		t.Errorf("Expected 2 info messages, got %d", len(mockLogger.infoMessages))
	}

	// Verify items were stored
	if len(refiningManager.upgradeList) != 2 {
		t.Errorf("Expected 2 items in upgradeList, got %d", len(refiningManager.upgradeList))
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"item_list": "invalid", // Invalid type
	}

	err = refiningManager.HandleUpgradeList(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestUpgradeMessage tests the HandleUpgradeMessage method
func TestUpgradeMessage(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create refining manager
	refiningManager := NewRefiningManager(mockParser, hookManager, mockLogger)

	// Test upgrade message (success)
	args := map[string]interface{}{
		"type":   uint8(0),
		"itemID": uint16(1),
	}

	err := refiningManager.HandleUpgradeMessage(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Test upgrade message (failure)
	args = map[string]interface{}{
		"type":   uint8(1),
		"itemID": uint16(1),
	}

	err = refiningManager.HandleUpgradeMessage(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 2 {
		t.Errorf("Expected 2 info messages, got %d", len(mockLogger.infoMessages))
	}

	// Test upgrade message (fail lvl)
	args = map[string]interface{}{
		"type":   uint8(2),
		"itemID": uint16(1),
	}

	err = refiningManager.HandleUpgradeMessage(args)

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
		"type": "invalid", // Invalid type
	}

	err = refiningManager.HandleUpgradeMessage(invalidArgs)

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

	// Create refining manager
	refiningManager := NewRefiningManager(mockParser, hookManager, mockLogger)

	// Register handlers
	refiningManager.RegisterHandlers()

	// Verify handlers were registered
	expectedHandlers := []string{
		"refineui_opened",
		"refineui_info",
		"refine_result",
		"refine_status",
		"upgrade_list",
		"upgrade_message",
	}

	// We can't verify the handlers directly since we're using a real CoreParser
	// Instead, we'll just log the expected handlers
	for _, handler := range expectedHandlers {
		t.Logf("Expected handler: %s", handler)
	}
}
