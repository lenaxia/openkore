package core

import (
	"testing"
	"time"

	"github.com/lenaxia/goKore/network"
	"github.com/lenaxia/goKore/network/hooks"
)

func TestNewAccountManager(t *testing.T) {
	parser := NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()
	manager := NewAccountManager(parser, hookManager, logger)

	if manager == nil {
		t.Fatal("NewAccountManager() returned nil")
	}

	if manager.parser != parser {
		t.Error("manager.parser was not set correctly")
	}

	if manager.session == nil {
		t.Error("manager.session was not initialized")
	}

	if manager.session.State != AccountStateLoggedOut {
		t.Errorf("manager.session.State = %v, want %v", manager.session.State, AccountStateLoggedOut)
	}

	if manager.networkState != network.NotConnected {
		t.Errorf("manager.networkState = %v, want %v", manager.networkState, network.NotConnected)
	}
}

func TestRegisterHandlers(t *testing.T) {
	parser := NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()
	manager := NewAccountManager(parser, hookManager, logger)

	// Register handlers
	manager.RegisterHandlers()

	// Verify handlers were registered
	handlerNames := []string{
		"account_server_info",
		"received_characters_info",
		"login_error_game_login_server",
		"character_creation_successful",
		"received_character_ID_and_Map",
		"map_loaded",
	}

	for _, name := range handlerNames {
		if _, exists := parser.handlers[name]; !exists {
			t.Errorf("Handler %s was not registered", name)
		}
	}
}

func TestSetGetAccountID(t *testing.T) {
	parser := NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()
	manager := NewAccountManager(parser, hookManager, logger)

	// Set account ID
	accountID := uint32(12345)
	err := manager.SetAccountID(accountID)
	if err != nil {
		t.Fatalf("SetAccountID() returned error: %v", err)
	}

	// Get session and check account ID
	session := manager.GetSession()
	if session.AccountID != accountID {
		t.Errorf("session.AccountID = %v, want %v", session.AccountID, accountID)
	}

	// Try to set invalid account ID
	err = manager.SetAccountID(0)
	if err != ErrInvalidAccountID {
		t.Errorf("SetAccountID(0) returned error: %v, want %v", err, ErrInvalidAccountID)
	}
}

func TestSetGetCharID(t *testing.T) {
	parser := NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()
	manager := NewAccountManager(parser, hookManager, logger)

	// Set character ID
	charID := uint32(67890)
	err := manager.SetCharID(charID)
	if err != nil {
		t.Fatalf("SetCharID() returned error: %v", err)
	}

	// Get session and check character ID
	session := manager.GetSession()
	if session.CharID != charID {
		t.Errorf("session.CharID = %v, want %v", session.CharID, charID)
	}

	// Try to set invalid character ID
	err = manager.SetCharID(0)
	if err != ErrInvalidCharID {
		t.Errorf("SetCharID(0) returned error: %v, want %v", err, ErrInvalidCharID)
	}
}

func TestSetGetState(t *testing.T) {
	parser := NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()
	manager := NewAccountManager(parser, hookManager, logger)

	// Set state
	manager.SetState(AccountStateLoggedIn)

	// Get state
	state := manager.GetState()
	if state != AccountStateLoggedIn {
		t.Errorf("GetState() = %v, want %v", state, AccountStateLoggedIn)
	}

	// Get session and check state
	session := manager.GetSession()
	if session.State != AccountStateLoggedIn {
		t.Errorf("session.State = %v, want %v", session.State, AccountStateLoggedIn)
	}
}

func TestSetGetNetworkState(t *testing.T) {
	parser := NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()
	manager := NewAccountManager(parser, hookManager, logger)

	// Test all network states
	testCases := []struct {
		networkState int
		accountState AccountState
	}{
		{network.NotConnected, AccountStateLoggedOut},
		{network.ConnectedToMasterServer, AccountStateLoggingIn},
		{network.ConnectedToLoginServer, AccountStateLoggingIn},
		{network.ConnectedToCharServer, AccountStateLoggedIn},
		{network.InGame, AccountStateInGame},
	}

	for _, tc := range testCases {
		// Set network state
		manager.SetNetworkState(tc.networkState)

		// Get network state
		state := manager.GetNetworkState()
		if state != tc.networkState {
			t.Errorf("GetNetworkState() = %v, want %v", state, tc.networkState)
		}

		// Check that account state was updated
		accountState := manager.GetState()
		if accountState != tc.accountState {
			t.Errorf("GetState() = %v, want %v for network state %v", accountState, tc.accountState, tc.networkState)
		}
	}
}

