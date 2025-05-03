package core

import (
	"testing"
	"time"

	"github.com/lenaxia/goKore/network"
	"github.com/lenaxia/goKore/network/hooks"
)

func TestHandleConnectionRefused(t *testing.T) {
	parser := NewCoreParser("ServerType0", hooks.NewHookManager())
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()
	manager := NewAccountManager(parser, hookManager, logger)

	// Set initial state
	manager.SetNetworkState(network.ConnectedToLoginServer)
	manager.SetState(AccountStateLoggingIn)

	// Create test packet arguments
	args := map[string]interface{}{
		"error": byte(1), // Some error code
	}

	// Call handler
	err := manager.handleConnectionRefused(args)
	if err != nil {
		t.Fatalf("handleConnectionRefused() returned error: %v", err)
	}

	// Check that session was updated
	if manager.GetNetworkState() != network.NotConnected {
		t.Errorf("NetworkState = %v, want %v", manager.GetNetworkState(), network.NotConnected)
	}

	session := manager.GetSession()
	if session.State != AccountStateLoggedOut {
		t.Errorf("session.State = %v, want %v", session.State, AccountStateLoggedOut)
	}
}

func TestHandleMapLoadError(t *testing.T) {
	parser := NewCoreParser("ServerType0", hooks.NewHookManager())
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()
	manager := NewAccountManager(parser, hookManager, logger)

	// Set initial state
	manager.SetNetworkState(network.ConnectedToCharServer)
	manager.SetState(AccountStateSelectingChar)

	// Create test packet arguments
	args := map[string]interface{}{
		"error": byte(1), // Some error code
	}

	// Call handler
	err := manager.handleMapLoadError(args)
	if err != nil {
		t.Fatalf("handleMapLoadError() returned error: %v", err)
	}

	// Check that session was updated
	if manager.GetNetworkState() != network.ConnectedToCharServer {
		t.Errorf("NetworkState = %v, want %v", manager.GetNetworkState(), network.ConnectedToCharServer)
	}

	session := manager.GetSession()
	if session.State != AccountStateSelectingChar {
		t.Errorf("session.State = %v, want %v", session.State, AccountStateSelectingChar)
	}
}

func TestHandleReceivedSync(t *testing.T) {
	parser := NewCoreParser("ServerType0", hooks.NewHookManager())
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()
	manager := NewAccountManager(parser, hookManager, logger)

	// Set initial state
	manager.SetNetworkState(network.InGame)
	manager.SetState(AccountStateInGame)

	// Create test packet arguments
	args := map[string]interface{}{
		"time": uint32(12345),
	}

	// Call handler
	err := manager.handleReceivedSync(args)
	if err != nil {
		t.Fatalf("handleReceivedSync() returned error: %v", err)
	}

	// Check that session was updated
	if manager.GetNetworkState() != network.InGame {
		t.Errorf("NetworkState = %v, want %v", manager.GetNetworkState(), network.InGame)
	}

	session := manager.GetSession()
	if session.State != AccountStateInGame {
		t.Errorf("session.State = %v, want %v", session.State, AccountStateInGame)
	}
}

func TestHandleActorMovementInterrupted(t *testing.T) {
	parser := NewCoreParser("ServerType0", hooks.NewHookManager())
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()
	manager := NewAccountManager(parser, hookManager, logger)

	// Set initial state
	manager.SetNetworkState(network.InGame)
	manager.SetState(AccountStateInGame)

	// Create test packet arguments
	args := map[string]interface{}{
		"ID":   uint32(12345),
		"posX": uint16(100),
		"posY": uint16(200),
	}

	// Call handler
	err := manager.handleActorMovementInterrupted(args)
	if err != nil {
		t.Fatalf("handleActorMovementInterrupted() returned error: %v", err)
	}

	// Check that session was updated
	if manager.GetNetworkState() != network.InGame {
		t.Errorf("NetworkState = %v, want %v", manager.GetNetworkState(), network.InGame)
	}

	session := manager.GetSession()
	if session.State != AccountStateInGame {
		t.Errorf("session.State = %v, want %v", session.State, AccountStateInGame)
	}
}

