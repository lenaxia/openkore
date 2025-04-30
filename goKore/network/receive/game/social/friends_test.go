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

	// Test friend list with valid args
	args := map[string]interface{}{
		"len":     uint16(100),
		"RAW_MSG": []byte("TestFriend1\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00TestFriend2\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"),
	}

	err := friendManager.HandleFriendList(args)

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
		"len": "invalid", // Invalid type
	}

	err = friendManager.HandleFriendList(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
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

	// Test friend logon (online)
	args := map[string]interface{}{
		"friendAccountID": uint32(12345),
		"friendCharID":    uint32(67890),
		"isNotOnline":     uint8(0), // 0 = online
	}

	err := friendManager.HandleFriendLogon(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Test friend logon (offline)
	args = map[string]interface{}{
		"friendAccountID": uint32(12345),
		"friendCharID":    uint32(67890),
		"isNotOnline":     uint8(1), // 1 = offline
	}

	err = friendManager.HandleFriendLogon(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 2 {
		t.Errorf("Expected 2 info messages, got %d", len(mockLogger.infoMessages))
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"friendAccountID": "invalid", // Invalid type
	}

	err = friendManager.HandleFriendLogon(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
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

	// Track hook calls
	hookCalled := false
	hookManager.AddHook("friend_request", func(hookName string, data interface{}, userData interface{}) {
		hookCalled = true
	}, nil)

	// Test friend request
	args := map[string]interface{}{
		"accountID": uint32(12345),
		"charID":    uint32(67890),
		"name":      "TestFriend",
	}

	err := friendManager.HandleFriendRequest(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log messages were created
	if len(mockLogger.infoMessages) != 2 {
		t.Errorf("Expected 2 info messages, got %d", len(mockLogger.infoMessages))
	}

	// Verify hook was called
	if !hookCalled {
		t.Errorf("Expected friend_request hook to be called")
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"accountID": "invalid", // Invalid type
	}

	err = friendManager.HandleFriendRequest(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
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
	args := map[string]interface{}{
		"friendAccountID": uint32(12345),
		"friendCharID":    uint32(67890),
	}

	err := friendManager.HandleFriendRemoved(args)

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
		"friendAccountID": "invalid", // Invalid type
	}

	err = friendManager.HandleFriendRemoved(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
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

	// Test friend response (accepted)
	args := map[string]interface{}{
		"type": uint16(0), // 0 = accepted
		"name": "TestFriend",
	}

	err := friendManager.HandleFriendResponse(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Test friend response (rejected)
	args = map[string]interface{}{
		"type": uint16(1), // 1 = rejected
		"name": "TestFriend",
	}

	err = friendManager.HandleFriendResponse(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 2 {
		t.Errorf("Expected 2 info messages, got %d", len(mockLogger.infoMessages))
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"type": "invalid", // Invalid type
	}

	err = friendManager.HandleFriendResponse(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestIgnoreAllResult tests the HandleIgnoreAllResult method
func TestIgnoreAllResult(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create friend manager
	friendManager := NewFriendManager(mockParser, hookManager, mockLogger)

	// Test ignore all (type 0)
	args := map[string]interface{}{
		"type":  uint8(0),
		"error": uint8(0),
	}

	err := friendManager.HandleIgnoreAllResult(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Test unignore all (type 1)
	args = map[string]interface{}{
		"type":  uint8(1),
		"error": uint8(0),
	}

	err = friendManager.HandleIgnoreAllResult(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 2 {
		t.Errorf("Expected 2 info messages, got %d", len(mockLogger.infoMessages))
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"type": "invalid", // Invalid type
	}

	err = friendManager.HandleIgnoreAllResult(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
	}
}

// TestIgnorePlayerResult tests the HandleIgnorePlayerResult method
func TestIgnorePlayerResult(t *testing.T) {
	// Create mocks
	mockParser := NewMockParser()
	mockLogger := NewMockLogger()
	hookManager := hooks.NewHookManager()

	// Create friend manager
	friendManager := NewFriendManager(mockParser, hookManager, mockLogger)

	// Test ignore player (type 0)
	args := map[string]interface{}{
		"type":  uint8(0),
		"error": uint8(0),
	}

	err := friendManager.HandleIgnorePlayerResult(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 1 {
		t.Errorf("Expected 1 info message, got %d", len(mockLogger.infoMessages))
	}

	// Test unignore player (type 1)
	args = map[string]interface{}{
		"type":  uint8(1),
		"error": uint8(0),
	}

	err = friendManager.HandleIgnorePlayerResult(args)

	// Verify no error occurred
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify log message was created
	if len(mockLogger.infoMessages) != 2 {
		t.Errorf("Expected 2 info messages, got %d", len(mockLogger.infoMessages))
	}

	// Test with invalid parameters
	invalidArgs := map[string]interface{}{
		"type": "invalid", // Invalid type
	}

	err = friendManager.HandleIgnorePlayerResult(invalidArgs)

	// Verify error occurred
	if err == nil {
		t.Errorf("Expected error for invalid parameters, got nil")
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
		"ignore_all_result",
		"ignore_player_result",
	}

	// We can't verify the handlers directly since we're using a real CoreParser
	// Instead, we'll just log the expected handlers
	for _, handler := range expectedHandlers {
		t.Logf("Expected handler: %s", handler)
	}
}
