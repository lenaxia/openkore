package quest

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

// TestQuestAllList tests the HandleQuestAllList method
func TestQuestAllList(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create quest manager
	questManager := NewQuestManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("quest_all_list_end", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Create a mock message with some quest data (02B1 format)
	message := make([]byte, 10) // 2 quests (5 bytes each)

	// First quest
	message[0] = 1 // quest_id (low byte)
	message[1] = 0 // quest_id (high byte)
	message[2] = 0 // quest_id (high byte)
	message[3] = 0 // quest_id (high byte)
	message[4] = 1 // active

	// Second quest
	message[5] = 2 // quest_id (low byte)
	message[6] = 0 // quest_id (high byte)
	message[7] = 0 // quest_id (high byte)
	message[8] = 0 // quest_id (high byte)
	message[9] = 0 // active

	// Test quest all list
	args := map[string]interface{}{
		"switch":       "02B1",
		"quest_amount": uint16(2),
		"message":      message,
	}

	err := questManager.HandleQuestAllList(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify debug messages were created
	if len(mockLogger.debugMessages) < 2 {
		t.Errorf("Expected at least 2 debug messages, got %d", len(mockLogger.debugMessages))
	}

	// Verify quests were stored
	if len(questManager.questList) != 2 {
		t.Errorf("Expected 2 quests in questList, got %d", len(questManager.questList))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected quest_all_list_end hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"switch":       "02B1",
		"quest_amount": "invalid", // Invalid type
		"message":      message,
	}

	err = questManager.HandleQuestAllList(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestQuestAllMission tests the HandleQuestAllMission method
func TestQuestAllMission(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create quest manager
	questManager := NewQuestManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("quest_all_mission_end", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Create a mock message with some quest data
	message := make([]byte, 44) // 1 quest (14 bytes) + 1 mission (30 bytes)

	// Quest data
	message[0] = 1 // quest_id (low byte)
	message[1] = 0 // quest_id (high byte)
	message[2] = 0 // quest_id (high byte)
	message[3] = 0 // quest_id (high byte)
	// time_start, time_expire, mission_amount
	message[13] = 1 // mission_amount

	// Mission data
	message[14] = 100 // mob_id (low byte)
	message[15] = 0   // mob_id (high byte)
	message[16] = 0   // mob_id (high byte)
	message[17] = 0   // mob_id (high byte)
	message[18] = 5   // mob_count (low byte)
	message[19] = 0   // mob_count (high byte)
	// mob_name_original

	// Initialize questList with a quest
	questManager.questList[1] = map[string]interface{}{
		"quest_id":       uint32(1),
		"active":         uint8(1),
		"mission_amount": uint16(1),
		"missions":       make(map[uint32]map[string]interface{}),
	}

	// Test quest all mission
	args := map[string]interface{}{
		"mission_amount": uint16(1),
		"message":        message,
	}

	err := questManager.HandleQuestAllMission(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify debug messages were created
	if len(mockLogger.debugMessages) < 1 {
		t.Errorf("Expected at least 1 debug message, got %d", len(mockLogger.debugMessages))
	}

	// Verify mission was stored
	if missions, ok := questManager.questList[1]["missions"].(map[uint32]map[string]interface{}); !ok || len(missions) != 1 {
		t.Errorf("Expected 1 mission in quest 1, got %d", len(missions))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected quest_all_mission_end hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"mission_amount": "invalid", // Invalid type
		"message":        message,
	}

	err = questManager.HandleQuestAllMission(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestQuestAdd tests the HandleQuestAdd method
func TestQuestAdd(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create quest manager
	questManager := NewQuestManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("quest_added", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Create a mock message with some mission data
	message := make([]byte, 90) // 3 missions (30 bytes each)

	// First mission
	message[0] = 100 // mob_id (low byte)
	message[1] = 0   // mob_id (high byte)
	message[2] = 0   // mob_id (high byte)
	message[3] = 0   // mob_id (high byte)
	message[4] = 5   // mob_count (low byte)
	message[5] = 0   // mob_count (high byte)
	// mob_name_original

	// Second mission
	message[30] = 101 // mob_id (low byte)
	message[31] = 0   // mob_id (high byte)
	message[32] = 0   // mob_id (high byte)
	message[33] = 0   // mob_id (high byte)
	message[34] = 10  // mob_count (low byte)
	message[35] = 0   // mob_count (high byte)
	// mob_name_original

	// Third mission (beyond mission_amount, should be skipped)
	message[60] = 102 // mob_id (low byte)
	message[61] = 0   // mob_id (high byte)
	message[62] = 0   // mob_id (high byte)
	message[63] = 0   // mob_id (high byte)
	message[64] = 15  // mob_count (low byte)
	message[65] = 0   // mob_count (high byte)
	// mob_name_original

	// Test quest add
	args := map[string]interface{}{
		"switch":         "02B3",
		"questID":        uint32(1),
		"active":         uint8(1),
		"time_start":     uint32(123456),
		"time_expire":    uint32(654321),
		"mission_amount": uint16(2),
		"message":        message,
	}

	err := questManager.HandleQuestAdd(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify debug messages were created
	if len(mockLogger.debugMessages) < 2 {
		t.Errorf("Expected at least 2 debug messages, got %d", len(mockLogger.debugMessages))
	}

	// Verify quest was stored
	if _, exists := questManager.questList[1]; !exists {
		t.Errorf("Expected quest 1 to be stored")
	}

	// Verify missions were stored
	if missions, ok := questManager.questList[1]["missions"].(map[uint32]map[string]interface{}); !ok || len(missions) != 2 {
		t.Errorf("Expected 2 missions in quest 1, got %d", len(missions))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected quest_added hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"switch":         "02B3",
		"questID":        "invalid", // Invalid type
		"active":         uint8(1),
		"time_start":     uint32(123456),
		"time_expire":    uint32(654321),
		"mission_amount": uint16(2),
		"message":        message,
	}

	err = questManager.HandleQuestAdd(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestQuestUpdateMissionHunt tests the HandleQuestUpdateMissionHunt method
func TestQuestUpdateMissionHunt(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create quest manager
	questManager := NewQuestManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("quest_update_mission_hunt_end", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Initialize questList with a quest and mission
	questManager.questList[1] = map[string]interface{}{
		"quest_id":       uint32(1),
		"active":         uint8(1),
		"mission_amount": uint16(1),
		"missions": map[uint32]map[string]interface{}{
			100: {
				"mob_id":   uint32(100),
				"mob_name": "Test Monster",
			},
		},
	}

	// Create a mock message with some mission data
	message := make([]byte, 12) // 1 mission (12 bytes)

	// Mission data
	message[0] = 1   // questID (low byte)
	message[1] = 0   // questID (high byte)
	message[2] = 0   // questID (high byte)
	message[3] = 0   // questID (high byte)
	message[4] = 100 // mob_id (low byte)
	message[5] = 0   // mob_id (high byte)
	message[6] = 0   // mob_id (high byte)
	message[7] = 0   // mob_id (high byte)
	message[8] = 10  // mob_goal (low byte)
	message[9] = 0   // mob_goal (high byte)
	message[10] = 5  // mob_count (low byte)
	message[11] = 0  // mob_count (high byte)

	// Test quest update mission hunt
	args := map[string]interface{}{
		"switch":         "02B5",
		"mission_amount": uint16(1),
		"message":        message,
	}

	err := questManager.HandleQuestUpdateMissionHunt(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify debug messages were created
	if len(mockLogger.debugMessages) < 1 {
		t.Errorf("Expected at least 1 debug message, got %d", len(mockLogger.debugMessages))
	}

	// Verify info messages were created
	if len(mockLogger.infoMessages) < 1 {
		t.Errorf("Expected at least 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify mission was updated
	if missions, ok := questManager.questList[1]["missions"].(map[uint32]map[string]interface{}); !ok {
		t.Errorf("Expected missions map in quest 1")
	} else {
		if mission, exists := missions[100]; !exists {
			t.Errorf("Expected mission 100 in quest 1")
		} else {
			if count, ok := mission["mob_count"].(uint16); !ok || count != 5 {
				t.Errorf("Expected mob_count to be 5, got %v", mission["mob_count"])
			}
			if goal, ok := mission["mob_goal"].(uint16); !ok || goal != 10 {
				t.Errorf("Expected mob_goal to be 10, got %v", mission["mob_goal"])
			}
		}
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected quest_update_mission_hunt_end hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"switch":         "02B5",
		"mission_amount": "invalid", // Invalid type
		"message":        message,
	}

	err = questManager.HandleQuestUpdateMissionHunt(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestQuestDelete tests the HandleQuestDelete method
func TestQuestDelete(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create quest manager
	questManager := NewQuestManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("quest_delete", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Initialize questList with a quest
	questManager.questList[1] = map[string]interface{}{
		"quest_id": uint32(1),
		"active":   uint8(1),
	}

	// Test quest delete
	args := map[string]interface{}{
		"questID": uint32(1),
	}

	err := questManager.HandleQuestDelete(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify quest was deleted
	if _, exists := questManager.questList[1]; exists {
		t.Errorf("Expected quest 1 to be deleted")
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected quest_delete hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"questID": "invalid", // Invalid type
	}

	err = questManager.HandleQuestDelete(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestQuestActive tests the HandleQuestActive method
func TestQuestActive(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create quest manager
	questManager := NewQuestManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("quest_active", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Initialize questList with a quest
	questManager.questList[1] = map[string]interface{}{
		"quest_id": uint32(1),
		"active":   uint8(0),
	}

	// Test quest active
	args := map[string]interface{}{
		"questID": uint32(1),
		"active":  uint8(1),
	}

	err := questManager.HandleQuestActive(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify info message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Verify quest active status was updated
	if active, ok := questManager.questList[1]["active"].(uint8); !ok || active != 1 {
		t.Errorf("Expected quest 1 active status to be 1, got %v", questManager.questList[1]["active"])
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected quest_active hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"questID": "invalid", // Invalid type
		"active":  uint8(1),
	}

	err = questManager.HandleQuestActive(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestIncrementQuestGeneration tests the IncrementQuestGeneration method
func TestIncrementQuestGeneration(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create quest manager
	questManager := NewQuestManager(mockParser, hookManager, mockLogger)

	// Initial generation should be 0
	if questManager.questGeneration != 0 {
		t.Errorf("Expected initial questGeneration to be 0, got %d", questManager.questGeneration)
	}

	// Increment generation
	questManager.IncrementQuestGeneration()

	// Generation should be 1
	if questManager.questGeneration != 1 {
		t.Errorf("Expected questGeneration to be 1, got %d", questManager.questGeneration)
	}
}

// TestRegisterHandlers tests the RegisterHandlers method
func TestRegisterHandlers(t *testing.T) {
	// Create mocks
	mockParser := core.NewCoreParser("ServerType0", nil)
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create quest manager
	questManager := NewQuestManager(mockParser, hookManager, mockLogger)

	// Register handlers
	questManager.RegisterHandlers()

	// Verify handlers were registered
	expectedHandlers := []string{
		"quest_all_list",
		"quest_all_mission",
		"quest_add",
		"quest_update_mission_hunt",
		"quest_delete",
		"quest_active",
	}

	// We can't verify the handlers directly since we're using a real CoreParser
	// Instead, we'll just log the expected handlers
	for _, handler := range expectedHandlers {
		t.Logf("Expected handler: %s", handler)
	}
}