func TestHandleMapChange(t *testing.T) {
	parser := NewCoreParser("ServerType0", hooks.NewHookManager())
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()
	manager := NewAccountManager(parser, hookManager, logger)

	// Set initial state
	manager.SetNetworkState(network.InGame)
	manager.SetState(AccountStateInGame)

	// Create test packet arguments
	args := map[string]interface{}{
		"map": "new_map",
		"x":   uint16(100),
		"y":   uint16(200),
	}

	// Call handler
	err := manager.handleMapChange(args)
	if err != nil {
		t.Fatalf("handleMapChange() returned error: %v", err)
	}

	// Check that session was updated
	if manager.GetNetworkState() != network.InGame {
		t.Errorf("NetworkState = %v, want %v", manager.GetNetworkState(), network.InGame)
	}

	session := manager.GetSession()
	if session.MapName != "new_map" {
		t.Errorf("session.MapName = %v, want new_map", session.MapName)
	}
	if session.State != AccountStateInGame {
		t.Errorf("session.State = %v, want %v", session.State, AccountStateInGame)
	}
}

func TestHandleMapChanged(t *testing.T) {
	parser := NewCoreParser("ServerType0", hooks.NewHookManager())
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()
	manager := NewAccountManager(parser, hookManager, logger)

	// Set initial state
	manager.SetNetworkState(network.InGame)
	manager.SetState(AccountStateInGame)
	manager.session.MapName = "old_map"

	// Create test packet arguments
	args := map[string]interface{}{
		"map": "new_map",
		"x":   uint16(100),
		"y":   uint16(200),
	}

	// Call handler
	err := manager.handleMapChanged(args)
	if err != nil {
		t.Fatalf("handleMapChanged() returned error: %v", err)
	}

	// Check that session was updated
	if manager.GetNetworkState() != network.InGame {
		t.Errorf("NetworkState = %v, want %v", manager.GetNetworkState(), network.InGame)
	}

	session := manager.GetSession()
	if session.MapName != "new_map" {
		t.Errorf("session.MapName = %v, want new_map", session.MapName)
	}
	if session.State != AccountStateInGame {
		t.Errorf("session.State = %v, want %v", session.State, AccountStateInGame)
	}
}

func TestHandleQuitResponse(t *testing.T) {
	parser := NewCoreParser("ServerType0", hooks.NewHookManager())
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()
	manager := NewAccountManager(parser, hookManager, logger)

	// Test case 1: Successful disconnect (fail = 0)
	t.Run("SuccessfulDisconnect", func(t *testing.T) {
		// Set initial state
		manager.SetNetworkState(network.InGame)
		manager.SetState(AccountStateInGame)

		// Create test packet arguments
		args := map[string]interface{}{
			"fail": byte(0),
		}

		// Call handler
		err := manager.handleQuitResponse(args)
		if err != nil {
			t.Fatalf("handleQuitResponse() returned error: %v", err)
		}

		// Check that session was updated
		if manager.GetNetworkState() != network.NotConnected {
			t.Errorf("NetworkState = %v, want %v", manager.GetNetworkState(), network.NotConnected)
		}

		session := manager.GetSession()
		if session.State != AccountStateLoggedOut {
			t.Errorf("session.State = %v, want %v", session.State, AccountStateLoggedOut)
		}
	})

	// Test case 2: Failed disconnect (fail = 1)
	t.Run("FailedDisconnect", func(t *testing.T) {
		// Set initial state
		manager.SetNetworkState(network.InGame)
		manager.SetState(AccountStateInGame)

		// Create test packet arguments
		args := map[string]interface{}{
			"fail": byte(1),
		}

		// Call handler
		err := manager.handleQuitResponse(args)
		if err != nil {
			t.Fatalf("handleQuitResponse() returned error: %v", err)
		}

		// Check that session state remains unchanged
		if manager.GetNetworkState() != network.InGame {
			t.Errorf("NetworkState = %v, want %v", manager.GetNetworkState(), network.InGame)
		}

		session := manager.GetSession()
		if session.State != AccountStateInGame {
			t.Errorf("session.State = %v, want %v", session.State, AccountStateInGame)
		}
	})
}

