package banking

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

// TestBankingCheck tests the HandleBankingCheck method
func TestBankingCheck(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create banking manager
	bankingManager := NewBankingManager(mockParser, hookManager, mockLogger)
	bankingManager.SetCharZeny(1000) // Set character zeny for testing

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("banking_opened", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test banking check
	args := map[string]interface{}{
		"zeny": uint32(5000),
	}

	err := bankingManager.HandleBankingCheck(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify banking state was updated
	if !bankingManager.bankingOpened {
		t.Errorf("Expected bankingOpened to be true")
	}
	if bankingManager.bankingZeny != 5000 {
		t.Errorf("Expected bankingZeny to be 5000, got %d", bankingManager.bankingZeny)
	}

	// Verify info messages were created
	if len(mockLogger.infoMessages) != 4 {
		t.Errorf("Expected 4 info messages, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected banking_opened hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"zeny": "invalid", // Invalid type
	}

	err = bankingManager.HandleBankingCheck(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestBankingDeposit tests the HandleBankingDeposit method
func TestBankingDeposit(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create banking manager
	bankingManager := NewBankingManager(mockParser, hookManager, mockLogger)
	bankingManager.SetCharZeny(1000) // Set character zeny for testing

	// Track hook calls
	successHookCalled := false
	failedHookCalled := false
	hookManager.AddHook("banking_deposit_success", func(hookName string, data interface{}, userData interface{}) {
		successHookCalled = true
	}, nil)
	hookManager.AddHook("banking_deposit_failed", func(hookName string, data interface{}, userData interface{}) {
		failedHookCalled = true
	}, nil)

	// Test successful deposit
	args := map[string]interface{}{
		"reason":  uint16(0), // BDA_SUCCESS
		"money":   uint64(500),
		"balance": uint32(500), // New balance after deposit
	}

	err := bankingManager.HandleBankingDeposit(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify character zeny was updated
	if bankingManager.charZeny != 500 {
		t.Errorf("Expected charZeny to be 500, got %d", bankingManager.charZeny)
	}

	// Verify success message was created
	if len(mockLogger.successMessages) != 1 {
		t.Errorf("Expected 1 success message, got %d", len(mockLogger.successMessages))
	}

	// Verify success hook was called
	if !successHookCalled {
		t.Errorf("Expected banking_deposit_success hook to be called")
	}

	// Test failed deposit
	mockLogger = NewMockLogger()
	bankingManager.logger = mockLogger
	successHookCalled = false
	failedHookCalled = false

	args = map[string]interface{}{
		"reason":  uint16(1), // BDA_ERROR
		"money":   uint64(500),
		"balance": uint32(500),
	}

	err = bankingManager.HandleBankingDeposit(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify error message was created
	if len(mockLogger.errorMessages) != 1 {
		t.Errorf("Expected 1 error message, got %d", len(mockLogger.errorMessages))
	}

	// Verify failed hook was called
	if !failedHookCalled {
		t.Errorf("Expected banking_deposit_failed hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"reason":  "invalid", // Invalid type
		"money":   uint64(500),
		"balance": uint32(500),
	}

	err = bankingManager.HandleBankingDeposit(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestBankingWithdraw tests the HandleBankingWithdraw method
func TestBankingWithdraw(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create banking manager
	bankingManager := NewBankingManager(mockParser, hookManager, mockLogger)
	bankingManager.SetCharZeny(500) // Set character zeny for testing

	// Track hook calls
	successHookCalled := false
	failedHookCalled := false
	hookManager.AddHook("banking_withdraw_success", func(hookName string, data interface{}, userData interface{}) {
		successHookCalled = true
	}, nil)
	hookManager.AddHook("banking_withdraw_failed", func(hookName string, data interface{}, userData interface{}) {
		failedHookCalled = true
	}, nil)

	// Test successful withdraw
	args := map[string]interface{}{
		"reason":  uint16(0), // BWA_SUCCESS
		"money":   uint64(500),
		"balance": uint32(1000), // New balance after withdraw
	}

	err := bankingManager.HandleBankingWithdraw(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify character zeny was updated
	if bankingManager.charZeny != 1000 {
		t.Errorf("Expected charZeny to be 1000, got %d", bankingManager.charZeny)
	}

	// Verify success message was created
	if len(mockLogger.successMessages) != 1 {
		t.Errorf("Expected 1 success message, got %d", len(mockLogger.successMessages))
	}

	// Verify success hook was called
	if !successHookCalled {
		t.Errorf("Expected banking_withdraw_success hook to be called")
	}

	// Test failed withdraw
	mockLogger = NewMockLogger()
	bankingManager.logger = mockLogger
	successHookCalled = false
	failedHookCalled = false

	args = map[string]interface{}{
		"reason":  uint16(1), // BWA_NO_MONEY
		"money":   uint64(500),
		"balance": uint32(1000),
	}

	err = bankingManager.HandleBankingWithdraw(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify error message was created
	if len(mockLogger.errorMessages) != 1 {
		t.Errorf("Expected 1 error message, got %d", len(mockLogger.errorMessages))
	}

	// Verify failed hook was called
	if !failedHookCalled {
		t.Errorf("Expected banking_withdraw_failed hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"reason":  "invalid", // Invalid type
		"money":   uint64(500),
		"balance": uint32(1000),
	}

	err = bankingManager.HandleBankingWithdraw(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}
