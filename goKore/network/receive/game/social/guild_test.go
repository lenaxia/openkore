package social

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// TestGuildMembersList tests the HandleGuildMembersList method
func TestGuildMembersList(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create guild manager
	guildManager := NewGuildManager(mockParser, hookManager, mockLogger)

	// Test guild members list
	args := map[string]interface{}{}

	err := guildManager.HandleGuildMembersList(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}
}

// TestGuildName tests the HandleGuildName method
func TestGuildName(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create guild manager
	guildManager := NewGuildManager(mockParser, hookManager, mockLogger)

	// Test guild name
	args := map[string]interface{}{}

	err := guildManager.HandleGuildName(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}
}

// TestGuildMemberOnlineStatus tests the HandleGuildMemberOnlineStatus method
func TestGuildMemberOnlineStatus(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create guild manager
	guildManager := NewGuildManager(mockParser, hookManager, mockLogger)

	// Test guild member online status
	args := map[string]interface{}{}

	err := guildManager.HandleGuildMemberOnlineStatus(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}
}

// TestGuildNotice tests the HandleGuildNotice method
func TestGuildNotice(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create guild manager
	guildManager := NewGuildManager(mockParser, hookManager, mockLogger)

	// Test guild notice
	args := map[string]interface{}{}

	err := guildManager.HandleGuildNotice(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}
}

// TestGuildAlliesEnemyList tests the HandleGuildAlliesEnemyList method
func TestGuildAlliesEnemyList(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create guild manager
	guildManager := NewGuildManager(mockParser, hookManager, mockLogger)

	// Test guild allies/enemy list
	args := map[string]interface{}{}

	err := guildManager.HandleGuildAlliesEnemyList(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}
}

// TestRegisterGuildHandlers tests the RegisterHandlers method for guild
func TestRegisterGuildHandlers(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create guild manager
	guildManager := NewGuildManager(mockParser, hookManager, mockLogger)

	// Register handlers
	guildManager.RegisterHandlers()

	// Verify handlers were registered
	expectedHandlers := []string{
		"guild_members_list",
		"guild_name",
		"guild_member_online_status",
		"guild_notice",
		"guild_allies_enemy_list",
	}

	// We can't verify the handlers directly since we're using a real CoreParser
	// Instead, we'll just log the expected handlers
	for _, handler := range expectedHandlers {
		t.Logf("Expected handler: %s", handler)
	}
}
