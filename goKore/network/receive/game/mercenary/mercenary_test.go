package mercenary

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

// TestSlaveCalcPropertyHandler tests the SlaveCalcPropertyHandler function
func TestSlaveCalcPropertyHandler(t *testing.T) {
	// Test with aspd < 10
	slave := make(map[string]interface{})
	args := map[string]interface{}{
		"aspd": uint16(5),
	}

	SlaveCalcPropertyHandler(slave, args)

	// Verify attack_speed was calculated correctly
	attackSpeed, ok := slave["attack_speed"].(int)
	if !ok {
		t.Errorf("Expected attack_speed to be set")
	} else if attackSpeed != 190 { // 200 - 10
		t.Errorf("Expected attack_speed to be 190, got %d", attackSpeed)
	}

	// Test with aspd >= 10
	slave = make(map[string]interface{})
	args = map[string]interface{}{
		"aspd": uint16(100),
	}

	SlaveCalcPropertyHandler(slave, args)

	// Verify attack_speed was calculated correctly
	attackSpeed, ok = slave["attack_speed"].(int)
	if !ok {
		t.Errorf("Expected attack_speed to be set")
	} else if attackSpeed != 190 { // 200 - (100 / 10)
		t.Errorf("Expected attack_speed to be 190, got %d", attackSpeed)
	}

	// Test with invalid aspd
	slave = make(map[string]interface{})
	args = map[string]interface{}{
		"aspd": "invalid", // Invalid type
	}

	SlaveCalcPropertyHandler(slave, args)

	// Verify attack_speed was not set
	_, ok = slave["attack_speed"]
	if ok {
		t.Errorf("Expected attack_speed to not be set")
	}
}

// TestMercenaryInit tests the HandleMercenaryInit method
func TestMercenaryInit(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create mercenary manager
	mercManager := NewMercManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("mercenary_init", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test mercenary init
	args := map[string]interface{}{
		"ID":          uint32(12345),
		"name":        "TestMercenary",
		"level":       uint16(50),
		"hp":          uint32(1000),
		"maxhp":       uint32(1200),
		"sp":          uint32(100),
		"maxsp":       uint32(120),
		"atk":         uint16(100),
		"matk":        uint16(50),
		"hit":         uint16(80),
		"crit":        uint16(10),
		"def":         uint16(50),
		"mdef":        uint16(30),
		"flee":        uint16(70),
		"aspd":        uint16(190),
		"atk_range":   uint16(2),
		"expire_time": uint32(3600),
		"faith":       uint16(100),
		"calls":       uint32(5),
		"kills":       uint32(10),
	}

	err := mercManager.HandleMercenaryInit(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify mercenary info was updated
	if mercManager.mercInfo == nil {
		t.Errorf("Expected mercInfo to be created")
	} else {
		if mercManager.mercInfo.ID != 12345 {
			t.Errorf("Expected ID to be 12345, got %d", mercManager.mercInfo.ID)
		}
		if mercManager.mercInfo.Name != "TestMercenary" {
			t.Errorf("Expected name to be TestMercenary, got %s", mercManager.mercInfo.Name)
		}
		if mercManager.mercInfo.Level != 50 {
			t.Errorf("Expected level to be 50, got %d", mercManager.mercInfo.Level)
		}
		if mercManager.mercInfo.HP != 1000 {
			t.Errorf("Expected HP to be 1000, got %d", mercManager.mercInfo.HP)
		}
		if mercManager.mercInfo.MaxHP != 1200 {
			t.Errorf("Expected MaxHP to be 1200, got %d", mercManager.mercInfo.MaxHP)
		}
	}

	// Verify info messages were created
	if len(mockLogger.infoMessages) < 6 {
		t.Errorf("Expected at least 6 info messages, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected mercenary_init hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"ID":   "invalid", // Invalid type
		"name": "TestMercenary",
	}

	err = mercManager.HandleMercenaryInit(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestMercenaryOff tests the HandleMercenaryOff method
func TestMercenaryOff(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create mercenary manager
	mercManager := NewMercManager(mockParser, hookManager, mockLogger)
	mercManager.mercInfo = &MercInfo{
		ID:   12345,
		Name: "TestMercenary",
	}

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("mercenary_off", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test mercenary off
	args := map[string]interface{}{}

	err := mercManager.HandleMercenaryOff(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify mercenary was removed
	if mercManager.mercInfo != nil {
		t.Errorf("Expected mercInfo to be nil")
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected mercenary_off hook to be called")
	}

	// Test with no mercenary
	hookCalled = false
	mockLogger = NewMockLogger()
	mercManager.logger = mockLogger

	err = mercManager.HandleMercenaryOff(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify warning message was created
	if len(mockLogger.warningMessages) != 1 {
		t.Errorf("Expected 1 warning message, got %d", len(mockLogger.warningMessages))
	}

	// Verify hook was not called
	if hookCalled {
		t.Errorf("Expected mercenary_off hook to not be called")
	}
}
