package pet

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

// TestPetCaptureProcess tests the HandlePetCaptureProcess method
func TestPetCaptureProcess(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create pet manager
	petManager := NewPetManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("pet_capture_process", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test pet capture process
	args := map[string]interface{}{}

	err := petManager.HandlePetCaptureProcess(args)

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
		t.Errorf("Expected pet_capture_process hook to be called")
	}
}

// TestPetCaptureResult tests the HandlePetCaptureResult method
func TestPetCaptureResult(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create pet manager
	petManager := NewPetManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("pet_capture_result", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test successful pet capture
	args := map[string]interface{}{
		"success": uint8(1),
	}

	err := petManager.HandlePetCaptureResult(args)

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
		t.Errorf("Expected pet_capture_result hook to be called")
	}

	// Test failed pet capture
	hookCalled = false
	mockLogger = NewMockLogger()
	petManager.logger = mockLogger

	args = map[string]interface{}{
		"success": uint8(0),
	}

	err = petManager.HandlePetCaptureResult(args)

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
		t.Errorf("Expected pet_capture_result hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"success": "invalid", // Invalid type
	}

	err = petManager.HandlePetCaptureResult(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestPetEmotion tests the HandlePetEmotion method
func TestPetEmotion(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create pet manager
	petManager := NewPetManager(mockParser, hookManager, mockLogger)
	petManager.petInfo = &PetInfo{
		ID:   12345,
		Name: "TestPet",
	}

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("pet_emotion", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test pet emotion
	args := map[string]interface{}{
		"ID":   uint32(12345),
		"type": uint8(1),
	}

	err := petManager.HandlePetEmotion(args)

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
		t.Errorf("Expected pet_emotion hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"ID":   "invalid", // Invalid type
		"type": uint8(1),
	}

	err = petManager.HandlePetEmotion(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestPetEvolutionResult tests the HandlePetEvolutionResult method
func TestPetEvolutionResult(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create pet manager
	petManager := NewPetManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("pet_evolution_result", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test successful pet evolution
	args := map[string]interface{}{
		"result": uint8(6),
	}

	err := petManager.HandlePetEvolutionResult(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify success message was created
	if len(mockLogger.successMessages) != 1 {
		t.Errorf("Expected 1 success message, got %d", len(mockLogger.successMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected pet_evolution_result hook to be called")
	}

	// Test failed pet evolution
	hookCalled = false
	mockLogger = NewMockLogger()
	petManager.logger = mockLogger

	args = map[string]interface{}{
		"result": uint8(0),
	}

	err = petManager.HandlePetEvolutionResult(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify error message was created
	if len(mockLogger.errorMessages) != 1 {
		t.Errorf("Expected 1 error message, got %d", len(mockLogger.errorMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected pet_evolution_result hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"result": "invalid", // Invalid type
	}

	err = petManager.HandlePetEvolutionResult(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestPetFood tests the HandlePetFood method
func TestPetFood(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create pet manager
	petManager := NewPetManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("pet_food", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test successful feeding
	args := map[string]interface{}{
		"success": uint8(1),
		"foodID":  uint16(501),
	}

	err := petManager.HandlePetFood(args)

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
		t.Errorf("Expected pet_food hook to be called")
	}

	// Test failed feeding
	hookCalled = false
	mockLogger = NewMockLogger()
	petManager.logger = mockLogger

	args = map[string]interface{}{
		"success": uint8(0),
		"foodID":  uint16(501),
	}

	err = petManager.HandlePetFood(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify error message was created
	if len(mockLogger.errorMessages) != 1 {
		t.Errorf("Expected 1 error message, got %d", len(mockLogger.errorMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected pet_food hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"success": "invalid", // Invalid type
		"foodID":  uint16(501),
	}

	err = petManager.HandlePetFood(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestPetInfo tests the HandlePetInfo method
func TestPetInfo(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create pet manager
	petManager := NewPetManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("pet_info", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test pet info
	args := map[string]interface{}{
		"name":       "TestPet",
		"renameflag": uint8(1),
		"level":      uint16(10),
		"hungry":     uint16(50),
		"friendly":   uint16(100),
		"accessory":  uint16(0),
		"type":       uint16(1),
	}

	err := petManager.HandlePetInfo(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify pet info was updated
	if petManager.petInfo == nil {
		t.Errorf("Expected petInfo to be created")
	} else {
		if petManager.petInfo.Name != "TestPet" {
			t.Errorf("Expected name to be TestPet, got %s", petManager.petInfo.Name)
		}
		if !petManager.petInfo.RenameFlag {
			t.Errorf("Expected renameflag to be true")
		}
		if petManager.petInfo.Level != 10 {
			t.Errorf("Expected level to be 10, got %d", petManager.petInfo.Level)
		}
		if petManager.petInfo.Hungry != 50 {
			t.Errorf("Expected hungry to be 50, got %d", petManager.petInfo.Hungry)
		}
		if petManager.petInfo.Friendly != 100 {
			t.Errorf("Expected friendly to be 100, got %d", petManager.petInfo.Friendly)
		}
	}

	// Verify debug message was created
	if len(mockLogger.debugMessages) != 1 {
		t.Errorf("Expected 1 debug message, got %d", len(mockLogger.debugMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected pet_info hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"name": 123, // Invalid type
	}

	err = petManager.HandlePetInfo(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestPetInfo2 tests the HandlePetInfo2 method
func TestPetInfo2(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create pet manager
	petManager := NewPetManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("pet_info2", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test pet info2 (friendly update)
	args := map[string]interface{}{
		"type":  uint8(1),
		"ID":    uint32(12345),
		"value": uint32(100),
	}

	err := petManager.HandlePetInfo2(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify pet info was updated
	if petManager.petInfo == nil {
		t.Errorf("Expected petInfo to be created")
	} else {
		if petManager.petInfo.Friendly != 100 {
			t.Errorf("Expected friendly to be 100, got %d", petManager.petInfo.Friendly)
		}
	}

	// Verify debug message was created
	if len(mockLogger.debugMessages) != 1 {
		t.Errorf("Expected 1 debug message, got %d", len(mockLogger.debugMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected pet_info2 hook to be called")
	}

	// Test pet info2 (hungry update)
	hookCalled = false
	mockLogger = NewMockLogger()
	petManager.logger = mockLogger

	args = map[string]interface{}{
		"type":  uint8(2),
		"ID":    uint32(12345),
		"value": uint32(50),
	}

	err = petManager.HandlePetInfo2(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify pet info was updated
	if petManager.petInfo == nil {
		t.Errorf("Expected petInfo to be created")
	} else {
		if petManager.petInfo.Hungry != 50 {
			t.Errorf("Expected hungry to be 50, got %d", petManager.petInfo.Hungry)
		}
	}

	// Verify debug message was created
	if len(mockLogger.debugMessages) != 1 {
		t.Errorf("Expected 1 debug message, got %d", len(mockLogger.debugMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected pet_info2 hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"type":  "invalid", // Invalid type
		"ID":    uint32(12345),
		"value": uint32(100),
	}

	err = petManager.HandlePetInfo2(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}
