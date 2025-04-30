package buyingstore

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

// TestOpenBuyingStore tests the HandleOpenBuyingStore method
func TestOpenBuyingStore(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create buying store manager
	buyingStoreManager := NewBuyingStoreManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("open_buying_store", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test open buying store
	args := map[string]interface{}{
		"amount": uint32(5),
	}

	err := buyingStoreManager.HandleOpenBuyingStore(args)

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
		t.Errorf("Expected open_buying_store hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"amount": "invalid", // Invalid type
	}

	err = buyingStoreManager.HandleOpenBuyingStore(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestBuyerItems tests the HandleBuyerItems method
func TestBuyerItems(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create buying store manager
	buyingStoreManager := NewBuyingStoreManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("buyer_items", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Create test message
	msg := []byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, // Header (8 bytes)
		0x09, 0x0A, 0x0B, 0x0C, // Header (4 bytes)
		0x01, 0x00, 0x00, 0x00, // Total (4 bytes)
		0x64, 0x00, 0x00, 0x00, // Price (4 bytes)
		0x05, 0x00, // Amount (2 bytes)
		0x00,       // Type (1 byte)
		0x01, 0x02, // NameID (2 bytes)
	}

	// Test buyer items
	args := map[string]interface{}{
		"venderID": uint32(12345),
		"msg":      msg,
	}

	err := buyingStoreManager.HandleBuyerItems(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected buyer_items hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"venderID": "invalid", // Invalid type
		"msg":      msg,
	}

	err = buyingStoreManager.HandleBuyerItems(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestGetItemName tests the GetItemName method
func TestGetItemName(t *testing.T) {
	// Create buying store manager
	buyingStoreManager := NewBuyingStoreManager(nil, nil, nil)

	// Test getting item name
	itemName := buyingStoreManager.GetItemName(501)

	// Verify item name
	expectedName := "Item#501"
	if itemName != expectedName {
		t.Errorf("Expected item name to be %s, got %s", expectedName, itemName)
	}
}
