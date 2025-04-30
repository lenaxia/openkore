package auction

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
	successMessages []string
}

func NewMockLogger() *MockLogger {
	return &MockLogger{
		debugMessages:   []string{},
		infoMessages:    []string{},
		warningMessages: []string{},
		errorMessages:   []string{},
		successMessages: []string{},
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

// TestAuctionMySellStop tests the HandleAuctionMySellStop method
func TestAuctionMySellStop(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create auction manager
	auctionManager := NewAuctionManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("auction_my_sell_stop", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test auction my sell stop (ended auction)
	args := map[string]interface{}{
		"flag": uint8(0),
	}

	err := auctionManager.HandleAuctionMySellStop(args)

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
		t.Errorf("Expected auction_my_sell_stop hook to be called")
	}

	// Test auction my sell stop (unknown flag)
	hookCalled = false
	mockLogger = NewMockLogger()
	auctionManager.logger = mockLogger

	args = map[string]interface{}{
		"flag": uint8(10), // Unknown flag
	}

	err = auctionManager.HandleAuctionMySellStop(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify warning message was created
	if len(mockLogger.warningMessages) != 1 {
		t.Errorf("Expected 1 warning message, got %d", len(mockLogger.warningMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected auction_my_sell_stop hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"flag": "invalid", // Invalid type
	}

	err = auctionManager.HandleAuctionMySellStop(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestAuctionWindows tests the HandleAuctionWindows method
func TestAuctionWindows(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create auction manager
	auctionManager := NewAuctionManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("auction_windows", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test auction windows (opened)
	args := map[string]interface{}{
		"flag": uint8(0),
	}

	err := auctionManager.HandleAuctionWindows(args)

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
		t.Errorf("Expected auction_windows hook to be called")
	}

	// Test auction windows (closed)
	hookCalled = false
	mockLogger = NewMockLogger()
	auctionManager.logger = mockLogger

	args = map[string]interface{}{
		"flag": uint8(1),
	}

	err = auctionManager.HandleAuctionWindows(args)

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
		t.Errorf("Expected auction_windows hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"flag": "invalid", // Invalid type
	}

	err = auctionManager.HandleAuctionWindows(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestAuctionAddItem tests the HandleAuctionAddItem method
func TestAuctionAddItem(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create auction manager
	auctionManager := NewAuctionManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("auction_add_item", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test auction add item (success)
	args := map[string]interface{}{
		"fail": uint8(0),
		"ID":   uint32(12345),
	}

	err := auctionManager.HandleAuctionAddItem(args)

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
		t.Errorf("Expected auction_add_item hook to be called")
	}

	// Test auction add item (fail)
	hookCalled = false
	mockLogger = NewMockLogger()
	auctionManager.logger = mockLogger

	args = map[string]interface{}{
		"fail": uint8(1),
		"ID":   uint32(12345),
	}

	err = auctionManager.HandleAuctionAddItem(args)

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
		t.Errorf("Expected auction_add_item hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"fail": "invalid", // Invalid type
		"ID":   uint32(12345),
	}

	err = auctionManager.HandleAuctionAddItem(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestAuctionResult tests the HandleAuctionResult method
func TestAuctionResult(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create auction manager
	auctionManager := NewAuctionManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("auction_result", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test auction result (successful bid)
	args := map[string]interface{}{
		"flag": uint8(1),
	}

	err := auctionManager.HandleAuctionResult(args)

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
		t.Errorf("Expected auction_result hook to be called")
	}

	// Test auction result (won auction)
	hookCalled = false
	mockLogger = NewMockLogger()
	auctionManager.logger = mockLogger

	args = map[string]interface{}{
		"flag": uint8(6),
	}

	err = auctionManager.HandleAuctionResult(args)

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
		t.Errorf("Expected auction_result hook to be called")
	}

	// Test auction result (unknown flag)
	hookCalled = false
	mockLogger = NewMockLogger()
	auctionManager.logger = mockLogger

	args = map[string]interface{}{
		"flag": uint8(100), // Unknown flag
	}

	err = auctionManager.HandleAuctionResult(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify warning message was created
	if len(mockLogger.warningMessages) != 1 {
		t.Errorf("Expected 1 warning message, got %d", len(mockLogger.warningMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected auction_result hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"flag": "invalid", // Invalid type
	}

	err = auctionManager.HandleAuctionResult(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}
