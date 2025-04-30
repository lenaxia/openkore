package homunculus

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

// mockCurrentTime is used to mock the currentTime function for testing
func mockCurrentTime(t int64) {
	currentTime = func() int64 {
		return t
	}
}

// TestHomProperty tests the HandleHomProperty method
func TestHomProperty(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create homunculus manager
	homManager := NewHomManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("homunculus_property", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test homunculus property
	args := map[string]interface{}{
		"name":         "TestHomunculus",
		"modified":     uint8(3), // Named and Vaporized
		"level":        uint16(50),
		"hunger":       uint16(80),
		"intimacy":     uint16(100),
		"equip_id":     uint16(0),
		"atk":          uint16(100),
		"matk":         uint16(50),
		"hit":          uint16(80),
		"crit":         uint16(10),
		"def":          uint16(50),
		"mdef":         uint16(30),
		"flee":         uint16(70),
		"aspd":         uint16(190),
		"hp":           uint16(1000),
		"max_hp":       uint16(1200),
		"sp":           uint16(100),
		"max_sp":       uint16(120),
		"exp":          uint32(5000),
		"max_exp":      uint32(10000),
		"skill_points": uint16(3),
		"atk_range":    uint16(2),
	}

	err := homManager.HandleHomProperty(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify homunculus info was updated
	if homManager.homInfo == nil {
		t.Errorf("Expected homInfo to be created")
	} else {
		if homManager.homInfo.Name != "TestHomunculus" {
			t.Errorf("Expected name to be TestHomunculus, got %s", homManager.homInfo.Name)
		}
		if homManager.homInfo.Level != 50 {
			t.Errorf("Expected level to be 50, got %d", homManager.homInfo.Level)
		}
		if homManager.homInfo.Hunger != 80 {
			t.Errorf("Expected hunger to be 80, got %d", homManager.homInfo.Hunger)
		}
		if homManager.homInfo.Intimacy != 100 {
			t.Errorf("Expected intimacy to be 100, got %d", homManager.homInfo.Intimacy)
		}
		if !homManager.homInfo.State.Named {
			t.Errorf("Expected homunculus to be named")
		}
		if !homManager.homInfo.State.Vaporized {
			t.Errorf("Expected homunculus to be vaporized")
		}
	}

	// Verify info messages were created
	if len(mockLogger.infoMessages) < 2 {
		t.Errorf("Expected at least 2 info messages, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected homunculus_property hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"name": 123, // Invalid type
	}

	err = homManager.HandleHomProperty(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestHomInfo tests the HandleHomInfo method
func TestHomInfo(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create homunculus manager
	homManager := NewHomManager(mockParser, hookManager, mockLogger)
	homManager.homInfo = &HomInfo{
		ID:   12345,
		Dead: true,
	}

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("homunculus_info", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test homunculus info (pre-init)
	args := map[string]interface{}{
		"type":  uint8(0),
		"state": uint8(HO_PRE_INIT),
		"ID":    uint32(12345),
		"val":   uint32(0),
	}

	err := homManager.HandleHomInfo(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify homunculus info was updated
	if homManager.homInfo.Dead {
		t.Errorf("Expected homunculus to be not dead after pre-init")
	}

	// Verify debug message was created
	if len(mockLogger.debugMessages) < 1 {
		t.Errorf("Expected at least 1 debug message, got %d", len(mockLogger.debugMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected homunculus_info hook to be called")
	}

	// Test homunculus info (relationship changed)
	hookCalled = false
	args = map[string]interface{}{
		"type":  uint8(0),
		"state": uint8(HO_RELATIONSHIP_CHANGED),
		"ID":    uint32(12345),
		"val":   uint32(200),
	}

	err = homManager.HandleHomInfo(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify homunculus intimacy was updated
	if homManager.homInfo.Intimacy != 200 {
		t.Errorf("Expected intimacy to be 200, got %d", homManager.homInfo.Intimacy)
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected homunculus_info hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"state": "invalid", // Invalid type
		"ID":    uint32(12345),
		"val":   uint32(0),
	}

	err = homManager.HandleHomInfo(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestHomFood tests the HandleHomFood method
func TestHomFood(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create homunculus manager
	homManager := NewHomManager(mockParser, hookManager, mockLogger)
	homManager.homInfo = &HomInfo{
		Hunger: 50,
	}

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("homunculus_food", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test successful feeding
	args := map[string]interface{}{
		"success": uint8(1),
		"foodID":  uint16(501),
	}

	err := homManager.HandleHomFood(args)

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
		t.Errorf("Expected homunculus_food hook to be called")
	}

	// Test failed feeding
	hookCalled = false
	mockLogger = NewMockLogger()
	homManager.logger = mockLogger
	homManager.homInfo.Hunger = 11 // Critical hunger

	args = map[string]interface{}{
		"success": uint8(0),
		"foodID":  uint16(501),
	}

	err = homManager.HandleHomFood(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify error messages were created
	if len(mockLogger.errorMessages) != 2 {
		t.Errorf("Expected 2 error messages, got %d", len(mockLogger.errorMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected homunculus_food hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"success": "invalid", // Invalid type
		"foodID":  uint16(501),
	}

	err = homManager.HandleHomFood(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestEggList tests the HandleEggList method
func TestEggList(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create homunculus manager
	homManager := NewHomManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("egg_list", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Create test raw message
	// First 4 bytes are header, then 2 bytes per egg index
	rawMsg := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x00, 0x06, 0x00}

	// Test egg list
	args := map[string]interface{}{
		"RAW_MSG": rawMsg,
	}

	err := homManager.HandleEggList(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info messages were created
	if len(mockLogger.infoMessages) < 3 {
		t.Errorf("Expected at least 3 info messages, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected egg_list hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"RAW_MSG": "invalid", // Invalid type
	}

	err = homManager.HandleEggList(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}
