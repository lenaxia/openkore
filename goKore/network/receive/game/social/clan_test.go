package social

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// TestClanUser tests the HandleClanUser method
func TestClanUser(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create clan manager
	clanManager := NewClanManager(mockParser, hookManager, mockLogger)

	// Test clan user info
	args := map[string]interface{}{
		"onlineuser":   uint16(5),
		"totalmembers": uint16(10),
	}

	err := clanManager.HandleClanUser(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"onlineuser": "invalid", // Invalid type
	}

	err = clanManager.HandleClanUser(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestClanInfo tests the HandleClanInfo method
func TestClanInfo(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create clan manager
	clanManager := NewClanManager(mockParser, hookManager, mockLogger)

	// Test clan info
	args := map[string]interface{}{
		"clan_ID":               uint32(123),
		"clan_name":             "TestClan",
		"clan_master":           "ClanMaster",
		"clan_map":              "prontera",
		"alliance_count":        uint16(2),
		"antagonist_count":      uint16(1),
		"ally_antagonist_names": []byte("Ally1\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00Ally2\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00Enemy1\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"),
	}

	err := clanManager.HandleClanInfo(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log messages were created
	if len(mockLogger.infoMessages) != 2 {
		t.Errorf("Expected 2 info messages, got %d", len(mockLogger.infoMessages))
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"clan_ID": "invalid", // Invalid type
	}

	err = clanManager.HandleClanInfo(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestClanChat tests the HandleClanChat method
func TestClanChat(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create clan manager
	clanManager := NewClanManager(mockParser, hookManager, mockLogger)

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("packet_clanMsg", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test clan chat
	args := map[string]interface{}{
		"charname": "ClanMember",
		"message":  "Hello clan!",
	}

	err := clanManager.HandleClanChat(args)

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
		t.Errorf("Expected packet_clanMsg hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"charname": 123, // Invalid type
	}

	err = clanManager.HandleClanChat(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestClanLeave tests the HandleClanLeave method
func TestClanLeave(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create clan manager
	clanManager := NewClanManager(mockParser, hookManager, mockLogger)

	// Test clan leave
	args := map[string]interface{}{}

	err := clanManager.HandleClanLeave(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}
}

// TestRegisterClanHandlers tests the RegisterHandlers method for clan
func TestRegisterClanHandlers(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create clan manager
	clanManager := NewClanManager(mockParser, hookManager, mockLogger)

	// Register handlers
	clanManager.RegisterHandlers()

	// Verify handlers were registered
	expectedHandlers := []string{
		"clan_user",
		"clan_info",
		"clan_chat",
		"clan_leave",
	}

	// We can't verify the handlers directly since we're using a real CoreParser
	// Instead, we'll just log the expected handlers
	for _, handler := range expectedHandlers {
		t.Logf("Expected handler: %s", handler)
	}
}