func TestHandleSwitchCharacter(t *testing.T) {
	parser := NewCoreParser("ServerType0", hooks.NewHookManager())
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()
	manager := NewAccountManager(parser, hookManager, logger)

	// Test case 1: Valid switch character request (result = 1)
	t.Run("ValidSwitchCharacter", func(t *testing.T) {
		// Set initial state
		manager.SetNetworkState(network.InGame)
		manager.SetState(AccountStateInGame)

		// Add some test characters
		manager.session.Characters = []CharacterInfo{
			{
				CharID: 12345,
				Name:   "TestChar1",
			},
			{
				CharID: 67890,
				Name:   "TestChar2",
			},
		}
		manager.session.SelectedCharIndex = 0

		// Create test packet arguments
		args := map[string]interface{}{
			"result": byte(1),
		}

		// Call handler
		err := manager.handleSwitchCharacter(args)
		if err != nil {
			t.Fatalf("handleSwitchCharacter() returned error: %v", err)
		}

		// Check that session was updated
		if manager.GetNetworkState() != network.ConnectedToMasterServer {
			t.Errorf("NetworkState = %v, want %v", manager.GetNetworkState(), network.ConnectedToMasterServer)
		}

		session := manager.GetSession()
		if session.State != AccountStateLoggingIn {
			t.Errorf("session.State = %v, want %v", session.State, AccountStateLoggingIn)
		}

		// Check that character list was cleared
		if len(session.Characters) != 0 {
			t.Errorf("session.Characters length = %d, want 0", len(session.Characters))
		}

		if session.SelectedCharIndex != -1 {
			t.Errorf("session.SelectedCharIndex = %d, want -1", session.SelectedCharIndex)
		}
	})

	// Test case 2: Invalid switch character request (result != 1)
	t.Run("InvalidSwitchCharacter", func(t *testing.T) {
		// Set initial state
		manager.SetNetworkState(network.InGame)
		manager.SetState(AccountStateInGame)

		// Add some test characters
		manager.session.Characters = []CharacterInfo{
			{
				CharID: 12345,
				Name:   "TestChar1",
			},
			{
				CharID: 67890,
				Name:   "TestChar2",
			},
		}
		manager.session.SelectedCharIndex = 0

		// Create test packet arguments
		args := map[string]interface{}{
			"result": byte(0),
		}

		// Call handler
		err := manager.handleSwitchCharacter(args)
		if err != nil {
			t.Fatalf("handleSwitchCharacter() returned error: %v", err)
		}

		// Check that session state remains unchanged
		if manager.GetNetworkState() != network.InGame {
			t.Errorf("NetworkState = %v, want %v", manager.GetNetworkState(), network.InGame)
		}

		session := manager.GetSession()
		if session.State != AccountStateInGame {
			t.Errorf("session.State = %v, want %v", session.State, AccountStateInGame)
		}

		// Check that character list was not cleared
		if len(session.Characters) != 2 {
			t.Errorf("session.Characters length = %d, want 2", len(session.Characters))
		}

		if session.SelectedCharIndex != 0 {
			t.Errorf("session.SelectedCharIndex = %d, want 0", session.SelectedCharIndex)
		}
	})
}