func TestIsLoggedIn(t *testing.T) {
	parser := NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()
	manager := NewAccountManager(parser, hookManager, logger)

	// Test not logged in states
	notLoggedInStates := []AccountState{
		AccountStateUnknown,
		AccountStateLoggedOut,
		AccountStateLoggingIn,
	}

	for _, state := range notLoggedInStates {
		manager.SetState(state)
		if manager.IsLoggedIn() {
			t.Errorf("IsLoggedIn() = true, want false for state %v", state)
		}
	}

	// Test logged in states
	loggedInStates := []AccountState{
		AccountStateLoggedIn,
		AccountStateSelectingChar,
		AccountStateInGame,
	}

	for _, state := range loggedInStates {
		manager.SetState(state)
		if !manager.IsLoggedIn() {
			t.Errorf("IsLoggedIn() = false, want true for state %v", state)
		}
	}
}

func TestIsInGame(t *testing.T) {
	parser := NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()
	manager := NewAccountManager(parser, hookManager, logger)

	// Test not in game states
	notInGameStates := []AccountState{
		AccountStateUnknown,
		AccountStateLoggedOut,
		AccountStateLoggingIn,
		AccountStateLoggedIn,
		AccountStateSelectingChar,
	}

	for _, state := range notInGameStates {
		manager.SetState(state)
		if manager.IsInGame() {
			t.Errorf("IsInGame() = true, want false for state %v", state)
		}
	}

	// Test in game state
	manager.SetState(AccountStateInGame)
	if !manager.IsInGame() {
		t.Errorf("IsInGame() = false, want true for state %v", AccountStateInGame)
	}
}

func TestResetSession(t *testing.T) {
	parser := NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()
	manager := NewAccountManager(parser, hookManager, logger)

	// Set some session data
	manager.SetAccountID(12345)
	manager.SetCharID(67890)
	manager.SetState(AccountStateInGame)

	// Reset session
	manager.ResetSession()

	// Check that session was reset
	session := manager.GetSession()
	if session.AccountID != 0 {
		t.Errorf("session.AccountID = %v, want 0", session.AccountID)
	}
	if session.CharID != 0 {
		t.Errorf("session.CharID = %v, want 0", session.CharID)
	}
	if session.State != AccountStateLoggedOut {
		t.Errorf("session.State = %v, want %v", session.State, AccountStateLoggedOut)
	}
}

func TestUpdateLastPacketTime(t *testing.T) {
	parser := NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()
	manager := NewAccountManager(parser, hookManager, logger)

	// Get initial last packet time
	initialTime := manager.GetSession().LastPacketTime

	// Wait a bit
	time.Sleep(10 * time.Millisecond)

	// Update last packet time
	manager.UpdateLastPacketTime()

	// Check that last packet time was updated
	newTime := manager.GetSession().LastPacketTime
	if !newTime.After(initialTime) {
		t.Errorf("LastPacketTime was not updated")
	}
}

func TestIsSessionExpired(t *testing.T) {
	parser := NewCoreParser("ServerType0", nil)
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()
	manager := NewAccountManager(parser, hookManager, logger)

	// Session should not be expired initially
	if manager.IsSessionExpired(1 * time.Second) {
		t.Error("IsSessionExpired() = true, want false for new session")
	}

	// Set last packet time to a long time ago
	manager.session.LastPacketTime = time.Now().Add(-2 * time.Second)

	// Session should be expired now
	if !manager.IsSessionExpired(1 * time.Second) {
		t.Error("IsSessionExpired() = false, want true for expired session")
	}
}

func TestHandleAccountServerInfo(t *testing.T) {
	parser := NewCoreParser("ServerType0", hooks.NewHookManager())
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()
	manager := NewAccountManager(parser, hookManager, logger)

	// Create test packet arguments
	args := map[string]interface{}{
		"sessionID":  []byte{1, 2, 3, 4},
		"accountID":  []byte{5, 6, 7, 8},
		"sessionID2": []byte{9, 10, 11, 12},
		"accountSex": byte(1),
	}

	// Call handler
	err := manager.handleAccountServerInfo(args)
	if err != nil {
		t.Fatalf("handleAccountServerInfo() returned error: %v", err)
	}

	// Check that session was updated
	session := manager.GetSession()
	if len(session.SessionID) != 4 || session.SessionID[0] != 1 || session.SessionID[1] != 2 || session.SessionID[2] != 3 || session.SessionID[3] != 4 {
		t.Errorf("session.SessionID = %v, want [1 2 3 4]", session.SessionID)
	}
	if session.AccountID != 0x08070605 {
		t.Errorf("session.AccountID = %v, want %v", session.AccountID, 0x08070605)
	}
	if len(session.SessionID2) != 4 || session.SessionID2[0] != 9 || session.SessionID2[1] != 10 || session.SessionID2[2] != 11 || session.SessionID2[3] != 12 {
		t.Errorf("session.SessionID2 = %v, want [9 10 11 12]", session.SessionID2)
	}
	if session.Sex != 1 {
		t.Errorf("session.Sex = %v, want 1", session.Sex)
	}
	if session.State != AccountStateLoggedIn {
		t.Errorf("session.State = %v, want %v", session.State, AccountStateLoggedIn)
	}
}

