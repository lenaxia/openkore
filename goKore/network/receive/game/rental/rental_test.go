package rental

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

// TestRentalExpired tests the HandleRentalExpired method
func TestRentalExpired(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create rental manager
	rentalManager := NewRentalManager(mockParser, hookManager, mockLogger)

	// Set up inventory with a test item
	inventory := map[uint32]map[string]interface{}{
		12345: {
			"binID":  uint16(1),
			"nameID": uint16(501),
		},
	}
	rentalManager.SetInventory(inventory)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("rental_expired", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test rental expired for an item in inventory
	args := map[string]interface{}{
		"ID":     uint32(12345),
		"nameID": uint16(501),
	}

	err := rentalManager.HandleRentalExpired(args)

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
		t.Errorf("Expected rental_expired hook to be called")
	}

	// Verify item was removed from inventory
	if _, exists := rentalManager.inventory[12345]; exists {
		t.Errorf("Expected item to be removed from inventory")
	}

	// Test rental expired for an item not in inventory
	hookCalled = false
	mockLogger = NewMockLogger()
	rentalManager.logger = mockLogger

	args = map[string]interface{}{
		"ID":     uint32(54321),
		"nameID": uint16(502),
	}

	err = rentalManager.HandleRentalExpired(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was not called
	if hookCalled {
		t.Errorf("Expected rental_expired hook to not be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"ID":     "invalid", // Invalid type
		"nameID": uint16(501),
	}

	err = rentalManager.HandleRentalExpired(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestRentalTime tests the HandleRentalTime method
func TestRentalTime(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create rental manager
	rentalManager := NewRentalManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("rental_time", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test rental time
	args := map[string]interface{}{
		"nameID":  uint16(501),
		"seconds": uint32(3600), // 60 minutes
	}

	err := rentalManager.HandleRentalTime(args)

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
		t.Errorf("Expected rental_time hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"nameID":  "invalid", // Invalid type
		"seconds": uint32(3600),
	}

	err = rentalManager.HandleRentalTime(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestGetItemName tests the GetItemName method
func TestGetItemName(t *testing.T) {
	// Create rental manager
	rentalManager := NewRentalManager(nil, nil, nil)

	// Test getting item name
	itemName := rentalManager.GetItemName(501)

	// Verify item name
	expectedName := "Item#501"
	if itemName != expectedName {
		t.Errorf("Expected item name to be %s, got %s", expectedName, itemName)
	}
}