func TestHandleCharacterDeletionSuccessful(t *testing.T) {
	parser := NewCoreParser("ServerType0", hooks.NewHookManager())
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()
	manager := NewAccountManager(parser, hookManager, logger)

	// Test case 1: Character deletion with selected index
	t.Run("DeletionWithSelectedIndex", func(t *testing.T) {
		// Set initial state
		manager.SetNetworkState(network.ConnectedToCharServer)
		manager.SetState(AccountStateSelectingChar)

		// Add some test characters
		manager.session.Characters = []CharacterInfo{
			{
				CharID: 12345,
				Name:   "TestChar1",
			},
			{
				CharID: 67890,
				Name:   "TestChar2",
			},
			{
				CharID: 54321,
				Name:   "TestChar3",
			},
		}
		manager.session.SelectedCharIndex = 1 // Select the second character for deletion

		// Create test packet arguments (empty for this packet)
		args := map[string]interface{}{}

		// Call handler
		err := manager.handleCharacterDeletionSuccessful(args)
		if err != nil {
			t.Fatalf("handleCharacterDeletionSuccessful() returned error: %v", err)
		}

		// Check that session was updated
		if manager.GetNetworkState() != network.ConnectedToCharServer {
			t.Errorf("NetworkState = %v, want %v", manager.GetNetworkState(), network.ConnectedToCharServer)
		}

		session := manager.GetSession()
		if session.State != AccountStateSelectingChar {
			t.Errorf("session.State = %v, want %v", session.State, AccountStateSelectingChar)
		}

		// Check that character was deleted
		if len(session.Characters) != 2 {
			t.Errorf("session.Characters length = %d, want 2", len(session.Characters))
		}

		// Check that the correct character was deleted
		for _, char := range session.Characters {
			if char.Name == "TestChar2" {
				t.Errorf("Character 'TestChar2' should have been deleted but still exists")
			}
		}

		// Check that selected index was reset
		if session.SelectedCharIndex != -1 {
			t.Errorf("session.SelectedCharIndex = %d, want -1", session.SelectedCharIndex)
		}
	})

	// Test case 2: Character deletion without selected index
	t.Run("DeletionWithoutSelectedIndex", func(t *testing.T) {
		// Set initial state
		manager.SetNetworkState(network.ConnectedToCharServer)
		manager.SetState(AccountStateSelectingChar)

		// Add some test characters
		manager.session.Characters = []CharacterInfo{
			{
				CharID: 12345,
				Name:   "TestChar1",
			},
			{
				CharID: 54321,
				Name:   "TestChar3",
			},
		}
		manager.session.SelectedCharIndex = -1 // No character selected

		// Create test packet arguments (empty for this packet)
		args := map[string]interface{}{}

		// Call handler
		err := manager.handleCharacterDeletionSuccessful(args)
		if err != nil {
			t.Fatalf("handleCharacterDeletionSuccessful() returned error: %v", err)
		}

		// Check that session was updated
		if manager.GetNetworkState() != network.ConnectedToCharServer {
			t.Errorf("NetworkState = %v, want %v", manager.GetNetworkState(), network.ConnectedToCharServer)
		}

		session := manager.GetSession()
		if session.State != AccountStateSelectingChar {
			t.Errorf("session.State = %v, want %v", session.State, AccountStateSelectingChar)
		}

		// Check that character list remains unchanged
		if len(session.Characters) != 2 {
			t.Errorf("session.Characters length = %d, want 2", len(session.Characters))
		}

		// Check that selected index remains unchanged
		if session.SelectedCharIndex != -1 {
			t.Errorf("session.SelectedCharIndex = %d, want -1", session.SelectedCharIndex)
		}
	})
}

func TestHandleCharacterDeletionFailed(t *testing.T) {
	parser := NewCoreParser("ServerType0", hooks.NewHookManager())
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()
	manager := NewAccountManager(parser, hookManager, logger)

	// Set initial state
	manager.SetNetworkState(network.ConnectedToCharServer)
	manager.SetState(AccountStateInGame) // Intentionally set to a different state to verify it changes

	// Add some test characters
	manager.session.Characters = []CharacterInfo{
		{
			CharID: 12345,
			Name:   "TestChar1",
		},
		{
			CharID: 67890,
			Name:   "TestChar2",
		},
	}
	manager.session.SelectedCharIndex = 1 // Select a character for deletion

	// Create test packet arguments (empty for this packet)
	args := map[string]interface{}{}

	// Call handler
	err := manager.handleCharacterDeletionFailed(args)
	if err != nil {
		t.Fatalf("handleCharacterDeletionFailed() returned error: %v", err)
	}

	// Check that session was updated
	if manager.GetNetworkState() != network.ConnectedToCharServer {
		t.Errorf("NetworkState = %v, want %v", manager.GetNetworkState(), network.ConnectedToCharServer)
	}

	session := manager.GetSession()
	if session.State != AccountStateSelectingChar {
		t.Errorf("session.State = %v, want %v", session.State, AccountStateSelectingChar)
	}

	// Check that character list remains unchanged
	if len(session.Characters) != 2 {
		t.Errorf("session.Characters length = %d, want 2", len(session.Characters))
	}

	// Check that selected index was reset
	if session.SelectedCharIndex != -1 {
		t.Errorf("session.SelectedCharIndex = %d, want -1", session.SelectedCharIndex)
	}
}

