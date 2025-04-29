package social

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// TestPartyChat tests the HandlePartyChat method
func TestPartyChat(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create party manager
	partyManager := NewPartyManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("packet_partyMsg", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test party chat
	args := map[string]interface{}{
		"user":    "TestUser",
		"message": "Hello, party!",
	}

	err := partyManager.HandlePartyChat(args)

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
		t.Errorf("Expected packet_partyMsg hook to be called")
	}
}

// TestPartyExp tests the HandlePartyExp method
func TestPartyExp(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create party manager
	partyManager := NewPartyManager(mockParser, hookManager, mockLogger)

	// Test individual take exp distribution
	args := map[string]interface{}{
		"expOption":  uint8(0),
		"itemOption": uint8(0),
	}

	err := partyManager.HandlePartyExp(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log messages were created
	if len(mockLogger.infoMessages) != 3 { // EXP, item pickup, item division
		t.Errorf("Expected 3 info messages, got %d", len(mockLogger.infoMessages))
	}

	// Test even share exp distribution
	args = map[string]interface{}{
		"expOption":  uint8(1),
		"itemOption": uint8(1),
	}

	err = partyManager.HandlePartyExp(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log messages were created
	if len(mockLogger.infoMessages) != 6 { // 3 + 3 more
		t.Errorf("Expected 6 info messages, got %d", len(mockLogger.infoMessages))
	}
}

// TestPartyLeader tests the HandlePartyLeader method
func TestPartyLeader(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create party manager
	partyManager := NewPartyManager(mockParser, hookManager, mockLogger)

	// Test party leader change
	args := map[string]interface{}{
		"newLeaderID": uint32(12345),
	}

	err := partyManager.HandlePartyLeader(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}
}

// TestPartyHpInfo tests the HandlePartyHpInfo method
func TestPartyHpInfo(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create party manager
	partyManager := NewPartyManager(mockParser, hookManager, mockLogger)

	// Test party HP update
	args := map[string]interface{}{
		"ID":    uint32(12345),
		"hp":    uint32(800),
		"maxHp": uint32(1000),
	}

	err := partyManager.HandlePartyHpInfo(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify debug message was created
	if len(mockLogger.debugMessages) != 1 {
		t.Errorf("Expected 1 debug message, got %d", len(mockLogger.debugMessages))
	}
}

// TestPartyInvite tests the HandlePartyInvite method
func TestPartyInvite(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create party manager
	partyManager := NewPartyManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("party_invite", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test party invitation
	args := map[string]interface{}{
		"partyName": "TestParty",
	}

	err := partyManager.HandlePartyInvite(args)

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
		t.Errorf("Expected party_invite hook to be called")
	}
}

// TestPartyInviteResult tests the HandlePartyInviteResult method
func TestPartyInviteResult(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create party manager
	partyManager := NewPartyManager(mockParser, hookManager, mockLogger)

	// Test accepted invitation
	args := map[string]interface{}{
		"name":   "TestUser",
		"result": uint8(2), // ANSWER_JOIN_ACCEPT
	}

	err := partyManager.HandlePartyInviteResult(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Test rejected invitation
	args = map[string]interface{}{
		"name":   "TestUser",
		"result": uint8(1), // ANSWER_JOIN_REFUSE
	}

	err = partyManager.HandlePartyInviteResult(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify warning message was created
	if len(mockLogger.warningMessages) != 1 {
		t.Errorf("Expected 1 warning message, got %d", len(mockLogger.warningMessages))
	}
}

// TestPartyJoin tests the HandlePartyJoin method
func TestPartyJoin(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create party manager
	partyManager := NewPartyManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("packet_partyJoin", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test player joining party
	args := map[string]interface{}{
		"partyName": "TestParty",
		"ID":        uint32(12345),
		"name":      "TestUser",
		"isPlayer":  true,
	}

	err := partyManager.HandlePartyJoin(args)

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
		t.Errorf("Expected packet_partyJoin hook to be called")
	}

	// Test other member joining party
	args = map[string]interface{}{
		"partyName": "TestParty",
		"ID":        uint32(67890),
		"name":      "OtherUser",
		"isPlayer":  false,
	}

	err = partyManager.HandlePartyJoin(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 2 {
		t.Errorf("Expected 2 info messages, got %d", len(mockLogger.infoMessages))
	}
}

// TestPartyLvInfo tests the HandlePartyLvInfo method
func TestPartyLvInfo(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create party manager
	partyManager := NewPartyManager(mockParser, hookManager, mockLogger)

	// Test party level info update
	args := map[string]interface{}{
		"ID":    uint32(12345),
		"job":   uint16(4002), // Example job ID
		"level": uint16(99),   // Example level
	}

	err := partyManager.HandlePartyLvInfo(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify debug message was created
	if len(mockLogger.debugMessages) != 1 {
		t.Errorf("Expected 1 debug message, got %d", len(mockLogger.debugMessages))
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"ID":    "invalid", // Invalid type
		"job":   uint16(4002),
		"level": uint16(99),
	}

	err = partyManager.HandlePartyLvInfo(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestRegisterPartyHandlers tests the RegisterHandlers method for party
func TestRegisterPartyHandlers(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create party manager
	partyManager := NewPartyManager(mockParser, hookManager, mockLogger)

	// Register handlers
	partyManager.RegisterHandlers()

	// Verify handlers were registered
	expectedHandlers := []string{
		"party_chat",
		"party_exp",
		"party_leader",
		"party_hp_info",
		"party_invite",
		"party_invite_result",
		"party_location",
		"party_organize_result",
		"party_show_picker",
		"party_users_info",
		"party_dead",
		"partylv_info",
		"party_join",
		"party_allow_invite",
		"party_leave",
	}

	for _, handler := range expectedHandlers {
		if _, exists := mockParser.handlers[handler]; !exists {
			t.Errorf("Expected handler %s to be registered", handler)
		}
	}
}