func TestHandleReceivedCharactersInfo(t *testing.T) {
	parser := NewCoreParser("ServerType0", hooks.NewHookManager())
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()
	manager := NewAccountManager(parser, hookManager, logger)

	// Create test packet arguments
	args := map[string]interface{}{
		"total_slot":         byte(12),
		"premium_start_slot": byte(3),
		"premium_end_slot":   byte(5),
	}

	// Call handler
	err := manager.handleReceivedCharactersInfo(args)
	if err != nil {
		t.Fatalf("handleReceivedCharactersInfo() returned error: %v", err)
	}

	// Check that session was updated
	session := manager.GetSession()
	if session.CharacterSlots != 12 {
		t.Errorf("session.CharacterSlots = %v, want 12", session.CharacterSlots)
	}
	if session.PremiumStartSlot != 3 {
		t.Errorf("session.PremiumStartSlot = %v, want 3", session.PremiumStartSlot)
	}
	if session.PremiumEndSlot != 5 {
		t.Errorf("session.PremiumEndSlot = %v, want 5", session.PremiumEndSlot)
	}
	if session.State != AccountStateSelectingChar {
		t.Errorf("session.State = %v, want %v", session.State, AccountStateSelectingChar)
	}
}

func TestHandleLoginError(t *testing.T) {
	parser := NewCoreParser("ServerType0", hooks.NewHookManager())
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()
	manager := NewAccountManager(parser, hookManager, logger)

	// Set initial state
	manager.SetState(AccountStateLoggingIn)

	// Call handler
	err := manager.handleLoginError(map[string]interface{}{})
	if err != nil {
		t.Fatalf("handleLoginError() returned error: %v", err)
	}

	// Check that session was updated
	session := manager.GetSession()
	if session.State != AccountStateLoggedOut {
		t.Errorf("session.State = %v, want %v", session.State, AccountStateLoggedOut)
	}
}

func TestHandleCharacterCreationSuccessful(t *testing.T) {
	parser := NewCoreParser("ServerType0", hooks.NewHookManager())
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()
	manager := NewAccountManager(parser, hookManager, logger)

	// Set initial state
	manager.SetState(AccountStateLoggedIn)

	// Call handler
	err := manager.handleCharacterCreationSuccessful(map[string]interface{}{})
	if err != nil {
		t.Fatalf("handleCharacterCreationSuccessful() returned error: %v", err)
	}

	// Check that session was updated
	session := manager.GetSession()
	if session.State != AccountStateSelectingChar {
		t.Errorf("session.State = %v, want %v", session.State, AccountStateSelectingChar)
	}
}

func TestHandleReceivedCharacterIDAndMap(t *testing.T) {
	parser := NewCoreParser("ServerType0", hooks.NewHookManager())
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()
	manager := NewAccountManager(parser, hookManager, logger)

	// Create test packet arguments
	args := map[string]interface{}{
		"charID":  []byte{1, 2, 3, 4},
		"mapName": "prontera",
	}

	// Call handler
	err := manager.handleReceivedCharacterIDAndMap(args)
	if err != nil {
		t.Fatalf("handleReceivedCharacterIDAndMap() returned error: %v", err)
	}

	// Check that session was updated
	session := manager.GetSession()
	if session.CharID != 0x04030201 {
		t.Errorf("session.CharID = %v, want %v", session.CharID, 0x04030201)
	}
	if session.MapName != "prontera" {
		t.Errorf("session.MapName = %v, want prontera", session.MapName)
	}
	if session.State != AccountStateInGame {
		t.Errorf("session.State = %v, want %v", session.State, AccountStateInGame)
	}
}

func TestHandleMapLoaded(t *testing.T) {
	parser := NewCoreParser("ServerType0", hooks.NewHookManager())
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()
	manager := NewAccountManager(parser, hookManager, logger)

	// Set initial state
	manager.SetState(AccountStateSelectingChar)

	// Call handler
	err := manager.handleMapLoaded(map[string]interface{}{})
	if err != nil {
		t.Fatalf("handleMapLoaded() returned error: %v", err)
	}

	// Check that session was updated
	session := manager.GetSession()
	if session.State != AccountStateInGame {
		t.Errorf("session.State = %v, want %v", session.State, AccountStateInGame)
	}
}

func TestAccountStateString(t *testing.T) {
	testCases := []struct {
		state AccountState
		want  string
	}{
		{AccountStateUnknown, "Unknown"},
		{AccountStateLoggedOut, "LoggedOut"},
		{AccountStateLoggingIn, "LoggingIn"},
		{AccountStateLoggedIn, "LoggedIn"},
		{AccountStateSelectingChar, "SelectingChar"},
		{AccountStateInGame, "InGame"},
		{AccountState(99), "Invalid"},
	}

	for _, tc := range testCases {
		got := tc.state.String()
		if got != tc.want {
			t.Errorf("AccountState(%d).String() = %v, want %v", tc.state, got, tc.want)
		}
	}
}
