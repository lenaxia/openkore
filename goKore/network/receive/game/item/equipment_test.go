package item

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// TestEquipItem tests the HandleEquipItem method
func TestEquipItem(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create equipment manager
	equipmentManager := NewEquipmentManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("item_equipped", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test successful equip
	args := map[string]interface{}{
		"index":   uint16(1001),
		"type":    uint16(2), // Upper headgear
		"success": uint8(1),
	}

	err := equipmentManager.HandleEquipItem(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected item_equipped hook to be called")
	}

	// Test failed equip
	args = map[string]interface{}{
		"index":   uint16(1001),
		"type":    uint16(2), // Upper headgear
		"success": uint8(0),
	}

	err = equipmentManager.HandleEquipItem(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify error message was created
	if len(mockLogger.errorMessages) != 1 {
		t.Errorf("Expected 1 error message, got %d", len(mockLogger.errorMessages))
	}
}

// TestUnequipItem tests the HandleUnequipItem method
func TestUnequipItem(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create equipment manager
	equipmentManager := NewEquipmentManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("item_unequipped", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test successful unequip
	args := map[string]interface{}{
		"index":   uint16(1001),
		"type":    uint16(2), // Upper headgear
		"success": uint8(1),
	}

	err := equipmentManager.HandleUnequipItem(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected item_unequipped hook to be called")
	}

	// Test failed unequip
	args = map[string]interface{}{
		"index":   uint16(1001),
		"type":    uint16(2), // Upper headgear
		"success": uint8(0),
	}

	err = equipmentManager.HandleUnequipItem(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify error message was created
	if len(mockLogger.errorMessages) != 1 {
		t.Errorf("Expected 1 error message, got %d", len(mockLogger.errorMessages))
	}
}

// TestEquipItemSwitch tests the HandleEquipItemSwitch method
func TestEquipItemSwitch(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create equipment manager
	equipmentManager := NewEquipmentManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("item_equipped_switch", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test successful equip switch
	args := map[string]interface{}{
		"index":    uint16(1001),
		"position": uint32(2), // Upper headgear
		"success":  uint8(1),
	}

	err := equipmentManager.HandleEquipItemSwitch(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected item_equipped_switch hook to be called")
	}

	// Test failed equip switch
	args = map[string]interface{}{
		"index":    uint16(1001),
		"position": uint32(2), // Upper headgear
		"success":  uint8(0),
	}

	err = equipmentManager.HandleEquipItemSwitch(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify error message was created
	if len(mockLogger.errorMessages) != 1 {
		t.Errorf("Expected 1 error message, got %d", len(mockLogger.errorMessages))
	}
}

// TestEquipSwitchRunRes tests the HandleEquipSwitchRunRes method
func TestEquipSwitchRunRes(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create equipment manager
	equipmentManager := NewEquipmentManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("equipment_set_switched", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test successful equipment set switch
	args := map[string]interface{}{
		"success": uint8(1),
	}

	err := equipmentManager.HandleEquipSwitchRunRes(args)

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
		t.Errorf("Expected equipment_set_switched hook to be called")
	}

	// Test failed equipment set switch
	args = map[string]interface{}{
		"success": uint8(0),
	}

	err = equipmentManager.HandleEquipSwitchRunRes(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify error message was created
	if len(mockLogger.errorMessages) != 1 {
		t.Errorf("Expected 1 error message, got %d", len(mockLogger.errorMessages))
	}
}

// TestEquipSwitchLog tests the HandleEquipSwitchLog method
func TestEquipSwitchLog(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create equipment manager
	equipmentManager := NewEquipmentManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("equipment_switch_log", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test equipment switch log
	items := []map[string]interface{}{
		{
			"nameID": uint16(1001),
			"type":   uint8(4),
		},
		{
			"nameID": uint16(1002),
			"type":   uint8(4),
		},
	}

	args := map[string]interface{}{
		"items": items,
	}

	err := equipmentManager.HandleEquipSwitchLog(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log messages were created
	if len(mockLogger.infoMessages) != 3 { // Header + 2 items
		t.Errorf("Expected 3 info messages, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected equipment_switch_log hook to be called")
	}
}

// TestArrowEquipped tests the HandleArrowEquipped method
func TestArrowEquipped(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create equipment manager
	equipmentManager := NewEquipmentManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("arrow_equipped", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test arrow equipped
	args := map[string]interface{}{
		"index": uint16(1001),
	}

	err := equipmentManager.HandleArrowEquipped(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected arrow_equipped hook to be called")
	}
}

// TestRegisterEquipmentHandlers tests the RegisterHandlers method for equipment
func TestRegisterEquipmentHandlers(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create equipment manager
	equipmentManager := NewEquipmentManager(mockParser, hookManager, mockLogger)

	// Register handlers
	equipmentManager.RegisterHandlers()

	// Verify handlers were registered
	expectedHandlers := []string{
		"equip_item",
		"unequip_item",
		"equip_item_switch",
		"unequip_item_switch",
		"equip_switch_run_res",
		"equip_switch_log",
		"arrow_equipped",
		"arrow_none",
	}

	for _, handler := range expectedHandlers {
		if _, exists := mockParser.handlers[handler]; !exists {
			t.Errorf("Expected handler %s to be registered", handler)
		}
	}
}