func TestHandleCharDelete2Result(t *testing.T) {
	parser := NewCoreParser("ServerType0", hooks.NewHookManager())
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()
	manager := NewAccountManager(parser, hookManager, logger)

	// Test case 1: Successful deletion request
	t.Run("SuccessfulDeletionRequest", func(t *testing.T) {
		// Set initial state
		manager.SetNetworkState(network.ConnectedToCharServer)
		manager.SetState(AccountStateInGame) // Intentionally set to a different state to verify it changes

		// Add some test characters
		manager.session.Characters = []CharacterInfo{
			{
				CharID: 12345,
				Name:   "TestChar1",
			},
			{
				CharID: 67890,
				Name:   "TestChar2",
			},
		}
		manager.session.SelectedCharIndex = 1 // Select a character for deletion

		// Create test packet arguments
		currentTime := uint32(time.Now().Add(24 * time.Hour).Unix()) // Delete date is 24 hours from now
		args := map[string]interface{}{
			"result":     byte(1),
			"deleteDate": currentTime,
		}

		// Call handler
		err := manager.handleCharDelete2Result(args)
		if err != nil {
			t.Fatalf("handleCharDelete2Result() returned error: %v", err)
		}

		// Check that session was updated
		if manager.GetNetworkState() != network.ConnectedToCharServer {
			t.Errorf("NetworkState = %v, want %v", manager.GetNetworkState(), network.ConnectedToCharServer)
		}

		session := manager.GetSession()
		if session.State != AccountStateSelectingChar {
			t.Errorf("session.State = %v, want %v", session.State, AccountStateSelectingChar)
		}

		// Check that character list remains unchanged
		if len(session.Characters) != 2 {
			t.Errorf("session.Characters length = %d, want 2", len(session.Characters))
		}

		// Check that the character is marked for deletion
		if !session.Characters[1].IsDeleting {
			t.Errorf("Character should be marked for deletion")
		}

		// Check that the delete date is set
		expectedTime := time.Unix(int64(currentTime), 0)
		if !session.Characters[1].DeleteDate.Equal(expectedTime) {
			t.Errorf("DeleteDate = %v, want %v", session.Characters[1].DeleteDate, expectedTime)
		}
	})

	// Test case 2: Failed deletion request (already planned to be erased)
	t.Run("FailedDeletionRequest", func(t *testing.T) {
		// Set initial state
		manager.SetNetworkState(network.ConnectedToCharServer)
		manager.SetState(AccountStateInGame) // Intentionally set to a different state to verify it changes

		// Add some test characters
		manager.session.Characters = []CharacterInfo{
			{
				CharID: 12345,
				Name:   "TestChar1",
			},
			{
				CharID: 67890,
				Name:   "TestChar2",
			},
		}
		manager.session.SelectedCharIndex = 1 // Select a character for deletion

		// Create test packet arguments
		args := map[string]interface{}{
			"result":     byte(0), // 0 = already planned to be erased
			"deleteDate": uint32(0),
		}

		// Call handler
		err := manager.handleCharDelete2Result(args)
		if err != nil {
			t.Fatalf("handleCharDelete2Result() returned error: %v", err)
		}

		// Check that session was updated
		if manager.GetNetworkState() != network.ConnectedToCharServer {
			t.Errorf("NetworkState = %v, want %v", manager.GetNetworkState(), network.ConnectedToCharServer)
		}

		session := manager.GetSession()
		if session.State != AccountStateSelectingChar {
			t.Errorf("session.State = %v, want %v", session.State, AccountStateSelectingChar)
		}

		// Check that character list remains unchanged
		if len(session.Characters) != 2 {
			t.Errorf("session.Characters length = %d, want 2", len(session.Characters))
		}

		// Check that the character is not marked for deletion
		if session.Characters[1].IsDeleting {
			t.Errorf("Character should not be marked for deletion")
		}

		// Check that the delete date is not set
		if !session.Characters[1].DeleteDate.IsZero() {
			t.Errorf("DeleteDate should be zero, got %v", session.Characters[1].DeleteDate)
		}
	})
}

