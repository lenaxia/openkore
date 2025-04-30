package deal

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

// TestDealAddOther tests the HandleDealAddOther method
func TestDealAddOther(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create deal manager
	dealManager := NewDealManager(mockParser, hookManager, mockLogger)

	// Initialize current deal
	dealManager.currentDeal = &DealState{
		Name:       "TestPlayer",
		OtherItems: make(map[uint16]*DealItem),
		YourItems:  make(map[uint16]*DealItem),
	}

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("deal_add_other", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test deal add other (item)
	args := map[string]interface{}{
		"nameID":     uint32(501),
		"amount":     uint16(10),
		"identified": uint32(1),
		"broken":     uint8(0),
		"upgrade":    uint8(7),
		"cards":      []byte{1, 0, 0, 0, 2, 0, 0, 0},
		"options":    []byte{1, 0, 0, 0, 5},
	}

	err := dealManager.HandleDealAddOther(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify item was added to deal
	item, exists := dealManager.currentDeal.OtherItems[501]
	if !exists {
		t.Errorf("Expected item to be added to deal")
	} else {
		if item.Amount != 10 {
			t.Errorf("Expected item amount to be 10, got %d", item.Amount)
		}
		if !item.Identified {
			t.Errorf("Expected item to be identified")
		}
		if item.Broken {
			t.Errorf("Expected item to not be broken")
		}
		if item.Upgrade != 7 {
			t.Errorf("Expected item upgrade to be 7, got %d", item.Upgrade)
		}
		if len(item.Cards) != 2 {
			t.Errorf("Expected 2 cards, got %d", len(item.Cards))
		}
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected deal_add_other hook to be called")
	}

	// Test deal add other (zeny)
	hookCalled = false
	mockLogger = NewMockLogger()
	dealManager.logger = mockLogger

	args = map[string]interface{}{
		"nameID":     uint32(0),
		"amount":     uint16(1000),
		"identified": uint32(0),
		"broken":     uint8(0),
		"upgrade":    uint8(0),
		"cards":      []byte{},
		"options":    []byte{},
	}

	err = dealManager.HandleDealAddOther(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify zeny was added to deal
	if dealManager.currentDeal.OtherZeny != 1000 {
		t.Errorf("Expected other zeny to be 1000, got %d", dealManager.currentDeal.OtherZeny)
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected deal_add_other hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"nameID": "invalid", // Invalid type
		"amount": uint16(10),
	}

	err = dealManager.HandleDealAddOther(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestDealAddYou tests the HandleDealAddYou method
func TestDealAddYou(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create deal manager
	dealManager := NewDealManager(mockParser, hookManager, mockLogger)

	// Initialize current deal
	dealManager.currentDeal = &DealState{
		Name:           "TestPlayer",
		OtherItems:     make(map[uint16]*DealItem),
		YourItems:      make(map[uint16]*DealItem),
		LastItemAmount: 5,
	}

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("deal_you_added", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test deal add you (success)
	args := map[string]interface{}{
		"fail": uint8(0),
		"ID":   uint32(501),
	}

	err := dealManager.HandleDealAddYou(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify item was added to deal
	item, exists := dealManager.currentDeal.YourItems[501]
	if !exists {
		t.Errorf("Expected item to be added to deal")
	} else {
		if item.Amount != 5 {
			t.Errorf("Expected item amount to be 5, got %d", item.Amount)
		}
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected deal_you_added hook to be called")
	}

	// Test deal add you (fail)
	hookCalled = false
	mockLogger = NewMockLogger()
	dealManager.logger = mockLogger

	args = map[string]interface{}{
		"fail": uint8(1),
		"ID":   uint32(501),
	}

	err = dealManager.HandleDealAddYou(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify error message was created
	if len(mockLogger.errorMessages) != 1 {
		t.Errorf("Expected 1 error message, got %d", len(mockLogger.errorMessages))
	}

	// Verify hook was not called
	if hookCalled {
		t.Errorf("Expected deal_you_added hook to not be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"fail": "invalid", // Invalid type
		"ID":   uint32(501),
	}

	err = dealManager.HandleDealAddYou(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestDealBegin tests the HandleDealBegin method
func TestDealBegin(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create deal manager
	dealManager := NewDealManager(mockParser, hookManager, mockLogger)

	// Initialize incoming deal
	dealManager.incomingDeal = &IncomingDeal{
		ID:    12345,
		Name:  "TestPlayer",
		Level: 99,
	}

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("engaged_deal", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test deal begin (success)
	args := map[string]interface{}{
		"type": uint8(3),
	}

	err := dealManager.HandleDealBegin(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify current deal was initialized
	if dealManager.currentDeal == nil {
		t.Errorf("Expected current deal to be initialized")
	} else {
		if dealManager.currentDeal.Name != "TestPlayer" {
			t.Errorf("Expected current deal name to be TestPlayer, got %s", dealManager.currentDeal.Name)
		}
	}

	// Verify incoming deal was cleared
	if dealManager.incomingDeal != nil {
		t.Errorf("Expected incoming deal to be cleared")
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected engaged_deal hook to be called")
	}

	// Test deal begin (error)
	hookCalled = false
	mockLogger = NewMockLogger()
	dealManager.logger = mockLogger
	dealManager.currentDeal = nil

	errorHookCalled := false
	hookManager.AddHook("error_deal", func(hookName string, data interface{}, userData interface{}) {
		errorHookCalled = true
	}, nil)

	args = map[string]interface{}{
		"type": uint8(0),
	}

	err = dealManager.HandleDealBegin(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify error message was created
	if len(mockLogger.errorMessages) != 1 {
		t.Errorf("Expected 1 error message, got %d", len(mockLogger.errorMessages))
	}

	// Verify current deal was not initialized
	if dealManager.currentDeal != nil {
		t.Errorf("Expected current deal to not be initialized")
	}

	// Verify error hook was called
	if !errorHookCalled {
		t.Errorf("Expected error_deal hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"type": "invalid", // Invalid type
	}

	err = dealManager.HandleDealBegin(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestDealCancelled tests the HandleDealCancelled method
func TestDealCancelled(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create deal manager
	dealManager := NewDealManager(mockParser, hookManager, mockLogger)

	// Initialize deal state
	dealManager.currentDeal = &DealState{
		Name:       "TestPlayer",
		OtherItems: make(map[uint16]*DealItem),
		YourItems:  make(map[uint16]*DealItem),
	}
	dealManager.incomingDeal = &IncomingDeal{
		ID:    12345,
		Name:  "TestPlayer",
		Level: 99,
	}
	dealManager.outgoingDeal = &OutgoingDeal{
		ID: 12345,
	}

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("cancelled_deal", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test deal cancelled
	args := map[string]interface{}{}

	err := dealManager.HandleDealCancelled(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify deal state was cleared
	if dealManager.currentDeal != nil {
		t.Errorf("Expected current deal to be cleared")
	}
	if dealManager.incomingDeal != nil {
		t.Errorf("Expected incoming deal to be cleared")
	}
	if dealManager.outgoingDeal != nil {
		t.Errorf("Expected outgoing deal to be cleared")
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected cancelled_deal hook to be called")
	}
}

// TestDealComplete tests the HandleDealComplete method
func TestDealComplete(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create deal manager
	dealManager := NewDealManager(mockParser, hookManager, mockLogger)

	// Initialize deal state
	dealManager.currentDeal = &DealState{
		Name:       "TestPlayer",
		OtherItems: make(map[uint16]*DealItem),
		YourItems:  make(map[uint16]*DealItem),
	}
	dealManager.incomingDeal = &IncomingDeal{
		ID:    12345,
		Name:  "TestPlayer",
		Level: 99,
	}
	dealManager.outgoingDeal = &OutgoingDeal{
		ID: 12345,
	}

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("complete_deal", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test deal complete
	args := map[string]interface{}{}

	err := dealManager.HandleDealComplete(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify deal state was cleared
	if dealManager.currentDeal != nil {
		t.Errorf("Expected current deal to be cleared")
	}
	if dealManager.incomingDeal != nil {
		t.Errorf("Expected incoming deal to be cleared")
	}
	if dealManager.outgoingDeal != nil {
		t.Errorf("Expected outgoing deal to be cleared")
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected complete_deal hook to be called")
	}
}

// TestDealFinalize tests the HandleDealFinalize method
func TestDealFinalize(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create deal manager
	dealManager := NewDealManager(mockParser, hookManager, mockLogger)

	// Initialize current deal
	dealManager.currentDeal = &DealState{
		Name:       "TestPlayer",
		OtherItems: make(map[uint16]*DealItem),
		YourItems:  make(map[uint16]*DealItem),
	}

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("finalized_deal", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test deal finalize (other)
	args := map[string]interface{}{
		"type": uint8(1),
	}

	err := dealManager.HandleDealFinalize(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify other finalize was set
	if !dealManager.currentDeal.OtherFinalize {
		t.Errorf("Expected other finalize to be true")
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected finalized_deal hook to be called")
	}

	// Test deal finalize (you)
	hookCalled = false
	mockLogger = NewMockLogger()
	dealManager.logger = mockLogger

	args = map[string]interface{}{
		"type": uint8(0),
	}

	err = dealManager.HandleDealFinalize(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify your finalize was set
	if !dealManager.currentDeal.YourFinalize {
		t.Errorf("Expected your finalize to be true")
	}

	// Verify hook was not called
	if hookCalled {
		t.Errorf("Expected finalized_deal hook to not be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"type": "invalid", // Invalid type
	}

	err = dealManager.HandleDealFinalize(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestDealRequest tests the HandleDealRequest method
func TestDealRequest(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create deal manager
	dealManager := NewDealManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("incoming_deal", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test deal request
	args := map[string]interface{}{
		"user":  "TestPlayer",
		"ID":    uint32(12345),
		"level": uint16(99),
	}

	err := dealManager.HandleDealRequest(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info messages were created
	if len(mockLogger.infoMessages) != 2 {
		t.Errorf("Expected 2 info messages, got %d", len(mockLogger.infoMessages))
	}

	// Verify incoming deal was created
	if dealManager.incomingDeal == nil {
		t.Errorf("Expected incoming deal to be created")
	} else {
		if dealManager.incomingDeal.Name != "TestPlayer" {
			t.Errorf("Expected incoming deal name to be TestPlayer, got %s", dealManager.incomingDeal.Name)
		}
		if dealManager.incomingDeal.ID != 12345 {
			t.Errorf("Expected incoming deal ID to be 12345, got %d", dealManager.incomingDeal.ID)
		}
		if dealManager.incomingDeal.Level != 99 {
			t.Errorf("Expected incoming deal level to be 99, got %d", dealManager.incomingDeal.Level)
		}
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected incoming_deal hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"user":  123, // Invalid type
		"ID":    uint32(12345),
		"level": uint16(99),
	}

	err = dealManager.HandleDealRequest(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}
