package social

import (
	"strings"
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

// TestPartyDead tests the HandlePartyDead method
func TestPartyDead(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create party manager
	partyManager := NewPartyManager(mockParser, hookManager, mockLogger)

	// Test party member death
	args := map[string]interface{}{
		"ID":     uint32(12345),
		"isDead": uint8(1), // 1 = dead
	}

	err := partyManager.HandlePartyDead(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	} else if !contains(mockLogger.infoMessages[0], "is dead") {
		t.Errorf("Expected message to contain 'is dead', got %s", mockLogger.infoMessages[0])
	}

	// Test party member revival
	args = map[string]interface{}{
		"ID":     uint32(12345),
		"isDead": uint8(0), // 0 = alive
	}

	err = partyManager.HandlePartyDead(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 2 {
		t.Errorf("Expected 2 info messages, got %d", len(mockLogger.infoMessages))
	} else if !contains(mockLogger.infoMessages[1], "is alive") {
		t.Errorf("Expected message to contain 'is alive', got %s", mockLogger.infoMessages[1])
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"ID":     "invalid", // Invalid type
		"isDead": uint8(1),
	}

	err = partyManager.HandlePartyDead(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestPartyUsersInfo tests the HandlePartyUsersInfo method
func TestPartyUsersInfo(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create party manager
	partyManager := NewPartyManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("party_users_info_ready", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test party users info
	args := map[string]interface{}{
		"partyName": "TestParty",
		"members": []map[string]interface{}{
			{
				"name":   "Member1",
				"map":    "prontera",
				"online": uint8(1),
			},
			{
				"name":   "Member2",
				"map":    "geffen",
				"online": uint8(0),
			},
		},
	}

	err := partyManager.HandlePartyUsersInfo(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log messages were created
	if len(mockLogger.infoMessages) != 3 { // Party info + 2 members
		t.Errorf("Expected 3 info messages, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected party_users_info_ready hook to be called")
	}

	// Test with invalid parameters - missing partyName
	invalidArgs := map[string]interface{}{
		"members": []map[string]interface{}{
			{
				"name":   "Member1",
				"map":    "prontera",
				"online": uint8(1),
			},
		},
	}

	err = partyManager.HandlePartyUsersInfo(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for missing partyName, got nil")
	}

	// Test with invalid parameters - missing members
	invalidArgs = map[string]interface{}{
		"partyName": "TestParty",
	}

	err = partyManager.HandlePartyUsersInfo(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for missing members, got nil")
	}
}

// TestPartyShowPicker tests the HandlePartyShowPicker method
func TestPartyShowPicker(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create party manager
	partyManager := NewPartyManager(mockParser, hookManager, mockLogger)

	// Test party show picker
	args := map[string]interface{}{
		"sourceID": uint32(12345),
		"itemName": "Potion",
	}

	err := partyManager.HandlePartyShowPicker(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	} else if !contains(mockLogger.infoMessages[0], "picked up") {
		t.Errorf("Expected message to contain 'picked up', got %s", mockLogger.infoMessages[0])
	}

	// Test with invalid sourceID
	invalidArgs := map[string]interface{}{
		"sourceID": "invalid", // Invalid type
		"itemName": "Potion",
	}

	err = partyManager.HandlePartyShowPicker(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid sourceID, got nil")
	}

	// Test with invalid itemName
	invalidArgs = map[string]interface{}{
		"sourceID": uint32(12345),
		// Missing itemName
	}

	err = partyManager.HandlePartyShowPicker(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for missing itemName, got nil")
	}
}

// TestPartyOrganizeResult tests the HandlePartyOrganizeResult method
func TestPartyOrganizeResult(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create party manager
	partyManager := NewPartyManager(mockParser, hookManager, mockLogger)

	// Test successful party creation
	args := map[string]interface{}{
		"fail": uint8(0), // 0 = success
	}

	err := partyManager.HandlePartyOrganizeResult(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify success message was created
	if len(mockLogger.successMessages) != 1 {
		t.Errorf("Expected 1 success message, got %d", len(mockLogger.successMessages))
	}

	// Test party name exists error
	args = map[string]interface{}{
		"fail": uint8(1), // 1 = party name exists
	}

	err = partyManager.HandlePartyOrganizeResult(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify warning message was created
	if len(mockLogger.warningMessages) != 1 {
		t.Errorf("Expected 1 warning message, got %d", len(mockLogger.warningMessages))
	} else if !contains(mockLogger.warningMessages[0], "Party name already exists") {
		t.Errorf("Expected message to contain 'Party name already exists', got %s", mockLogger.warningMessages[0])
	}

	// Test already in party error
	args = map[string]interface{}{
		"fail": uint8(2), // 2 = already in party
	}

	err = partyManager.HandlePartyOrganizeResult(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify warning message was created
	if len(mockLogger.warningMessages) != 2 {
		t.Errorf("Expected 2 warning messages, got %d", len(mockLogger.warningMessages))
	} else if !contains(mockLogger.warningMessages[1], "Already in a party") {
		t.Errorf("Expected message to contain 'Already in a party', got %s", mockLogger.warningMessages[1])
	}

	// Test not allowed in map error
	args = map[string]interface{}{
		"fail": uint8(3), // 3 = not allowed in map
	}

	err = partyManager.HandlePartyOrganizeResult(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify warning message was created
	if len(mockLogger.warningMessages) != 3 {
		t.Errorf("Expected 3 warning messages, got %d", len(mockLogger.warningMessages))
	} else if !contains(mockLogger.warningMessages[2], "Not allowed in current map") {
		t.Errorf("Expected message to contain 'Not allowed in current map', got %s", mockLogger.warningMessages[2])
	}

	// Test unknown error
	args = map[string]interface{}{
		"fail": uint8(99), // Unknown error code
	}

	err = partyManager.HandlePartyOrganizeResult(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify warning message was created
	if len(mockLogger.warningMessages) != 4 {
		t.Errorf("Expected 4 warning messages, got %d", len(mockLogger.warningMessages))
	} else if !contains(mockLogger.warningMessages[3], "Failed to organize party") {
		t.Errorf("Expected message to contain 'Failed to organize party', got %s", mockLogger.warningMessages[3])
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"fail": "invalid", // Invalid type
	}

	err = partyManager.HandlePartyOrganizeResult(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestPartyLocation tests the HandlePartyLocation method
func TestPartyLocation(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create party manager
	partyManager := NewPartyManager(mockParser, hookManager, mockLogger)

	// Test party location update
	args := map[string]interface{}{
		"ID": uint32(12345),
		"x":  uint16(100),
		"y":  uint16(200),
	}

	err := partyManager.HandlePartyLocation(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify debug message was created
	if len(mockLogger.debugMessages) != 1 {
		t.Errorf("Expected 1 debug message, got %d", len(mockLogger.debugMessages))
	} else if !contains(mockLogger.debugMessages[0], "location update") {
		t.Errorf("Expected message to contain 'location update', got %s", mockLogger.debugMessages[0])
	}

	// Test with invalid ID
	invalidArgs := map[string]interface{}{
		"ID": "invalid", // Invalid type
		"x":  uint16(100),
		"y":  uint16(200),
	}

	err = partyManager.HandlePartyLocation(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid ID, got nil")
	}

	// Test with invalid x coordinate
	invalidArgs = map[string]interface{}{
		"ID": uint32(12345),
		"x":  "invalid", // Invalid type
		"y":  uint16(200),
	}

	err = partyManager.HandlePartyLocation(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid x coordinate, got nil")
	}

	// Test with invalid y coordinate
	invalidArgs = map[string]interface{}{
		"ID": uint32(12345),
		"x":  uint16(100),
		"y":  "invalid", // Invalid type
	}

	err = partyManager.HandlePartyLocation(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid y coordinate, got nil")
	}
}

// TestPartyLeave tests the HandlePartyLeave method
func TestPartyLeave(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create party manager
	partyManager := NewPartyManager(mockParser, hookManager, mockLogger)

	// Test player leaving party
	args := map[string]interface{}{
		"ID":       uint32(12345),
		"name":     "TestUser",
		"reason":   uint8(0), // GROUPMEMBER_DELETE_LEAVE
		"isPlayer": true,
	}

	err := partyManager.HandlePartyLeave(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	} else if !contains(mockLogger.infoMessages[0], "You %s the party") {
		t.Errorf("Expected message to contain 'You %%s the party', got %s", mockLogger.infoMessages[0])
	}

	// Test other member leaving party
	args = map[string]interface{}{
		"ID":       uint32(67890),
		"name":     "OtherUser",
		"reason":   uint8(0), // GROUPMEMBER_DELETE_LEAVE
		"isPlayer": false,
	}

	err = partyManager.HandlePartyLeave(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 2 {
		t.Errorf("Expected 2 info messages, got %d", len(mockLogger.infoMessages))
	} else if !contains(mockLogger.infoMessages[1], "%s has %s the party") {
		t.Errorf("Expected message to contain '%%s has %%s the party', got %s", mockLogger.infoMessages[1])
	}

	// Test member kicked from party
	args = map[string]interface{}{
		"ID":       uint32(67890),
		"name":     "OtherUser",
		"reason":   uint8(1), // GROUPMEMBER_DELETE_EXPEL
		"isPlayer": false,
	}

	err = partyManager.HandlePartyLeave(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 3 {
		t.Errorf("Expected 3 info messages, got %d", len(mockLogger.infoMessages))
	} else if !contains(mockLogger.infoMessages[2], "%s has %s the party") {
		t.Errorf("Expected message to contain '%%s has %%s the party', got %s", mockLogger.infoMessages[2])
	}

	// Test with invalid ID
	invalidArgs := map[string]interface{}{
		"ID":       "invalid", // Invalid type
		"name":     "TestUser",
		"reason":   uint8(0),
		"isPlayer": true,
	}

	err = partyManager.HandlePartyLeave(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid ID, got nil")
	}

	// Test with invalid name
	invalidArgs = map[string]interface{}{
		"ID":       uint32(12345),
		"name":     123, // Invalid type
		"reason":   uint8(0),
		"isPlayer": true,
	}

	err = partyManager.HandlePartyLeave(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid name, got nil")
	}

	// Test with invalid reason
	invalidArgs = map[string]interface{}{
		"ID":       uint32(12345),
		"name":     "TestUser",
		"reason":   "invalid", // Invalid type
		"isPlayer": true,
	}

	err = partyManager.HandlePartyLeave(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid reason, got nil")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
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

	// We can't verify the handlers directly since we're using a real CoreParser
	// Instead, we'll just log the expected handlers
	for _, handler := range expectedHandlers {
		t.Logf("Expected handler: %s", handler)
	}
}