func TestHandleCharDelete2AcceptResult(t *testing.T) {
	parser := NewCoreParser("ServerType0", hooks.NewHookManager())
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()
	manager := NewAccountManager(parser, hookManager, logger)

	// Test case 1: Successful deletion acceptance
	t.Run("SuccessfulDeletionAcceptance", func(t *testing.T) {
		// Set initial state
		manager.SetNetworkState(network.ConnectedToCharServer)
		manager.SetState(AccountStateInGame) // Intentionally set to a different state to verify it changes

		// Add some test characters
		manager.session.Characters = []CharacterInfo{
			{
				CharID: 12345,
				Name:   "TestChar1",
			},
			{
				CharID: 67890,
				Name:   "TestChar2",
			},
		}
		manager.session.SelectedCharIndex = 1 // Select a character for deletion

		// Create test packet arguments
		// Character ID 67890 in little-endian: 0x10932 = [0x32, 0x09, 0x01, 0x00]
		charID := []byte{0x32, 0x09, 0x01, 0x00}
		args := map[string]interface{}{
			"charID": charID,
			"result": byte(1), // 1 = success
		}

		// Call handler
		err := manager.handleCharDelete2AcceptResult(args)
		if err != nil {
			t.Fatalf("handleCharDelete2AcceptResult() returned error: %v", err)
		}

		// Check that session was updated
		if manager.GetNetworkState() != network.ConnectedToCharServer {
			t.Errorf("NetworkState = %v, want %v", manager.GetNetworkState(), network.ConnectedToCharServer)
		}

		session := manager.GetSession()
		if session.State != AccountStateSelectingChar {
			t.Errorf("session.State = %v, want %v", session.State, AccountStateSelectingChar)
		}

		// Check that character was deleted
		if len(session.Characters) != 1 {
			t.Errorf("session.Characters length = %d, want 1", len(session.Characters))
		}

		// Check that the correct character was deleted
		for _, char := range session.Characters {
			if char.Name == "TestChar2" {
				t.Errorf("Character 'TestChar2' should have been deleted but still exists")
			}
		}

		// Check that selected index was reset
		if session.SelectedCharIndex != -1 {
			t.Errorf("session.SelectedCharIndex = %d, want -1", session.SelectedCharIndex)
		}
	})

	// Test case 2: Failed deletion acceptance
	t.Run("FailedDeletionAcceptance", func(t *testing.T) {
		// Set initial state
		manager.SetNetworkState(network.ConnectedToCharServer)
		manager.SetState(AccountStateInGame) // Intentionally set to a different state to verify it changes

		// Add some test characters
		manager.session.Characters = []CharacterInfo{
			{
				CharID: 12345,
				Name:   "TestChar1",
			},
		}
		manager.session.SelectedCharIndex = 0 // Select a character for deletion

		// Create test packet arguments
		// Character ID 12345 in little-endian: 0x3039 = [0x39, 0x30, 0x00, 0x00]
		charID := []byte{0x39, 0x30, 0x00, 0x00}
		args := map[string]interface{}{
			"charID": charID,
			"result": byte(5), // 5 = birthday does not match
		}

		// Call handler
		err := manager.handleCharDelete2AcceptResult(args)
		if err != nil {
			t.Fatalf("handleCharDelete2AcceptResult() returned error: %v", err)
		}

		// Check that session was updated
		if manager.GetNetworkState() != network.ConnectedToCharServer {
			t.Errorf("NetworkState = %v, want %v", manager.GetNetworkState(), network.ConnectedToCharServer)
		}

		session := manager.GetSession()
		if session.State != AccountStateSelectingChar {
			t.Errorf("session.State = %v, want %v", session.State, AccountStateSelectingChar)
		}

		// Check that character list remains unchanged
		if len(session.Characters) != 1 {
			t.Errorf("session.Characters length = %d, want 1", len(session.Characters))
		}

		// Check that selected index was reset
		if session.SelectedCharIndex != -1 {
			t.Errorf("session.SelectedCharIndex = %d, want -1", session.SelectedCharIndex)
		}
	})
}

