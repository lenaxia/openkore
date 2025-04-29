package social

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// TestNewHandler tests the NewHandler function
func TestNewHandler(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create handler
	handler := NewHandler(mockParser, hookManager, mockLogger)

	// Verify handler was created
	if handler == nil {
		t.Errorf("Expected handler to be created")
	}

	// Verify sub-managers were created
	if handler.chatManager == nil {
		t.Errorf("Expected chat manager to be created")
	}

	if handler.partyManager == nil {
		t.Errorf("Expected party manager to be created")
	}

	if handler.friendManager == nil {
		t.Errorf("Expected friend manager to be created")
	}

	if handler.guildManager == nil {
		t.Errorf("Expected guild manager to be created")
	}
}

// TestRegisterAllSocialHandlers tests the RegisterHandlers method for the main handler
func TestRegisterAllSocialHandlers(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create handler
	handler := NewHandler(mockParser, hookManager, mockLogger)

	// Register handlers
	handler.RegisterHandlers()

	// Verify handlers were registered
	// We can't directly check the handlers, but we can verify that the
	// sub-managers' RegisterHandlers methods were called by checking
	// if the expected handlers were registered

	// Expected chat handlers
	chatHandlers := []string{
		"self_chat",
		"private_message",
		"private_message_sent",
		"system_chat",
		"local_broadcast",
		"chat_created",
		"chat_info",
		"chat_users",
		"chat_join_result",
		"chat_modified",
		"chat_newowner",
		"chat_user_join",
		"chat_user_leave",
		"chat_removed",
		"whisper_list",
	}

	// Expected party handlers
	partyHandlers := []string{
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

	// Expected friend handlers
	friendHandlers := []string{
		"friend_list",
		"friend_logon",
		"friend_request",
		"friend_removed",
		"friend_response",
	}

	// Expected guild handlers
	guildHandlers := []string{
		"guild_members_list",
		"guild_name",
		"guild_member_online_status",
		"guild_notice",
		"guild_allies_enemy_list",
	}

	// Check if all expected handlers were registered
	for _, handlerName := range chatHandlers {
		if _, exists := mockParser.handlers[handlerName]; !exists {
			t.Errorf("Expected chat handler %s to be registered", handlerName)
		}
	}

	for _, handlerName := range partyHandlers {
		if _, exists := mockParser.handlers[handlerName]; !exists {
			t.Errorf("Expected party handler %s to be registered", handlerName)
		}
	}

	for _, handlerName := range friendHandlers {
		if _, exists := mockParser.handlers[handlerName]; !exists {
			t.Errorf("Expected friend handler %s to be registered", handlerName)
		}
	}

	for _, handlerName := range guildHandlers {
		if _, exists := mockParser.handlers[handlerName]; !exists {
			t.Errorf("Expected guild handler %s to be registered", handlerName)
		}
	}
}

// TestGetChatManager tests the GetChatManager method
func TestGetChatManager(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create handler
	handler := NewHandler(mockParser, hookManager, mockLogger)

	// Get chat manager
	chatManager := handler.GetChatManager()

	// Verify chat manager was returned
	if chatManager == nil {
		t.Errorf("Expected chat manager to be returned")
	}

	// Verify it's the same instance
	if chatManager != handler.chatManager {
		t.Errorf("Expected returned chat manager to be the same instance")
	}
}

// TestGetPartyManager tests the GetPartyManager method
func TestGetPartyManager(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create handler
	handler := NewHandler(mockParser, hookManager, mockLogger)

	// Get party manager
	partyManager := handler.GetPartyManager()

	// Verify party manager was returned
	if partyManager == nil {
		t.Errorf("Expected party manager to be returned")
	}

	// Verify it's the same instance
	if partyManager != handler.partyManager {
		t.Errorf("Expected returned party manager to be the same instance")
	}
}

// TestGetFriendManager tests the GetFriendManager method
func TestGetFriendManager(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create handler
	handler := NewHandler(mockParser, hookManager, mockLogger)

	// Get friend manager
	friendManager := handler.GetFriendManager()

	// Verify friend manager was returned
	if friendManager == nil {
		t.Errorf("Expected friend manager to be returned")
	}

	// Verify it's the same instance
	if friendManager != handler.friendManager {
		t.Errorf("Expected returned friend manager to be the same instance")
	}
}

// TestGetGuildManager tests the GetGuildManager method
func TestGetGuildManager(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create handler
	handler := NewHandler(mockParser, hookManager, mockLogger)

	// Get guild manager
	guildManager := handler.GetGuildManager()

	// Verify guild manager was returned
	if guildManager == nil {
		t.Errorf("Expected guild manager to be returned")
	}

	// Verify it's the same instance
	if guildManager != handler.guildManager {
		t.Errorf("Expected returned guild manager to be the same instance")
	}
}
