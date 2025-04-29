package social

import (
	"testing"

	"github.com/lenaxia/goKore/network/hooks"
)

// TestFriendList tests the HandleFriendList method
func TestFriendList(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create friend manager
	friendManager := NewFriendManager(mockParser, hookManager, mockLogger)

	// Test friend list
	args := map[string]interface{}{}

	err := friendManager.HandleFriendList(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}
}

// TestFriendLogon tests the HandleFriendLogon method
func TestFriendLogon(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create friend manager
	friendManager := NewFriendManager(mockParser, hookManager, mockLogger)

	// Test friend logon
	args := map[string]interface{}{}

	err := friendManager.HandleFriendLogon(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}
}

// TestFriendRequest tests the HandleFriendRequest method
func TestFriendRequest(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create friend manager
	friendManager := NewFriendManager(mockParser, hookManager, mockLogger)

	// Test friend request
	args := map[string]interface{}{}

	err := friendManager.HandleFriendRequest(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}
}

// TestFriendRemoved tests the HandleFriendRemoved method
func TestFriendRemoved(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create friend manager
	friendManager := NewFriendManager(mockParser, hookManager, mockLogger)

	// Test friend removed
	args := map[string]interface{}{}

	err := friendManager.HandleFriendRemoved(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}
}

// TestFriendResponse tests the HandleFriendResponse method
func TestFriendResponse(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create friend manager
	friendManager := NewFriendManager(mockParser, hookManager, mockLogger)

	// Test friend response
	args := map[string]interface{}{}

	err := friendManager.HandleFriendResponse(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}
}

// TestRegisterFriendHandlers tests the RegisterHandlers method for friend
func TestRegisterFriendHandlers(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create friend manager
	friendManager := NewFriendManager(mockParser, hookManager, mockLogger)

	// Register handlers
	friendManager.RegisterHandlers()

	// Verify handlers were registered
	expectedHandlers := []string{
		"friend_list",
		"friend_logon",
		"friend_request",
		"friend_removed",
		"friend_response",
	}

	for _, handler := range expectedHandlers {
		if _, exists := mockParser.handlers[handler]; !exists {
			t.Errorf("Expected handler %s to be registered", handler)
		}
	}
}