func TestHandleCharDelete2CancelResult(t *testing.T) {
	parser := NewCoreParser("ServerType0", hooks.NewHookManager())
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()
	manager := NewAccountManager(parser, hookManager, logger)

	// Test case 1: Successful deletion cancellation
	t.Run("SuccessfulDeletionCancellation", func(t *testing.T) {
		// Set initial state
		manager.SetNetworkState(network.ConnectedToCharServer)
		manager.SetState(AccountStateInGame) // Intentionally set to a different state to verify it changes

		// Add some test characters with one scheduled for deletion
		deleteDate := time.Now().Add(24 * time.Hour)
		manager.session.Characters = []CharacterInfo{
			{
				CharID: 12345,
				Name:   "TestChar1",
			},
			{
				CharID:     67890,
				Name:       "TestChar2",
				DeleteDate: deleteDate,
				IsDeleting: true,
			},
		}
		manager.session.SelectedCharIndex = 1 // Select the character scheduled for deletion

		// Create test packet arguments
		args := map[string]interface{}{
			"result": byte(1), // 1 = success
		}

		// Call handler
		err := manager.handleCharDelete2CancelResult(args)
		if err != nil {
			t.Fatalf("handleCharDelete2CancelResult() returned error: %v", err)
		}

		// Check that session was updated
		if manager.GetNetworkState() != network.ConnectedToCharServer {
			t.Errorf("NetworkState = %v, want %v", manager.GetNetworkState(), network.ConnectedToCharServer)
		}

		session := manager.GetSession()
		if session.State != AccountStateSelectingChar {
			t.Errorf("session.State = %v, want %v", session.State, AccountStateSelectingChar)
		}

		// Check that character list remains unchanged
		if len(session.Characters) != 2 {
			t.Errorf("session.Characters length = %d, want 2", len(session.Characters))
		}

		// Check that the character is no longer scheduled for deletion
		if session.Characters[1].IsDeleting {
			t.Errorf("Character should not be marked for deletion")
		}

		// Check that the delete date is cleared
		if !session.Characters[1].DeleteDate.IsZero() {
			t.Errorf("DeleteDate should be zero, got %v", session.Characters[1].DeleteDate)
		}

		// Check that selected index was reset
		if session.SelectedCharIndex != -1 {
			t.Errorf("session.SelectedCharIndex = %d, want -1", session.SelectedCharIndex)
		}
	})

	// Test case 2: Failed deletion cancellation
	t.Run("FailedDeletionCancellation", func(t *testing.T) {
		// Set initial state
		manager.SetNetworkState(network.ConnectedToCharServer)
		manager.SetState(AccountStateInGame) // Intentionally set to a different state to verify it changes

		// Add some test characters with one scheduled for deletion
		deleteDate := time.Now().Add(24 * time.Hour)
		manager.session.Characters = []CharacterInfo{
			{
				CharID: 12345,
				Name:   "TestChar1",
			},
			{
				CharID:     67890,
				Name:       "TestChar2",
				DeleteDate: deleteDate,
				IsDeleting: true,
			},
		}
		manager.session.SelectedCharIndex = 1 // Select the character scheduled for deletion

		// Create test packet arguments
		args := map[string]interface{}{
			"result": byte(2), // 2 = database error
		}

		// Call handler
		err := manager.handleCharDelete2CancelResult(args)
		if err != nil {
			t.Fatalf("handleCharDelete2CancelResult() returned error: %v", err)
		}

		// Check that session was updated
		if manager.GetNetworkState() != network.ConnectedToCharServer {
			t.Errorf("NetworkState = %v, want %v", manager.GetNetworkState(), network.ConnectedToCharServer)
		}

		session := manager.GetSession()
		if session.State != AccountStateSelectingChar {
			t.Errorf("session.State = %v, want %v", session.State, AccountStateSelectingChar)
		}

		// Check that character list remains unchanged
		if len(session.Characters) != 2 {
			t.Errorf("session.Characters length = %d, want 2", len(session.Characters))
		}

		// Check that the character is still scheduled for deletion
		if !session.Characters[1].IsDeleting {
			t.Errorf("Character should still be marked for deletion")
		}

		// Check that the delete date is still set
		if session.Characters[1].DeleteDate != deleteDate {
			t.Errorf("DeleteDate = %v, want %v", session.Characters[1].DeleteDate, deleteDate)
		}

		// Check that selected index was reset
		if session.SelectedCharIndex != -1 {
			t.Errorf("session.SelectedCharIndex = %d, want -1", session.SelectedCharIndex)
		}
	})
}

func TestHandleSyncRequestEx(t *testing.T) {
	parser := NewCoreParser("ServerType0", hooks.NewHookManager())
	hookManager := hooks.NewHookManager()
	logger := NewMockLogger()
	manager := NewAccountManager(parser, hookManager, logger)

	// Set initial state
	manager.SetNetworkState(network.InGame)
	manager.SetState(AccountStateInGame)

	// Create test packet arguments with a switch ID
	args := map[string]interface{}{
		"switch": "07E5", // Example packet ID
	}

	// Set up the sync_ex_reply table
	manager.syncExReplyTable = map[string]string{
		"07E5": "07E6", // Example mapping
	}

	// Call handler
	err := manager.handleSyncRequestEx(args)
	if err != nil {
		t.Fatalf("handleSyncRequestEx() returned error: %v", err)
	}

	// Check that session was updated
	if manager.GetNetworkState() != network.InGame {
		t.Errorf("NetworkState = %v, want %v", manager.GetNetworkState(), network.InGame)
	}

	session := manager.GetSession()
	if session.State != AccountStateInGame {
		t.Errorf("session.State = %v, want %v", session.State, AccountStateInGame)
	}
}
